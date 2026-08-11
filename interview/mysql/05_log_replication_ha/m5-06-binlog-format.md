# M5-06 binlog 的三种格式 Statement、Row、Mixed 各有什么优缺点？生产环境怎么选？


可视化页面：[m5-06-binlog-format.html](m5-06-binlog-format.html)
分类：日志、复制与高可用

难度：L2/L3

优先级：P1

关键词：binlog_format、Statement、Row、Mixed、主从复制、binlog_row_image、数据恢复

复习状态：已成稿

来源题号：LC100 M4-06

## 问题

`binlog_format` 的 `STATEMENT`、`ROW`、`MIXED` 分别是什么？为什么生产环境通常推荐 `ROW`？

## 先讲人话

binlog 可以用三种方式记录一次修改：

1. 只记“我执行了哪条 SQL”。
2. 详细记“哪一行从什么值变成什么值”。
3. 平时记 SQL，遇到危险 SQL 再自动改成记行变化。

越详细越安全，但日志越大。

## 前置概念

| 概念 | 小白解释 |
| --- | --- |
| Statement Based Replication | 基于 SQL 语句复制，从库重新执行 SQL。 |
| Row Based Replication | 基于行变更复制，从库按行应用变化。 |
| Mixed | MySQL 自动在 Statement 和 Row 之间选择。 |
| 不确定函数 | 如 `UUID()`、`RAND()`、部分 `NOW()` 场景，主从重新执行可能不同。 |
| `binlog_row_image` | Row 格式下记录哪些列，影响 binlog 体积。 |

## 30 秒短答

`STATEMENT` 记录原始 SQL，日志量小，但遇到不确定函数、触发器、存储过程等场景可能导致主从结果不一致。`ROW` 记录每一行的具体变更，一致性最好，也方便恢复和审计，但日志量大。`MIXED` 由 MySQL 判断安全 SQL 用 Statement，不安全时切 Row，兼顾空间和一致性，但行为不如 Row 直观。

生产一般推荐 `ROW`，再配合 `binlog_row_image=MINIMAL` 降低日志体积。

## 三种格式对比

| 格式 | 记录内容 | 优点 | 缺点 | 适合场景 |
| --- | --- | --- | --- | --- |
| `STATEMENT` | 原始 SQL | 日志小、可读性强 | 主从可能因不确定 SQL 不一致 | 简单 SQL、旧系统，不推荐新项目默认用 |
| `ROW` | 每行变更 | 一致性最好、恢复精确 | 日志量大 | 生产主从、数据恢复、审计 |
| `MIXED` | 自动选择 | 兼顾空间和安全 | 行为不够透明 | 想保守节省空间但接受复杂度 |

## 面试时可以直接这样回答

binlog 有三种格式。

`Statement` 记录 SQL 原文，比如一条 `update user set score = score + 1 where city = 'SZ'`。优点是日志量小，一条 SQL 影响多少行都只记一条。缺点是从库要重新执行 SQL，如果 SQL 中有 `UUID()`、`RAND()` 这类不确定函数，或者执行环境不同，就可能主从结果不一致。

`Row` 记录每一行的变更，主库哪一行从什么值变成什么值，从库就应用同样的变更，所以一致性最好。缺点是大批量更新时 binlog 可能很大。

`Mixed` 是 MySQL 自动选择，安全 SQL 用 Statement，不安全 SQL 用 Row。它折中但也增加不确定性。

生产环境我会优先选 Row，因为主从一致性和恢复准确性更重要。如果担心日志体积，可以配合 `binlog_row_image=MINIMAL`，只记录必要列和修改列。

## `binlog_row_image`

| 取值 | 含义 | 特点 |
| --- | --- | --- |
| `FULL` | 记录修改前后所有列 | 信息最完整，日志最大。 |
| `MINIMAL` | 只记录必要标识列和被修改列 | 日志更小，生产常用于减小 Row 体积。 |
| `NOBLOB` | 类似 FULL，但未修改 BLOB/TEXT 不记录 | 对大字段友好。 |

## 原理拆解

### 1. 为什么 Statement 可能不安全

例如：

```sql
insert into token(id, value) values(1, uuid());
```

主库执行一次，生成一个 UUID；从库重新执行，可能生成另一个 UUID。主从数据就不一样。

### 2. 为什么 Row 更适合恢复

Row 记录的是行级变化。做数据订正或回放时，可以更明确知道哪一行被改成了什么，而不是重新理解 SQL 的语义。

### 3. Row 的代价是什么

一条 SQL 如果影响 10 万行，Row 格式可能记录 10 万行变更，binlog 体积、网络传输和从库回放压力都会增加。

## 结合项目怎么讲

可以这样说：

> 如果业务依赖主从复制、binlog 订阅或按时间点恢复，我会优先选 Row。因为数据一致性和可恢复性比节省一点日志空间更重要。对于大批量更新，会通过分批执行、控制事务大小、监控 binlog 增长和主从延迟来降低风险。

待补充：

- 当前项目是否有 Canal/binlog consumer。
- 是否遇到过大事务导致主从延迟。
- 是否使用 `binlog_row_image=MINIMAL`。

## 容易说错的点

1. 说 Row 没有缺点。Row 的缺点是日志量大。
2. 说 Statement 一定不一致。它是某些场景不安全，不是所有 SQL 都不安全。
3. 忽略 `binlog_row_image`。
4. 把 binlog 格式和 redo log 类型混在一起。

## 高频追问

### 追问 1：为什么 RC 隔离级别更推荐 Row？

Statement 在 RC 下更容易因为执行环境和可见数据差异导致主从不一致，Row 直接复制行变更，安全性更高。

### 追问 2：Row 格式大事务有什么风险？

binlog 暴增、传输变慢、从库回放慢、主从延迟扩大。处理方式是分批、限速、避开高峰。

### 追问 3：Mixed 是不是最优？

不一定。Mixed 让 MySQL 自动判断格式，减少部分日志量，但也让行为不够直观。生产为了确定性，通常更偏向 Row。

## 记忆口诀

```text
Statement 记 SQL，Row 记行变；
Row 更稳，日志更大，MINIMAL 来减肥。
```

