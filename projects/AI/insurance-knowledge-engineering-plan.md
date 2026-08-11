# Shopee Insurance 知识工程化总体规划设计

> 来源：根据用户提供的截图整理。当前版本不是逐字 OCR，而是保留文档的核心架构、边界和可复习表达。原始链接、作者、落地状态待补充。

## 一句话定位

这是一个面向 Shopee Insurance 的 **AI 可用知识底座规划**：把分散在代码仓库、Confluence、PRD、TD、交付证据、业务知识和 Agent 反馈里的信息，按层级沉淀、索引、治理，并通过统一的 Knowledge Service 提供给 Skill / Agent 使用。

## 目标与边界

目标：

- 建立 Insurance 项目级知识体系，让 Agent 能更稳定地理解业务、仓库、需求和交付上下文。
- 把已有知识从“散落文档”整理为可检索、可追溯、可治理的结构化知识。
- 支撑研发、需求分析、TD 生成、代码审查、Smoke Evidence、上线物料整理等工作流。

边界：

- 不是替代 Git、Confluence 或现有 PRD/TD 文档。
- 不是让 Agent 随意扫描所有资料，而是通过明确来源、owner、metadata 和更新时间进行受控检索。
- 不把敏感数据、凭证、客户隐私直接写入知识库。

## 总体架构

```text
开发与交付资料
├── Git / Markdown / 代码仓库
├── Project KB / Repo KB
├── Domain KB
├── Overall KB
└── Confluence / PRD / TD / Coding Report / AICR / Smoke Evidence

知识服务层
└── Insurance Knowledge Service
    ├── 索引与检索
    ├── metadata 管理
    ├── source / owner / 更新时间追踪
    ├── 权限与可见性控制
    └── 知识质量反馈

使用层
├── Skills / Agent
├── PRD / TD / Code Review / Smoke Evidence 辅助
└── Harness Feedback Loop
```

## 核心组件职责

| 组件 | 主要职责 |
| --- | --- |
| Project KB / Repo KB | 沉淀单个代码仓库的架构、目录、运行方式、owner、常见改动点和注意事项。 |
| Domain KB | 沉淀某个业务域的规则、术语、核心流程、上下游依赖和历史背景。 |
| Overall KB | 沉淀 Insurance 全局视角，包括跨项目关系、通用规范、知识地图和统一原则。 |
| Confluence | 保存正式 PRD、TD、Coding Report、AICR、Smoke Evidence、发布说明等交付文档。 |
| Insurance Knowledge Service | 统一接入 Git、Markdown、Confluence 和 metadata，提供检索、索引、权限、owner 和更新状态管理。 |
| Skills / Agent | 通过 Knowledge Service 获取上下文，用于需求分析、技术方案、代码生成、排查和交付检查。 |
| Harness Feedback Loop | 收集 Agent 使用知识时暴露的问题，反推补齐 Project KB、Domain KB、Overall KB 或 Confluence。 |

## 知识使用路径

| 场景 | 优先路径 |
| --- | --- |
| 当前 repo 下 Code / Agent 任务 | 先读当前 repo 的 Project KB，再结合代码。 |
| 当前 repo 架构或运行方式不清楚 | 查 Project KB / Repo KB。 |
| 涉及业务域、术语、规则或跨模块背景 | 查 Domain KB。 |
| 涉及多个 repo、多个 domain 或团队级规范 | 查 Overall KB。 |
| 需要 PRD、TD、AICR、Smoke Evidence 等正式证据 | 查 Confluence 或通过 Knowledge Service 定位原始文档。 |
| 知识缺失、过期、冲突或 owner 不清楚 | 进入 Harness Feedback Loop，形成补文档或修正知识的任务。 |

关键原则：

- Project KB / Repo KB 解决“这个仓库怎么改、怎么跑、谁负责”。
- Domain KB 解决“这个业务为什么这样做、规则是什么”。
- Overall KB 解决“Insurance 全局怎么组织、跨系统怎么理解”。
- Knowledge Service 负责把这些知识组织起来，而不是把所有内容复制成一份新的大文档。

## 治理与保护边界

知识治理重点：

- 每条知识需要有来源、owner、适用范围和更新时间。
- Project KB / Domain KB / Overall KB 应以 Markdown + Git 管理为主，保持 review、版本追踪和变更责任。
- Confluence 继续承载正式交付文档和评审证据。
- Knowledge Service 维护索引、metadata、权限和检索质量。

保护边界：

- 不写入 token、cookie、Bearer、账号密码等凭证。
- 不沉淀客户敏感数据或未脱敏生产数据。
- 不把 Agent 推断结果直接当事实，必须保留来源或标记为 `待确认`。
- 不让 L1 / L2 / L3 知识层互相污染：仓库事实、业务规则、团队原则要分层维护。

## 生命周期闭环

```text
需求 / PRD / TD / 代码 / Review / 测试 / 上线
        ↓
项目知识沉淀到 Project KB
        ↓
业务规则沉淀到 Domain KB
        ↓
跨项目原则沉淀到 Overall KB
        ↓
Knowledge Service 建索引并提供检索
        ↓
Skill / Agent 在真实任务中使用
        ↓
Harness Feedback Loop 收集问题
        ↓
反向修复 KB、metadata、owner、权限和文档缺口
```

这个闭环的重点是：知识库不是一次性整理，而是在真实 Agent 任务中不断验证、修正和补充。

## 成功标准

- Skill / Agent 能稳定找到当前 repo、业务域和交付材料的可信来源。
- 当前 repo 的 Project KB 能回答架构、运行、owner、常见改动点和注意事项。
- Agent 生成 PRD 分析、TD、Coding Report、Smoke Evidence 时能引用明确来源。
- 知识缺失、过期、冲突时能进入反馈闭环，而不是长期散落在聊天记录里。
- L1 / L2 / L3 知识层边界清楚，避免把仓库事实、业务规则和团队原则混在一起。

## 关键原则

1. 知识源头仍然是 Markdown + Git、Confluence 和正式交付文档。
2. Confluence 是正式交付和评审记录，不能被随意替代。
3. Project KB 贴近代码仓库，负责 repo 内可执行、可维护的知识。
4. Insurance Knowledge Service 是索引、检索和治理层，不是新的大而全文档仓库。
5. 所有知识都要能追溯 source、owner、更新时间和适用范围。
6. Agent 使用知识后的错误、缺口和不确定性，需要进入 Harness Feedback Loop。

## 面试表达

### 30 秒短答

我参与理解的是一个 Insurance 知识工程化方案，本质是给 AI Agent 建一个可用的项目知识底座。它把 repo 级知识、业务域知识、团队级知识和 Confluence 交付文档分层管理，再通过 Knowledge Service 做索引、metadata、owner 和权限治理，让 Agent 在做 PRD 分析、TD、代码生成、审查和交付检查时能查到可信来源。

### 1 到 2 分钟版本

这个方案不是单纯做文档库，而是解决 Agent 在真实研发任务里上下文不稳定的问题。它把知识分成三层：Project KB / Repo KB 管单仓库事实，比如架构、运行方式、owner 和常见改动点；Domain KB 管业务规则和术语；Overall KB 管 Insurance 全局关系和统一原则。正式 PRD、TD、Coding Report、AICR、Smoke Evidence 仍然放在 Confluence。

中间的 Insurance Knowledge Service 负责把这些来源连接起来，维护索引、metadata、source、owner、更新时间和权限。Agent 不直接乱翻所有材料，而是通过这个服务按场景查可信知识。最后通过 Harness Feedback Loop 收集“查不到、查错、知识过期、owner 不清楚”等问题，再反向修复 KB 和文档。

这个设计的价值是把隐性的项目经验变成可检索、可追溯、可持续维护的团队知识，让 AI Agent 从一次性问答变成可以嵌入研发流程的工具。

## 待补充

- 原始文档链接或 HTML 文件。
- 作者、创建时间、当前评审状态。
- Insurance Knowledge Service 的具体技术选型。
- KB metadata 字段定义。
- 是否已有 MVP、接入了哪些 repo、哪些 Skill / Agent 已经使用。
- 个人在其中的实际贡献：需求分析、架构设计、文档整理、实现、评审或试点验证。
