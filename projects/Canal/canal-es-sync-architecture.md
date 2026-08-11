# Canal / ES 同步架构梳理

> 可视化图解：[Canal / ES 同步架构图解](./canal-es-sync-architecture.html)
>
> 所属项目索引：[Canal 项目材料](./README.md)
>
> Insurance 总览：[Insurance 部门业务材料](../Insurance/README.md)

## 专项补充文档

如果只想背重点，可以直接看这三篇：

- [Canal / ES 核心知识点速记](./canal-es-key-knowledge.md)：binlog、instance、Kafka、adapter、sync script、mapping、ETL、DLQ。
- [Canal DLQ / 同步阻塞优化专项](./canal-dlq-blocking-optimization.md)：异常分类、Bulk 部分失败、DLQ、重试、补数和监控。
- [Canal 重点架构：RDS 切换、ETL、ES 宽表](./canal-rds-etl-architecture-notes.md)：RDS 扩容切换、ETL 分批、宽表设计和回滚。

## 一、Canal 在 Insurance 里做什么？一句话定位

**Canal 是 Insurance 的 MySQL -> ES 近实时同步链路**。它解析 MySQL binlog，把 Product、Order、Policy、Promotion、Marketing 等业务库的变更同步到 Elasticsearch，支撑 O-BFF / Admin 的复杂查询、宽表、分表聚合、报表和运营分析。

面试开场：

> Insurance 后台很多查询不能直接扫 MySQL，尤其是保单、订单、促销、用户、人群这些数据会分库分表，页面还需要复杂条件查询和报表。我们用 Canal 订阅 MySQL binlog，经 Kafka 投递给 Canal Adapter，Adapter 按 ES sync script 把数据写入 Elasticsearch。Canal 的价值是把 MySQL 事实数据同步成适合查询的 ES 视图，但事实源仍然是各 owner 服务的 MySQL。

个人边界：

> Canal 是组内数据同步基础设施。我能讲清配置、同步链路、RDS 扩容切换和排障思路；个人实际做过的 script、切换单和排障记录需要按事实补充。

## 项目路径速查

| 项目 | 本机路径 | 这里看什么 |
| --- | --- | --- |
| `canal-adapter` | `/Users/si.chen/GolandProjects/canal-adapter` | Adapter 启动、配置、sync script、模板工具、RDS 扩容物料。 |
| O-BFF | `/Users/si.chen/GolandProjects/insurance-operator-bff` | Admin 查询和 ES adapter 调用入口。 |
| 业务 owner 服务 | Product / Order / Policy / Promotion / Marketing 等 | MySQL 事实源和业务写入方。 |

常用入口：

| 想看什么 | 推荐路径 |
| --- | --- |
| live adapter 配置 | `/Users/si.chen/GolandProjects/canal-adapter/conf/live/application.yml` |
| sync script 模板 | `/Users/si.chen/GolandProjects/canal-adapter/tools/template/` |
| policy 分表模板 | `/Users/si.chen/GolandProjects/canal-adapter/tools/template/policy_{region}.policy_n_tab_{index}.yml` |
| 生成模板工具 | `/Users/si.chen/GolandProjects/canal-adapter/tools/template_to_yml/generate_es_config.go` |
| RDS 扩容 TD | `/Users/si.chen/GolandProjects/canal-adapter/docs/id_vn_rds_expansion_canal_td.md` |
| 执行清单 | `/Users/si.chen/GolandProjects/canal-adapter/docs/id_vn_rds_expansion_canal_checklist.md` |
| ETL 聚合问题 | `/Users/si.chen/GolandProjects/canal-adapter/docs/canal-adapter-etl-group-by-fix.md` |

## 二、核心概念

| 概念 | 人话解释 | 面试要点 |
| --- | --- | --- |
| MySQL binlog | MySQL 记录数据增删改的日志。 | Canal 的数据来源，不是主动扫表。 |
| Canal Server / Instance | 模拟 MySQL slave，订阅某个库/表的 binlog。 | Instance 对应一个目标库或一组表，RDS 扩容时经常要新增 instance。 |
| Kafka topic | Canal Server 到 Adapter 的中间消息通道。 | Insurance 常见 topic 如 `live-id.tp.spin_es`、`live-vn.tp.spin_es`。 |
| Canal Adapter | 消费 Kafka，按配置把数据写到 ES。 | 关注 datasource、sync script、groupId、DLQ、lag。 |
| ES sync script | 一份 YAML，定义 dataSourceKey、destination、groupId、ES index、SQL、etlCondition、commitBatch。 | 这是从 MySQL 映射到 ES 的核心配置。 |
| ES Mapping | ES index 的字段类型定义。 | Mapping 要先建好；字段类型错了可能要重建索引。 |
| ETL | 存量同步，用 SQL 把历史数据刷进 ES。 | 增量同步和存量 ETL 的问题不完全一样。 |
| DLQ | Dead Letter Queue，异常消息兜底队列。 | 避免异常消息一直阻塞 partition。 |

## 三、完整数据流

```text
业务服务写 MySQL
  -> MySQL 产生 binlog
  -> Canal Instance 订阅并解析 binlog
  -> 投递到 Kafka topic
  -> Canal Adapter 消费 topic
  -> 根据 sync script 找 datasource 和 SQL mapping
  -> 回查 MySQL / 组装 ES document
  -> 写入 Elasticsearch index
  -> O-BFF / Admin / Marketing-Data 查询 ES
```

### Stage 1: 业务服务写 MySQL

Product、Order、Policy、Promotion、Marketing 等 owner 服务负责真实业务写入。Canal 不改变业务事实，只消费变更。

面试说法：

> Canal 同步的是 owner 服务 MySQL 的变化。排查时不能只看 ES，必须先确认 owner MySQL 有没有正确数据。

### Stage 2: Canal Instance 解析 binlog

Canal 模拟 MySQL slave 向 MySQL master/slave 发 dump 请求，接收 binary log event，再解析数据变化。

面试说法：

> Instance 是否启动、位点是否正确、binlog 是否开启、ROW format 是否符合要求，都会影响同步。

### Stage 3: Kafka 解耦

Insurance 的同步链路通常是 Canal Server -> Kafka -> Canal Adapter。

面试说法：

> Kafka 把 server 和 adapter 解耦，也让我们能通过 lag 看 adapter 是否消费正常。但如果异常消息一直失败，也可能阻塞 partition，所以要关注 DLQ 和失败窗口。

### Stage 4: Adapter 执行 sync script

典型 YAML：

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

关键字段：

| 字段 | 作用 |
| --- | --- |
| `dataSourceKey` | 指向 `application.yml` 里的 MySQL datasource。 |
| `destination` | Canal destination 或 Kafka topic。 |
| `groupId` | Adapter 消费组，Insurance 常见是 `canal-adapter`。 |
| `_index` | ES index 名。 |
| `_id` / `pk` | ES document 唯一 ID。 |
| `sql` | 回查 MySQL 并组装 ES 文档的 SQL。 |
| `etlCondition` | 存量 ETL 分批条件。 |
| `commitBatch` | 每批提交大小。 |

### Stage 5: ES 提供查询视图

Admin/O-BFF 通过 ES 做宽表、分表聚合、复杂查询和报表。

常见场景：

- Policy Management：保单宽表、多表字段合并、分表统一查询。
- Promotion：促销活动效果统计、用户行为分析。
- Marketing：用户营销数据、人群数据、group item 查询。
- Claim / Order：复杂列表和导出。

## 四、配置物料怎么讲

### 1. `application.yml`

`conf/live/application.yml` 里包含：

- `mode: kafka`
- Kafka consumer properties
- `deadLetterQueue` 配置
- `srcDataSources`，包括 `canal` 管理库和大量业务 datasource
- Prometheus `/metrics` 配置

面试说法：

> Adapter 运行时要先知道自己消费哪个 Kafka、有哪些 datasource、DLQ 怎么配置，再加载 sync script 执行映射。

### 2. sync script 模板

`tools/template/` 下维护按模块、地区、分表展开的 YAML 模板，例如：

- `policy_{region}.policy_n_tab_{index}.yml`
- `order_{region}.item_n_tab_{index}.yml`
- `marketing_{region}.user_marketing_data_tab.yml`
- `promotion_{region}.user_voucher_tab.yml`

生成工具：

- `/Users/si.chen/GolandProjects/canal-adapter/tools/template_to_yml/generate_es_config.go`

### 3. mapping 关系管理

本地工具文档里抽象了 `SyncScriptMapping`：

```text
SyncScriptMapping
├── region
├── syncScript
├── indexName
├── datasourceKey
├── tableName
├── moduleName
└── isDeleted
```

面试说法：

> Canal 不是只维护一堆散落的 YAML，实际要能回答某个 region、某个 module、某张表同步到哪个 ES index，使用哪个 datasource 和 script。

## 五、RDS 扩容 / 迁移切换流程

RDS 扩容时，`db host`、`db name`、分库数量可能变化。Canal 侧要同时处理 Instance、datasource、sync script 和回滚。

典型流程：

```text
准备新 RDS 清单
  -> 生成新 Canal Instance，初始关闭
  -> 预置新 ES sync script，is_deleted = 1
  -> 更新 adapter application.yml datasource
  -> AT 先完整验证
  -> 停旧 Instance
  -> 启用新 script，软删旧 script
  -> 发布 adapter 配置
  -> 打开新 Instance
  -> 观察 Kafka lag / DLQ / ES 写入 / 抽样数据
  -> AT 通过后再 EQ5/Main
```

关键口径：

- topic 可以保持不变，例如 `live-id.tp.spin_es`、`live-vn.tp.spin_es`。
- 新 script 预置时用 `is_deleted = 1`，切换时再启用。
- 旧 script 不硬删，而是软删并标记 `old_[region]_db`。
- 不要默认用 Adapter `syncSwitch` 停整个 region topic；它粒度太大，可能影响不在本轮切换范围的模块。
- AT 完整闭环通过后再 EQ5/Main。

## 六、ETL 和增量同步的区别

| 维度 | 增量同步 | 存量 ETL |
| --- | --- | --- |
| 数据来源 | binlog event | SQL 按条件扫描历史数据 |
| 触发方式 | 业务写入自动触发 | 人工/任务发起 |
| 常见问题 | instance 未启动、Kafka lag、script 报错、mapping 类型错 | 分批条件错、聚合被覆盖、SQL 超时、历史数据缺失 |
| 排障重点 | 位点、topic、adapter 日志、DLQ | `etlCondition`、业务主键、commitBatch、SQL 性能 |

ETL 聚合问题典型案例：

- SQL 用 `group_concat` / `min` / `max` 等聚合。
- 如果没有 `GROUP BY`，一批数据可能被聚合成一行。
- 加了 `GROUP BY` 后，如果按表主键 `id` 分批，而业务主键跨批次，会出现后批覆盖前批。
- 长期方案是按业务主键分批，或小表不分批。

面试说法：

> 增量同步和存量 ETL 的正确性边界不同。增量通常根据 binlog 事件回查完整业务行；ETL 如果分批字段选错，就可能把同一个业务文档拆到不同批次，导致 ES 文档覆盖或不完整。

## 七、排障路径

### ES 查不到数据

```text
1. 查 owner MySQL：业务数据是否真的存在
2. 查 binlog / RDS slave：主从延迟、binlog、ROW format
3. 查 Canal Instance：是否 running、位点是否推进
4. 查 Kafka：topic、lag、partition 是否阻塞
5. 查 Adapter：日志、DLQ、script parse、datasource 连接
6. 查 ES：index、mapping、document id、查询 DSL
```

### 某地区有数据，某地区没有

常见原因：

- 该地区 Canal Instance 未启动。
- 该地区 datasource 或 script 配置错误。
- ES mapping 在不同环境/地区不一致。
- 查询 DSL 使用的字段类型不适配。

### 查询报 `all shards failed`

常见原因：

- ES mapping 缺字段或字段类型不对。
- 对 text 字段用 term 精确查询。
- range 查询字段不是 date/number。

### Kafka lag 持续增长

常见原因：

- Adapter 消费异常。
- 某条消息一直失败阻塞 partition。
- 下游 ES 写入慢。
- datasource 回查慢或连接失败。

处理思路：

- 看 adapter 日志和 DLQ。
- 判断是否某个 script 或某个 table 导致。
- 必要时暂停目标 instance/script，避免影响其他模块。

## 八、设计哲学

### 1. MySQL 是事实源，ES 是查询视图

ES 用来加速 Admin 查询和分析，但不负责业务状态正确性。业务事实要回到 Product/Order/Policy/Marketing 等 owner 服务。

### 2. 配置先于代码

Canal 的核心工作很多不是 Java 代码，而是 instance、datasource、sync script、mapping、topic、groupId、operator、is_deleted 的一致性。

### 3. 切换必须可回滚

RDS 扩容或迁移时，新旧 instance、script、datasource 要可切换。旧配置不要直接删除，回滚 SQL 和旧 `application.yml` 要提前准备。

### 4. Mapping 要谨慎

字段类型错了往往不能直接改，可能要重建索引和重刷历史数据。上线前先确认 ES mapping，再发布 sync script。

### 5. ETL 要关注业务主键

分批字段不等于表主键。只按 `id` 分批可能破坏业务文档完整性，尤其是聚合型宽表。

## 九、面试叙事顺序

1. **定位**：Canal 是 MySQL 到 ES 的近实时同步链路，支撑 Admin 查询和报表。
2. **原理**：Canal 模拟 MySQL slave 解析 binlog，经 Kafka 给 Adapter。
3. **配置模型**：Instance、topic、datasource、sync script、ES mapping、groupId。
4. **Insurance 使用场景**：Policy/Order/Promotion/Marketing 的宽表、分表聚合、复杂查询。
5. **切换经验**：RDS 扩容要新增 instance、预置 script、发布 datasource、AT 验证、再 Main。
6. **排障**：MySQL -> Instance -> Kafka -> Adapter -> ES -> DSL 分层看。
7. **边界**：ES 不是事实源，不能用 ES 单独判断业务数据是否存在。

## 十、常见追问

### Canal 为什么比直接查 MySQL 更适合 Admin 查询？

Admin 查询经常跨表、跨分片、条件复杂、需要报表统计。直接扫 MySQL 会影响业务库，也很难做宽表查询。Canal 把业务变更同步到 ES，Admin 可以用 ES 做复杂查询和聚合。

### Canal 同步是强一致吗？

不是。它是近实时/最终一致链路，中间有 binlog、Kafka、Adapter、ES 写入和 mapping。业务事实仍然以 MySQL owner 服务为准。

### ES mapping 为什么重要？

查询能力取决于字段类型。比如 text 字段用 term 精确查询会查不到，date/number 类型错了 range 查询会报错。已有字段类型不能随便改，严重时要重建索引。

### sync script 写 SQL 要注意什么？

主表要清晰，关联条件要稳定，ES `_id` 要能唯一表示业务文档。`etlCondition` 要和业务主键匹配；复杂聚合要避免跨批次覆盖。

### RDS 扩容时 Canal 要改什么？

要新增 Canal Instance，更新 Adapter datasource，预置新 ES sync script，切换时启用新脚本并软删旧脚本，发布 `application.yml`，打开新 instance，验证 Kafka lag、DLQ、ES 写入和抽样数据。

### 同步阻塞怎么处理？

先定位阻塞在 Instance、Kafka、Adapter、DLQ 还是 ES。异常消息如果持续失败，要分类处理，必要时进入 DLQ 或暂停目标脚本，不能让一个错误消息长期阻塞整个 partition。

## 十一、可讲项目经验模板

稳妥版：

> 我参与/熟悉 Insurance 的 Canal/ES 同步链路。这个链路把 owner 服务 MySQL 的 binlog 同步到 ES，支撑 O-BFF/Admin 的复杂查询和报表。我能从 Instance、Kafka topic、Adapter datasource、ES sync script、mapping 和 ETL 几个层面排查问题。比如 RDS 扩容时，要先准备新 instance 和 datasource，预置新 script 为 disabled，切换时启用新 script、软删旧 script，并观察 Kafka lag、DLQ、Adapter 日志和 ES 抽样数据。

具体个人项目：`待补充`。

## 十二、优化项重点记忆

这些是 Canal/ES 相关 Confluence 里最值得背的优化项，面试问“同步链路怎么优化 / 怎么避免阻塞 / 怎么设计 ES”时优先讲。

| 优化项 | 解决的问题 | 面试怎么讲 |
| --- | --- | --- |
| 同步阻塞 DLQ 优化 | 一条坏数据影响整个批次，ACK 后又可能永久丢失。 | 异常分类处理：可重试异常阻塞或重试；脚本/mapping/数据问题写 DLQ；好数据继续写 ES。 |
| ETL 阶段和 ES Bulk 阶段分开 | 提交 ES 前不知道具体失败记录，提交后可以从 BulkResponse 找失败记录。 | ETL 异常对失败批次逐条重试定位；ES Bulk 异常解析 response 提取失败记录，无需重复提交整批。 |
| DLQ 手动重试和时间窗口补数 | 问题修复后需要恢复异常消息，最终失败还要知道补数范围。 | 修复 mapping/script 后手动消费 DLQ；最终失败上报时间窗口，用于后续全量刷数补齐。 |
| DLQ 监控指标 | 异常进入死信后如果不可见，仍然难排障。 | 记录 `canal_dlq_produce_total`、retry success、final failure，label 包括 destination、config_file、es_index、exception_type。 |
| ES index 设计规范 | 滥建 index、字段类型随意、dynamic keyword 会造成容量和查询问题。 | 只有分库分表非分片键查询、聚合查询、复杂模糊查询等场景才建 ES；字段最小化，mapping 先行。 |
| 分片容量评估 | shard 太大影响同步和搜索，太多增加集群负担。 | 业务类单分片建议控制在 20G 级别，单分片文档数不高于 1 亿，副本和 refresh interval 按业务评估。 |
| 模板生成工具 | 分表很多，手写成百上千份 YAML 容易错。 | 用模板 + 生成脚本展开 region/index，统一 dataSourceKey、destination、SQL 和 etlCondition。 |
| ETL 深分页优化 | 大表用 limit offset 会越来越慢，拖慢全量同步。 | 用 ID、时间或业务字段预先分段，再多次触发 ETL，避免深分页。 |
| 聚合宽表按业务主键分批 | 按表主键 `id` 分批时，同一业务文档跨批会被后批覆盖。 | 聚合型 ETL 要按业务主键分批；小表可不分批，大表要选有索引、区分度高的业务键。 |
| RDS 扩容可回滚切换 | 新旧 RDS、instance、datasource、script 同时变化，风险高。 | 新 instance 初始关闭，新 script 初始 `is_deleted=1`；AT 验证后再 Main；旧 script/datasource 保留回滚。 |

一句话总结：

> Canal 的优化主线是“不让坏数据拖垮整条同步链路”：异常分类、DLQ、Bulk 部分失败解析、监控可见、ETL 分批正确、ES mapping 先行、RDS 切换可回滚。

## 十三、资料来源

- Confluence：`ES 使用分享`，页面 ID `2347394775`。
- Confluence：`同步阻塞优化`，页面 ID `3045066306`。
- Confluence：`ES以及Canal同步工具维护说用`，页面 ID `1857867981`。
- Confluence：`ES设计规范`，页面 ID `2083164892`。
- 本地 README：`/Users/si.chen/GolandProjects/canal-adapter/README.md`
- live 配置：`/Users/si.chen/GolandProjects/canal-adapter/conf/live/application.yml`
- 模板：`/Users/si.chen/GolandProjects/canal-adapter/tools/template/`
- 模板生成工具：`/Users/si.chen/GolandProjects/canal-adapter/tools/template_to_yml/generate_es_config.go`
- mapping 工具说明：`/Users/si.chen/GolandProjects/canal-adapter/tools/README_MAPPING_EXTRACTOR.md`
- RDS 扩容 TD：`/Users/si.chen/GolandProjects/canal-adapter/docs/id_vn_rds_expansion_canal_td.md`
- RDS 扩容清单：`/Users/si.chen/GolandProjects/canal-adapter/docs/id_vn_rds_expansion_canal_checklist.md`
- ETL 聚合问题修复：`/Users/si.chen/GolandProjects/canal-adapter/docs/canal-adapter-etl-group-by-fix.md`
