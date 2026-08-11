# Marketing-Data 面试专项笔记

## 面试定位

Marketing-Data 是 Insurance Marketing 的离线数据服务。它从 Marketing 主服务里拆出来，专门处理离线人群拉取、ES 查询预览、文件/Hive/Insight 数据源、batch history 和 Kafka 推送。

一句话：

> Marketing 负责计划执行，Marketing-Data 负责离线数据抽取和推送；它把重查询、文件处理、批次状态和数据源差异从主执行引擎里拆出来。

更稳的个人口径：

> 这是我们组 Marketing 域的一部分。我熟悉它和 Marketing Engine 的边界、`OfflineDataFetchAndPush` 主链路、batch 幂等、DataSource handler 插件化、ES/S3/CSV/Hive 数据处理和执行历史查询。个人具体改动需要补 PR/Jira。

个人证据：`待补充`。

## 1 分钟回答

Marketing-Data 是 Marketing 的离线数据支撑服务。Marketing 主服务负责运营计划执行，比如 Trigger、Processor、Handler、发券、通知、PNAR；但定时离线人群或大批量数据抽取不适合放在主执行链路里，所以拆出 Marketing-Data。

主入口是 `OfflineDataFetchAndPush`。Marketing 会把 `plan_id`、`batch_id`、`dest_topic`、`data_source`、`group_condition`、`user_group_id` 等信息传过来。Marketing-Data 先做参数校验和 batch 幂等检查，如果这个 `batch_id` 已经有执行记录就直接返回；否则先写 `offline_plan_history_tab` 初始化记录，再异步选择对应 DataSource handler 执行。handler 会做 `BeforeProcessData`、`ProcessData`、`AfterProcessData`，从 ES、文件、Hive 或 Insight 等数据源拉数据，处理后推到 Kafka，最后更新 batch history 和监控。

所以面试里可以讲它的价值：隔离重数据处理、支持多数据源扩展、保证批次幂等和执行可追踪，让 Marketing Engine 更专注于计划编排和业务动作执行。

## Marketing 与 Marketing-Data 的区别

| 维度 | Marketing | Marketing-Data |
| --- | --- | --- |
| 定位 | 运营计划主执行引擎 | 离线数据抽取和数据预处理服务 |
| 核心问题 | 什么时间、什么人群、执行什么营销动作 | 从哪里拉人群、怎么批量处理、怎么推给 Marketing |
| 典型模块 | Engine、RuleCenter、Processor、Handler、Plan、Notification、Voucher、PNAR | `OfflineDataFetchAndPush`、`QueryOfflinePlanHistory`、`QueryEsList`、DataSource handler、offline history |
| 主要风险 | 重复触达、规则漏过滤、handler 幂等、活动状态 | 大查询超时、文件异常、重复推送、batch 状态丢失、数据源扩展 |
| 面试讲法 | 配置化执行流水线 | 异步离线数据管道 |

## 主链路怎么讲

```text
Marketing Scheduler / Plan Engine
  -> Marketing-Data.InternalService.OfflineDataFetchAndPush
  -> 参数校验 plan_id / batch_id / dest_topic / data_source / group_condition / user_group_id
  -> batch_id 幂等检查
  -> 写 offline_plan_history_tab，状态 init
  -> 异步执行
  -> 根据 data_source 找 handler
  -> BeforeProcessData
  -> ProcessData: ES / File / Hive / Insight 拉取、过滤、去重、推 Kafka
  -> AfterProcessData: 文件路径记录等收尾
  -> 更新 history 为 success / failed
  -> 上报 plan_batch_extract_total_insurance 监控
  -> Marketing 消费数据并继续执行营销动作
```

## 核心模块

| 模块 | 作用 | 代码证据 |
| --- | --- | --- |
| Internal Service | `OfflineDataFetchAndPush` 和 `QueryOfflinePlanHistory`，负责内部调用入口和参数校验。 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/service/marketing_data_internal_service.go` |
| Admin Service | `QueryEsList`，用于 Admin 侧 ES 条件预览，限制 `size <= 100`。 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/service/marketing_data_admin_service.go` |
| OfflineDataFetchManager | batch 幂等、初始化 history、异步执行 handler、更新结果、上报监控。 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/manager/impl/offline_data_fetch_manager_impl.go` |
| DataSource Handler | `BeforeProcessData`、`ProcessData`、`AfterProcessData`、`GetStaticDetail` 四个阶段，按 data_source 自注册。 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/handler/offline_data_fetch_handler.go` |
| Handler 实现 | ES、File、Hive File、Insight File 等数据源实现。 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/handler/offline_data_fetch/impl/` |
| MQ Producer | handler 通过 `PushWithTopic` 把抽取结果推到目标 topic。 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/mq/producer/` |
| ES Repo | 查询 ES，支持按 size 查询和 scroll batch。 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/repo/elasticsearch/` |
| History Repo | 读写 `offline_plan_history_tab`，支持按 batch_id 查询和按 plan_id 查询最新 batch。 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/repo/repo_dao/impl/offline_plan_history_repo_dao_impl.go` |
| Monitor | 上报 plan/batch 抽取结果、数量、耗时和 estimate 信息。 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/monitor/register.go`、`/Users/si.chen/GolandProjects/insurance-marketing-data/src/monitor/report.go` |

## 关键设计点

### 1. batch_id 幂等

代码里 `FetchAndPushData` 会先用 `batch_id` 查 `offline_plan_history_tab`：

```text
如果 batch_id 已存在 -> 直接 return
如果不存在 -> 创建 init history -> 异步执行
```

面试说法：

> 离线抽取可能因为超时、重试或调度重复触发被调用多次。`batch_id` 幂等可以避免同一批数据重复抽取、重复推送 Kafka，降低重复触达用户的风险。

注意：

- 代码里也有 SQL 增加 `uk_batch_id` 的痕迹，说明 batch_id 是重要唯一键。
- 幂等不是只靠代码判断，DB 唯一约束能兜底并发插入。

### 2. 先写 history，再异步执行

执行前先写 `offline_plan_history_tab`，初始状态是 `init`；异步处理结束后更新为 success 或 failed，并保存统计和错误详情。

面试说法：

> 这样做是为了让异步任务可观测。即使任务后面失败，也能通过 batch history 查到开始时间、结束时间、状态、总数、去重数和错误信息。

### 3. DataSource handler 插件化

`handler.InjectInitHandlerFunc(dataSource, function)` 把每种数据源注册到 map，Manager 根据 `req.DataSource` 找初始化函数。

面试说法：

> 新增数据源时，不需要改主流程，只要实现统一 handler 接口并注册 data_source。主流程只关心 Before/Process/After 三阶段和统计结果。

已有实现：

- `es_data_fetch_impl.go`
- `file_data_fetch_impl.go`
- `hive_file_data_fetch_impl.go`
- `insight_file_data_fetch_impl.go`

### 4. QueryEsList 为什么限制 size

Admin API `QueryEsList` 会校验 `index`、`dsl`、`size`，并限制 `size <= 100`。

面试说法：

> 这个接口主要用于 Admin 侧预览查询结果，不是大批量抽取入口。限制 size 可以避免用户在页面上触发大查询，影响 ES 和服务稳定性。真正离线抽取应该走 `OfflineDataFetchAndPush` 的异步链路。

### 5. ES 不是事实源

Marketing-Data 可能从 ES 查人群，也可能从文件、Hive、Insight 拉数据。ES 更适合查询视图，不应被说成唯一事实源。

面试说法：

> 如果人群数据查不到，我会先确认数据源类型，再分别看 owner DB、Canal/ES 同步、文件路径/S3、Hive/Insight 产物和 DSL 转换，而不是只看 ES。

## 常见追问

### 为什么 Marketing 要拆出 Marketing-Data？

短答：

> 因为运营计划执行和离线数据抽取的变化频率、资源消耗和失败模式不一样。Marketing 负责业务执行，Marketing-Data 负责重数据处理，拆开后主链路更稳定，也方便单独扩展和排障。

展开：

- 离线人群抽取可能很重，容易超时或占资源。
- 文件/Hive/Insight/ES 的异常和 Marketing Handler 的异常不是一类问题。
- batch history 和数据抽取监控适合独立沉淀。
- 新增数据源时，用 handler 扩展，不影响 Marketing 主引擎。

### `OfflineDataFetchAndPush` 具体做什么？

短答：

> 它接收 Marketing 传来的计划批次和数据源条件，做参数校验、batch 幂等、初始化 history，然后异步拉取数据并推送到目标 topic，最后更新结果和监控。

展开：

- 必填字段：`plan_id`、`batch_id`、`dest_topic`、`data_source`、`group_condition`、`user_group_id`。
- 幂等：`batch_id` 已存在则直接返回。
- 执行：根据 `data_source` 找 handler。
- 结果：写 `total`、`unique_total`、`err_detail`、开始/结束时间、耗时。

### 如果同一个 batch 被重复调用怎么办？

短答：

> 先按 `batch_id` 查询 history，存在就直接返回；DB 层还有唯一索引可以兜底并发重复插入。

展开：

- 解决调度重试、RPC 重试、人工重复触发。
- 避免重复推 Kafka。
- 方便用 batch_id 查执行结果。

### 新增一种数据源怎么做？

短答：

> 实现 `OfflineDataFetchHandler` 接口，然后用 `InjectInitHandlerFunc(dataSource, initFunc)` 注册。主流程不用改，只根据 `data_source` 找 handler。

展开：

- `BeforeProcessData` 做参数和数据源前置检查。
- `ProcessData` 做实际拉取、转换、过滤、推送。
- `AfterProcessData` 做收尾。
- `GetStaticDetail` 返回统计信息。

### 大批量 ES 查询怎么处理？

短答：

> Admin 预览用 `QueryEsList`，限制 size；批量抽取走异步链路。大结果集要看 scroll batch size、DSL 转换、超时、ES mapping 和错误监控。

展开：

- `QueryEsList` 不适合大批量。
- ES handler/Repo 负责真正数据抽取。
- 如果 ES 查询失败，先看 DSL 转换、自定义 tag handler、mapping 类型和 ES 错误。

### 怎么排查离线任务失败？

短答：

> 用 `batch_id` 查 `QueryOfflinePlanHistory`，看状态、错误、总数、去重数和执行时间；再根据 `data_source` 看 ES/S3/CSV/Hive/Insight 具体错误。

展开：

- 参数错误：看 service 校验。
- unsupported data_source：看 handler 是否注册。
- ES 错误：看 DSL、index、mapping、scroll。
- 文件错误：看 S3 session、路径重复、CSV 字段、文件完整性。
- 推送/消费问题：看 dest_topic、Kafka、Marketing 消费端。

### Marketing-Data 会直接执行 PNAR 吗？

短答：

> 不建议这么说。PNAR 或通知动作属于 Marketing Handler 的执行范围，Marketing-Data 主要提供离线人群数据，把抽取结果推给 Marketing。

展开：

- Marketing-Data 解决“人群从哪里来”。
- Marketing Engine 解决“拿到人群后执行什么动作”。
- 面试时要把数据管道和业务动作分开。

## 可以包装成项目经历的说法

### 稳妥版

> 我参与/熟悉 Marketing 离线人群链路。这个链路里，Marketing Engine 生成计划批次，把数据源、用户组条件和目标 topic 传给 Marketing-Data；Marketing-Data 负责 batch 幂等、异步抽取、数据源 handler 处理、Kafka 推送和历史记录。这个设计把重数据处理从主引擎拆开，降低了营销执行链路的耦合和失败影响。

具体项目名、PR、Jira：`待补充`。

### STAR 模板

背景：

> 运营计划需要定时触达某类用户，例如保单到期提醒、PNAR、发券或通知场景。人群可能来自 ES、文件或离线平台，数据量比较大。

任务：

> 我需要支持离线人群抽取链路，保证同一批计划不会重复抽取/重复推送，并且任务执行结果可查询。

行动：

> 我先确认 Marketing 传入的 `plan_id`、`batch_id`、`data_source`、`group_condition` 和 `dest_topic`；在 Marketing-Data 侧通过 `batch_id` 做幂等，先写 batch history，再异步执行对应 DataSource handler；处理完成后更新 total、unique_total、err_detail 和执行时间，并上报监控。

结果：

> 让离线人群抽取从主执行链路中解耦出来，Marketing 可以通过 batch history 判断执行结果，并降低重复触达风险。具体效果和指标：`待补充`。

## 风险边界

| 不建议说 | 风险 | 推荐替换 |
| --- | --- | --- |
| Marketing-Data 是营销主引擎 | 它主要是离线数据服务 | Marketing 是执行引擎，Marketing-Data 是离线数据支撑 |
| Marketing-Data 直接执行发券/通知/PNAR | 容易混淆 Handler 所属服务 | Marketing-Data 抽取人群并推给 Marketing，具体营销动作由 Marketing 执行 |
| ES 是唯一事实源 | ES 是查询视图，可能延迟 | 根据 data_source 判断事实源；ES 场景也要回查 owner DB/Canal |
| 我个人 owner 了 Marketing-Data 全部 | 面试官会追 PR、事故和模块 owner | 我参与/熟悉 Marketing-Data 离线数据链路，个人具体任务按事实补充 |

## 代码和资料证据

- Internal Service：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/service/marketing_data_internal_service.go`
- Admin Service：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/service/marketing_data_admin_service.go`
- Manager 主流程：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/manager/impl/offline_data_fetch_manager_impl.go`
- Handler 接口：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/handler/offline_data_fetch_handler.go`
- Handler 实现：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/handler/offline_data_fetch/impl/`
- ES Repo：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/repo/elasticsearch/`
- History Repo：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/repo/repo_dao/impl/offline_plan_history_repo_dao_impl.go`
- 监控：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/monitor/register.go`
- `batch_id` 唯一索引痕迹：`/Users/si.chen/GolandProjects/insurance-marketing-data/resource/sql/dev-1.2.60/ddl.sql`
