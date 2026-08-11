# Marketing-Data 架构梳理

> 可视化图解：[Marketing-Data 架构图解](./marketing-data-architecture.html)
>
> 所属项目索引：[Marketing-Data 项目材料](./README.md)
>
> Marketing 主项目：[Marketing 项目材料](../Marketing/README.md)

## 专项补充文档

如果只想背深挖点，可以直接看这三篇：

- [Marketing-Data 核心知识深挖](./marketing-data-key-knowledge-deep-dive.md)：OfflineDataFetchAndPush、batch history、服务边界和深挖速答。
- [Marketing-Data Handler / pre-filter 深挖](./marketing-data-handler-prefilter-deep-dive.md)：DataSource handler、ES/File/Hive/Insight、pre-filter 实时预处理。
- [Marketing-Data 可靠性和一致性深挖](./marketing-data-reliability-deep-dive.md)：batch 幂等、history、Kafka publish 风险、排障和不能夸大的点。

## 一、Marketing-Data 是什么？一句话定位

**Marketing-Data 是 Marketing 的离线/批量数据服务**。Marketing Engine 负责计划执行，Marketing-Data 负责离线人群拉取、ES 查询预览、S3/CSV/Hive/Insight 文件处理、实时消息预过滤、batch history 和 Kafka 推送。

面试开场：

> Marketing-Data 是从 Marketing 主服务拆出来的数据抽取服务。拆分原因是运营数据提取功能相对稳定，但执行期间不适合被频繁发版打断；同时大查询、文件处理、Hive/Insight 产物读取和 batch history 更适合独立维护。它通过 `OfflineDataFetchAndPush` 接收 Marketing 的计划批次，做参数校验和 `batch_id` 幂等，写 `offline_plan_history_tab`，再异步选择 DataSource handler 拉数、过滤、去重并推 Kafka，最后更新执行结果和监控。

个人边界：

> Marketing-Data 是组内 Marketing 域的一部分。我能讲清它和 Marketing Engine 的边界、离线抽取主链路和排障方式；个人具体负责的 PR/Jira 需要按事实补充。

## 项目路径速查

| 项目 | 本机路径 | 这里看什么 |
| --- | --- | --- |
| `insurance-marketing-data` | `/Users/si.chen/GolandProjects/insurance-marketing-data` | Marketing-Data 主服务，离线拉数、pre-filter、history、MQ 都在这里。 |
| `insurance-marketing-data-api` | `/Users/si.chen/GolandProjects/insurance-marketing-data-api` | Internal/Admin API proto。 |
| `insurance-marketing` | `/Users/si.chen/GolandProjects/insurance-marketing` | Marketing 主引擎，消费数据后执行 Handler。 |

常用代码入口：

| 想看什么 | 推荐路径 |
| --- | --- |
| `OfflineDataFetchAndPush` 参数校验 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/service/marketing_data_internal_service.go` |
| `QueryEsList` Admin 预览 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/service/marketing_data_admin_service.go` |
| 离线拉数编排 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/manager/impl/offline_data_fetch_manager_impl.go` |
| DataSource handler 接口 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/handler/offline_data_fetch_handler.go` |
| ES/File/Hive/Insight 实现 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/handler/offline_data_fetch/impl/` |
| ES Scroll 查询 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/repo/elasticsearch/impl/common_info_impl_es.go` |
| Kafka producer | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/mq/producer/impl/producer_impl.go` |
| pre-filter consumer | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/mq/consumer/impl/consumer_impl.go` |
| Prometheus 监控 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/monitor/register.go` |

## 二、为什么从 Marketing 拆出来

Confluence 技术设计里给出的背景是：运营业务的数据提取功能相对稳定，任务执行期间不宜被打断，重试可能触发数据重新提取，所以从 `insurance-marketing` 中抽出来独立维护。

拆分后的目标：

| 目标 | 面试解释 |
| --- | --- |
| 离线/实时数据提取独立维护 | 数据抽取和营销动作执行的变化频率不同，拆开能降低主引擎耦合。 |
| 离线执行结果可查看 | 通过 `offline_plan_history_tab` 记录状态、耗时、总数、去重数、错误详情。 |
| go-schedule 切换 | Marketing 原来的定时任务切到 Marketing-Data，避免主服务承担重数据处理。 |
| pre-filter 独立 | 实时消息预过滤规则迁移到 Marketing-Data，减少 Marketing 主服务负担。 |

一句话：

> Marketing 管“什么时候对哪些人做什么动作”，Marketing-Data 管“这些人从哪里来、怎么批量抽取、怎么推给 Marketing”。

## 三、核心模型：离线抽取请求

`OfflineDataFetchAndPushReq` 可以抽象成：

```text
OfflineDataFetchAndPushReq
├── plan_id          -> 哪个 Marketing Plan
├── batch_id         -> 本次执行批次，幂等键
├── dest_topic       -> 推给 Marketing 的 Kafka topic
├── data_source      -> es / file / hive / insight 等
├── group_condition  -> 数据源查询/文件配置/过滤条件 JSON
├── user_group_id    -> Marketing 用户组
└── estimate         -> 预估数量、耗时、owner，用于监控规则
```

必填校验：

- `plan_id > 0`
- `batch_id > 0`
- `dest_topic` 非空
- `data_source` 非空
- `group_condition` 非空
- `user_group_id > 0`

## 四、完整数据流：五阶段流水线

```text
┌──────────────┐    ┌──────────────┐    ┌────────────────┐    ┌─────────────┐    ┌─────────────┐
│ Marketing    │ -> │ Service 校验 │ -> │ Manager 幂等   │ -> │ Handler 拉数 │ -> │ Kafka/History│
│ Plan/Schedule│    │ 参数/API     │    │ history/async  │    │ 过滤/去重    │    │ monitor      │
└──────────────┘    └──────────────┘    └────────────────┘    └─────────────┘    └─────────────┘
```

### Stage 1: Service 参数校验

Internal API：

- `OfflineDataFetchAndPush`
- `QueryOfflinePlanHistory`

Admin API：

- `QueryEsList`

`QueryEsList` 校验 `index`、`dsl`、`size`，且 `size <= 100`。它是 Admin 预览接口，不是大批量抽取入口。

### Stage 2: Manager 做 batch 幂等和 history 初始化

`FetchAndPushData` 的关键逻辑：

```text
1. 用 batch_id 查询 offline_plan_history_tab
2. 如果已有记录，直接 return
3. 初始化 PlanExecutionDetail，状态 init
4. 写 offline_plan_history_tab
5. 异步执行 handler
6. 根据结果更新 success/failed、耗时、统计、错误
7. 上报 Prometheus
```

面试抓手：

> `batch_id` 是离线任务的幂等键。调度重试、RPC 超时重试或人工重复触发时，已执行过的 batch 不应该重复拉取和重复推 Kafka，否则可能造成重复触达。

### Stage 3: DataSource handler 自注册

统一接口：

```go
type OfflineDataFetchHandler interface {
    BeforeProcessData()
    ProcessData()
    AfterProcessData()
    GetStaticDetail() *model_ext.OffLinePlanStatistic
}
```

注册方式：

```text
handler.InjectInitHandlerFunc(dataSource, initFunc)
```

已有实现：

| data source | 代表文件 | 主要职责 |
| --- | --- | --- |
| ES | `es_data_fetch_impl.go` | 解析 DSL，scroll 拉取 ES，按 key 去重，推 Kafka。 |
| File | `file_data_fetch_impl.go` | 读取 S3/CSV 文件，按 mapping 转字段，规则过滤，日期过滤，去重，推 Kafka。 |
| Hive File | `hive_file_data_fetch_impl.go` | 处理 Hive 产出的文件，通常跳过表头并校验文件。 |
| Insight File | `insight_file_data_fetch_impl.go` | 读取 Insight 推送文件，并校验推送时间窗口。 |

设计亮点：

> 新增一种数据源时，实现统一 handler 接口并注册 data_source，主流程不需要改。

### Stage 4: ES / File 处理细节

ES 处理：

```text
group_condition -> FeatureSearchVo / DSL
  -> custom tag handler 补充动态条件
  -> QueryTotalCount 做 max length 检查
  -> QueryWithScroll 分批拉取
  -> deduplicationKeys + BloomFilter 去重
  -> FormatOfflineRawResult
  -> PushWithTopic
```

File 处理：

```text
group_condition -> FileConditionDetail
  -> S3 获取文件对象
  -> 本地临时文件读取 CSV
  -> mapping 映射字段
  -> 正则规则过滤
  -> DateFilters 日期窗口过滤
  -> AppendData 补充字段
  -> BloomFilter 去重
  -> PushWithTopic
  -> recordFilePath 防重复处理
```

### Stage 5: 结果、监控和查询

history 记录：

```text
offline_plan_history_tab
├── plan_id
├── batch_id
├── result_status: init / success / failed
├── execution_second
├── detail: total / filter_total / unique_total / err_detail
├── gmt_start_time
└── gmt_end_time
```

监控：

- `extract_length_counter`
- `plan_batch_extract_total_insurance`

监控 labels 包括：

- namespace、service
- planId、batchId
- result、count、duration
- minCount、maxCount、maxDuration、owner、ruleCode

面试说法：

> Marketing-Data 不只是把数据推到 Kafka，还会把每个 batch 的结果落 history 并上报监控。这样 Marketing 或排障人员可以通过 `batch_id` 查到一次离线抽取到底成功、失败、耗时多久、抽到多少人、错误是什么。

## 五、实时 pre-filter 链路

Marketing-Data 不只处理离线批量，也承接实时消息预过滤：

```text
source topic
  -> reliable consumer
  -> pre_filter_tab 规则缓存
  -> RegexpMatch 匹配消息
  -> FillVariable 组装目标消息
  -> 补 product_type / user_id
  -> PushWithTopic 到 dest topic
```

关键表：

- `pre_filter_tab`
- `reliable_consumer_handler_tab`
- `reliable_event_handler_tab`

面试说法：

> pre-filter 的作用是把外部实时消息先按规则过滤和字段补齐，再转发给 Marketing 使用。这样 Marketing 主引擎只消费已经整理过的营销消息。

## 六、目录结构

主服务路径：`/Users/si.chen/GolandProjects/insurance-marketing-data`

```text
src/
├── service/                       # Internal/Admin service 入口
├── manager/                       # OfflineDataFetchManager 编排
├── handler/                       # DataSource handler 接口和实现
│   └── offline_data_fetch/impl/   # es / file / hive / insight
├── biz/                           # plan/pre-filter/product/account 等业务能力
├── repo/                          # DB repo、ES repo、S3 repo
├── mq/                            # producer / consumer
├── monitor/                       # Prometheus + RC 上报
├── model/ model_ext/              # history、plan message、condition DTO
├── task/                          # pre-filter cache update task
├── event/                         # reliable-event 初始化
└── config/ constant/ util/        # 配置、常量、工具
```

## 七、设计哲学

### 1. 重数据处理从主引擎拆出

离线人群、大 ES 查询、文件读取、Hive/Insight 产物都可能耗时长、失败模式复杂。拆出 Marketing-Data 后，Marketing Engine 更专注于计划执行。

### 2. batch 幂等优先

离线任务重复执行的后果是重复推送和重复触达。`batch_id` + history 查询 + DB 唯一约束是关键防线。

### 3. 数据源插件化

ES、File、Hive、Insight 都实现同一套 handler 生命周期。新增数据源时扩展 handler，不改主流程。

### 4. 执行结果可观测

异步任务必须可查结果。history + monitor 让排障从 `batch_id` 开始，而不是只查日志。

### 5. Admin 预览和批量抽取分离

`QueryEsList` 限制 `size <= 100`，适合页面预览。大批量应该走异步抽取链路。

## 八、面试叙事顺序

1. **定位**：Marketing-Data 是 Marketing 的离线/批量数据支撑服务。
2. **为什么拆**：数据提取稳定但任务不宜被主服务发版打断，大查询和文件处理适合独立。
3. **主链路**：`OfflineDataFetchAndPush` -> 校验 -> `batch_id` 幂等 -> history -> async -> handler -> Kafka -> history/monitor。
4. **handler 设计**：ES/File/Hive/Insight 自注册，统一 Before/Process/After。
5. **可靠性**：batch history、唯一索引、Prometheus、错误详情、max length 风控。
6. **边界**：Marketing-Data 管人群数据，Marketing Engine 管具体营销动作。
7. **排障**：用 batch_id 查 history，再按 data_source 分层排查 ES/S3/CSV/Hive/Kafka。

## 九、常见追问

### Marketing-Data 和 Marketing 的边界是什么？

Marketing 是计划执行引擎，负责触发、筛选、Handler 执行。Marketing-Data 是数据支撑服务，负责离线人群和预处理消息。不要把 PNAR/发券/通知说成 Marketing-Data 执行。

### `batch_id` 为什么重要？

它是同一计划批次的幂等键。离线任务可能重试或重复触发，`batch_id` 已存在就直接返回，避免重复抽取、重复推 Kafka 和重复触达。

### 新增 DataSource 怎么做？

实现 `OfflineDataFetchHandler`，在 `init()` 中注册 `data_source -> initFunc`，主流程根据 `req.DataSource` 找 handler。新增逻辑集中在 handler 内部。

### ES 查询失败怎么排？

先看 group_condition 和 DSL 转换，再看 custom tag handler、index/mapping、ES scroll、max length 检查和异常码。Admin 预览和离线批量不是同一个入口。

### 文件数据失败怎么排？

看 S3 session、文件路径、文件完整性、CSV 字段 mapping、日期过滤、路径重复记录和具体异常。文件处理后要记录路径，避免重复处理同一个文件。

### 推 Kafka 失败会怎样？

当前 producer 里单条 `Publish` 失败会记录 error，不直接让整个任务失败。面试时可以说这点需要结合业务可接受性看：如果要求强一致，需要补充失败重试或 DLQ 机制；当前文档不夸大为“每条都强保证”。

## 十、可讲项目经验模板

稳妥版：

> 我参与/熟悉 Marketing 离线人群链路。Marketing Engine 生成 plan batch 后，把 `plan_id`、`batch_id`、`data_source`、`group_condition`、`dest_topic` 传给 Marketing-Data；Marketing-Data 用 `batch_id` 做幂等，先写 history，再异步选择 ES/File/Hive/Insight handler 拉数，处理后推 Kafka 给 Marketing，最后更新 success/failed、数量、耗时和错误详情。这个设计把大查询和文件处理从主执行引擎里拆出来，也让离线任务可追踪。

具体个人项目：`待补充`。

## 十一、优化项重点记忆

这些点来自 Confluence 技术设计和本地代码，适合面试里回答“为什么拆服务 / 怎么保证离线任务稳定 / 怎么优化大数据抽取”。

| 优化项 | 解决的问题 | 面试怎么讲 |
| --- | --- | --- |
| 从 Marketing 拆出 Marketing-Data | 运营数据提取功能稳定，但任务执行期间不宜被频繁发版打断；重试会重新提取数据。 | 把重数据处理从主执行引擎拆开，Marketing 专注计划执行，Marketing-Data 专注数据抽取。 |
| go-schedule 切换到 Marketing-Data | 定时离线人群任务如果留在 Marketing，会增加主服务负担。 | 定时任务和预处理迁移到 Marketing-Data，降低 Marketing 主链路压力。 |
| `batch_id` 幂等 + 唯一索引 | 调度/RPC/人工重试可能导致同一批次重复抽取和重复推送。 | 先查 history，存在直接 return；DB 唯一约束兜底并发重复插入。 |
| 先写 history，再异步执行 | 异步任务如果只打日志，失败后难追踪。 | 初始化 `offline_plan_history_tab`，执行后更新 success/failed、数量、耗时和错误。 |
| DataSource handler 插件化 | ES、文件、Hive、Insight 逻辑差异大，硬写分支难维护。 | 统一 Before/Process/After 生命周期，新增数据源只新增 handler 并注册。 |
| `QueryEsList size <= 100` | Admin 页面预览如果允许大查询，会拖慢 ES 和服务。 | 预览接口限量，大批量抽取必须走异步离线链路。 |
| BloomFilter 去重 | 同一人群可能从 ES/文件中重复出现，重复触达风险高。 | 按 `deduplicationKeys` 做内存去重，统计 total 和 unique_total。 |
| MaxLength + RC/Prometheus | 人群过大可能造成触达风险和系统压力。 | handler 处理前可做最大数量检查，超限上报并阻断。 |
| pre-filter 迁移和开关 | 实时消息预处理如果和主引擎耦合，升级风险高。 | pre-filter 规则迁移到 Marketing-Data，并支持配置开关便于切换和回退。 |

风险提醒：

> 当前 producer 单条 Kafka publish 失败会记录 error，但不会直接让整个离线任务失败。面试里不要夸成“每条消息强一致保证”；可以说这是后续可优化点，需要结合 DLQ、重试或失败落表来增强。

一句话总结：

> Marketing-Data 的优化主线是把“重、慢、可重试、可观测”的数据抽取能力从 Marketing 主执行链路里拆出来，并用 batch 幂等、history、handler 插件化和监控降低重复触达与排障成本。

## 十二、资料来源

- Confluence：`[v1.2.X]insurance-marketing-data 技术设计`，页面 ID `1613567933`。
- 代码：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/service/marketing_data_internal_service.go`
- 代码：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/manager/impl/offline_data_fetch_manager_impl.go`
- 代码：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/handler/offline_data_fetch_handler.go`
- 代码：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/handler/offline_data_fetch/impl/es_data_fetch_impl.go`
- 代码：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/handler/offline_data_fetch/impl/file_data_fetch_impl.go`
- 代码：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/mq/consumer/impl/consumer_impl.go`
- 代码：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/mq/producer/impl/producer_impl.go`
- 代码：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/monitor/register.go`
