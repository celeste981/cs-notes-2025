# Insurance 部门业务材料

> Projects 总入口：[../README.md](../README.md)

这个目录用于记录 Insurance 部门的组内业务、上下游服务和面试可讲口径。它和 `projects/Marketing/` 的区别是：这里不是单个项目深挖，而是把整个 Insurance 后端业务地图串起来，方便面试时回答“你们组负责什么”“你了解哪些系统”“跨服务怎么协作”。

## 目录定位

| 内容 | 放在这里吗 | 说明 |
| --- | --- | --- |
| Insurance 部门业务地图、跨服务关系、整体面试口径 | 是 | 这是本目录的核心用途。 |
| O-BFF / Marketing / Marketing-Data / Canal 的一句话定位和协作关系 | 是 | 用于回答“我们组整体负责什么”。 |
| 某个具体服务的完整架构、模块深挖、优化点 | 否 | 放到对应项目目录：`../Marketing/`、`../O-BFF/`、`../Marketing-Data/`、`../Canal/`。 |
| 个人真实 PR、Jira、上线和排障记录 | 视情况 | 跨项目证据可放这里；单服务证据优先放对应项目 README。 |

## 当前文档

- [Insurance 其他业务总览](other-business-overview.md)：按服务梳理 OBFF、Marketing、Marketing-Data、Canal 以及其他核心业务服务。
- [Insurance 业务地图 HTML](insurance-business-map.html)：用一页静态 HTML 看入口层、核心服务、数据链路和面试讲法。
- [业务流程地图](business-flow-map.md)：把 C 端、B 端、物流、Admin、Marketing、Marketing-Data、Canal 的端到端链路串起来。
- [服务深挖笔记](service-deep-dive-notes.md)：从 Confluence 和本地代码提取各服务边界、关键模块和追问点。
- [O-BFF 面试专项笔记](o-bff-interview-notes.md)：专门准备 O-BFF 的定位、分层、审批、批量、数据修复、proxy 和追问回答。
- [Marketing-Data 面试专项笔记](marketing-data-interview-notes.md)：专门准备 Marketing-Data 的离线数据链路、batch 幂等、DataSource handler、ES 查询和追问回答。
- [资料来源索引](source-index.md)：记录已查到的 Confluence 页面和本地代码路径，方便后续继续补。
- [面试口径与边界](interview-speaking-notes.md)：准备“我们组负责 / 我参与支持 / 我能讲清链路”的表达方式，并标出容易被追问的风险点。

## 重点复习顺序

1. 先看 `other-business-overview.md` 的“服务地图”和“核心链路”。
2. 打开 `insurance-business-map.html`，先把整体架构层次和几条主链路看顺。
3. 再看 `business-flow-map.md`，准备端到端流程题。
4. 最后看 `interview-speaking-notes.md` 的 30 秒和 1 分钟口述版。
5. O-BFF 和 Marketing-Data 被追问时，优先看专项笔记：
   - [O-BFF 面试专项笔记](o-bff-interview-notes.md)
   - [Marketing-Data 面试专项笔记](marketing-data-interview-notes.md)
6. 如果面试官追问更细节，再回到对应项目目录：
   - Marketing 深挖：[../Marketing/README.md](../Marketing/README.md)
   - O-BFF 代码：`/Users/si.chen/GolandProjects/insurance-operator-bff`
   - Marketing-Data 代码：`/Users/si.chen/GolandProjects/insurance-marketing-data`
   - Canal 代码：`/Users/si.chen/GolandProjects/canal-adapter`

## 信息来源

- Confluence：`OBFF 功能分析`、Marketing / O-BFF / Promotion 相关需求分析和 TD 检索结果。
- Confluence：`B-BFF 功能分析`、`CBFF功能分析`、`LBFF功能分析`、`EC Buyer业务梳理`、`[Marketing] 运营计划架构`、`insurance-marketing-data 技术设计`、`ES 使用分享`。
- 本地 Project KB：`insurance-operator-bff/project-kb/`、`insurance-marketing/project-kb/`。
- 本地代码：`insurance-operator-bff`、`insurance-marketing`、`insurance-marketing-data`、`canal-adapter`、`insurance-business-bff`、`insurance-consumer-bff`、`insurance-product`、`insurance-order`、`insurance-policy`、`insurance-promotion`、`insurance-account`、`insurance-external-aggregator`。

## 真实性边界

- 可以说“我所在组负责这些业务域，我参与维护/支持/开发其中的后台和营销链路”。
- 不建议说“我个人独立负责所有服务”，除非你能拿出具体 PR、排障记录或上线材料。
- 当前缺少你个人在 OBFF、Marketing-Data、Canal 上的具体任务证据，先统一标记为 `待补充`。
