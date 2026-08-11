# Go 基础与类型系统

本目录整理 Go 基础类型系统相关面试题。重点不是背源码字段名，而是能在面试里判断代码输出、解释底层原因，并说明工程上如何避坑。

## 今晚学习清单

| 编号 | 问题 | 难度 | 文件 |
| --- | --- | --- | --- |
| G1-01 | slice 扩容机制：append 之后旧切片还能用吗？ | L2/L3 | [g1-01-slice-growth.md](g1-01-slice-growth.md) |
| G1-02 | map 底层实现：哈希冲突如何解决？扩容是怎么做的？ | L2/L3 | [g1-02-map-internals.md](g1-02-map-internals.md) |
| G1-03 | defer 的执行顺序和陷阱：defer、return、函数返回值的执行时序 | L2 | [g1-03-defer-return-order.md](g1-03-defer-return-order.md) |
| G1-04 | for range 的复制陷阱：遍历时修改元素为什么不生效？ | L2 | [g1-04-for-range-copy.md](g1-04-for-range-copy.md) |

## 共同主线

这四题都在考一个核心能力：

```text
不要只看 Go 语法表象，要能判断变量背后到底拷贝了什么、共享了什么、什么时候延迟执行。
```

| 题目 | 面试关键 |
| --- | --- |
| slice | slice header 是值拷贝，底层数组可能共享。 |
| map | map 底层是哈希表，扩容是渐进式的，不能并发读写。 |
| defer | defer 延迟执行，但参数立即求值；return 不是一步完成。 |
| for range | 循环变量通常是值拷贝，Go 1.22 前闭包捕获有经典坑。 |
