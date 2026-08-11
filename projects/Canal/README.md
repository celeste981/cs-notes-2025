# Canal 项目材料

这个目录用于沉淀 Insurance Canal / ES 同步链路的项目背景、架构梳理、配置物料、迁移切换和面试表达。复习时优先抓住一条主线：`MySQL binlog` -> `Canal Instance` -> `Kafka topic` -> `Canal Adapter` -> `ES sync script` -> `Elasticsearch` -> `Admin/O-BFF 查询`。

## 当前文档

- [Canal / ES 同步架构梳理](canal-es-sync-architecture.md)：核心项目文档，适合回答 Canal 原理、Insurance ES 使用场景、sync script、RDS 扩容切换、ETL、DLQ 和排障。
- [Canal / ES 同步架构图解](canal-es-sync-architecture.html)：上面架构文档的可视化版本。
- [Canal DLQ / 同步阻塞优化专项](canal-dlq-blocking-optimization.md)：专门背同步阻塞、异常分类、Bulk 部分失败、DLQ、重试和补数。
- [Canal / ES 核心知识点速记](canal-es-key-knowledge.md)：专门背 binlog、instance、Kafka、adapter、sync script、mapping、ETL、DLQ 等基础知识。
- [Canal 重点架构：RDS 切换、ETL、ES 宽表](canal-rds-etl-architecture-notes.md)：专门背 RDS 扩容切换、ETL 分批、宽表设计和回滚口径。
- [Insurance 总览](../Insurance/README.md)：从整体业务地图看 Canal 与 O-BFF、Marketing、Marketing-Data 的关系。

## 本机项目路径

| 项目 | 本机路径 | 面试定位 |
| --- | --- | --- |
| `canal-adapter` | `/Users/si.chen/GolandProjects/canal-adapter` | MySQL binlog 到 ES 的同步基础设施，支持 adapter、sync script、datasource、ETL、RDS 扩容切换和回滚物料。 |

## 代码和物料入口速查

| 想看什么 | 推荐路径 |
| --- | --- |
| Canal README | `/Users/si.chen/GolandProjects/canal-adapter/README.md` |
| live adapter 配置 | `/Users/si.chen/GolandProjects/canal-adapter/conf/live/application.yml` |
| ES6/ES7 adapter | `/Users/si.chen/GolandProjects/canal-adapter/es6x/`, `/Users/si.chen/GolandProjects/canal-adapter/es7x/` |
| sync script 模板 | `/Users/si.chen/GolandProjects/canal-adapter/tools/template/` |
| 模板生成工具 | `/Users/si.chen/GolandProjects/canal-adapter/tools/template_to_yml/generate_es_config.go` |
| mapping 提取工具说明 | `/Users/si.chen/GolandProjects/canal-adapter/tools/README_MAPPING_EXTRACTOR.md` |
| RDS 扩容迁移 TD | `/Users/si.chen/GolandProjects/canal-adapter/docs/id_vn_rds_expansion_canal_td.md` |
| RDS 扩容执行清单 | `/Users/si.chen/GolandProjects/canal-adapter/docs/id_vn_rds_expansion_canal_checklist.md` |
| ETL 聚合问题修复 | `/Users/si.chen/GolandProjects/canal-adapter/docs/canal-adapter-etl-group-by-fix.md` |

## 面试复习顺序

1. 先看 [Canal / ES 同步架构图解](canal-es-sync-architecture.html)，记住从 MySQL 到 ES 的完整链路。
2. 再看 [Canal / ES 同步架构梳理](canal-es-sync-architecture.md)，补充配置、切换和排障。
3. 单独背三个专项：[核心知识点](canal-es-key-knowledge.md)、[DLQ 阻塞优化](canal-dlq-blocking-optimization.md)、[RDS/ETL/宽表架构](canal-rds-etl-architecture-notes.md)。
4. 回到 [Insurance 总览](../Insurance/README.md)，理解 Canal 为什么支撑 Admin 查询、Marketing 人群和 O-BFF 报表。

## 待补充事实

- 你实际做过的 ES sync script、mapping、RDS 扩容、切换、回滚、补数或同步阻塞排障。
- 某个真实 index 的 mapping、sync script、数据源、topic、instance 名。
- 具体 PR、SWP/DBPortal/Jira、上线记录、验证截图。
