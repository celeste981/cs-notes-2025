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

## 公司服务视角

| 主题 | 说明 | 文件 |
| --- | --- | --- |
| Shopee Insurance 服务如何使用 MySQL | 把 MySQL 八股里的客户端、连接器、长连接、主从、分表、事务映射到内部服务实践 | [company-service-db-usage.md](01_architecture_storage_engine/company-service-db-usage.md) |

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
