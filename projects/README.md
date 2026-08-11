# Projects 项目复习总入口

这个目录专门放项目经历、业务背景、系统架构、优化点、排障故事和面试讲法。复习时不要直接在文件列表里翻，先按下面三层路径走：

> 可视化导航：[Projects 项目复习结构图](project-map.html)

1. `Insurance/`：部门级业务地图，回答“你们组做什么、上下游怎么协作”。
2. `Marketing/`、`O-BFF/`、`Marketing-Data/`、`Canal/`：具体服务深挖，回答“你负责/熟悉的系统怎么设计、怎么优化、怎么排障”。
3. `AI/`：AI Agent、知识工程化、MCP/Skill 等补充项目，作为加分项。

## 目录结构

| 目录 | 定位 | 先看什么 | 适合回答的问题 |
| --- | --- | --- | --- |
| [Insurance](Insurance/README.md) | Insurance 部门全局业务地图和跨服务关系。 | `insurance-business-map.html`、`other-business-overview.md` | 你们部门负责什么？一次保险业务链路怎么走？这些服务之间怎么协作？ |
| [Marketing](Marketing/README.md) | Marketing Engine 主服务，重点是 Plan、Trigger、Processor、Handler、RuleCenter 和用户组匹配。 | `marketing-engine-architecture.html`、`marketing-engine-architecture.md` | 营销系统怎么设计？如何做人群筛选？如何防重复触达？有什么优化点？ |
| [O-BFF](O-BFF/README.md) | Insurance Admin / Operator BFF，重点是后台聚合、配置化组件、审批、批量任务、Task Center。 | `o-bff-architecture.html`、`o-bff-approval-task-center-architecture.html` | BFF 为什么存在？审批中心/任务中心怎么设计？配置化组件怎么做？ |
| [Marketing-Data](Marketing-Data/README.md) | Marketing 离线/批量数据服务，重点是数据拉取、batch history、handler、Kafka 推送。 | `marketing-data-architecture.html`、`marketing-data-architecture.md` | 离线人群怎么供数？如何保证幂等和一致性？Kafka 推送失败怎么处理？ |
| [Canal](Canal/README.md) | MySQL binlog 到 ES 的同步链路，重点是 Canal、Kafka、Adapter、sync script、ETL、DLQ。 | `canal-es-sync-architecture.html`、`canal-es-sync-architecture.md` | ES 数据怎么同步？binlog 到 ES 中间有哪些环节？同步阻塞怎么排障？ |
| [AI](AI/README.md) | AI Agent / 知识工程化补充项目。 | `voucher-operations-agent-design.md`、`insurance-knowledge-engineering-plan.md` | 你怎么用 AI 工程化提升团队知识沉淀？如何把 Agent 做成生产系统？ |

## 推荐复习路线

### 30 分钟快速过项目

1. 打开 [Insurance 业务地图](Insurance/insurance-business-map.html)，先把入口层、核心业务服务、数据链路和后台服务记住。
2. 打开 [Marketing Engine 架构图解](Marketing/marketing-engine-architecture.html)，记住 `Plan = Trigger + User Group + Handler`。
3. 打开 [O-BFF 架构图解](O-BFF/o-bff-architecture.html)，记住 `Admin 请求 = 拦截治理 + 配置化组件 + 聚合编排 + 中心化后台能力`。
4. 打开 [Marketing-Data 架构图解](Marketing-Data/marketing-data-architecture.html)，记住 `离线抽取 = 幂等 + history + handler + Kafka + monitor`。
5. 打开 [Canal / ES 同步架构图解](Canal/canal-es-sync-architecture.html)，记住 `MySQL binlog -> Canal -> Kafka -> Adapter -> ES -> Admin 查询`。

### 面试前重点背

| 优先级 | 文档 | 为什么要背 |
| --- | --- | --- |
| P0 | [Marketing Engine 架构梳理](Marketing/marketing-engine-architecture.md) | 主项目，最容易被要求详细讲。 |
| P0 | [Marketing 用户名单 Redis / Roaring Bitmap 深挖](Marketing/marketing-user-list-redis-roaring-bitmap-deep-dive.md) | 有明确技术优化点，适合讲缓存、Redis、数据结构和性能。 |
| P0 | [O-BFF 配置化组件深挖](O-BFF/o-bff-config-components-deep-dive.md) | 你提到 O-BFF 做了一整套配置化组件，适合包装成后台效率提升项目。 |
| P0 | [Approval Center / Task Center 架构说明](O-BFF/o-bff-approval-task-center-architecture.md) | 审批中心、任务中心很容易被追问架构边界。 |
| P1 | [Marketing-Data 可靠性和一致性深挖](Marketing-Data/marketing-data-reliability-deep-dive.md) | 适合回答离线任务、幂等、重试、对账。 |
| P1 | [Canal DLQ / 同步阻塞优化专项](Canal/canal-dlq-blocking-optimization.md) | 适合回答基础设施和线上排障。 |
| P1 | [Marketing 稳定性和排障深挖](Marketing/marketing-reliability-troubleshooting-deep-dive.md) | 准备稳定性追问和事故排查题。 |

## 文件放置规则

为了避免后面继续变乱，新增材料按这些规则放：

| 材料类型 | 放哪里 | 命名方式 |
| --- | --- | --- |
| 部门整体业务、跨服务流程、统一面试口径 | `projects/Insurance/` | `business-flow-map.md`、`interview-speaking-notes.md` 这类总览名 |
| 某个服务的主架构文档 | `projects/<Project>/` | `<project>-architecture.md`，需要图解时同名 `.html` |
| 某个模块或技术点深挖 | `projects/<Project>/` | `<project>-<topic>-deep-dive.md` |
| 优化点集合 | `projects/<Project>/` | `<project>-technical-optimization-points.md` |
| 线上排障或 STAR 故事 | `projects/<Project>/` | `<incident-or-topic>-star.md`，需要图解时同名 `.html` |
| 资料来源和已查证路径 | 优先放对应项目 README；跨项目来源放 `Insurance/source-index.md` | 不要散落在临时文件里 |

## Markdown 和 HTML 规则

- 有 HTML 图解时，Markdown 顶部要放同名 HTML 链接。
- HTML 和 Markdown 尽量同名，例如 `marketing-user-list-redis-roaring-bitmap-deep-dive.md` 对应 `marketing-user-list-redis-roaring-bitmap-deep-dive.html`。
- HTML 必须能直接本地打开，样式和脚本内联，不依赖 CDN。
- 不要为了“看起来丰富”给每篇文档都建 HTML；只有架构、链路、缓存结构、生命周期、排障流程这类视觉化收益明显的主题才建。

## 内容边界

- `Insurance/` 讲“全局地图”和“跨服务口径”，不承载某个服务的全部深挖。
- `Marketing/`、`O-BFF/`、`Marketing-Data/`、`Canal/` 才是具体项目的主目录。
- 每个项目目录里的 `README.md` 是该项目唯一入口；新增文档后要同步更新对应 README。
- 不确定的个人贡献、线上指标、事故结论统一写 `待补充`，不要写成已确认事实。

## 本机代码路径速查

| 服务 | 本机代码路径 |
| --- | --- |
| Marketing | `/Users/si.chen/GolandProjects/insurance-marketing` |
| Marketing API | `/Users/si.chen/GolandProjects/insurance-marketing-api` |
| Marketing-Data | `/Users/si.chen/GolandProjects/insurance-marketing-data` |
| Marketing-Data API | `/Users/si.chen/GolandProjects/insurance-marketing-data-api` |
| O-BFF | `/Users/si.chen/GolandProjects/insurance-operator-bff` |
| Canal Adapter | `/Users/si.chen/GolandProjects/canal-adapter` |

## 后续整理清单

- `待补充`：把你真实做过的 PR、Jira、上线记录和排障记录补到各项目 README 的“待补充事实”里。
- `待补充`：给每个 P0 项目准备一个 1 分钟口述版和一个 5 分钟深挖版。
- `待补充`：如果后续新增 Product、Order、Policy、Promotion 等业务，只在内容超过 3 篇文档时单独建目录；否则先放在 `Insurance/other-business-overview.md`。
