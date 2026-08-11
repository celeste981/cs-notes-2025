# Marketing-Data 项目材料

这个目录用于沉淀 Insurance Marketing-Data 的项目背景、架构梳理、离线数据链路和面试表达。复习时优先抓住一条主线：`Marketing Plan` -> `OfflineDataFetchAndPush` -> `batch history` -> `DataSource handler` -> `Kafka` -> `Marketing Engine`。

## 当前文档

- [Marketing-Data 架构梳理](marketing-data-architecture.md)：核心项目文档，适合回答服务拆分、离线人群、batch 幂等、DataSource handler、ES/文件/Hive/Insight、pre-filter 和监控。
- [Marketing-Data 架构图解](marketing-data-architecture.html)：上面架构文档的可视化版本。
- [Marketing-Data 核心知识深挖](marketing-data-key-knowledge-deep-dive.md)：专门背 OfflineDataFetchAndPush、batch history、服务边界和深挖速答。
- [Marketing-Data Handler / pre-filter 深挖](marketing-data-handler-prefilter-deep-dive.md)：专门背 DataSource handler、ES/File/Hive/Insight、pre-filter 实时预处理。
- [Marketing-Data 可靠性和一致性深挖](marketing-data-reliability-deep-dive.md)：专门背 batch 幂等、history、Kafka publish 风险、排障和不能夸大的点。
- [Insurance 总览中的 Marketing-Data 面试专项](../Insurance/marketing-data-interview-notes.md)：偏面试问答和防守边界。

## 本机项目路径

| 项目 | 本机路径 | 面试定位 |
| --- | --- | --- |
| `insurance-marketing-data` | `/Users/si.chen/GolandProjects/insurance-marketing-data` | Marketing 的离线/批量数据服务，负责离线数据拉取、ES/S3/CSV/Hive/Insight、Kafka 推送和离线执行记录。 |
| `insurance-marketing-data-api` | `/Users/si.chen/GolandProjects/insurance-marketing-data-api` | Marketing-Data API/proto 仓库。 |
| `insurance-marketing` | `/Users/si.chen/GolandProjects/insurance-marketing` | Marketing 主服务，消费 Marketing-Data 推送的数据后执行具体营销动作。 |

## 代码入口速查

| 想看什么 | 推荐路径 |
| --- | --- |
| Internal Service | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/service/marketing_data_internal_service.go` |
| Admin Service | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/service/marketing_data_admin_service.go` |
| OfflineDataFetchManager | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/manager/impl/offline_data_fetch_manager_impl.go` |
| Handler 接口 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/handler/offline_data_fetch_handler.go` |
| Handler 实现 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/handler/offline_data_fetch/impl/` |
| ES Repo | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/repo/elasticsearch/` |
| MQ Producer / Consumer | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/mq/` |
| 监控 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/monitor/` |

## 面试复习顺序

1. 先看 [Marketing-Data 架构图解](marketing-data-architecture.html)，记住 `离线抽取 = 幂等 + history + handler + Kafka + monitor`。
2. 再看 [Marketing-Data 架构梳理](marketing-data-architecture.md)，补充代码路径、核心流程和追问。
3. 单独背三个深挖专项：[核心知识](marketing-data-key-knowledge-deep-dive.md)、[Handler/pre-filter](marketing-data-handler-prefilter-deep-dive.md)、[可靠性一致性](marketing-data-reliability-deep-dive.md)。
4. 最后看 [Marketing-Data 面试专项](../Insurance/marketing-data-interview-notes.md)，准备 1 分钟回答和防守边界。

## 待补充事实

- 你在 Marketing-Data 上实际负责过的 data source、ES 查询、文件处理、batch history 或 Kafka 推送问题。
- 某个真实 Plan 的 `plan_id`、`batch_id`、`dest_topic`、`data_source` 和执行结果。
- 具体 PR、Jira、上线记录、排障记录。
