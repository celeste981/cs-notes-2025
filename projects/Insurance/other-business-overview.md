# Insurance 其他业务总览

[HTML 业务地图](insurance-business-map.html) | [业务流程地图](business-flow-map.md) | [服务深挖笔记](service-deep-dive-notes.md) | [资料来源索引](source-index.md)

## 一句话定位

Insurance 后端不是单一服务，而是一组围绕“保险产品售卖、保单履约、运营后台、营销触达、数据同步”的服务群。面试时可以把它讲成四层：

```text
入口层：C-BFF / B-BFF / L-BFF / O-BFF / Unified Gateway
核心业务层：Product / Order / Policy / Account / Promotion / Marketing / Marketing-Data / Risk Control
外部集成层：EA / UEA / 通知 / 第三方 provider
数据支撑层：MySQL / Redis / ES / Kafka / Canal / S3
```

## 面试主线

> 我所在的 Insurance 后端组主要支撑保险业务的后台运营和营销域。面试里我会重点讲 O-BFF、Marketing、Marketing-Data 和 Canal：O-BFF 是 Admin 后台聚合入口，Marketing 负责运营计划和触达，Marketing-Data 负责离线人群和批量数据，Canal 负责 MySQL 到 ES 等搜索/报表数据的同步。其他 Product、Order、Policy、Promotion、Account 等服务是核心业务 owner，我需要理解它们的边界和调用契约。

## 核心服务地图

| 服务 | 面试定位 | 关键能力 | 代码路径 |
| --- | --- | --- | --- |
| O-BFF / `insurance-operator-bff` | Insurance Admin 后台聚合层 | 后台 gRPC API、配置化组件、审批、批量操作、任务中心、数据修复、ES 查询、下游服务编排 | `/Users/si.chen/GolandProjects/insurance-operator-bff` |
| Marketing / `insurance-marketing` | 运营自动化与营销触达服务 | Plan Engine、人群、通知、发券、PNAR、RuleCenter、活动和个性化展示 | `/Users/si.chen/GolandProjects/insurance-marketing` |
| Marketing-Data / `insurance-marketing-data` | Marketing 的离线/批量数据服务 | 离线人群拉取、ES 查询、S3/CSV/Hive/Insight 文件数据、批次历史、Kafka 推送 | `/Users/si.chen/GolandProjects/insurance-marketing-data` |
| Canal / `canal-adapter` | 数据同步基础设施 | MySQL binlog 同步、Adapter、ES sync script、RDS 扩容切换、回滚物料 | `/Users/si.chen/GolandProjects/canal-adapter` |
| Product | 产品域 owner | 产品配置、产品规则、计划、价格/定价相关产品数据 | `/Users/si.chen/GolandProjects/insurance-product` |
| Order | 订单域 owner | 保险订单、order item、支付、订单和保单关联 | `/Users/si.chen/GolandProjects/insurance-order` |
| Policy | 保单域 owner | 保单生命周期、保单查询、核保、理赔通知、续保/取消等子能力 | `/Users/si.chen/GolandProjects/insurance-policy` |
| Promotion / Voucher | 促销域 owner | 促销活动、折扣、优惠券、库存、预算、活动状态 | `/Users/si.chen/GolandProjects/insurance-promotion` |
| Account | 账户域 owner | SPIN account、Shopee / ShopeePay 用户映射、accountId 查询 | `/Users/si.chen/GolandProjects/insurance-account` |
| Risk Control | 风控与清结算相关能力 | 风险告警、对账、清结算、保费重算等 | `/Users/si.chen/GolandProjects/insurance-risk-control` |
| EA / UEA | 外部系统适配层 | 第三方保险供应商、通知、SeaTalk、Transify、ShopeePay、Wiz 等 | `/Users/si.chen/GolandProjects/insurance-external-aggregator`, `/Users/si.chen/GolandProjects/insurance-unified-external-aggregator` |

## 入口 BFF 补充

| BFF | 上游 | 核心模型 | 适合面试展开 |
| --- | --- | --- | --- |
| C-BFF | App / H5 C 端用户 | 产品展示、报价、预检、订单、支付、投保、保单和理赔查询 | “C 端购买闭环”和“多地区/多产品线适配”。 |
| B-BFF | MP / DP / CL 等渠道 | 插件化工作流：Rollout、Auth、DataFilling、PresaleCheck、PremiumsCalc、Marketing、Priority、OrderCreation/Update | “为什么用插件化”以及“MP/CL/DP 流程差异”。 |
| L-BFF | OFG / OMS / OA / ODO 消息 | Consumer + Adaptor Chain + Filter Chain + State Machine + Processor | “事件驱动物流保险”和“状态机/幂等/短路过滤”。 |
| O-BFF | Insurance Admin | Service/Manager/Integrate + Approval/Batch/Task/DataFix | “后台高风险操作怎么控风险”。 |

## OBFF

### 先讲人话

O-BFF 是 Insurance Admin 后台的后端聚合层。Admin 页面不直接调用 Product、Order、Policy、Marketing 等所有服务，而是先打到 O-BFF，由 O-BFF 做参数校验、权限/脱敏/审批拦截、业务编排和下游 RPC 调用。

### 关键能力

- Service 层：注册 gRPC handler，适配 Admin 请求。
- Manager/Biz 层：做后台业务编排。
- Integrate 层：封装 Product、Order、Policy、Marketing、Promotion、Voucher、Risk Control 等下游 RPC。
- Approval Center：审批流程定义、实例、审批动作和审批事件。
- Batch Operate / Task Center：批量上传、审批、异步执行、结果下载。
- Data Fix Center：数据修复需要审批，通过可靠事件异步执行。
- ES 查询：后台复杂检索和运营分析。

### 代码入口

| 想看什么 | 路径 |
| --- | --- |
| 系统架构 KB | `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/System_Architecture.md` |
| 跨应用关系 | `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/Cross_App_Overview.md` |
| Service | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/service/` |
| Manager | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/manager/` |
| Integrate | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/integrate/` |
| Approval Center | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/approval_center/` |
| Batch Operate Center | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/batch_operate_center/` |
| Task Center | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/task_center/` |
| Data Fix Center | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/data_fix_center/` |
| 配置化组件地图 | `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/system/Admin_Platform_Component_Map.md` |
| Assembly 配置 ADR | `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/decisions/ADR-0002-admin-current-assembly-configuration.md` |

## Marketing

### 先讲人话

Marketing 是运营自动化服务。运营在 Admin 上配置 Plan，系统按“触发 -> 筛选用户 -> 执行动作”的流水线自动执行通知、发券、PNAR、活动展示等营销动作。

### 关键能力

- Plan Engine：统一承接定时、MQ、实时事件。
- User Group Processor：用 Processor 链筛选目标用户。
- Handler：执行通知、发券、PNAR、Chat、数据清洗等动作。
- RuleCenter：复用通用条件判断。
- Basic 能力：Group Center、Notify Center、Reward Center。
- 下游依赖：Account、Product、Promotion/Voucher、Policy、Order、C-BFF、UEA 等。

### 代码入口

| 想看什么 | 路径 |
| --- | --- |
| 深挖材料 | `/Users/si.chen/Desktop/prep/projects/Marketing/README.md` |
| Engine | `/Users/si.chen/GolandProjects/insurance-marketing/src/engine/` |
| RuleCenter | `/Users/si.chen/GolandProjects/insurance-marketing/src/basic/rule-center/` |
| Group / Notify / Reward | `/Users/si.chen/GolandProjects/insurance-marketing/src/basic/` |
| 下游 RPC client | `/Users/si.chen/GolandProjects/insurance-marketing/src/common/integrate/` |

## Marketing-Data

### 先讲人话

Marketing-Data 是 Marketing 的离线数据助手。Marketing Engine 需要批量圈人时，不适合在主服务里直接做大查询和文件处理，所以由 Marketing-Data 负责查 ES、读取 S3/CSV/Hive/Insight 数据、生成批次执行记录，并把用户数据推到 Kafka 给 Marketing Engine 消费。

### 关键能力

- Internal API：`OfflineDataFetchAndPush`、`QueryOfflinePlanHistory`。
- Admin API：`QueryEsList`，支持 Admin 侧预览 ES 查询结果。
- OfflineDataFetchManager：做 batch_id 幂等、初始化 history、异步执行、更新执行结果、上报监控。
- Data Handler 插件：按 `DataSource` 自注册不同数据源处理器。
- 数据源实现：`es_data_fetch_impl.go`、`file_data_fetch_impl.go`、`hive_file_data_fetch_impl.go`、`insight_file_data_fetch_impl.go`。

### 代码入口

| 想看什么 | 路径 |
| --- | --- |
| Internal Service | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/service/marketing_data_internal_service.go` |
| Admin Service | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/service/marketing_data_admin_service.go` |
| 离线拉数编排 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/manager/impl/offline_data_fetch_manager_impl.go` |
| Handler 接口 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/handler/offline_data_fetch_handler.go` |
| Handler 实现 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/handler/offline_data_fetch/impl/` |
| ES/S3 数据层 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/repo/` |

## Canal

### 先讲人话

Canal 负责把 MySQL 的变更同步到下游数据系统，Insurance 里常见用途是 MySQL -> Kafka/Adapter -> Elasticsearch，让 Admin 查询、报表和搜索类能力能从 ES 读到近实时数据。

### 关键能力

- Canal Instance：订阅某个 MySQL / RDS 分片的 binlog。
- Adapter datasource：配置回查数据源，用于 ETL 和同步。
- ES sync script：定义从源表到 ES index 的字段映射、过滤条件和脚本版本。
- 切换/回滚：RDS 扩容或迁移时，预置新 instance 和新 sync script，再按地区/模块启用新脚本、软删旧脚本。

### 代码和物料入口

| 想看什么 | 路径 |
| --- | --- |
| Canal README | `/Users/si.chen/GolandProjects/canal-adapter/README.md` |
| RDS 扩容迁移 TD | `/Users/si.chen/GolandProjects/canal-adapter/docs/id_vn_rds_expansion_canal_td.md` |
| Canal 物料目录 | `/Users/si.chen/GolandProjects/canal-adapter/docs/canal_materials/` |
| ES adapter | `/Users/si.chen/GolandProjects/canal-adapter/es7x/`, `/Users/si.chen/GolandProjects/canal-adapter/es6x/` |
| Launcher | `/Users/si.chen/GolandProjects/canal-adapter/launcher/` |
| 工具脚本 | `/Users/si.chen/GolandProjects/canal-adapter/tools/` |

## 典型核心链路

### Admin 后台配置链路

```text
Insurance Admin
  -> O-BFF Service / Interceptor
  -> Manager / Biz
  -> 下游 Product / Policy / Order / Marketing / Promotion
  -> DB / ES / Config / 审批 / 批量任务
```

面试抓手：O-BFF 的价值是统一后台入口，减少 Admin 直接对接多个后端服务的复杂度。

### Marketing 计划执行链路

```text
Admin 配置 Plan
  -> Marketing Plan Engine
  -> Trigger: Cron / MQ / Event
  -> User Group Processor 筛选
  -> RuleCenter 校验
  -> Handler 执行通知 / 发券 / PNAR
```

面试抓手：配置化、策略模式、自注册、事件驱动、统一入口。

### 离线人群链路

```text
Marketing Engine 定时计划
  -> Marketing-Data OfflineDataFetchAndPush
  -> ES / S3 / CSV / Hive / Insight 数据源
  -> 批次 history + 监控
  -> Kafka topic
  -> Marketing Engine 消费执行
```

面试抓手：把大查询和批量数据处理从主执行链路拆出去，降低 Marketing 主服务压力。

### ES 同步链路

```text
MySQL binlog
  -> Canal Instance
  -> Kafka topic / Adapter
  -> ES sync script 映射
  -> Elasticsearch index
  -> O-BFF / Admin 查询和报表
```

面试抓手：ES 不是业务事实源，事实源仍然是 MySQL owner 服务；ES 用于查询加速和分析展示。

### C 端购买链路

```text
App / H5
  -> C-BFF
  -> Product / Promotion / Voucher / Risk
  -> Order / Payment
  -> Policy
  -> EA / UEA
  -> Document / Notification / Marketing
```

面试抓手：C-BFF 做用户侧聚合，复杂点在产品线、地区差异、支付回调幂等、投保失败补偿。

### B 端 MP 嵌入式链路

```text
MP PDP / Checkout
  -> B-BFF
  -> Product / Promotion
  -> ODO Kafka
  -> B-BFF Consumer
  -> Order / Policy
  -> ODP 查询保单
```

面试抓手：B-BFF 用插件化工作流承接不同渠道和市场差异，创单与投保更新可拆 topic 避免互相阻塞。

### 物流保险链路

```text
OFG / OMS 消息
  -> L-BFF Consumer
  -> Adaptor Chain 补全数据
  -> Filter Chain 校验规则
  -> State Machine 推进状态
  -> Processor 调用 Core / Order / Policy
```

面试抓手：物流保险主要靠事件驱动，要关注可靠消费、幂等、状态机恢复和过滤器短路。

## Confluence / KB 参考

| 来源 | 用途 |
| --- | --- |
| `OBFF 功能分析` Confluence：`https://confluence.shopee.io/pages/viewpage.action?pageId=2968425075` | OBFF 功能、审批、批量、任务、数据修复、ES 查询等能力总结。 |
| `B-BFF 功能分析` Confluence：`https://confluence.shopee.io/pages/viewpage.action?pageId=2968424010` | B-BFF 插件工作流、MP/DP/CL 聚合能力。 |
| `CBFF功能分析` Confluence：`https://confluence.shopee.io/pages/viewpage.action?pageId=2968423411` | C 端产品展示、报价、下单、支付、投保、保单和理赔闭环。 |
| `LBFF功能分析` Confluence：`https://confluence.shopee.io/pages/viewpage.action?pageId=2968424743` | 物流保险 Consumer、Adaptor、Filter、State Machine、Processor 模型。 |
| `EC Buyer业务梳理` Confluence：`https://confluence.shopee.io/pages/viewpage.action?pageId=2578909398` | MP/CL/DP 业务流程、topic、接口和市场差异。 |
| `[Marketing] 运营计划架构` Confluence：`https://confluence.shopee.io/pages/viewpage.action?pageId=1825819385` | Marketing 与 Marketing-Data 拆分、运营计划能力、数据模型和性能监控。 |
| `ES 使用分享` Confluence：`https://confluence.shopee.io/pages/viewpage.action?pageId=2347394775` | ES 在 Admin 的用途、Canal 同步原理、mapping 和常见排障。 |
| `SPIN Global Service Map` 本地 KB | Insurance 服务边界和全局服务地图。 |
| `insurance-marketing/project-kb/Cross_App_Overview.md` | Marketing 上下游关系和依赖 owner。 |
| `canal-adapter/docs/id_vn_rds_expansion_canal_td.md` | Canal RDS 扩容、ES sync script、切换/回滚物料。 |

## 待补充

- 你在 OBFF 上实际做过的 Admin 页面/API/批量任务/审批需求。
- 你在 Marketing-Data 上实际做过的 DataSource、ES 查询或批次执行案例。
- 你在 Canal 上实际做过的 sync script、RDS 扩容、ES mapping 或回滚案例。
- 每个项目可以防守的 PR、Jira、Confluence TD 或排障记录。
