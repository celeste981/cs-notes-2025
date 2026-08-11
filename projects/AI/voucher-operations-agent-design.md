# Voucher Operations Agent：生产级 Agent Backend 项目设计

> 状态：求职作品项目规划，尚未宣称已上线或已产生业务指标。  
> 目标：证明“Senior Backend Engineer 能把 Agent 从 Demo 做成 Production System”。

## 一句话定位

这是一个面向 Voucher 运营场景的受控 Agent Backend：它通过有权限边界的业务 Tool、可持久化工作流和知识检索，协助运营人员查询、分析和生成配置草稿；所有高风险写操作需要人工审批，并支持审计、重试、超时、恢复、评估和成本治理。

## 为什么选这个题目

- 直接连接 Promotion / Voucher 业务背景，避免做与简历脱节的通用聊天机器人。
- 能同时展示业务理解、Golang/微服务经验迁移、Agent Workflow 与生产工程设计。
- 任务既有查询/分析，也有带风险的配置草稿/提交，天然适合讲权限、审批、幂等和审计。

## 用户任务与边界

| 用户任务 | Agent 能做什么 | 风险控制 |
| --- | --- | --- |
| 查询配置、规则、库存、Budget | 查询业务 Tool，汇总结构化结果与来源 | 行级/字段级权限；敏感字段脱敏 |
| 生成 Voucher 配置草稿 | 按明确模板生成草稿和校验清单 | 仅生成 Draft，不直接生效 |
| 检测配置冲突 | 调用规则校验 Tool，解释冲突及建议 | 校验结果保留规则版本与证据 |
| 分析发券失败、核销异常 | 查询任务/日志/指标的脱敏聚合数据，提出排查路径 | 不向模型暴露凭证、PII 或任意生产日志 |
| 查询运营知识 | RAG 检索受治理的业务文档并附来源 | 只检索用户可见的知识；外部内容视为不可信数据 |
| 提交配置变更 | 仅在审批通过后异步执行 | Human Approval、幂等键、审计和可回滚/补偿策略 |

非目标：直接开放数据库写入、由模型决定权限、让模型自由执行任意 SQL / Shell、以聊天 UI 作为核心亮点。

## 架构

```text
Operator / Admin UI
        |
Agent API (FastAPI 或 Go Gateway)
        |
AuthN/AuthZ + Request Policy + Rate Limit
        |
Agent Orchestrator / Workflow Engine
  |       |          |            |
State    LLM       Tool Registry  Knowledge Service
Store    Gateway     |             |
  |       |      Voucher Tools     RAG Retrieval
  |       |      Rule/Budget/      (source, owner,
  |       |      Inventory/Task     version, permission)
  |       |
Async Worker / Queue / Scheduler
        |
Approval Service + Audit Log + Trace/Eval/Metric
```

### 核心组件

| 组件 | 职责 | 面试可讲点 |
| --- | --- | --- |
| Agent API | 接收任务、鉴权、创建 task id、返回同步结果或异步句柄 | API 边界、限流、幂等、SSE/轮询取舍 |
| Orchestrator | 控制 Agent Loop 与确定性 Workflow；校验 tool call；管理预算与终止条件 | LLM 是决策器，不是权限与执行器 |
| State Store | 保存 task、step、tool result 摘要、approval 状态、checkpoint | 长任务恢复、可回放、避免重复执行 |
| Tool Registry | 定义 schema、权限、risk level、timeout、retry policy、owner | Tool schema / 版本治理 / 最小权限 |
| Knowledge Service | 检索运营规则和 SOP，并返回来源、版本、权限信息 | RAG 质量、引用、知识治理 |
| Async Worker | 执行慢查询、批量校验、提交与补偿 | timeout、retry、DLQ、状态机 |
| Approval & Audit | 审批高风险动作、记录谁在何时批准了什么 | HITL、合规、责任可追溯 |
| Observability & Eval | Trace、指标、日志、回放、离线评估与成本统计 | 不能只看 HTTP 200，要看任务正确性 |

## Workflow 与状态机

优先将稳定路径做成确定性 Workflow；只有“需要判断下一步查什么、如何归纳”的环节交给 LLM。这样能减少随机性，也更容易测试和恢复。

```text
RECEIVED
-> AUTHORIZED
-> PLANNING
-> TOOL_RUNNING <-> WAITING_RETRY
-> NEED_APPROVAL -> APPROVED -> ASYNC_EXECUTING
-> SUCCEEDED / FAILED / CANCELLED / EXPIRED
```

关键状态字段建议：`task_id`、`request_id`、`user_id`、`tenant/scope`、`workflow_version`、`status`、`current_step`、`checkpoint`、`idempotency_key`、`approval_id`、`tool_call_id`、`cost_budget`、`trace_id`。

### 失败与恢复

- 每个副作用 Tool 必须有 `idempotency_key`；恢复任务时先查询该 key 的历史结果。
- 超时不等于失败：Worker 应区分“请求未发出、已发出但结果未知、已完成未回写”。
- 可重试错误按退避策略重试；不可重试错误返回结构化 `error_code + suggestion` 给 Workflow。
- 重试耗尽进入 `FAILED` 或 `MANUAL_REVIEW`，不要让 Agent 无限循环。
- 所有状态转换记录事件；借助 checkpoint 从最近安全节点继续，而不是从头再次执行。

## Tool 设计原则

Tool 不只是一个函数名；它是一份受治理的后端契约。每个 Tool 至少应具有：输入/输出 schema、owner、版本、风险级别、权限要求、超时、重试策略、幂等要求、审计字段和错误码。

示例：

| Tool | 类型 | 风险级别 | 是否需审批 |
| --- | --- | --- | --- |
| `get_voucher_config` | Read | Low | 否 |
| `check_voucher_conflicts` | Read / Compute | Low | 否 |
| `query_budget_and_inventory` | Read | Medium（数据范围） | 否，需权限 |
| `generate_voucher_draft` | Draft | Medium | 否，输出不生效 |
| `submit_voucher_change` | Write | High | 是 |
| `analyze_issue` | Read / Analysis | Medium | 否，使用脱敏数据 |

MCP 可作为 Tool 接入协议：MCP Server 提供能力发现与调用，但不会绕过本系统的鉴权、审批、审计和执行策略。业务 Tool 的最终执行权仍在 Orchestrator / Tool Gateway。

## RAG 与 Memory

### RAG

- 知识来源：Voucher 规则、运营 SOP、配置说明、常见异常 Runbook；数据需脱敏并有 `source/owner/version/update_time/permission` metadata。
- 检索路径：query rewrite（可选）→ permission filter → retrieval → rerank（可选）→ 引用来源生成。
- 质量检查：无可信来源时明确回答“不确定，需要人工确认”，不让模型补全事实。
- 评估指标：检索命中率、引用覆盖率、答案正确率、过期知识命中率、P95 延迟。

### Memory

- 短期：当前任务的 workflow state 和必要 tool observation。
- 长期偏好：用户授权后保存的稳定展示偏好，不能混入业务事实。
- 事实知识：进入受治理 RAG，不存入随意的对话摘要。
- 过程记录：用于审计与评估，访问需受权限与保留期策略约束。

## 权限、安全与审计

- 将 LLM 输出和检索到的外部文本都视为不可信输入；先 schema 校验，再做业务校验和权限校验。
- RBAC/ABAC 应同时约束用户、租户、业务范围、字段与动作；不能因模型“建议”而提升权限。
- Read Tool 也要最小化返回数据，默认脱敏并限制数据范围。
- Write Tool 强制 HITL：展示变更 diff、影响范围、规则校验结果和审批人；审批绑定 `task_id + payload hash + version`，避免审批后篡改。
- Audit Log 至少记录 actor、request、approval、tool、参数摘要、结果、时间、trace id、策略与版本；避免记录凭证或原始 PII。
- 防 Prompt Injection：把检索文本标注为 data，不允许其中指令改变系统策略或触发写操作。

## 观测、评估与成本

### Trace / Metrics

- 每个任务串联 `trace_id`；每一 step 记录模型、prompt 版本、工具、耗时、状态变化和结果摘要。
- 过程指标：task completion rate、step count、tool success rate、retry/fallback rate、approval latency、async queue latency。
- 质量指标：配置冲突检出率、失败归因准确率、引用覆盖率、人工纠正率。
- 体验/成本：TTFT、端到端 P50/P95、input/output token、cost per task、超预算终止率。

### Evaluation

- 建立脱敏 Golden Set：配置冲突、预算不足、库存不足、重复发券、规则不匹配、发券/核销异常等案例。
- 分层评测：Tool contract unit test → Workflow integration test → Task-level offline eval → 人工抽样评审。
- Prompt、Tool schema、模型版本均需版本化；变化先跑回归集，再灰度。

### 成本治理

- 按任务类型设置最大 step、最大 token、总耗时和模型预算。
- 能用确定性规则/普通 API 完成的步骤，不调用 LLM。
- 大工具结果先服务端聚合/摘要；RAG 按需检索，避免无差别塞 context。
- 高成本模型只用于复杂规划/归纳；简单路由、字段抽取或降级答复使用更低成本模型或规则。

## MVP 实施顺序

1. **最小闭环**：只读 `get_voucher_config`、`check_voucher_conflicts`、RAG 查询、Agent trace；用模拟/脱敏数据完成端到端任务。
2. **状态与异步**：任务表、step 事件、checkpoint、超时/重试、长任务查询接口与恢复逻辑。
3. **受控写操作**：草稿生成、变更 diff、Human Approval、幂等提交与 Audit Log。
4. **生产化证据**：Golden Set、评估报告、Dashboard、故障注入、模型 fallback、成本预算。

## 1 到 2 分钟面试表达

我会做一个 Voucher Operations Agent，而不是普通聊天机器人。它服务运营侧的真实任务：查询配置和库存、检查规则冲突、分析发券或核销异常、生成配置草稿。核心设计是把 LLM 限定在规划和归纳上，所有业务操作都通过受治理的 Tool Gateway 执行。

系统会把任务建模成可持久化的 Workflow，每一步都有状态、超时、重试和 checkpoint。涉及写配置时，Agent 只能生成草稿和影响分析，真正提交必须经过权限校验和 Human Approval；审批会绑定具体 payload hash，避免审批后内容被替换。所有工具调用、状态转换、模型成本和审批记录会串到同一个 trace 里。

知识部分不是简单向量检索：每个文档带来源、owner、版本、更新时间和权限，回答要附引用。评估上会准备 Voucher 冲突、预算不足、库存不足和异常排查等脱敏案例，分别测工具契约、Workflow 和最终任务成功率。这个项目想证明的是，我能把已有的分布式后端经验迁移到 Agent 的可靠性、安全和可运维性上。

## 待补充

- 可公开或可本地模拟的 Voucher 字段、规则、状态和异常案例。
- 技术选型：Python/FastAPI 与 Go 服务的具体边界，工作流引擎、队列、存储与向量检索组件。
- MVP 的代码仓库、部署方式、演示脚本和真实评估数据。
- 个人实际实现内容、完成时间、可量化结果；完成前不得写入简历为线上成果。
