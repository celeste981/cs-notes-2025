# Marketing-Data 可靠性和一致性深挖

> 返回：[Marketing-Data 架构梳理](./marketing-data-architecture.md)
>
> 图解：[Marketing-Data 架构图解](./marketing-data-architecture.html)
>
> 项目索引：[Marketing-Data 项目材料](./README.md)

## 一、这篇重点背什么

Marketing-Data 可能被深挖的可靠性问题：

- batch 幂等为什么重要。
- history 为什么要先写。
- Kafka publish 失败怎么处理。
- ES / 文件 / Hive 异常怎么定位。
- 重复触达风险怎么控制。
- 当前链路哪些地方不能夸成强一致。

## 二、可靠性防线

```mermaid
flowchart LR
  A["请求校验"] --> B["batch_id 幂等"]
  B --> C["DB 唯一约束兜底"]
  C --> D["先写 init history"]
  D --> E["异步 handler 执行"]
  E --> F["过滤 / 去重"]
  F --> G["Kafka 推送"]
  G --> H["更新 history + monitor"]
```

## 三、batch_id 幂等深挖

### 为什么重要

离线抽取可能重复触发：

- Marketing 调用超时后重试。
- Scheduler 重复调度。
- 人工重复触发。
- 网络抖动导致调用方不确定是否成功。

如果没有幂等，后果可能是：

- 重复拉取人群。
- 重复推 Kafka。
- Marketing 重复消费后重复触达。
- history 里一批数据多条记录，排障混乱。

回答模板：

> `batch_id` 是离线抽取的业务幂等键。Manager 先查 history，存在就直接返回；DB 唯一索引可以兜底并发重复插入。这样同一批不会重复抽取和重复推送。

## 四、先写 history 的价值

| 如果不先写 | 先写 history 后 |
| --- | --- |
| 异步任务失败后只剩日志。 | 可以按 batch_id 查到 init/failed。 |
| 调用方不知道是否开始执行。 | 有明确执行记录。 |
| 错误统计分散。 | total、unique_total、err_detail、duration 可沉淀。 |
| 排障只能靠机器日志。 | Admin/接口可查询批次历史。 |

面试说法：

> 异步任务一定要先记账再执行。这样即使 handler 崩了、ES 超时、文件异常，也能通过 batch history 找到批次状态和错误详情。

## 五、Kafka publish 失败怎么防守

当前文档口径要保守：

> 代码里单条 Publish 失败会记录 error，不应该夸成每条消息都有强一致保证。如果业务要求更强，需要补充 DLQ、重试、失败落表或发送结果对账。

### 面试回答

| 问题 | 稳妥回答 |
| --- | --- |
| Publish 失败会不会导致任务失败？ | 当前要看代码实现和 history 错误记录，不夸大为强保证。 |
| 怎么增强？ | producer 重试、DLQ、失败落表、按 batch 对账、Marketing 消费确认。 |
| 如果部分消息失败怎么办？ | 记录失败详情，结合 batch_id 做补偿或重推。 |
| 怎么避免重复？ | batch_id 幂等 + 消息业务 key + Marketing 消费端幂等。 |

## 六、按 data_source 排障

| data_source | 排障路径 |
| --- | --- |
| ES | batch history -> DSL -> index/mapping -> scroll -> Canal 同步 -> owner MySQL。 |
| File/S3 | batch history -> 文件路径 -> header/mapping -> 文件大小/编码 -> 过滤规则。 |
| Hive | batch history -> Hive 产物是否生成 -> 文件路径 -> 字段格式。 |
| Insight | batch history -> Insight 产物 -> 字段语义 -> 文件读取。 |
| pre-filter | source topic -> 规则缓存 -> regexp -> FillVariable -> dest topic。 |

## 七、重复触达风险

Marketing-Data 不直接触达用户，但它会影响 Marketing 是否重复触达。

风险来源：

- 同一 batch 重复推 Kafka。
- 同一用户在数据源里重复出现。
- 去重 key 设计不合理。
- Kafka publish 失败后人工重推导致重复。
- Marketing 消费端没有幂等。

防守：

```text
batch_id 幂等
  -> handler 去重
  -> message key / user key
  -> Marketing 消费端幂等
  -> Handler 执行记录
```

## 八、监控和 history 怎么讲

应该能回答：

- 这批是否开始执行。
- 当前是 success 还是 failed。
- 抽取总数是多少。
- 去重后数量是多少。
- 过滤掉多少。
- 错误详情是什么。
- 执行耗时多久。
- 推送到了哪个 topic。

面试说法：

> Marketing-Data 的可观测性入口不是只查日志，而是 batch history + monitor。排障从 batch_id 开始，再按 data_source 分层看。

## 九、不能夸大的点

| 不要这么说 | 更稳说法 |
| --- | --- |
| “Kafka 推送一定不会丢。” | “当前 producer 失败会记录错误；强一致需要额外 DLQ/重试/落表。” |
| “batch_id 可以解决所有重复。” | “batch_id 解决批次重复，用户级重复还要 handler 去重和消费端幂等。” |
| “ES 查不到就是没人群。” | “先看 ES DSL/mapping/Canal，再回到 owner MySQL。” |
| “Marketing-Data 执行营销动作。” | “Marketing-Data 准备人群，Marketing Handler 执行动作。” |

## 十、面试 1 分钟回答

> Marketing-Data 的可靠性重点是 batch 幂等、history 可观测和数据源分层排障。`batch_id` 用来防止同一计划批次重复抽取和重复推 Kafka；任务开始前先写 init history，异步执行后更新 success/failed、数量、耗时和错误详情。不同 data_source 的问题要分开看，ES 看 DSL/mapping/scroll/Canal，文件看路径和字段 mapping，pre-filter 看 topic 和规则缓存。Kafka publish 失败这块要保守表达，不能说强一致，后续增强可以用 DLQ、重试和失败落表。

## 十一、代码和资料证据

- Manager：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/manager/impl/offline_data_fetch_manager_impl.go`
- History Repo：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/repo/repo_dao/impl/offline_plan_history_repo_dao_impl.go`
- Producer：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/mq/producer/impl/producer_impl.go`
- Monitor：`/Users/si.chen/GolandProjects/insurance-marketing-data/src/monitor/register.go`
- `batch_id` 唯一索引痕迹：`/Users/si.chen/GolandProjects/insurance-marketing-data/resource/sql/dev-1.2.60/ddl.sql`
