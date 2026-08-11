# Marketing 项目材料

这个目录用于沉淀 Insurance Marketing 相关的项目背景、架构梳理、排障故事和面试表达。复习时优先抓住一条主线：`Admin 配置 Plan` -> `Marketing Engine 触发/筛选/执行` -> `Marketing-Data 批量供数` -> `O-BFF/Admin 聚合入口`。

## 当前文档

- [Marketing Engine 架构梳理](marketing-engine-architecture.md)：核心项目文档，适合回答架构设计、规则引擎、事件驱动、配置化、Processor/Handler 插件化。
- [Marketing Engine 架构图解](marketing-engine-architecture.html)：上面架构文档的可视化版本，适合先建立整体图像。
- [Marketing Engine 核心知识深挖](marketing-engine-key-knowledge-deep-dive.md)：专门背 Plan、Trigger、Processor、Handler、RuleCenter 和服务边界。
- [Marketing Processor / Handler / RuleCenter 深挖](marketing-processor-handler-rulecenter-deep-dive.md)：专门背筛选链、执行动作、RuleCenter 边界、重复触达防守。
- [Marketing 稳定性和排障深挖](marketing-reliability-troubleshooting-deep-dive.md)：专门背 OOMKilled、大批量人群、Kafka、Handler 失败和配置风险。
- [Marketing PlanDispatch 生命周期深挖](marketing-plan-dispatch-lifecycle-deep-dive.md)：专门背 PlanRawMessage、PlanDispatch、Trigger 分支、DoExecuteHandler 和异常记录。
- [Marketing User Group Processor 链深挖](marketing-user-group-processor-chain-deep-dive.md)：专门背 Processor 注册、三类 Processor、ConditionType、Delay 和筛选优化。
- [Marketing 用户名单 Redis / Roaring Bitmap 深挖](marketing-user-list-redis-roaring-bitmap-deep-dive.md)：专门背用户名单缓存、Redis 二级 Roaring Bitmap、LocalCache、RemoteLocalCache 和缓存一致性。
- [Marketing 用户名单 Redis / Roaring Bitmap 图解](marketing-user-list-redis-roaring-bitmap-deep-dive.html)：用图表记忆方案选型、high/low 拆分、Redis 二级结构和生命周期。
- [Marketing Handler 执行和风险控制深挖](marketing-handler-risk-control-deep-dive.md)：专门背 Handler 接口、Notify、Voucher、PNAR、HandlerRules 和资损/骚扰风险。
- [Marketing Consumer / Batch History 可靠性深挖](marketing-consumer-batch-history-reliability-deep-dive.md)：专门背 consumer 注册、topic 并发限制、batch history 对账和异常状态。
- [Marketing 技术优化点和深挖清单](marketing-technical-optimization-points.md)：按优化点汇总配置化、并发、批次对账、防重复、可观测性和配置风险。
- [Marketing 容器 OOMKilled 排障](oom-killed-troubleshooting-star.md)：稳定性和线上排障故事，适合回答容器内存、Go 服务排障、如何区分 panic 和 OOM。
- [OOMKilled 排障图解](oom-killed-troubleshooting-star.html)：上面排障故事的证据链图解，适合复习 `exit code 137` 和 `usage > limit`。

## 本机项目路径

| 项目 | 本机路径 | 当前分支 | 面试定位 |
| --- | --- | --- | --- |
| `insurance-marketing` | `/Users/si.chen/GolandProjects/insurance-marketing` | `release-v1.3.254` | Marketing 主服务，承载 Plan Engine、RuleCenter、Group/Notify/Reward 等营销能力。 |
| `insurance-marketing-data` | `/Users/si.chen/GolandProjects/insurance-marketing-data` | `release-v1.3.255` | 批量/离线人群数据服务，负责离线数据拉取、ES/S3/CSV、Kafka 推送和离线执行记录。 |
| `insurance-marketing-api` | `/Users/si.chen/GolandProjects/insurance-marketing-api` | `release-v1.3.266` | Marketing 主服务 API/proto 仓库。 |
| `insurance-marketing-data-api` | `/Users/si.chen/GolandProjects/insurance-marketing-data-api` | `master` | Marketing-Data API/proto 仓库。 |
| `insurance-operator-bff` | `/Users/si.chen/GolandProjects/insurance-operator-bff` | `master-ai` | O-BFF / Admin 聚合入口，负责后台查询、配置、审批、批量任务和下游 RPC。 |

> 注意：真正改代码或复盘线上问题时，要以当次需求分支、发版分支和运行环境为准。

## 代码入口速查

| 想看什么 | 推荐路径 |
| --- | --- |
| Engine 统一入口 `PlanDispatch` | `/Users/si.chen/GolandProjects/insurance-marketing/src/engine/internal/manager/impl/plan_manager_impl.go` |
| PlanManager 接口定义 | `/Users/si.chen/GolandProjects/insurance-marketing/src/engine/internal/manager/plan_manager.go` |
| User Group Processor 链式匹配 | `/Users/si.chen/GolandProjects/insurance-marketing/src/engine/internal/biz/impl/user_group_biz_impl.go` |
| 用户名单 Redis / Roaring Bitmap | `/Users/si.chen/GolandProjects/insurance-marketing/src/basic/group-center/` |
| Processor 三类接口 | `/Users/si.chen/GolandProjects/insurance-marketing/src/engine/internal/processor/user_group/` |
| Handler 插件 | `/Users/si.chen/GolandProjects/insurance-marketing/src/engine/internal/handler/` |
| RuleCenter | `/Users/si.chen/GolandProjects/insurance-marketing/src/basic/rule-center/` |
| Engine Kafka consumer | `/Users/si.chen/GolandProjects/insurance-marketing/src/engine/consumer/inner/default_consumer.go` |
| Marketing-Data 离线拉数编排 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/manager/impl/offline_data_fetch_manager_impl.go` |
| Marketing-Data 服务入口 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/service/` |
| O-BFF Admin 聚合入口 | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/service/`, `/Users/si.chen/GolandProjects/insurance-operator-bff/src/manager/`, `/Users/si.chen/GolandProjects/insurance-operator-bff/src/integrate/` |

## 面试复习顺序

1. 先看 [Marketing Engine 架构图解](marketing-engine-architecture.html)，记住 `Plan = Trigger + User Group + Handler`。
2. 再看 [Marketing Engine 架构梳理](marketing-engine-architecture.md)，补充代码路径、设计亮点和可追问点。
3. 单独背模块深挖：[PlanDispatch 生命周期](marketing-plan-dispatch-lifecycle-deep-dive.md)、[Processor 链](marketing-user-group-processor-chain-deep-dive.md)、[用户名单 Redis / Roaring Bitmap](marketing-user-list-redis-roaring-bitmap-deep-dive.md)、[Handler 风险控制](marketing-handler-risk-control-deep-dive.md)。
4. 再背可靠性和优化：[Consumer / Batch History](marketing-consumer-batch-history-reliability-deep-dive.md)、[技术优化点](marketing-technical-optimization-points.md)、[稳定性排障](marketing-reliability-troubleshooting-deep-dive.md)。
5. 最后看 [Marketing 容器 OOMKilled 排障](oom-killed-troubleshooting-star.md)，准备一个稳定性案例。

## 待补充事实

- 你在 Marketing Engine 里实际负责过的具体需求、PR 或模块。
- Marketing-Data 与 Engine 在具体 Plan 上的 topic、消息体字段和一次完整执行案例。
- O-BFF 到 Marketing 的真实 Admin 页面路由、Gateway/Assembly 配置和接口名。
- OOMKilled 案例最终是否扩容、根因是否定位到某个任务/consumer/handler。
