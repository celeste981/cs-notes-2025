# O-BFF 架构梳理

> 可视化图解：[O-BFF 架构图解](./o-bff-architecture.html)
>
> 所属项目索引：[O-BFF 项目材料](./README.md)
>
> Insurance 总览：[Insurance 部门业务材料](../Insurance/README.md)

## 专项补充文档

如果只想背深挖点，可以直接看这三篇：

- [O-BFF 核心知识深挖](./o-bff-key-knowledge-deep-dive.md)：分层、边界、Interceptor、下游契约、排障路径。
- [O-BFF 配置化组件深挖](./o-bff-config-components-deep-dive.md)：Assembly Proxy/Import/Export、Mask、Approval Mapping、RateConfig、配置化边界。
- [O-BFF 审批 / 批量 / 任务 / 数据修复深挖](./o-bff-governance-centers-deep-dive.md)：Approval、Batch Operate、Task Center、Data Fix 的关系和可靠性追问。
- [Approval Center / Task Center 架构说明](./o-bff-approval-task-center-architecture.md)：审批中心和任务中心的职责边界、协作流程和图解。

## 一、O-BFF 是什么？一句话定位

**O-BFF / Operator BFF 是 Insurance Admin 的后台聚合和治理层**。Admin 页面统一调用 O-BFF，O-BFF 再编排 Product、Order、Policy、Marketing、Promotion、Risk Control、Account 等下游服务，同时承载审批、批量任务、数据修复、操作日志、脱敏、proxy 和 ES 查询等后台通用能力。

面试开场：

> 我们组的 O-BFF 本质上是 Insurance Admin 的后台运营聚合层。它不是简单 RPC 转发，而是把后台高风险操作收口到统一链路里：请求先进 interceptor 做权限、脱敏、操作日志、审批信息和重复请求控制，再进入 Service/Manager/Integrate 编排下游服务；高风险修改走 Approval Center，批量导入导出走 Batch Operate / Task Center，数据修复走 Data Fix Center。

个人边界：

> O-BFF 是组内业务。我能讲清整体架构和链路，但个人具体负责的页面、接口、PR 和上线记录需要按事实补充。

## 项目路径速查

| 项目 | 本机路径 | 这里看什么 |
| --- | --- | --- |
| `insurance-operator-bff` | `/Users/si.chen/GolandProjects/insurance-operator-bff` | O-BFF 主服务，Admin API、审批、批量、任务、数据修复和下游 RPC 聚合都在这里。 |
| `insurance-product` | `/Users/si.chen/GolandProjects/insurance-product` | 产品、计划、规则、保费、产品工厂等下游 owner。 |
| `insurance-order` | `/Users/si.chen/GolandProjects/insurance-order` | 订单、支付、账单、invoice 等下游 owner。 |
| `insurance-policy` | `/Users/si.chen/GolandProjects/insurance-policy` | 保单、理赔、取消、续保等下游 owner。 |
| `insurance-marketing` | `/Users/si.chen/GolandProjects/insurance-marketing` | 营销计划、人群、推荐、PNAR、发券/通知等下游能力。 |

常用代码入口：

| 想看什么 | 推荐路径 |
| --- | --- |
| 系统总架构 | `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/System_Architecture.md` |
| 下游依赖关系 | `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/Cross_App_Overview.md` |
| Interceptor 注册顺序 | `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/system/interceptor/Agent.md` |
| 下游 API 契约 | `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/tech/Admin_Downstream_API_Contract.md` |
| 监控告警 | `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/tech/OBFF_Monitoring_Alerts.md` |
| Service/Manager/Biz/Integrate | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/service/`, `/src/manager/`, `/src/biz/`, `/src/integrate/` |

## 二、核心抽象：后台操作 Operation

面试时可以把 O-BFF 的核心对象抽象成一次后台操作：

```text
Admin Operation
├── Who: operator / permission / SOUP role
├── What: query / create / update / import / export / data-fix
├── Risk: 是否高风险，是否需要审批
├── Execution: 同步 RPC / proxy / task / reliable event
├── Audit: 操作日志 / 更新历史 / 安全日志 / 审批历史
└── Result: 页面 response / batch result / download file / error detail
```

这个抽象能解释为什么 O-BFF 不是普通 CRUD：

- 后台操作要知道谁在操作、有没有权限。
- 高风险操作要先审批。
- 批量和长耗时操作要异步任务化。
- 结果要可追踪、可下载、可审计。
- 下游服务才是业务事实 owner，O-BFF 做入口和编排。

## 三、完整数据流：四层链路

```text
┌──────────────┐    ┌────────────────┐    ┌────────────────────┐    ┌───────────────┐
│ Admin Portal │ -> │ Interceptor 链 │ -> │ Service/Manager/Biz│ -> │ Downstream/RDB │
└──────────────┘    └────────────────┘    └────────────────────┘    └───────────────┘
         │                    │                       │                       │
         │                    │                       ├─ Approval Center       │
         │                    │                       ├─ Batch Operate         │
         │                    │                       ├─ Task Center           │
         │                    │                       ├─ Data Fix Center       │
         │                    │                       └─ ES / Repo / Event     │
```

### Stage 1: Interceptor 链做公共治理

`src/interceptor/a_initialator.go` 注册了 gRPC server interceptor。Project KB 里记录的顺序是：

| 顺序 | Wrapper | 作用 |
| --- | --- | --- |
| 1 | `NewRepeatedHandlerWrapper` | 重复请求处理，降低重复提交风险。 |
| 2 | `NewCommonParamFillingHandlerWrapper` | 填充公共参数。 |
| 3 | `NewOperateLogHandlerWrapper` | 记录操作日志。 |
| 4 | `NewUpdateHistoryHandlerWrapper` | 记录更新历史。 |
| 5 | `NewMaskHandlerWrapper` | 响应脱敏。 |
| 6 | `approve_info_interceptor.NewApproveInfoHandlerWrapper` | 注入审批信息上下文。 |
| 7 | `NewRpcProxyHandlerWrapper` | 旧版 RPC proxy。 |
| 8 | `NewRpcNewProxyHandlerWrapper` | 新版 RPC proxy。 |
| 9 | `NewRpcOrchestrationHandlerWrapper` | RPC 编排代理。 |
| 10 | `NewSecurityLogsRetentionHandlerWrapper` | 安全日志留存。 |
| 11 | `NewPrintPanicStackWrapper` | panic 堆栈打印。 |

面试说法：

> Interceptor 是 O-BFF 的后台治理入口。权限、日志、脱敏、重复请求、审批上下文和 proxy 这类横切能力不应该散落在每个 service 方法里，所以通过 wrapper 链统一处理。

### Stage 2: Service / Manager / Biz / Integrate 分层编排

```text
Service
  -> 接 Admin API，做协议适配
Manager / Biz
  -> 编排后台业务，组合多个下游调用或中心模块
Integrate
  -> 封装 Product / Order / Policy / Marketing 等下游 RPC
Repo / ES DAO
  -> 访问 O-BFF 自身表、任务表、ES 查询视图
```

面试说法：

> Service 不应该写复杂业务，Manager/Biz 做编排，Integrate 隔离下游 API。这样下游 proto 或调用方式变化时，不会把影响扩散到页面入口。

### Stage 3: 后台中心模块承接高风险能力

| 中心模块 | 面试定位 | 代表路径 |
| --- | --- | --- |
| Approval Center | 通用审批/审核工作流，管理流程定义、节点推进、审批实例和审批事件。 | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/approval_center/` |
| Batch Operate Center | 批量操作记录的上传、列表、下载、审批提交和审批结果处理。 | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/batch_operate_center/` |
| Task Center | 配置驱动批量任务框架，支持 CSV/文件处理、S3、回调、监控和状态维护。 | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/task_center/` |
| Data Fix Center | 数据修复工具，支持创建修复记录、审批、可靠事件异步执行、字段元数据和辅助查询。 | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/data_fix_center/` |
| Operator Executer Center | 把批量操作、审批结果和批处理引擎衔接起来。 | `/Users/si.chen/GolandProjects/insurance-operator-bff/src/operator_executer_center/` |

### Stage 4: 下游 owner 服务负责领域事实

O-BFF 会调很多下游，但面试时要强调边界：

| 下游 | O-BFF 使用场景 | 事实边界 |
| --- | --- | --- |
| Product | 产品、计划、规则、声明、定价、产品工厂 | 产品配置和费率语义归 Product。 |
| Policy | 保单详情、理赔、取消/续保、保单对账、数据修复 | 保单状态机和保单事实归 Policy。 |
| Order | 订单、支付、账单、invoice、financial bill | 支付、金额、订单状态归 Order。 |
| Account | accountId、平台用户映射、用户查询 | 用户映射归 Account。 |
| Promotion / Voucher | 活动、券、预算、优惠、白名单 | 发券、预算、库存归 Promotion/Voucher。 |
| Marketing | 营销推荐、人群、活动产品配置、referral | 营销计划和人群语义归 Marketing。 |

## 四、核心流程

### 流程 1：批量操作审批流程

```text
Admin 上传文件/提交批量操作
  -> O-BFF 校验文件和操作类型
  -> 创建 batch operate 记录
  -> 发起 Approval Center 审批
  -> 审批通过
  -> Operator Executer / Task Center 异步执行
  -> 调下游批量 API 或执行文件任务
  -> 写任务状态、成功/失败数量、错误明细
  -> Admin 下载结果
```

面试抓手：

- 批量操作不能同步卡住页面。
- 高风险批量修改需要审批。
- 失败要定位到单条数据，常见字段是 `success_count`、`failed_count`、`error_data`、`error_msg`。
- 下游 API 要支持批量请求、局部失败、幂等和错误定位。

### 流程 2：Task Center 导入/导出流程

```text
Admin 发起导入/导出
  -> 读取 batch_task_config
  -> 按 remote_type 选择执行路径
  -> 下载/解析 CSV 或分页/scroll 查询
  -> 调下游 RPC / assembly / JSON RPC
  -> 上传结果到 S3
  -> 写 task_record / callback / monitor
```

面试抓手：

- Task Center 是配置驱动的批处理引擎。
- 它和 Batch Operate 的边界不同：Batch Operate 偏运营记录和审批，Task Center 偏文件处理、远程调用和任务状态。
- 大导出要看分页、scroll/search_after、稳定排序和超时。

### 流程 3：数据修复流程

```text
Admin 创建 data fix 记录
  -> 填修复模块、字段、业务 ID、修复明细
  -> 提交审批
  -> 审批通过
  -> reliable-event 异步执行
  -> 调 Policy / UserCenter / Spex / BatchOperate 等下游
  -> 更新修复状态和失败项
  -> 支持查询和失败重试
```

面试抓手：

- 数据修复不能说成“直接改库”。
- 正确说法是工具化、审批化、异步化、可追踪。
- owner 服务仍负责自己的数据语义，O-BFF 负责入口和流程。

### 流程 4：Proxy / Orchestration 流程

```text
Admin 通用接口
  -> rpc_proxy / rpc_new_proxy / orchestration interceptor
  -> 读取 proxy/orchestration 配置
  -> 前后置处理
  -> 下游 RPC
  -> response 加工
```

面试抓手：

- 简单标准接口可以走 proxy，减少重复样板代码。
- 复杂编排、审批、批量、资损链路不能只靠 proxy，应该显式 Manager/Biz 编排。
- 下游 API 要满足 Admin 契约：字段、分页、operator、时间、批量返回结构。

### 流程 5：ES 查询和报表流程

```text
Admin 列表/报表查询
  -> O-BFF 构造查询条件
  -> ES DAO / adapter 查询 ES
  -> 转换成 Admin 展示字段
  -> 返回页面或导出任务
```

面试抓手：

- ES 是查询视图，不是事实源。
- 数据不一致时，先看 owner MySQL，再看 Canal/Adapter/Kafka/ES mapping，最后看查询 DSL。
- 大结果集导出需要 scroll/search_after 或稳定分页。

## 五、配置化组件体系：Admin Assembly 能力

这是 O-BFF 面试里可以重点讲的一块。它不是只有手写接口，而是把 Admin 后台常见能力做成了一整套配置化组件：标准转发、批量导入、批量导出、审批映射、字段脱敏、操作记录、任务限流和文件字段转换。

一句话：

> 我们把 Admin 常见后台能力沉淀成配置化组件。简单标准场景通过 Assembly 配置接入，复杂高风险场景再写 Manager/Biz 编排，这样能降低重复开发，也能把审批、审计、脱敏、批量任务这些治理能力统一起来。

### 组件总览

| 配置化组件 | 主要解决什么 | 关键配置/入口 | 运行时链路 | 面试怎么说 |
| --- | --- | --- | --- | --- |
| Assembly Proxy / New Proxy | 标准 Admin 接口不想重复写转发代码。 | `admin_assembly_tab`、`proxy_config`、`admin-proxy-config`。 | `rpc_proxy` / `rpc_new_proxy` interceptor -> proxy handler -> 下游 RPC。 | 标准查询/更新接口走配置化转发，保留前后置 handler 扩展点。 |
| Assembly Import | CSV 批量导入、字段转换、下游批量更新。 | `batch_task_config`，`remote_type=6`，`admin-import-config`。 | 上传文件 -> Task Center -> 解析/转换 -> 下游 RPC -> 结果文件。 | 导入不是写死代码，而是配置字段映射、批次大小、并发和错误列。 |
| Assembly Export | 后台导出、分页/scroll 查询、S3 下载记录。 | `batch_task_config`，`remote_type=7`，`admin-export-config`。 | 创建导出任务 -> Page/Scroll 查询 -> 文件生成 -> S3/download record。 | 大导出要配置下载策略、游标字段、终止条件和超时保护。 |
| Approval Mapping | 标准接口也能接审批，不要求每个接口手写审批逻辑。 | `audit_biz_type_mapping`、审批配置、Assembly 配置。 | Proxy 后识别审批映射 -> 创建 batch/approval record -> 审批通过再执行。 | 高风险操作通过配置接入审批，避免绕过审计。 |
| Mask 脱敏组件 | Admin response 中 PII 字段需要统一脱敏。 | config-center `[mask.field].mask_fields`、`admin-mask-config`。 | mask interceptor -> 根据接口/字段配置脱敏 response。 | 姓名、手机号、邮箱、地址等敏感字段不在业务代码里零散处理。 |
| Operate History | 后台操作要有可追踪历史。 | operation log/update history interceptor、本地操作历史表。 | 请求进入 interceptor -> 记录操作者、动作、对象和更新前后信息。 | 支持审计、排障和事故追溯。 |
| RateConfig / Circuit | 批量任务不能无限并发打下游。 | `RateConfig`：batch size、并发、QPS、timeout、circuit。 | Task Center 执行前读取配置 -> 控制切批、并发、超时和熔断。 | 这是批量导入导出的稳定性保护。 |
| AssemblyExtraConfig | 文件列、数组字段、错误列、唯一键、下载策略要配置化。 | `file_transform_list`、`unique_field_map`、`download_strategy`、scroll/page 字段。 | 文件解析/结果生成/导出分页阶段读取配置。 | 支持不同 Admin 页面复用同一套批处理框架。 |

### 配置化接入流程

```text
Admin 新需求
  -> 判断是否是标准 proxy / import / export / mask / approval 场景
  -> 选择对应 Admin Assembly 组件
  -> 通过技能或 AssemblyTool 生成配置
  -> 写入配置表或 config-center
  -> O-BFF interceptor / Task Center / Approval Center 按配置执行
  -> 通过日志、任务记录、审批记录和监控验证
```

### 配置化和手写代码的边界

| 场景 | 推荐方式 | 原因 |
| --- | --- | --- |
| 标准查询、简单更新、字段映射清晰 | Assembly Proxy | 接入快，少写重复 glue code。 |
| CSV 导入、导出、可用统一批处理模型表达 | Assembly Import / Export | 文件、任务、S3、结果明细、监控都已有框架。 |
| 需要审批但业务执行标准 | Approval Mapping + Proxy/Batch | 审批治理统一，减少漏审风险。 |
| 涉及复杂状态机、资损、跨多个下游强编排 | 手写 Manager/Biz | 需要明确业务顺序、幂等、补偿和异常处理。 |
| 下游契约不满足分页、幂等、局部失败结构 | 先改下游契约或标记待确认 | 配置化不是兜底方案，契约不稳定会把风险转移到 O-BFF。 |

面试表达：

> 我会先判断这个 Admin 能力能不能被标准配置承接。比如纯转发走 Assembly Proxy，导入导出走 Task Center 的 Assembly Import/Export，PII 字段走 Mask 配置，高风险写操作走 Approval Mapping。只有当它涉及复杂领域状态机、资损风险或跨服务强编排时，才落到 Manager/Biz 手写逻辑。

## 六、目录结构

主服务路径：`/Users/si.chen/GolandProjects/insurance-operator-bff`

```text
src/
├── interceptor/                  # 请求公共治理：权限、日志、脱敏、proxy、安全日志
├── service/                      # O-BFF gRPC service 入口
├── manager/                      # 后台业务编排
├── biz/                          # 可复用业务逻辑
├── integrate/                    # 下游 RPC client
├── repo/                         # DB/ES repo
├── approval_center/              # 审批中心
├── batch_operate_center/         # 批量操作中心
├── task_center/                  # 批量任务中心
├── data_fix_center/              # 数据修复中心
├── operator_executer_center/     # 批量操作到任务执行的串联层
├── event/                        # 审批事件、可靠事件
├── processor/                    # 后台处理器/产品检查器
└── common/ model/ config/ util/  # 支撑模型和工具
```

## 七、设计哲学

### 1. 后台治理集中化

权限、日志、脱敏、审批、批量、数据修复、任务状态这些能力放在 O-BFF 统一治理，避免每个 Admin 页面重复实现。

### 2. BFF 只做聚合和编排

O-BFF 不应该抢 Product、Order、Policy 的领域事实。它适合做字段适配、下游调用组合、后台流程控制和审计。

### 3. 高风险操作审批化

后台写操作、批量导入、数据修复都可能产生资损或用户影响，所以要先审批，再异步执行，并保留完整记录。

### 4. 长耗时任务异步化

批量导入导出、数据修复、报表下载不应同步阻塞 Admin 页面。Task Center 用 task record、S3、callback、monitor 做结果管理。

### 5. 通用接口配置化，复杂接口代码化

Proxy / orchestration 适合标准接口。涉及审批、批量、状态机、资金或复杂聚合时，仍然要写清楚 Manager/Biz 编排。

## 八、面试叙事顺序

1. **定位**：O-BFF 是 Insurance Admin 的后台聚合和治理层。
2. **请求链路**：Admin -> interceptor -> service -> manager/biz -> integrate/repo -> downstream。
3. **核心能力**：审批、批量、任务、数据修复、proxy、ES 查询。
4. **配置化组件**：Assembly Proxy、Import、Export、Mask、Approval Mapping、Operate History。
5. **批量流程**：上传/审批/异步执行/结果下载/错误明细。
6. **数据修复**：工具化、审批化、可靠事件异步执行。
7. **边界意识**：下游 owner 服务负责领域事实，O-BFF 负责 Admin 场景治理。
8. **排障思路**：先看入口路由和 O-BFF 日志，再看中心模块记录，最后看下游 owner、ES/Canal 或配置。

## 九、常见追问

### O-BFF 和普通 BFF 有什么区别？

普通 BFF 更偏页面聚合。O-BFF 是后台运营 BFF，除了聚合下游数据，还要处理审批、批量、数据修复、脱敏、操作日志和审计。后台操作风险更高，所以治理能力更重。

### 批量操作和 Task Center 有什么区别？

Batch Operate Center 偏“运营提交的一次批量操作记录”，包括上传、审批、状态、撤回、下载。Task Center 偏“配置驱动的批处理执行引擎”，负责文件、S3、远程调用、回调和监控。两者会通过审批和执行器串起来。

### 为什么不能直接在 O-BFF 改数据？

因为数据 owner 是下游领域服务。O-BFF 直接改库会绕过状态机、校验、审计和领域约束。正确做法是 O-BFF 提供后台入口和审批流程，执行时调用 owner 服务提供的修复或更新接口。

### 下游 API 不满足 Admin 契约怎么办？

TD 里继续产出设计，但要在待确认事项写清差异、影响和期望下游改造；批量导入/导出配置不能直接发布，直到下游接口满足契约或业务明确接受兜底方案。

### ES 查询出问题怎么排？

先确认 owner MySQL 有没有数据，再看 Canal instance、Kafka/Adapter、ES mapping、sync script，最后看查询 DSL。不能只看 ES 判断业务数据不存在。

## 十、可讲项目经验模板

稳妥版：

> 我参与/熟悉 Insurance Admin 后台运营链路。O-BFF 负责把 Admin 页面操作统一接入到后台治理链路里：先通过 interceptor 处理权限、日志、脱敏和重复请求，再由 service/manager/integrate 调下游服务。我们还沉淀了一套配置化组件，标准接口走 Assembly Proxy，批量导入导出走 Assembly Import/Export 和 Task Center，敏感字段走 Mask 配置，高风险操作通过 Approval Mapping 接审批。对于更复杂的场景，再用 Manager/Biz 做显式编排。我的重点是能把后台操作从“页面请求”讲到“配置化接入、审批、异步执行、结果追踪和下游 owner 边界”。

具体个人项目：`待补充`。

## 十一、优化项重点记忆

这些点适合面试官问“你们系统做过哪些优化 / 为什么这样设计 / 怎么保证后台操作稳定”时展开。

| 优化项 | 解决的问题 | 面试怎么讲 |
| --- | --- | --- |
| 配置化组件体系 | Admin 页面重复开发 proxy、导入、导出、脱敏、审批接入，开发成本高且治理容易漏。 | 用 Assembly Proxy / Import / Export / Mask / Approval Mapping 把标准能力配置化，复杂资损链路再手写 Manager/Biz。 |
| Interceptor 链集中治理 | 权限、脱敏、日志、重复请求、审批信息散落在业务代码里会难维护。 | O-BFF 把公共治理收敛到 interceptor，业务 service 更专注场景编排。 |
| Approval Center 标准化 | 高风险后台操作如果直接执行，资损和审计风险高。 | 审批实例、节点推进、审批事件让批量操作和数据修复先审批再执行。 |
| Batch Operate + Task Center 拆层 | 批量上传、审批记录、文件执行、回调、结果下载混在一起会复杂。 | Batch Operate 管运营记录和审批，Task Center 管配置驱动执行、文件、S3、回调和监控。 |
| 下游 API 契约检查 | 下游接口分页、operator、错误结构不统一，会导致 Admin 接入成本高。 | 新增导入/导出前检查分页/游标、批量返回、`error_id`、`error_msg`、幂等和死锁风险。 |
| 大结果集导出约束 | MySQL 深分页、ES 大查询容易超时或重复/丢数。 | 导出设计要明确 Page/Scroll/SearchAfter、稳定排序字段和终止条件。 |
| 监控 label 规范 | 高基数 label 或敏感信息进 Prometheus 会造成监控和合规问题。 | 指标保留 namespace/service/method/status 等低基数字段，不放 userId、policyId、requestId、手机号、原始错误。 |
| ES 同步扁平化和冗余键 | ES 同步阶段跨库 join 或推导关联键不稳定，补偿复杂。 | 新增 ES 同步需求时推动源数据冗余稳定关联键，ES 索引保持扁平结构。 |

一句话总结：

> O-BFF 的优化重点不是单点性能，而是后台操作治理：公共能力集中化、高风险动作审批化、批量任务异步化、下游契约标准化、监控可观测化。

## 十二、资料来源

- Confluence：`OBFF 功能分析`，页面 ID `2968425075`。
- 本地 KB：`/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/System_Architecture.md`
- 本地 KB：`/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/Cross_App_Overview.md`
- 本地 KB：`/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/system/interceptor/Agent.md`
- 本地 KB：`/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/system/approval_center/Agent.md`
- 本地 KB：`/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/system/batch_operate_center/Agent.md`
- 本地 KB：`/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/system/task_center/Agent.md`
- 本地 KB：`/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/system/data_fix_center/Agent.md`
- 本地 KB：`/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/tech/Admin_Downstream_API_Contract.md`
- 本地 KB：`/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/tech/OBFF_Monitoring_Alerts.md`
- 本地 KB：`/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/decisions/ADR-0001-es-sync-flat-index-and-redundant-keys.md`
- 本地 KB：`/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/decisions/ADR-0002-admin-current-assembly-configuration.md`
- 本地 KB：`/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/system/Admin_Platform_Component_Map.md`
