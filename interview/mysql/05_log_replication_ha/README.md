# MySQL 05 日志、复制与高可用

> 来源：LC100 MySQL 日志专题，已按本地 MySQL 编号规则重排为 `M5-xx`。  
> 重点：redo log、undo log、binlog、WAL、两阶段提交、崩溃恢复、binlog 格式、Online DDL。

可视化总览：[m5-log-recovery-map.html](m5-log-recovery-map.html)

## 单题可视化

- [M5-01 三大日志](m5-01-redo-undo-binlog.html)
- [M5-02 WAL 机制](m5-02-wal.html)
- [M5-03 两阶段提交](m5-03-two-phase-commit.html)
- [M5-04 刷盘参数](m5-04-flush-params.html)
- [M5-05 崩溃恢复](m5-05-crash-recovery.html)
- [M5-06 binlog 格式](m5-06-binlog-format.html)
- [M5-07 Online DDL](m5-07-online-ddl.html)

## 题目清单

| 编号 | 问题 | 难度 | 优先级 | 文件 |
| --- | --- | --- | --- | --- |
| M5-01 | redo log、undo log 和 binlog 各自的作用是什么？它们之间是什么关系？ | L2/L3 | P0 | [m5-01-redo-undo-binlog.md](m5-01-redo-undo-binlog.md) |
| M5-02 | WAL 机制的原理是什么？为什么顺序写日志比随机写数据页快？ | L2/L3 | P0 | [m5-02-wal.md](m5-02-wal.md) |
| M5-03 | 两阶段提交是怎么保证 redo log 和 binlog 一致性的？崩溃恢复时怎么判断？ | L3 | P0 | [m5-03-two-phase-commit.md](m5-03-two-phase-commit.md) |
| M5-04 | `innodb_flush_log_at_trx_commit` 和 `sync_binlog` 的各种组合有什么含义？怎么选？ | L2/L3 | P0 | [m5-04-flush-params.md](m5-04-flush-params.md) |
| M5-05 | MySQL 崩溃恢复的完整流程是怎样的？ | L3 | P0 | [m5-05-crash-recovery.md](m5-05-crash-recovery.md) |
| M5-06 | binlog 的三种格式 Statement、Row、Mixed 各有什么优缺点？生产环境怎么选？ | L2/L3 | P1 | [m5-06-binlog-format.md](m5-06-binlog-format.md) |
| M5-07 | Online DDL 的实现原理是什么？大表加字段用 pt-osc 和 gh-ost 有什么区别？ | L3 | P1 | [m5-07-online-ddl.md](m5-07-online-ddl.md) |

## 这组题怎么复习

1. 先看 `M5-01`，把三类日志的职责边界背清楚。
2. 再看 `M5-02`，理解为什么 MySQL 不直接每次都刷数据页。
3. 重点攻 `M5-03` 和 `M5-05`，这是最容易被追问崩溃点推演的部分。
4. `M5-04` 用来连接生产参数选择，必须能讲数据安全和性能的取舍。
5. `M5-06` 和 `M5-07` 偏生产实践，重点讲主从一致性、大表 DDL 风险和回滚方案。

## 一句话主线

```text
undo 管回滚和 MVCC，redo 管崩溃恢复，binlog 管复制和归档；
WAL 让提交只等顺序日志刷盘；
两阶段提交让 redo 和 binlog 对同一事务保持一致；
崩溃恢复时用 redo 前滚、binlog 判定 prepare 事务、undo 回滚未提交事务。
```

## 高频追问地图

| 追问方向 | 要能答出 |
| --- | --- |
| 为什么有 binlog 还要 redo log？ | binlog 是 Server 层逻辑日志，不 crash-safe，不能高效恢复脏页。 |
| 为什么有 redo log 还要 binlog？ | redo 循环写且 InnoDB 独有，不能做长期归档、主从复制和 PITR。 |
| redo prepare 后宕机怎么办？ | 看 binlog 是否有完整 XID，有则提交，无则回滚。 |
| 双 1 配置为什么慢？ | 每次提交至少涉及 redo fsync 和 binlog fsync，安全性换延迟。 |
| Online DDL 为什么还会卡住？ | DDL 需要 MDL 写锁，长事务持有 MDL 读锁会阻塞 DDL，DDL 排队后又阻塞后续 DML。 |
