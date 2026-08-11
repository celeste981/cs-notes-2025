# Agent / LLM 题目索引


可视化页面：[agent-topic-index.html](agent-topic-index.html)
> 来源：[`../lc100/lc100-topic-index.md`](../lc100/lc100-topic-index.md) 中的 LLM / Agent 部分。  
> 题量：68 道。  
> 用法：先扫标题，遇到讲不清的题再回到 `agent-key-points.md` 或 LC100 原索引补细节。

## LLM 应用工程

| ID | 难度 | 主题 | 复习状态 |
| --- | --- | --- | --- |
| Q1-0a | L2 | System/User/Assistant/Tool 四种 Message Role 的设计意图 | 未整理 |
| Q1-0b | L2/L3 | Few-shot Prompting 的工程边界，什么时候不如 fine-tuning | 未整理 |
| Q1-0c | L2/L3 | Reasoning Models 和标准 Chat Models 在 Agent 中的区别 | 未整理 |
| Q1-1 | L2/L3 | Context Engineering 与 Prompt Engineering 的区别 | 已成稿 |
| Q1-2 | L3 | Context Window 满了之后的处理策略和代价 | 已成稿 |
| Q1-3 | L2/L3 | Structured Output 的价值和 schema 约束 | 未整理 |
| Q1-4 | L3 | Agent Streaming 响应的工程挑战 | 未整理 |
| Q1-5 | L2/L3 | Token 成本如何影响 Agent 架构 | 已成稿 |
| Q1-6 | L3 | 跨模型工具调用兼容设计 | 未整理 |
| Q1-7 | L2 | Temperature 和采样参数的工程实践 | 未整理 |

## Agent 推理与规划

| ID | 难度 | 主题 | 复习状态 |
| --- | --- | --- | --- |
| Q2-0a | L2 | Agent 和 Chain 的本质区别 | 已成稿 |
| Q2-0b | L2/L3 | Observation 为什么是必须的 | 已成稿 |
| Q2-0c | L2 | Agent Loop 基本结构 | 已成稿 |
| Q2-1 | L2/L3 | ReAct 模式的局限和替代优化 | 未整理 |
| Q2-2 | L3 | CoT、ToT、Self-Consistency 的差异和代价 | 未整理 |
| Q2-3 | L2/L3 | Plan-and-Execute vs Dynamic Re-planning | 未整理 |
| Q2-4 | L3 | Reflection 机制如何避免“假反思” | 未整理 |
| Q2-5 | L2/L3 | Agent 无限循环防护 | 已成稿 |
| Q2-6 | L3 | Human-in-the-Loop 如何暂停和恢复 | 未整理 |
| Q2-7 | L2/L3 | Agent Memory 的记忆衰退问题 | 未整理 |
| Q2-8 | L3 | 高可靠场景的决策确认机制 | 未整理 |

## 工具系统工程

| ID | 难度 | 主题 | 复习状态 |
| --- | --- | --- | --- |
| Q3-0a | L2/L3 | Function Calling 完整链路 | 已成稿 |
| Q3-0b | L2/L3 | Tool result 给代码和给 LLM 的区别 | 已成稿 |
| Q3-1 | L3 | 工具从 10 个增长到 100 个时的架构演进 | 未整理 |
| Q3-2 | L2/L3 | 高质量 Tool Schema 设计原则 | 已成稿 |
| Q3-3 | L3 | 并行工具调用和边界条件 | 未整理 |
| Q3-4 | L2/L3 | 工具调用错误处理 Pipeline | 已成稿 |
| Q3-5 | L3 | Code Execution 工具安全设计 | 未整理 |
| Q3-6 | L2/L3 | 工具级 RBAC / ABAC 权限设计 | 未整理 |
| Q3-7 | L3 | 工具调用可观测系统和关键指标 | 未整理 |

## 记忆与知识系统

| ID | 难度 | 主题 | 复习状态 |
| --- | --- | --- | --- |
| Q4-0a | L2/L3 | Embedding 的本质和语义相似来源 | 未整理 |
| Q4-0b | L2/L3 | Agentic RAG 与传统 RAG Pipeline 的区别 | 已成稿 |
| Q4-1 | L3 | 生产级 Agent Memory 架构 | 已成稿 |
| Q4-2 | L3 | RAG 高阶调优 Checklist | 已成稿 |
| Q4-3 | L3 | Hybrid Search 与 RRF | 未整理 |
| Q4-4 | L2/L3 | 中文文档 Chunking 策略 | 未整理 |
| Q4-5 | L3 | GraphRAG 原理和适用场景 | 未整理 |
| Q4-6 | L2/L3 | RAG 知识库实时更新策略 | 未整理 |
| Q4-7 | L3 | RAG 评估体系 | 已成稿 |
| Q4-8 | L3 | RAG vs Fine-tuning vs 结合 | 已成稿 |

## Multi-Agent 系统

| ID | 难度 | 主题 | 复习状态 |
| --- | --- | --- | --- |
| Q5-0a | L2/L3 | 为什么需要多个 Agent | 未整理 |
| Q5-1 | L3 | Agent 框架选型 | 未整理 |
| Q5-2 | L3 | Supervisor vs Peer-to-Peer | 未整理 |
| Q5-3 | L3 | LangGraph State 设计 | 未整理 |
| Q5-4a | L2/L3 | MCP 解决什么问题，如何实现 MCP Server | 未整理 |
| Q5-4b | L2/L3 | A2A 与 MCP 的关系 | 未整理 |
| Q5-5 | L3 | Multi-Agent 调试方法论 | 未整理 |
| Q5-6 | L3 | Agent Handoff 设计 | 未整理 |
| Q5-7 | L2/L3 | Multi-Agent 成本控制 | 未整理 |
| Q5-8 | L3 | Multi-Agent 整体性能评估 | 未整理 |

## 生产工程

| ID | 难度 | 主题 | 复习状态 |
| --- | --- | --- | --- |
| Q6-0a | L2/L3 | LLM 调用和传统 API 调用的本质不同 | 未整理 |
| Q6-0b | L2/L3 | 多模态能力带来的工程挑战 | 未整理 |
| Q6-1 | L3 | Agent Observability 体系 | 已成稿 |
| Q6-2 | L3 | Agent Eval 体系 | 已成稿 |
| Q6-3 | L3 | Agent 可靠性工程 | 未整理 |
| Q6-4 | L2/L3 | 生产延迟优化和 Tradeoff | 未整理 |
| Q6-5 | L2/L3 | Agent 幂等性设计 | 未整理 |
| Q6-6 | L3 | Prompt 版本管理和 A/B 测试 | 未整理 |
| Q6-7 | L3 | 成本监控和预算控制 | 已成稿 |
| Q6-8 | L3 | 灰度发布和回滚机制 | 未整理 |
| Q6-9 | L2/L3 | Circuit Breaker 和降级策略 | 未整理 |

## 安全与合规

| ID | 难度 | 主题 | 复习状态 |
| --- | --- | --- | --- |
| Q7-0a | L2/L3 | LLM 无法区分数据和指令的安全困境 | 已成稿 |
| Q7-1 | L2/L3 | Prompt Injection 攻击类型和防御纵深 | 未整理 |
| Q7-2 | L3 | AI Agent 完整安全防护层次 | 已成稿 |
| Q7-3 | L2/L3 | Guardrails 的实现层次 | 未整理 |
| Q7-4 | L3 | PII 和数据隐私处理方案 | 已成稿 |
| Q7-5 | L3 | Jailbreak 防御 | 未整理 |
| Q7-6 | L2/L3 | 高风险行业部署 Agent 的合规要求 | 未整理 |

