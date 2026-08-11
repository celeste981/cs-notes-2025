# LC100 面试复习资料

> 来源：https://lc100.pages.dev/  
> 整理日期：2026-05-26  
> 说明：这里把站点中的结构化面试题整理成复习入口、完整题目索引和重难点可视化页。内容按面试复习重新组织，项目经历相关部分统一保留为 `待补充`，避免写入不真实的生产细节。

## 文件

- [lc100-topic-index.md](./lc100-topic-index.md)：完整题目索引，共 132 道。
- [key-difficult-points.md](./key-difficult-points.md)：重难点 Markdown 版。
- [key-difficult-points.html](./key-difficult-points.html)：重难点可交互 HTML 版，可直接用浏览器打开。
- [../mysql/05_log_replication_ha/README.md](../mysql/05_log_replication_ha/README.md)：MySQL 日志专题已单独详细展开。
- [../agent/README.md](../agent/README.md)：Agent/LLM 题库和项目表达已单独整理。

## 模块分布

| 模块 | 题量 |
| --- | ---: |
| LLM / Agent | 68 |
| Golang | 22 |
| MySQL | 7 |
| Redis | 35 |

## 分类分布

| 分类 | 模块 | 题量 | L3 相关题量 |
| --- | --- | ---: | ---: |
| LLM 应用工程 | LLM / Agent | 10 | 8 |
| Agent 推理与规划 | LLM / Agent | 11 | 9 |
| 工具系统工程 | LLM / Agent | 9 | 9 |
| 记忆与知识系统 | LLM / Agent | 10 | 10 |
| Multi-Agent 系统 | LLM / Agent | 10 | 10 |
| 生产工程 | LLM / Agent | 11 | 11 |
| 安全与合规 | LLM / Agent | 7 | 7 |
| 内存管理与 GC | Golang | 6 | 6 |
| 接口与泛型 | Golang | 5 | 3 |
| 错误处理与工程实践 | Golang | 5 | 1 |
| 性能优化与系统设计 | Golang | 6 | 5 |
| 日志系统与崩溃恢复 | MySQL | 7 | 7 |
| 持久化机制 | Redis | 6 | 3 |
| 缓存设计与一致性 | Redis | 7 | 4 |
| 分布式锁与事务 | Redis | 7 | 4 |
| 性能优化与线程模型 | Redis | 7 | 3 |
| 应用场景与数据建模 | Redis | 8 | 5 |

## 建议复习顺序

1. 先看 [key-difficult-points.html](./key-difficult-points.html)，用筛选按钮按模块过一遍主线。
2. 再看 [lc100-topic-index.md](./lc100-topic-index.md)，把不熟的题目标为 `需强化`。
3. 每个重难点按“先讲人话 -> 再讲原理 -> 面试怎么说”准备 1 到 2 分钟口述版。
4. 结合自己的项目经历补齐 `待补充`：业务场景、服务规模、排查过程、优化结果，不能编造指标。
