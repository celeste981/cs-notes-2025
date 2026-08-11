# O-BFF 项目材料

这个目录用于沉淀 Insurance O-BFF / Operator BFF 的项目背景、架构梳理、后台治理能力和面试表达。复习时优先抓住一条主线：`Insurance Admin` -> `O-BFF Interceptor` -> `Service/Manager/Integrate` -> `Approval/Batch/Task/DataFix` -> `下游 owner 服务`。

## 当前文档

- [O-BFF 架构梳理](o-bff-architecture.md)：核心项目文档，适合回答后台 BFF、配置化组件、审批、批量任务、数据修复、proxy/orchestration、ES 查询和下游聚合。
- [O-BFF 架构图解](o-bff-architecture.html)：上面架构文档的可视化版本，适合先建立整体图像。
- [O-BFF 核心知识深挖](o-bff-key-knowledge-deep-dive.md)：专门背分层、边界、Interceptor、下游契约和排障路径。
- [O-BFF 配置化组件深挖](o-bff-config-components-deep-dive.md)：专门背 Assembly Proxy/Import/Export、Mask、Approval Mapping、RateConfig 和配置化边界。
- [O-BFF 审批 / 批量 / 任务 / 数据修复深挖](o-bff-governance-centers-deep-dive.md)：专门背 Approval、Batch Operate、Task Center、Data Fix 的关系和可靠性追问。
- [Approval Center / Task Center 架构说明](o-bff-approval-task-center-architecture.md)：专门解释审批中心和任务中心的架构、职责边界、协作流程和面试讲法。
- [Approval Center / Task Center 架构图解](o-bff-approval-task-center-architecture.html)：上面说明文档的可视化版本。
- [Insurance 总览中的 O-BFF 面试专项](../Insurance/o-bff-interview-notes.md)：偏面试问答和防守边界。

## 本机项目路径

| 项目 | 本机路径 | 面试定位 |
| --- | --- | --- |
| `insurance-operator-bff` | `/Users/si.chen/GolandProjects/insurance-operator-bff` | Insurance Admin 后台聚合入口，承接审批、批量、任务、数据修复、ES 查询和下游服务编排。 |

> 注意：之前的历史 clone 路径不要再使用；当前只记录 `/Users/si.chen/GolandProjects/insurance-operator-bff`。

## 代码入口速查

| 想看什么 | 推荐路径 |
| --- | --- |
| 系统架构 KB | `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/System_Architecture.md` |
| 跨应用关系 | `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/Cross_App_Overview.md` |
| Interceptor 链 | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/interceptor/` |
| Service 入口 | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/service/` |
| Manager 编排 | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/manager/` |
| 下游 RPC 封装 | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/integrate/` |
| Approval Center | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/approval_center/` |
| Batch Operate Center | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/batch_operate_center/` |
| Task Center | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/task_center/` |
| Data Fix Center | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/data_fix_center/` |
| 配置化组件地图 | `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/system/Admin_Platform_Component_Map.md` |
| 当前 Assembly 配置 ADR | `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/decisions/ADR-0002-admin-current-assembly-configuration.md` |

## 面试复习顺序

1. 先看 [O-BFF 架构图解](o-bff-architecture.html)，记住 `Admin 请求 = 拦截治理 + 配置化组件 + 聚合编排 + 中心化后台能力`。
2. 再看 [O-BFF 架构梳理](o-bff-architecture.md)，补充代码路径、核心流程、设计亮点和追问答案。
3. 单独看 [Approval Center / Task Center 架构图解](o-bff-approval-task-center-architecture.html)，把“审批管能不能做，任务管怎么执行”讲顺。
4. 单独背三个深挖专项：[核心知识](o-bff-key-knowledge-deep-dive.md)、[配置化组件](o-bff-config-components-deep-dive.md)、[治理中心](o-bff-governance-centers-deep-dive.md)。
5. 最后看 [O-BFF 面试专项](../Insurance/o-bff-interview-notes.md)，准备 1 分钟回答和防守边界。

## 待补充事实

- 你在 O-BFF 上实际负责过的 Admin 页面/API/批量任务/审批需求。
- 具体 PR、Jira、上线记录、排障记录。
- 某个真实 Admin 页面到 O-BFF service / manager / integrate 的接口名。
