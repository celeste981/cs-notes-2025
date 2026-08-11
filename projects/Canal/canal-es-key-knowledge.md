# Canal / ES 核心知识点速记

> 返回：[Canal / ES 同步架构梳理](./canal-es-sync-architecture.md)
>
> 图解：[Canal / ES 同步架构图解](./canal-es-sync-architecture.html)
>
> 项目索引：[Canal 项目材料](./README.md)

## 一、面试先背这 10 个点

1. Canal 同步的是 MySQL binlog，不是主动定时扫表。
2. ES 是查询视图，不是业务事实源。
3. Canal 链路是最终一致，不是强一致。
4. Adapter 运行依赖 `application.yml` 里的 Kafka、DLQ、datasource。
5. ES sync script 决定从哪张表回查、写哪个 index、用哪个 `_id`。
6. Mapping 要先设计，字段类型错了可能要重建索引。
7. ETL 是存量同步，增量同步来自 binlog，两者问题不同。
8. 聚合宽表 ETL 不能随便按表主键分批，可能跨批覆盖。
9. RDS 扩容要同时处理 instance、datasource、script、回滚。
10. 同步异常按 MySQL -> Instance -> Kafka -> Adapter -> ES -> DSL 分层排查。

## 二、核心概念表

| 概念 | 人话解释 | 面试怎么说 |
| --- | --- | --- |
| binlog | MySQL 的变更日志。 | Canal 订阅 binlog，所以要先确认 MySQL 事实源和 binlog 事件。 |
| Canal Instance | 订阅某个库/表变更的 Canal 单元。 | RDS 切换、扩容时常涉及新增或启停 instance。 |
| Kafka topic | Server 和 Adapter 中间的消息通道。 | 看 lag 可以判断 Adapter 是否消费正常。 |
| Canal Adapter | 消费消息并写 ES 的组件。 | Adapter 读取 datasource 和 sync script，回查 MySQL 后组装 ES document。 |
| sync script | MySQL 到 ES 的 YAML 映射。 | 关注 `dataSourceKey`、`destination`、`groupId`、`_index`、`_id`、`sql`、`etlCondition`。 |
| ES mapping | ES 字段类型。 | mapping 错了会导致写入失败、查询失败或要重建索引。 |
| ETL | 历史存量补刷。 | ETL 的分批字段、业务主键、SQL 性能决定正确性和效率。 |
| DLQ | Dead Letter Queue，死信队列。 | 用来承接异常消息，避免单条坏数据阻塞整条链路。 |

## 三、sync script 必背字段

```yaml
dataSourceKey: policy_{region}
destination: {env}-{region}.tp.spin_es
groupId: canal-adapter
esMapping:
  _index: spin_policy_{region}_es
  _type: _doc
  _id: id
  upsert: true
  sql: "select ... from policy_n_tab_{index} a"
  etlCondition: "where a.id >= {} and a.id < {}"
  commitBatch: 3000
```

| 字段 | 重点 |
| --- | --- |
| `dataSourceKey` | 要和 `application.yml` 里的 datasource 对得上。 |
| `destination` | 要和 Canal/Kafka topic 对得上。 |
| `groupId` | 决定 adapter 消费组。 |
| `_index` | ES index 或 alias 名。 |
| `_id` | 必须稳定唯一，否则会覆盖错文档。 |
| `sql` | 回查 MySQL 并组装 document。 |
| `etlCondition` | 存量 ETL 的分批条件。 |
| `commitBatch` | 每批提交数量，过大过小都可能有风险。 |

## 四、ES 设计规范怎么讲

### 什么时候适合建 ES index？

适合：

- 跨分库分表查询。
- 非分片键复杂查询。
- 多条件组合筛选。
- Admin 列表、报表、运营分析。
- 需要近实时查询视图，但事实源仍在 MySQL。

不适合：

- 强一致交易判断。
- 只按主键查一条数据。
- 直接承载业务状态机。
- 替代 owner 服务的领域模型。

### Mapping 重点

| 规则 | 为什么 |
| --- | --- |
| 字段最小化 | 字段越多，mapping、存储和查询成本越高。 |
| keyword 按需使用 | 不需要精确查询/聚合的字符串不要都做 keyword。 |
| 日期和金额类型要准确 | 类型错了 range 查询或聚合会出问题。 |
| JSON 字段要 null-safe | 下游解析 object 时不能因为 null 崩掉。 |
| mapping 先于 script | 先确认 index 和字段类型，再发布同步脚本。 |

### Shard 和容量口径

保守面试说法：

> ES index 设计要提前评估数据量、字段数、查询方式和 shard。业务类单 shard 不宜无限变大，副本和 refresh interval 要按查询实时性和写入压力评估。这里不要夸具体线上容量，除非有真实数据证据。

## 五、增量同步 vs ETL

| 维度 | 增量同步 | ETL |
| --- | --- | --- |
| 来源 | binlog event | 历史 SQL 扫描 |
| 触发 | 业务写入自动触发 | 人工/任务触发 |
| 关注 | instance、Kafka lag、Adapter、DLQ | 分批字段、SQL 性能、覆盖问题 |
| 风险 | 阻塞、mapping、脚本异常 | 深分页、跨批覆盖、漏刷 |

面试说法：

> 增量同步的核心是消费 binlog 并及时写 ES；ETL 是历史数据补刷，最容易出问题的是分批条件和业务主键。如果聚合宽表按明细表 id 分批，同一个业务 document 可能跨批被覆盖。

## 六、常见错误说法

| 不要这么说 | 更稳的说法 |
| --- | --- |
| “ES 里没有就是业务没有。” | “先查 owner MySQL，ES 只是查询视图。” |
| “Canal 是强一致。” | “Canal 是近实时/最终一致链路。” |
| “异常都进 DLQ 就好了。” | “要区分 ES 集群异常和单条数据/脚本异常。” |
| “mapping 可以随时改。” | “字段类型变更可能要重建索引和补刷。” |
| “ETL 按 id 分批就行。” | “聚合宽表要看业务主键，不能跨批覆盖。” |

## 七、1 分钟回答模板

> Canal/ES 是 Insurance 的数据同步和查询视图链路。业务服务写 MySQL 后，Canal Instance 订阅 binlog，经 Kafka 投递给 Adapter。Adapter 根据 sync script 找 datasource 回查 MySQL，并把结果写入 ES index。Admin/O-BFF 用 ES 做复杂列表、报表和分表聚合查询。但 ES 不是事实源，链路也是最终一致，所以排障要从 owner MySQL、binlog、instance、Kafka lag、Adapter、DLQ、ES mapping、查询 DSL 一层层看。

## 八、资料来源

- Confluence：`ES 使用分享`，页面 ID `2347394775`。
- Confluence：`ES设计规范`，页面 ID `2083164892`。
- Confluence：`ES以及Canal同步工具维护说用`，页面 ID `1857867981`。
- 本地模板：`/Users/si.chen/GolandProjects/canal-adapter/tools/template/`
- 本地配置：`/Users/si.chen/GolandProjects/canal-adapter/conf/live/application.yml`
