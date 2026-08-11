# Agent 面试复习

这个目录用于整理 LLM Agent、RAG、Tool Calling、Memory、Eval、Observability、安全合规，以及用户自己的 AI/知识工程项目表达。

## 当前材料

| 文件 | 用途 |
| --- | --- |
| [agent-topic-index.md](agent-topic-index.md) | 从 LC100 中拆出的 Agent/LLM 题目索引，共 68 道。 |
| [agent-key-points.md](agent-key-points.md) | Agent 高频重难点整理，适合先复习主线。 |
| [insurance-knowledge-engineering-interview.md](insurance-knowledge-engineering-interview.md) | 把 `projects/AI` 的知识工程项目转成可面试表达。 |

## 可视化页面

- [Agent 高频重难点可视化](agent-key-points.html)
- [Agent 题库矩阵可视化](agent-topic-index.html)
- [Insurance 知识工程化项目可视化](insurance-knowledge-engineering-interview.html)

## 推荐复习顺序

1. 先看 [agent-key-points.md](agent-key-points.md)，把主线讲顺。
2. 再看 [agent-topic-index.md](agent-topic-index.md)，挑 L3 题补追问。
3. 最后看 [insurance-knowledge-engineering-interview.md](insurance-knowledge-engineering-interview.md)，把八股和自己的项目经历连接起来。

## 复习主线

```text
Context Engineering
-> Agent Loop / ReAct / Planning
-> Tool Calling / Tool Schema / Tool Result
-> RAG / Memory / Knowledge Service
-> Multi-Agent / Handoff / MCP
-> Observability / Eval / Cost / Reliability
-> Prompt Injection / PII / Guardrails / HITL
```

## 面试回答原则

- 先说工程问题，不要一上来堆框架名。
- 每个概念都要落到：为什么需要、怎么实现、失败场景、如何监控。
- 项目经历只说确认过的事实；不确定的统一标记为 `待补充`。
- 如果没真正落地过生产系统，可以说“我参与的是方案设计/理解/整理，核心实现和指标待补充”，不要包装成已上线成果。
