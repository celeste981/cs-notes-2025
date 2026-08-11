# M5-04 innodb_flush_log_at_trx_commit 和 sync_binlog 的各种组合有什么含义？怎么选？


可视化页面：[m5-04-flush-params.html](m5-04-flush-params.html)
分类：日志、复制与高可用

难度：L2/L3

优先级：P0

关键词：刷盘策略、fsync、redo log、binlog、双 1 配置、RPO、性能取舍

复习状态：已成稿

来源题号：LC100 M4-04

## 问题

`innodb_flush_log_at_trx_commit` 和 `sync_binlog` 分别控制什么？生产环境怎么选？

## 先讲人话

这两个参数本质上是在问：事务提交时，日志到底要不要立刻“真正写进磁盘”。

写得越勤，越安全，但每次提交都要等磁盘同步，性能更差。写得越懒，吞吐更高，但机器断电或 OS 崩溃时可能丢最近一小段事务。

## 前置概念

| 概念 | 小白解释 |
| --- | --- |
| write | 写到操作系统 page cache，不一定已经落到物理磁盘。 |
| fsync | 要求操作系统把文件数据同步到磁盘，成本更高。 |
| MySQL 崩溃 | MySQL 进程挂了，但操作系统还在。 |
| OS 崩溃/断电 | 操作系统和 page cache 都可能丢失。 |
| RPO | Recovery Point Objective，最多能接受丢多少数据。 |

## 30 秒短答

`innodb_flush_log_at_trx_commit` 控制 redo log 的写盘策略，最安全是 `1`，每次事务提交都 write + fsync；`sync_binlog` 控制 binlog 的刷盘策略，最安全也是 `1`，每次提交都 fsync。

生产核心业务通常用“双 1”：`innodb_flush_log_at_trx_commit=1` 且 `sync_binlog=1`，安全性最高但性能较差。非核心、可容忍少量丢失的场景可以放宽，例如 redo 设 `2`、binlog 设 `100` 或更大，但要明确 RPO 风险。

## 参数含义

### `innodb_flush_log_at_trx_commit`

| 取值 | 提交时做什么 | MySQL 进程崩溃 | OS 崩溃/断电 | 特点 |
| --- | --- | --- | --- | --- |
| `1` | redo 写入 OS cache 并 fsync | 不丢已提交事务 | 通常不丢已提交事务 | 最安全，最慢 |
| `2` | redo 写入 OS cache，不每次 fsync | 通常不丢 | 可能丢最近约 1 秒 | 性能和安全折中 |
| `0` | 提交时不立即 write，由后台约每秒写和刷 | 可能丢最近约 1 秒 | 可能丢最近约 1 秒 | 性能更高，风险更大 |

### `sync_binlog`

| 取值 | 含义 | 风险 |
| --- | --- | --- |
| `1` | 每次事务提交都 fsync binlog | 最安全，性能成本高 |
| `0` | 交给 OS 决定何时刷盘 | OS 崩溃可能丢 binlog |
| `N` | 每 N 次事务提交 fsync 一次 | 可能丢最近 N 个事务左右的 binlog |

## 面试时可以直接这样回答

这两个参数分别控制 redo log 和 binlog 的持久化强度。

`innodb_flush_log_at_trx_commit=1` 表示每次事务提交都把 redo log 写入 OS page cache 并 fsync 到磁盘，是最安全的。`=2` 表示每次提交只写到 OS page cache，不每次 fsync，所以 MySQL 进程崩溃通常不丢，但 OS 崩溃或断电可能丢。`=0` 表示提交时不立即写 redo，由后台线程大约每秒处理，MySQL 进程崩溃也可能丢最近事务。

`sync_binlog=1` 表示每次提交都 fsync binlog；`=0` 表示交给 OS；`=N` 表示每 N 次提交刷一次。

如果是核心交易、订单、支付类业务，我会倾向双 1，牺牲一些性能换数据安全。如果是日志、埋点、临时中间结果这类可重放或可容忍少量丢失的数据，可以适当放宽，但要和业务确认 RPO。

## 组合选择

| 场景 | 推荐组合 | 说明 |
| --- | --- | --- |
| 金融/订单/库存核心链路 | `innodb_flush_log_at_trx_commit=1` + `sync_binlog=1` | 数据安全优先。 |
| 普通业务主库 | 优先双 1，压力大再评估优化 | 不要一开始为了性能牺牲一致性。 |
| 日志/埋点/可重放数据 | `2` + `100~1000` 可评估 | 接受少量丢失换吞吐。 |
| 开发/测试环境 | 可放宽 | 不代表生产也能这么配。 |

## 原理拆解

### 1. 为什么双 1 慢

双 1 下每个事务提交至少要关注两类同步：

1. redo log fsync。
2. binlog fsync。

虽然 group commit 能缓解，但磁盘同步仍然是写入延迟的重要来源。

### 2. 为什么 `=2` 时 MySQL 崩溃通常不丢 redo

`=2` 已经把 redo 写到 OS page cache。MySQL 进程挂掉时，操作系统还在，page cache 后续仍可能刷盘。

但如果是断电、内核 panic 或机器重启，page cache 可能丢失。

### 3. 为什么 `sync_binlog` 也重要

只保证 redo 安全还不够。binlog 如果丢了，主库本地可能恢复出数据，但从库和基于 binlog 的恢复链路没有这笔事务，仍然会产生一致性风险。

## 结合项目怎么讲

可以这样说：

> 我不会只从性能角度改这两个参数。会先和业务确认这张表的数据是否可重建、最多能接受丢多少秒或多少笔事务，然后再结合 QPS、提交延迟、磁盘 fsync 指标做选择。核心交易库默认双 1，非核心异步数据才考虑放宽。

待补充：

- 你的业务是否有订单/支付/库存这类强一致数据。
- 是否知道生产 MySQL 的这两个参数配置。
- 是否有磁盘 fsync 延迟监控。

## 容易说错的点

1. 说 `write` 等于落盘。真正落盘通常要 `fsync`。
2. 说 `innodb_flush_log_at_trx_commit=2` 在任何崩溃下都不丢。
3. 只调 redo 参数，不管 binlog 参数。
4. 不谈业务 RPO，直接说某个配置一定最好。

## 高频追问

### 追问 1：为什么双 1 也不是绝对 100% 安全？

还依赖磁盘、文件系统、缓存策略、硬件是否真正遵守 fsync 语义。但在数据库参数层面，双 1 是最强持久化配置。

### 追问 2：怎么优化双 1 的性能？

优先考虑 group commit、批量提交、减少小事务频率、优化磁盘能力，而不是直接放松安全参数。

### 追问 3：如果 sync_binlog=100，最多丢多少？

粗略理解为可能丢最近一批未 fsync 的 binlog，接近最多 100 个事务，但实际还受 OS 刷盘时机影响。

## 记忆口诀

```text
redo 看 innodb_flush，binlog 看 sync_binlog；
双 1 最安全，放宽先问 RPO。
```

