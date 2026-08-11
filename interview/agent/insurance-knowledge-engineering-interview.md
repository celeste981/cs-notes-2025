# Insurance 知识工程化项目：Agent 面试表达


可视化页面：[insurance-knowledge-engineering-interview.html](insurance-knowledge-engineering-interview.html)
来源材料：[projects/AI/insurance-knowledge-engineering-plan.md](../../projects/AI/insurance-knowledge-engineering-plan.md)

## 项目分类

- 类型：项目经验 / Agent 应用工程 / 知识工程 / RAG / 内部研发效率工具
- 主题：Project KB、Domain KB、Overall KB、Knowledge Service、Harness Feedback Loop
- 风险：当前材料还有较多 `待补充`，不能包装成完整上线成果。

## 面试官可能怎么问

1. 你做过 Agent 或 RAG 相关项目吗？
2. 你怎么解决 Agent 上下文不稳定、查不到项目知识的问题？
3. Knowledge Base 和普通文档库有什么区别？
4. 怎么保证 Agent 用到的是可信、不过期、可追溯的知识？
5. 如果 Agent 查错知识或生成错误答案，你怎么闭环修复？

## 30 秒短答

我参与理解和整理过一个 Insurance 知识工程化方案，目标是给 AI Agent 建一个可检索、可追溯、可治理的项目知识底座。它不是替代 Confluence 或代码仓库，而是把 repo 级知识、业务域知识、团队全局知识分层管理，再通过 Knowledge Service 提供检索、metadata、owner、权限和更新状态，让 Agent 在 PRD 分析、TD 生成、代码审查和交付检查时能获取可信上下文。

## 1 到 2 分钟版本

这个项目的背景是：Agent 在真实研发任务里经常因为上下文不足而不稳定。比如它不知道当前 repo 的架构、业务术语、历史决策、Confluence 里的正式 PRD/TD，或者拿到的是过期知识。这个方案就是把这些分散知识工程化。

架构上分三层知识。Project KB / Repo KB 负责单仓库事实，比如目录结构、运行方式、owner、常见改动点；Domain KB 负责业务规则、术语、上下游流程；Overall KB 负责 Insurance 全局视角和跨系统原则。正式交付材料，比如 PRD、TD、Coding Report、AICR 和 Smoke Evidence，仍然保留在 Confluence。

中间通过 Insurance Knowledge Service 做统一索引和治理，包括 source、owner、更新时间、适用范围、权限和检索质量。Agent 不直接无序扫描所有资料，而是按任务场景通过 Knowledge Service 获取可信上下文。最后通过 Harness Feedback Loop 收集 Agent 查不到、查错、知识过期、owner 不清楚等问题，再反向修复 KB 和文档。

这个项目的价值是把隐性项目经验变成可维护的团队知识，让 Agent 从一次性问答变成可以嵌入研发流程的工具。

## STAR 结构

### S：背景

Insurance 项目资料分散在代码仓库、Confluence、PRD、TD、交付证据和聊天记录里。Agent 在做需求分析、技术设计、代码生成或排查时，容易因为上下文缺失而回答不稳定。

### T：任务

设计一个可供 Agent 使用的知识底座，让它能按场景检索可信知识，并能追踪来源、owner、更新时间和权限边界。

### A：行动

可讲行动：

1. 将知识分层为 Project KB / Repo KB、Domain KB、Overall KB 和 Confluence 正式文档。
2. 设计 Knowledge Service 作为索引、检索和治理层，而不是复制所有文档。
3. 定义知识元数据：source、owner、更新时间、适用范围、权限、可信度。
4. 设计 Agent 使用路径：先查 repo 事实，再查业务规则，再查全局原则和正式交付材料。
5. 通过 Harness Feedback Loop 把 Agent 使用中暴露的问题反向沉淀为知识修复任务。

### R：结果

当前结果需要谨慎表达：

- 可以说：形成了清晰的知识分层、使用路径和治理原则。
- 可以说：为 Agent 在需求分析、TD、代码生成、Review、Smoke Evidence 场景中获取上下文提供了方案。
- 待补充：是否已有 MVP、接入哪些 repo、真实命中率/效率提升、上线状态。

## 和 Agent 八股怎么连接

| 八股问题 | 项目连接 |
| --- | --- |
| Context Engineering | Knowledge Service 控制 Agent 看到哪些上下文，而不是把所有文档塞进 prompt。 |
| RAG | Project KB、Domain KB、Confluence 是检索源，Knowledge Service 管索引和 metadata。 |
| Memory | Project KB 像长期项目记忆，Harness Feedback Loop 像从失败中更新记忆。 |
| Observability / Eval | 需要记录 Agent 查了什么 source、是否命中、答案是否引用可信来源。 |
| Security / Compliance | 不把凭证、客户隐私、未脱敏生产数据写入知识库；按权限控制可见性。 |
| Agent Reliability | 通过可信来源和反馈闭环降低幻觉和过期知识风险。 |

## 面试追问与回答要点

### 追问 1：这和普通文档库有什么区别？

普通文档库只是存文档。这个方案面向 Agent 使用，需要额外管理 source、owner、更新时间、适用范围、权限和检索质量。重点不是“多存一份文档”，而是让 Agent 能按任务拿到可信上下文。

### 追问 2：为什么不把所有 Confluence 都塞给 Agent？

因为 context window 有限制，全部塞进去成本高、噪音大，还容易引入过期或无关信息。更合理的是按任务检索相关知识，并保留引用来源。

### 追问 3：知识过期怎么办？

每条知识要有 owner 和更新时间；检索时可以根据时间和来源置信度排序；Agent 使用中发现冲突或过期，通过 Harness Feedback Loop 生成修复任务。

### 追问 4：怎么评估这个系统有没有用？

可以看任务完成率、检索命中率、引用覆盖率、人工纠错率、Agent 因上下文不足失败的比例、PRD/TD/Review 生成耗时变化。但这些具体指标目前需要 `待补充`。

### 追问 5：你个人贡献是什么？

必须按事实回答。当前可以准备两版：

如果只是参与整理：

> 我主要参与的是方案理解和结构化整理，把截图/文档里的知识分层、治理原则和 Agent 使用路径整理成可复用材料。

如果确实参与设计或落地：

> 待补充：具体负责的模块、实现内容、接入的 repo、效果指标。

## 风险说法

不要这样说：

- “我做了一个完整的 RAG 平台。”除非确实实现并上线。
- “Agent 准确率提升了 xx%。”除非有评估数据。
- “所有项目知识都自动同步。”除非有真实同步链路。

更稳的说法：

- “我参与过知识工程化方案的设计/整理。”
- “这个方案目标是解决 Agent 上下文不稳定和知识不可追溯的问题。”
- “落地状态和具体指标我需要结合实际项目补充。”

## 待补充

- 原始 PRD / 方案链接。
- Knowledge Service 是否已经实现。
- 接入了哪些仓库和 Confluence 空间。
- 检索技术选型：向量检索、关键词检索、Hybrid Search、Rerank。
- metadata 字段规范。
- 权限控制方案。
- Eval 指标和实际效果。
- 个人实际贡献。

