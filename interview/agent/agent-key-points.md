# Agent 高频重难点整理


可视化页面：[agent-key-points.html](agent-key-points.html)
> 这份笔记把 Agent 面试的主线压缩成可以复述的版本。更完整题目见 [agent-topic-index.md](agent-topic-index.md)。

## 1. Context Engineering

### 先讲人话

Agent 的表现不只取决于 prompt 写得好不好，而取决于模型在当前请求里能看到哪些信息：系统规则、用户问题、历史对话、工具结果、检索文档、用户权限、业务状态。Context Engineering 就是在管理这些信息。

### 再讲原理

Context Window 有长度限制，信息还会互相干扰。系统要决定：

- 哪些信息必须固定保留，比如 system prompt、工具 schema、任务目标。
- 哪些信息可以摘要，比如历史对话。
- 哪些信息按需检索，比如知识库文档。
- 哪些信息必须过滤，比如 PII、外部网页里的恶意指令。

常见策略有 Sliding Window、Progressive Summarization、Selective Retrieval、Context Distillation、Map-Reduce。

### 面试怎么说

我会把 context 当成一层数据架构来设计，而不是只拼字符串。核心目标是让模型看到足够、准确、低噪音、可追溯的信息，同时控制 token 成本。系统 prompt 和工具 schema 不被压缩，历史对话可以摘要，长文档走检索和分段摘要，工具结果要先去噪再注入。

## 2. Agent Loop 与 Observation

### 先讲人话

Chain 是固定流程，Agent 是边执行边看结果再决定下一步。Observation 就是 Agent 看真实世界的步骤，没有 Observation，Agent 只是在猜。

### 再讲原理

典型循环：

```text
User Input
-> LLM 决策
-> Tool Call
-> Tool 执行
-> Observation 注入 context
-> LLM 再决策
-> 直到输出最终答案
```

因为模型有非确定性，所以必须有防护：

- `max_steps`：最大执行步数。
- 重复动作检测：相同 action 连续出现时中止。
- 工具失败熔断：同一工具连续失败后换策略。
- 目标达成检查：不要只靠模型自己说结束。
- 总时间和成本预算：防止无限消耗。

### 面试怎么说

我会记录每一步的 thought/action/tool_result/latency/cost，并设置步骤预算、时间预算和失败熔断。简单任务不走完整 ReAct；工具之间没有依赖时用并行工具调用；结果过大时先摘要再回填给模型。

## 3. Function Calling 与 Tool Schema

### 先讲人话

LLM 不是真的执行函数。它只是生成一个结构化请求：调用哪个工具、参数是什么。真正执行函数、鉴权、校验、重试的是外层 Orchestrator。

### 再讲原理

完整链路：

```text
工具 schema 注入
-> LLM 选择工具并生成 arguments
-> Orchestrator 校验 JSON 和业务参数
-> 执行工具
-> 工具结果整理成 tool message
-> LLM 基于结果继续推理
```

高质量 Tool Schema 要说明：

- 工具什么时候用，什么时候不要用。
- 参数类型、范围、必填约束。
- 错误码和可恢复建议。
- 返回值要给 LLM 看摘要，不要直接塞大 JSON。

### 面试怎么说

我会把 LLM 生成的 tool call 当成不可信输入处理：先做 schema 校验、权限校验和业务参数校验，再执行工具。错误结果用结构化 `error_code + message + suggestion` 返回，方便 Agent 重新规划。

## 4. RAG、Memory 与 Knowledge Service

### 先讲人话

RAG 不是“丢进向量库搜一下”。生产里要解决知识从哪里来、怎么切块、怎么召回、怎么更新、怎么判断答案是否引用了可信来源。

### 再讲原理

传统 RAG 通常是固定 Pipeline：

```text
query -> retrieve -> rerank -> generate
```

Agentic RAG 会让 Agent 自己决定是否需要检索、查哪个源、是否换 query、是否多轮检索。

Memory 需要分层：

- 短期记忆：当前任务上下文。
- 长期偏好：用户稳定偏好。
- 事实知识：业务规则、项目文档。
- 过程经验：过去任务成功/失败经验。

每条知识要有 source、owner、更新时间、适用范围和置信度。

### 面试怎么说

我会先建立评估集，再调检索。指标包括命中率、引用覆盖率、答案正确率、幻觉率和延迟。知识更新频繁时优先 RAG；输出格式高度固定、知识变化小的时候才考虑 fine-tuning。

## 5. Observability、Eval 与成本

### 先讲人话

Agent 不像普通接口，只看 HTTP 200 没意义。它可能每一步都成功，但最终答案是错的。所以要监控过程，也要评估结果。

### 再讲原理

关键指标：

| 维度 | 指标 |
| --- | --- |
| 执行过程 | step_count、tool_call_count、retry_count、fallback_count |
| 工具质量 | tool_success_rate、tool_latency、tool_error_code |
| 模型成本 | input tokens、output tokens、cost per task |
| 用户体验 | TTFT、总延迟、任务完成率 |
| 答案质量 | 正确率、引用覆盖率、幻觉率、人工评分 |

Eval 应该分层：Unit Eval 测 prompt/tool schema，Task Eval 测单任务，System Eval 测完整 Agent 流程。

### 面试怎么说

我会把每次 Agent 执行记录成 trace，能回放每一步 context、工具输入输出和模型决策。Prompt 要版本化，变更走 A/B 测试和灰度；成本按任务类型设预算，超过预算就降级或停止。

## 6. 安全：Prompt Injection、PII、HITL

### 先讲人话

LLM 的根本安全问题是：它很难天然区分“数据”和“指令”。外部网页、邮件、文档里的一句话，可能被模型误当成用户命令。

### 再讲原理

防护要分层：

- 输入层：清洗、分类、PII 检测。
- Context 层：把外部内容明确标记为 untrusted data。
- Tool 层：最小权限、参数校验、敏感操作审批。
- Policy 层：高风险操作必须 HITL。
- 输出层：脱敏、合规检查、引用来源。
- 审计层：记录谁触发、模型看到什么、工具做了什么。

### 面试怎么说

我不会只靠 prompt 说“不要被注入”。工程上会把 tool result、网页内容、用户上传文件都当作不可信数据处理。涉及发邮件、删数据、转账、改生产配置这类操作时，需要人工确认和审计日志。

