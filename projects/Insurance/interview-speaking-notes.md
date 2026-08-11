# 面试口径与边界

## 推荐总口径

面试时建议用“业务域负责”而不是“我个人独立负责所有服务”。更稳的说法是：

> 我在 Insurance 后端组，主要覆盖后台运营和营销域相关系统。我们组负责 O-BFF、Marketing、Marketing-Data 以及 ES/Canal 同步相关链路。我参与这些系统的开发、维护和排障时，通常会从 Admin 配置入口、Marketing 执行链路、离线数据和 ES 同步几个方向去看问题。

如果面试官继续问“哪些是你个人做的”，要回到事实：

> 我个人最熟的是 Marketing Engine 和相关排障案例；OBFF、Marketing-Data、Canal 属于组内业务，我能讲清楚架构和链路，也参与过相关需求/排障/物料整理。具体个人 PR 和上线记录我会按实际经历展开。

个人具体经历：`待补充`。

## 30 秒口述版

我所在的 Insurance 后端组主要负责保险后台运营和营销域。核心系统有四块：第一是 O-BFF，作为 Insurance Admin 的后台聚合入口；第二是 Marketing，负责运营计划、用户筛选、通知、发券和 PNAR；第三是 Marketing-Data，负责离线人群拉取、ES/文件数据处理和批次执行；第四是 Canal/ES 同步，用于把 MySQL 变更同步到 ES，支撑后台查询和报表。我的面试重点会讲 Marketing Engine，同时能把 OBFF、Marketing-Data、Canal 的上下游链路串起来。

## 1 分钟口述版

我在 Insurance 后端组，业务上主要围绕后台运营和营销域。这个域不是一个单服务，而是一条从 Admin 配置到营销执行、再到数据查询的链路。

入口是 O-BFF，它面向 Insurance Admin，把产品、订单、保单、营销、风控、券等多个下游服务聚合成后台 API，同时承载审批、批量操作、任务中心和数据修复。营销执行侧是 Marketing，它的核心是 Plan Engine，通过 Trigger、User Group Processor 和 Handler 组成配置化运营流水线。定时或离线人群场景会用 Marketing-Data，它负责从 ES、S3、CSV、Hive 等数据源批量拉人群，记录 batch history，再推 Kafka 给 Marketing 执行。后台查询和报表还依赖 Canal，把 MySQL binlog 同步到 ES。

所以我讲这个业务时，会从 O-BFF 的 Admin 聚合、Marketing 的配置化执行、Marketing-Data 的离线数据、Canal 的 ES 同步四块展开。

## 2 分钟扩展版

如果面试官问我们 Insurance 后端整体架构，我会按入口层、核心业务层、数据支撑层和外部集成层讲。

入口层有 C-BFF、B-BFF、L-BFF 和 O-BFF。C-BFF 面向 C 端用户，聚合产品展示、保费计算、订单支付、投保、保单和理赔；B-BFF 面向 MP、DP、CL 等渠道，核心是插件化工作流，把授权、数据填充、售前校验、保费计算、营销和订单创建串起来；L-BFF 面向物流保险，消费 OFG/OMS 等消息，通过适配器链补数据、过滤器链校验规则、状态机推进流程；O-BFF 面向 Admin 后台，统一承接审批、批量、任务中心、数据修复和下游编排。

核心业务层是 Product、Order、Policy、Promotion、Account、Risk Control。Product 管产品、计划、保障和费率；Order 管订单、支付和状态机；Policy 管投保、保单、理赔、取消和续保；Promotion/Voucher 管活动、券、库存和价格；Account 管统一 accountId；EA/UEA 作为防腐层对接保司和外部系统。数据支撑层里，Marketing-Data 负责离线人群和批量数据，Canal/ES 支撑 Admin 查询和报表。

我个人面试重点还是后台运营和营销域，但会把这些上下游链路讲清楚，说明我理解整个 Insurance 系统怎么协作。

## 分项目说法

| 项目 | 可以怎么说 | 追问时要能答 |
| --- | --- | --- |
| O-BFF | 我们组的后台聚合入口，承接 Insurance Admin，负责审批、批量、任务、数据修复和下游服务编排。 | Service/Manager/Integrate 分层，审批和批量操作流，为什么需要 BFF。 |
| Marketing | 我重点熟悉的运营自动化服务，核心是 Plan Engine，支持定时、MQ、实时事件触发。 | Plan 模型、Processor 链、Handler、RuleCenter、PNAR/发券/通知。 |
| Marketing-Data | Marketing 的离线数据服务，负责批量拉人群、查 ES、处理文件数据、记录 batch history 并推 Kafka。 | `OfflineDataFetchAndPush`、batch 幂等、DataSource handler、ES/S3/CSV/Hive。 |
| Canal | 组内的数据同步链路，负责 MySQL 到 ES 的同步配置、RDS 扩容切换和回滚物料。 | Canal Instance、Adapter datasource、ES sync script、切换/回滚顺序。 |
| B-BFF | B 端渠道聚合层，对接 MP/DP/CL，用插件化工作流处理售前、保费、营销和订单。 | Workflow、plugin selector、MP/CL/DP 流程差异、可靠消费。 |
| C-BFF | C 端用户侧聚合层，处理产品展示、报价、支付、投保、保单和理赔。 | motor/accident_health/sales 模块、多地区 handler、支付回调和投保补偿。 |
| L-BFF | 物流保险编排层，靠 OFG/OMS 消息驱动。 | Adaptor Chain、Filter Chain、State Machine、Processor、幂等和重试。 |
| Product / Order / Policy / Promotion / Account / EA | 核心上下游服务，我理解边界和调用契约。 | 每个服务的事实源、状态机、外部依赖和对 Marketing/O-BFF 的影响。 |

## 可以强调的技术点

- BFF 聚合模式：前端只对接 O-BFF，O-BFF 屏蔽下游服务差异。
- 配置化引擎：Marketing 用 Plan 抽象运营活动，减少运营活动发版。
- 策略模式 / 自注册：Processor、Handler、Rule、DataSource handler 都能用 code/name 查找实现。
- 异步和可靠性：Kafka、reliable-event、batch history、失败重试、幂等。
- ES 查询与 Canal：MySQL 是事实源，ES 是查询/分析视图，Canal 负责变更同步。
- 审批和批量任务：后台高风险操作要审批、异步执行、可追踪、可下载结果。

## 风险说法

| 不建议说 | 风险 | 推荐替换 |
| --- | --- | --- |
| 我个人独立负责 OBFF、Marketing、Marketing-Data、Canal 全部系统 | 面试官追 PR、owner、线上事故，很容易穿帮 | 我所在组负责这些系统，我参与/熟悉其中的后台和营销链路 |
| Canal 是业务系统 | Canal 更像数据同步基础设施 | Canal 是支撑 Admin 查询和报表的 MySQL -> ES 同步链路 |
| ES 是数据事实源 | ES 可能延迟或丢同步 | MySQL owner 服务是事实源，ES 用于查询和分析 |
| Marketing-Data 是 Marketing 的替代服务 | 它是离线数据助手，不承载主执行引擎 | Marketing-Data 负责离线批量数据，Marketing Engine 负责计划执行 |

## 追问准备

### 你们组到底负责什么？

> 我们组主要负责 Insurance 后台运营和营销域，包括 O-BFF 的 Admin 后台聚合、Marketing 的运营计划引擎、Marketing-Data 的离线人群数据，以及 Canal/ES 同步相关的数据支撑。其他 Product、Order、Policy、Promotion 等是核心业务服务，我们会作为调用方理解它们的契约和边界。

### O-BFF 和 Marketing 的区别？

> O-BFF 是后台入口和聚合层，面向 Admin 页面，负责把后台操作编排到各个下游服务。Marketing 是营销业务服务，负责运营计划、人群筛选、通知、发券和活动展示。简单说，O-BFF 解决“后台怎么操作”，Marketing 解决“运营动作怎么自动执行”。

### O-BFF 为什么不是简单转发？

> 因为 Admin 后台操作有治理要求。O-BFF 除了调用下游，还要统一做权限、脱敏、操作日志、审批信息、重复请求校验、RPC proxy、批量任务、数据修复和 ES 查询。简单查询可以接近 proxy，但高风险修改和批量导入导出必须有审批、任务状态和错误明细。

### O-BFF 批量导入导出怎么讲？

> Admin 上传文件或发起导出后，O-BFF 会创建 batch/task 记录，做字段校验、去重、切批，必要时先走审批，再由 Task Center 异步执行。导入时下游最好返回每条失败数据的 `error_id` 和 `error_msg`；导出时要有稳定分页或 scroll，避免大数据量翻页丢数或重复。

### O-BFF 怎么避免变成大泥球？

> 边界要清楚：O-BFF 做 Admin 场景的聚合、适配和治理，不沉淀核心领域事实。Product、Order、Policy、Promotion 等 owner 服务负责自己的业务规则和数据事实；O-BFF 通过 Integrate 调用它们，并把审批、批量、日志、脱敏这类后台通用能力放在中心模块或 interceptor。

### Marketing 和 Marketing-Data 的区别？

> Marketing 是主业务执行服务，核心是 Plan Engine。Marketing-Data 是离线数据服务，负责批量拉人群、处理 ES/S3/CSV/Hive 等数据源、记录批次历史，再把结果推给 Marketing。这样可以把大查询和批量文件处理从主执行链路拆出来。

### Marketing-Data 的主链路怎么讲？

> Marketing 生成一个离线计划批次，把 `plan_id`、`batch_id`、`data_source`、`group_condition`、`dest_topic` 等传给 Marketing-Data。Marketing-Data 先校验参数，再用 `batch_id` 做幂等；如果没执行过就写 batch history，然后异步选择对应 DataSource handler 拉数、处理、推 Kafka，最后更新 success/failed、总数、去重数、错误详情和监控。

### Marketing-Data 为什么需要 batch_id 幂等？

> 离线任务可能因为调度重试、RPC 超时或人工重复触发被调用多次。`batch_id` 是同一批计划的唯一标识，先查 history，存在就直接返回，可以避免重复拉取和重复推 Kafka，降低重复触达用户的风险。

### Marketing-Data 怎么扩展新的数据源？

> 它有统一的 `OfflineDataFetchHandler` 接口，包含 `BeforeProcessData`、`ProcessData`、`AfterProcessData` 和 `GetStaticDetail`。新增数据源时实现这个接口，再按 `data_source` 注册初始化函数；主流程只负责根据 `data_source` 找 handler，不需要改整体 pipeline。

### Marketing-Data 会执行 PNAR 或发券吗？

> 不建议这么说。Marketing-Data 主要解决“人群从哪里来、怎么批量抽取和推送”，PNAR、发券、通知等具体营销动作属于 Marketing Engine / Handler 的执行范围。Marketing-Data 把离线人群结果推给 Marketing，Marketing 再继续执行营销动作。

### Canal 为什么重要？

> Admin 后台很多查询和报表不能直接扫 MySQL，所以会把 MySQL 变更通过 Canal 同步到 ES。Canal 这层要关注 instance、datasource、ES sync script、切换和回滚。它的核心价值是让后台查询快，但业务事实仍然以 MySQL owner 服务为准。

### 你怎么理解 B-BFF 插件化？

> B-BFF 要对接 MP、DP、CL 等渠道，不同地区、不同产品、不同渠道的售前校验、授权、保费、营销和订单逻辑都不一样。如果全部写在一个流程里会非常难维护，所以它把流程拆成 Rollout、Auth、DataFilling、PresaleCheck、PremiumsCalc、Marketing、Priority、OrderCreation 等插件，再用 plugin selector 和 workflow 编排。这样新增市场或策略时，优先新增插件或配置选择器。

### L-BFF 为什么需要状态机？

> 物流保险不是用户同步点击下单，而是 OFG/OMS 等外部物流状态驱动。消息可能重复、延迟或乱序，所以需要状态机判断当前状态能不能流转到下一状态，配合可靠消费和幂等，避免重复创建保险订单或非法确认订单。

### ES 查询出问题怎么排？

> 我会分层排查：先确认业务 MySQL 里有没有数据，再看 Canal instance 是否正常消费 binlog，Kafka/Adapter 是否阻塞，ES mapping 是否缺字段或类型错误，最后看查询 DSL 是否用了错误的 term/match/range。ES 是查询视图，不是事实源，所以不能只看 ES 判断业务数据不存在。

## 待补充个人证据

- OBFF：具体做过哪个 Admin 页面、接口、批量任务、审批或数据修复需求。
- Marketing：具体做过哪个 Plan、Processor、Handler、Rule 或 PNAR/发券/通知需求。
- Marketing-Data：具体做过哪个 DataSource、ES 查询、文件拉取、batch history 或 Kafka 推送问题。
- Canal：具体做过哪个 ES sync script、RDS 扩容、mapping、切换、回滚或数据同步排障。
