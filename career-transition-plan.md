# 三个月职业转型路线图：Senior Backend + AI Agent Engineering

> 求职定位：**Senior Backend Engineer with AI Agent Engineering capabilities**。  
> 目标城市：深圳；重点公司：腾讯。  
> 时间范围：约 3 个月。本文是求职准备的决策基线；任何新增学习任务先判断是否直接服务于简历、面试、项目、系统设计或投递。

## 一句话策略

不把已有后端积累推倒重来，而是在 `Financial Products + Complex Business Systems + Golang Backend` 的底座上，补齐能把 AI Agent 从 Demo 做到生产系统的后端工程能力。

## 目标岗位与能力重叠

| 方向 | 核心卖点 | 已有基础 | 3 个月要补齐的部分 |
| --- | --- | --- | --- |
| AI Agent Backend Engineer | 能把 Agent 做成可控、可靠、可观测的 Production System | Golang、MySQL、Redis、Kafka、gRPC、微服务、后台系统、金融产品复杂业务 | Python/FastAPI、LLM API、Tool Calling、MCP、Workflow、State、RAG、Eval、Agent 可观测性与治理 |
| Senior Backend Engineer | 复杂金融业务和分布式后端的设计、交付与稳定性能力 | Shopee Financial Products、Insurance、Promotion/Voucher、运营平台、分布式系统 | 项目深挖、系统设计、可靠性/一致性表达、算法高频题、简历与面试表达 |

两类岗位共享约 70% 的能力基础：系统设计、异步任务、状态机、幂等、重试、超时、审计、权限、观测、成本与稳定性。AI Agent 是在这些基础能力上的新应用层，而不是转向模型训练。

## 明确边界

### 高优先级

- 后端核心：Golang、MySQL、Redis、Kafka、分布式一致性、微服务、高可用、高并发、性能与可观测性。
- Agent 工程：LLM API Integration、Tool Calling、MCP、LangGraph/Workflow、State Persistence、Memory、RAG、Evaluation、Observability、Long-running Task、Recovery、Permission/Audit、Model Fallback、Cost Control。
- 求职产物：两份简历、一个可讲可演示的项目、系统设计题库、项目故事、目标 JD 清单、模拟面试与正式投递。

### 低优先级 / 非目标

不以 LLM Algorithm Engineer 为目标，因此不系统投入模型训练、深度学习数学、PyTorch 训练、Transformer 从零实现、SFT、RLHF、PPO/GRPO、CUDA、训练基础设施或纯推理优化。只有在帮助理解 Agent Backend 的最小范围内才按需补充。

## 三个月节奏与交付物

| 阶段 | 核心任务 | 必须产物 | 完成信号 |
| --- | --- | --- | --- |
| 第 1-2 周：Go 生存底座 | 语言基础、并发、context、测试；梳理一个真实项目 | Go 高频短答、一个项目请求链路、弱项清单 | 面对 Go 基础追问不再完全失语 |
| 第 3-4 周：数据与可靠性 | MySQL、Redis、Kafka、幂等/重试/DLQ；完成一个可靠性故事 | 后端专题笔记、一个 STAR 故事、JD 关键词 | 能把原理与项目连接起来 |
| 第 5-6 周：项目与系统设计 | 系统设计、gRPC/HTTP、性能与观测；Mock 与 Resume B | 5 分钟项目讲稿、系统设计图、Resume B v1 | 能守住 Senior Backend 主叙事 |
| 第 7-8 周：Agent MVP | 两个 Tool、RAG、状态、失败处理、Trace 与小评估集 | 可运行 MVP、Demo、README、评估样例 | 能诚实展示一段 Agent 工程能力 |
| 第 9-10 周：面试化与简历 | 项目收口、后端串讲、系统设计、Resume A/B、JD 匹配 | 两份简历 v1、Mock 反馈、目标 JD 列表 | 不依赖代码也能讲清项目与取舍 |
| 第 11-12 周：投递与迭代 | 小批量投递、Mock、真实反馈复盘 | 投递记录、面试复盘、弱项修订 | 进入可持续的投递—面试—复盘循环 |

逐周安排、每周验收和计划收缩规则见 [reviews/three-month-weekly-plan.md](reviews/three-month-weekly-plan.md)。

## 时间分配

- 55%：Backend Interview Preparation。
- 25%：AI Agent Backend MVP。
- 20%：算法高频题、JD 调研、简历、Mock Interview 与投递。

算法题只服务于 Backend / Agent Backend 面试；不按算法岗位的训练量投入。

## 项目主线：可信的 Agent Backend MVP

项目主线不是必须做 Voucher Agent，而是完成一个你能守住的 Agent 后端最小闭环。推荐方案和两周止损规则见 [projects/AI/agent-backend-mvp-plan.md](projects/AI/agent-backend-mvp-plan.md)。

业务场景默认采用中性的“运营知识与任务诊断助手”；仅在你能用脱敏/模拟数据解释 Voucher 规则和边界时，才将它包装为 Voucher 场景。现有的 [Voucher Operations Agent 设计](projects/AI/voucher-operations-agent-design.md) 保留为可选的扩展蓝图，不是必须实现的清单。

核心交付只需要证明五件事：受控 Tool Calling、带引用的 RAG、可持久化任务状态、有限的失败恢复、Trace + 小型评估集。不要把 UI、Multi-Agent、真实写操作或框架数量当作目标。

## 两份简历的叙事边界

| 简历 | 主叙事 | 项目呈现原则 |
| --- | --- | --- |
| Resume A：AI Agent Backend Engineer | 后端工程师把 Agent 做成生产系统的能力 | Agent Backend MVP 为加分项目；强调已实现的 Tool、Workflow、State、RAG、Eval 或 Observability；不把 Voucher 或未实现的生产能力写成事实 |
| Resume B：Senior Backend Engineer | Shopee Financial Products 中复杂业务与分布式后端经验 | Insurance、Promotion/Voucher、运营平台为主；AI Agent 作为工程能力加分项，不替代后端主身份 |

每条简历 bullet 都应可经受追问：问题规模、约束、个人职责、技术选择、结果证据。没有确认的指标和线上成果统一写 `待补充`，不编造。

## 每周复盘问题

1. 这周投入是否直接提高了目标岗位成功率？若没有，下一周降级或停止。
2. 是否新增了可写进简历、可在面试讲清楚、可实际演示的产物？
3. Agent MVP 是否增加了生产工程证据，而不只是增加 UI 或 prompt？
4. 后端核心弱项是否已转成可复述的项目故事或系统设计答案？
5. 从 JD / Mock 里出现的共性缺口是什么，下一周如何用最小投入补齐？

## 时间紧张时的优先级护栏

建议按每周 9-11 小时制定计划，设置 6-7 小时的最低保底，而不是以可透支的时间安排任务；持续性优先于短期冲刺。每周先锁定：2 个后端面试专题、1 个项目产物、1 次简历/JD/Mock 动作。其余内容全部视为可选。详细节奏见 [reviews/three-month-weekly-plan.md](reviews/three-month-weekly-plan.md)。

当以下情况发生时，应立即收缩范围：

- 项目两周内没有可运行 Demo：停止新增功能，先完成最小闭环与讲稿。
- 同一 Agent 概念连续学习两次仍没有代码、笔记或面试答案产出：暂停继续学习该概念。
- 投递期开始后：项目每周维护不超过 2-3 小时，其余时间投向 JD、简历、模拟面试和真实面试复盘。

## 当前待补充

- Resume A / Resume B 的真实经历素材、项目职责和可公开范围：已建立 [代码证据草稿](resume/resume-evidence-draft.md)，待本人确认职责与结果。
- Voucher 业务中可以用于本地模拟的规则、异常类型和脱敏样例：已在代码证据草稿中给出静态模拟范围。
- 每周可稳定投入的小时数，以及第 1 周的具体起始日期：标准每周 9-11 小时，低能量周允许降至 2-4 小时；第 1 周为 2026-08-10 至 2026-08-16。
- 深圳目标公司与目标岗位 JD 清单：已完成 [初步市场扫描](resume/shenzhen-role-market-scan.md)，真实 JD 清单待下周收集 10 条后建立。
