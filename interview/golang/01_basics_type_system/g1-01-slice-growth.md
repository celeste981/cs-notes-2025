# G1-01 slice 扩容机制：append 之后旧切片还能用吗？

分类：基础与类型系统

材料类型：interview question / knowledge topic

难度：L2/L3

优先级：P0

关键词：slice、append、len、cap、底层数组、扩容、三下标切片

复习状态：已成稿

来源：https://lc100.pages.dev/go

## 问题

slice 扩容机制是什么？`append` 之后旧切片还能不能继续用？

这道题常见考法是给一段 `append` 代码，让你判断输出结果。真正考察的是：

```text
slice 自己只是一个 header；
header 会被拷贝；
底层数组可能被多个 slice 共享；
append 是否扩容决定新旧 slice 是否互相影响。
```

## 先讲人话

slice 可以理解成一张“便利贴”，上面写着三件事：

```text
从哪个底层数组开始看
当前能看到几个元素 len
最多还能往后用多少容量 cap
```

把 slice 赋值给另一个变量，不是复制整段数组，而是复制这张便利贴。两张便利贴可能仍然指向同一块底层数组。

所以 `append` 后旧切片能不能继续用，关键看有没有扩容：

```text
cap 够：继续写原来的底层数组，新旧 slice 可能互相覆盖。
cap 不够：runtime 分配新数组并拷贝数据，新旧 slice 分离。
```

## 前置概念

| 概念 | 小白解释 |
| --- | --- |
| array | 固定长度数组，长度是类型的一部分，比如 `[3]int` 和 `[4]int` 是不同类型。 |
| slice | 对底层数组的一段视图，自己包含指针、长度、容量。 |
| len | 当前 slice 能访问的元素个数。 |
| cap | 从 slice 起点到底层数组末尾还能容纳的元素个数。 |
| append | 向 slice 追加元素，返回新的 slice header。必须接住返回值。 |
| 扩容 | cap 不够时，runtime 分配更大的底层数组，把旧元素拷贝过去。 |

## 30 秒短答

slice 底层不是直接存所有元素，而是一个 header，包含底层数组指针、长度和容量。slice 赋值或函数传参时，拷贝的是这个 header，不是底层数组。

`append` 后旧切片是否受影响，取决于容量是否足够。如果 `append` 没有超过 `cap`，新旧 slice 仍然共享同一个底层数组，后续修改可能互相影响；如果超过 `cap` 触发扩容，runtime 会分配新数组并拷贝旧数据，新 slice 指向新数组，旧 slice 仍指向旧数组，两者就分开了。

面试判断口诀：看到 `append`，先看 `cap` 够不够。

## 1-2 分钟标准回答

slice 本质上是一个很小的描述结构，通常可以理解为 `Data`、`Len`、`Cap` 三元组。`Data` 指向底层数组，`Len` 表示当前长度，`Cap` 表示容量。slice 变量本身是值类型，所以赋值、传参时复制的是这个 header，但 header 里的 `Data` 仍可能指向同一个底层数组。

`append` 的行为分两种。如果当前 slice 的容量足够，`append` 会直接复用原底层数组，只是返回一个新的 slice header，长度变大。这个时候新旧 slice 共享同一块底层数组，如果它们写到同一个位置，就会互相覆盖。很多代码输出题就是考这个。

如果容量不够，`append` 会触发扩容。Go runtime 会申请一个更大的底层数组，把旧数据拷贝过去，然后返回指向新数组的 slice。此时旧 slice 仍然指向旧数组，新 slice 指向新数组，两者后续修改互不影响。

工程上不要依赖 `append` 后旧 slice 的共享关系。如果想强制新 slice 和旧 slice 分离，可以用三下标切片限制容量，例如 `s[:len(s):len(s)]`，让下一次 `append` 必然扩容。

## 代码例子

```go
func main() {
	s := make([]int, 3, 5)
	s[0], s[1], s[2] = 1, 2, 3

	a := append(s, 4) // cap=5 足够，不扩容
	b := append(s, 5) // cap=5 足够，不扩容

	fmt.Println(a) // [1 2 3 5]，不是 [1 2 3 4]
	fmt.Println(b) // [1 2 3 5]
}
```

为什么 `a` 不是 `[1 2 3 4]`？

```text
a 和 b 都从 s append；
cap 足够，所以两次 append 都写同一个底层数组 index=3；
b 后执行，把 a 写进去的 4 覆盖成 5。
```

如果想让它们互不影响：

```go
s := make([]int, 3, 5)
s[0], s[1], s[2] = 1, 2, 3

base := s[:len(s):len(s)] // len=3, cap=3
a := append(base, 4)      // 必然扩容
b := append(base, 5)      // 必然扩容

fmt.Println(a) // [1 2 3 4]
fmt.Println(b) // [1 2 3 5]
```

## 原理拆解

### 1. slice header 是值拷贝

可以把 slice 想象成：

```go
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
```

真实代码里不要直接操作 `reflect.SliceHeader` 做转换，这里只是帮助理解。

当你写：

```go
s2 := s1
```

复制的是 `Data/Len/Cap` 这三个字段。`s1` 和 `s2` 的 header 是两份，但 `Data` 可能指向同一个底层数组。

### 2. append 必须接返回值

`append` 可能返回一个新的 header，甚至可能指向新的底层数组。

```go
s = append(s, x)
```

这是正确写法。不要只写：

```go
append(s, x) // 编译不通过：append 的结果没有使用
```

### 3. 扩容策略不用死背，但要知道趋势

Go 1.18 之后的扩容不再是简单的“小于 1024 翻倍，大于 1024 增长 25%”。大体规律是：

```text
小 slice：接近翻倍，减少频繁扩容。
大 slice：增长比例逐渐靠近 1.25 倍，减少内存浪费。
```

面试一般不要求背源码公式。更重要的是能说清楚：扩容会分配新数组并拷贝旧数据，因此可能带来内存分配和拷贝成本。

## 结合我的经历

待补充。

可以结合 Go 项目里的这些场景讲：

- 批量组装请求参数、订单列表、用户列表时，提前 `make([]T, 0, n)` 预分配容量，减少扩容。
- 对外返回 slice 前，如果不希望调用方继续共享内部底层数组，可以主动 copy。
- 在循环里大量 append 时，关注是否造成额外内存分配和 GC 压力。

## 常见追问

| 追问 | 回答要点 |
| --- | --- |
| 如何让 `append` 后的新旧 slice 互不影响？ | 用三下标切片 `s[:len(s):len(s)]` 把 cap 限制为 len，下一次 append 会触发扩容。也可以显式 `copy`。 |
| nil slice 和空 slice 有什么区别？ | nil slice 的指针是 nil，len/cap 都是 0；空 slice len/cap 也是 0，但指针可能指向 runtime 的零地址。两者 append 行为一致，但 JSON 序列化不同：nil 通常是 `null`，空 slice 是 `[]`。 |
| 为什么建议预分配 slice 容量？ | 减少扩容次数，降低内存分配、数据拷贝和 GC 压力。 |

## 易错点

- 以为 slice 赋值会复制底层数组。
- 忘记 `append` 可能复用底层数组，导致新旧 slice 互相覆盖。
- 忘记接住 `append` 的返回值。
- 过度死背扩容倍数，忽略真正要判断的是 `cap` 是否足够。

## 复习要点

- slice 是 header，不是数组本体。
- `append` 是否扩容，决定是否共享底层数组。
- 判断输出题先看 `len/cap`。
- 想隔离底层数组，用三下标切片或 `copy`。
