# Golang 面试八股整理

目标：把 Go 高频面试题整理成“先能听懂，再能复述，最后能接追问、结合项目”的复习材料。

当前阶段：先整理基础与类型系统。今晚学习的是 G1-01 到 G1-04 四道题。

## 分类目录

1. `01_basics_type_system`：基础与类型系统
   - slice、map、defer、for range、struct 内存对齐、string、make/new、值传递、零值行为

## 题目清单

| 编号 | 问题 | 分类 | 难度 | 文件 | 复习状态 |
| --- | --- | --- | --- | --- | --- |
| G1-01 | slice 扩容机制：append 之后旧切片还能用吗？ | 基础与类型系统 | L2/L3 | [g1-01-slice-growth.md](01_basics_type_system/g1-01-slice-growth.md) | 已成稿 |
| G1-02 | map 底层实现：哈希冲突如何解决？扩容是怎么做的？ | 基础与类型系统 | L2/L3 | [g1-02-map-internals.md](01_basics_type_system/g1-02-map-internals.md) | 已成稿 |
| G1-03 | defer 的执行顺序和陷阱：defer、return、函数返回值的执行时序 | 基础与类型系统 | L2 | [g1-03-defer-return-order.md](01_basics_type_system/g1-03-defer-return-order.md) | 已成稿 |
| G1-04 | for range 的复制陷阱：遍历时修改元素为什么不生效？ | 基础与类型系统 | L2 | [g1-04-for-range-copy.md](01_basics_type_system/g1-04-for-range-copy.md) | 已成稿 |

## 建议复习顺序

1. 先看 `slice`，理解 Go 里“看起来像引用，实际上传的是 header 拷贝”的思想。
2. 再看 `map`，重点记住桶、overflow bucket、渐进式扩容、并发不安全。
3. 然后看 `defer`，把 `return 赋值 -> defer 执行 -> 函数返回` 背熟。
4. 最后看 `for range`，把值拷贝、闭包捕获、Go 1.22 前后差异讲清楚。

## 每道题的整理格式

1. `先讲人话`：先用直觉解释问题。
2. `前置概念`：补齐必须懂的基础词。
3. `30 秒短答`：面试开口先给结论。
4. `1-2 分钟标准回答`：可直接复述。
5. `原理拆解`：应对追问。
6. `结合项目怎么讲`：不编造经历，缺少细节标为 `待补充`。
7. `常见追问` 和 `易错点`：用于考前快速复盘。
