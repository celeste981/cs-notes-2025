# MySQL 面试八股整理

目标：先把每道题按模块放好，再把答案整理成“先能听懂，再能复述，最后能应对追问、结合项目”的形式。

当前阶段：优先补齐 MySQL 高频八股。

推荐模板：[`../../templates/mysql-question.md`](../../templates/mysql-question.md)

## 分类目录

1. `01_architecture_storage_engine`：基础架构与存储引擎
   - Server 层、存储引擎层、连接器、分析器、优化器、执行器、InnoDB
2. `02_index`：索引
   - B+ 树、聚簇索引、二级索引、回表、覆盖索引、最左前缀、索引失效
3. `03_transaction_mvcc_lock`：事务、MVCC 与锁
   - ACID、隔离级别、Read View、undo log、间隙锁、临键锁、死锁
4. `04_sql_optimization`：SQL 优化与执行计划
   - explain、慢查询、扫描行数、Join 优化、分页优化、排序优化
5. `05_log_replication_ha`：日志、复制与高可用
   - redo log、binlog、undo log、主从复制、半同步、故障恢复
6. `06_sharding`：分库分表
   - 水平拆分、垂直拆分、分片键、全局 ID、跨库事务、扩容迁移
7. `07_troubleshooting`：排障与线上问题
   - 连接数打满、慢 SQL、锁等待、CPU 飙高、主从延迟、数据不一致

## 题目清单

| 编号 | 问题 | 分类 | 难度 | 文件 |
| --- | --- | --- | --- | --- |
| M1-01 | MySQL 一条查询语句从客户端发送到返回结果，中间经过了哪些步骤？ | 基础架构与存储引擎 | L2/L3 | [M1-01-query-flow.md](01_architecture_storage_engine/M1-01-query-flow.md) |
| M1-02 | 一条 UPDATE 语句的完整执行流程是什么？涉及哪些日志和缓冲组件？ | 基础架构与存储引擎 | L2/L3 | [M1-02-update-flow.md](01_architecture_storage_engine/M1-02-update-flow.md) |
| M1-03 | InnoDB 和 MyISAM 有什么本质区别？为什么 InnoDB 成为默认引擎？ | 基础架构与存储引擎 | L2/L3 | [m1-03-innodb-vs-myisam.md](01_architecture_storage_engine/m1-03-innodb-vs-myisam.md) |
| M1-04 | Buffer Pool 的工作原理是什么？为什么不用简单的 LRU 而用改进版本？ | 基础架构与存储引擎 | L2/L3 | [m1-04-buffer-pool-lru.md](01_architecture_storage_engine/m1-04-buffer-pool-lru.md) |
| M1-05 | InnoDB 的数据页结构是怎样的？一页能存多少条记录？行格式有什么区别？ | 基础架构与存储引擎 | L3 | [m1-05-innodb-page-row-format.md](01_architecture_storage_engine/m1-05-innodb-page-row-format.md) |
| M1-06 | Change Buffer 是什么？为什么唯一索引不能用 Change Buffer？ | 基础架构与存储引擎 | L2/L3 | [m1-06-change-buffer-unique-index.md](01_architecture_storage_engine/m1-06-change-buffer-unique-index.md) |
| M4-01 | 深分页慢查询怎么优化？为什么要改成游标分页并加复合索引？ | SQL 优化与执行计划 | L2/L3 | [m4-01-deep-pagination-index-optimization.md](04_sql_optimization/m4-01-deep-pagination-index-optimization.md) |
| M5-01 | redo log、undo log 和 binlog 各自的作用是什么？它们之间是什么关系？ | 日志、复制与高可用 | L2/L3 | [m5-01-redo-undo-binlog.md](05_log_replication_ha/m5-01-redo-undo-binlog.md) |
| M5-02 | WAL 机制的原理是什么？为什么顺序写日志比随机写数据页快？ | 日志、复制与高可用 | L2/L3 | [m5-02-wal.md](05_log_replication_ha/m5-02-wal.md) |
| M5-03 | 两阶段提交是怎么保证 redo log 和 binlog 一致性的？崩溃恢复时怎么判断？ | 日志、复制与高可用 | L3 | [m5-03-two-phase-commit.md](05_log_replication_ha/m5-03-two-phase-commit.md) |
| M5-04 | innodb_flush_log_at_trx_commit 和 sync_binlog 的各种组合有什么含义？怎么选？ | 日志、复制与高可用 | L2/L3 | [m5-04-flush-params.md](05_log_replication_ha/m5-04-flush-params.md) |
| M5-05 | MySQL 崩溃恢复的完整流程是怎样的？ | 日志、复制与高可用 | L3 | [m5-05-crash-recovery.md](05_log_replication_ha/m5-05-crash-recovery.md) |
| M5-06 | binlog 的三种格式 Statement、Row、Mixed 各有什么优缺点？生产环境怎么选？ | 日志、复制与高可用 | L2/L3 | [m5-06-binlog-format.md](05_log_replication_ha/m5-06-binlog-format.md) |
| M5-07 | Online DDL 的实现原理是什么？大表加字段用 pt-osc 和 gh-ost 有什么区别？ | 日志、复制与高可用 | L3 | [m5-07-online-ddl.md](05_log_replication_ha/m5-07-online-ddl.md) |
| M6-01 | 做过数据库扩容吗？分库扩容怎么做？遇到过什么坑？ | 分库分表 | L3/L4 | [m6-01-sharding-expansion-migration-case.md](06_sharding/m6-01-sharding-expansion-migration-case.md) |

## 公司服务视角

| 主题 | 说明 | 文件 |
| --- | --- | --- |
| Shopee Insurance 服务如何使用 MySQL | 把 MySQL 八股里的客户端、连接器、长连接、主从、分表、事务映射到内部服务实践 | [company-service-db-usage.md](01_architecture_storage_engine/company-service-db-usage.md) |

## 可视化页面

| 模块 | 页面 |
| --- | --- |
| MySQL 术语速查 | [glossary.html](glossary.html) |
| 查询流程 | [M1-01-query-flow.html](01_architecture_storage_engine/M1-01-query-flow.html) |
| UPDATE 流程 | [M1-02-update-flow.html](01_architecture_storage_engine/M1-02-update-flow.html) |
| InnoDB vs MyISAM | [m1-03-innodb-vs-myisam.html](01_architecture_storage_engine/m1-03-innodb-vs-myisam.html) |
| Buffer Pool LRU | [m1-04-buffer-pool-lru.html](01_architecture_storage_engine/m1-04-buffer-pool-lru.html) |
| 数据页与行格式 | [m1-05-innodb-page-row-format.html](01_architecture_storage_engine/m1-05-innodb-page-row-format.html) |
| Change Buffer | [m1-06-change-buffer-unique-index.html](01_architecture_storage_engine/m1-06-change-buffer-unique-index.html) |
| 深分页优化 | [m4-01-deep-pagination-index-optimization.html](04_sql_optimization/m4-01-deep-pagination-index-optimization.html) |
| 日志与崩溃恢复总览 | [m5-log-recovery-map.html](05_log_replication_ha/m5-log-recovery-map.html) |
| 分库扩容迁移 | [m6-01-sharding-expansion-migration-case.html](06_sharding/m6-01-sharding-expansion-migration-case.html) |

## 每道题的整理格式

1. 先给 `先讲人话`，用直觉解释这个问题在问什么。
2. 补 `前置概念`，把必须懂的术语讲清楚。
3. 再给 `30 秒短答`，保证面试开口有结论。
4. 再给 `1-2 分钟标准回答`，方便直接复述。
5. 用图或步骤拆开原理。
6. 补 `结合项目怎么讲`，把八股落到公司服务实践。
7. 最后准备 2 到 3 个高频追问和容易说错的点。

## 近期建议顺序

1. 先看 [`glossary.md`](glossary.md)：把常见术语混个脸熟。
2. `01_architecture_storage_engine`：先掌握查询流程、Server 层与 InnoDB 分工、Buffer Pool。
3. `02_index`：B+ 树、聚簇索引、二级索引、回表、覆盖索引、最左前缀、索引失效。
4. `03_transaction_mvcc_lock`：ACID、隔离级别、MVCC、Read View、undo log、行锁、间隙锁、死锁。
5. `04_sql_optimization`：explain、慢 SQL、Join、排序、分页、索引选择。
6. `05_log_replication_ha`：redo log、binlog、两阶段提交、主从复制、主从延迟。
7. `06_sharding`：分库分表、分片键、全局 ID、跨库事务、扩容迁移。
8. `07_troubleshooting`：连接打满、锁等待、CPU 飙高、慢 SQL、数据不一致。

## 编号规则

- `M1-xx`：基础架构与存储引擎。
- `M2-xx`：索引。
- `M3-xx`：事务、MVCC 与锁。
- `M4-xx`：SQL 优化。
- `M5-xx`：日志、复制与高可用。
- `M6-xx`：分库分表。
- `M7-xx`：排障与线上问题。
