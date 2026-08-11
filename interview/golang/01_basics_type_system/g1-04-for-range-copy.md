# G1-04 for range 的复制陷阱：遍历时修改元素为什么不生效？

分类：基础与类型系统

材料类型：interview question / knowledge topic

难度：L2

优先级：P0

关键词：for range、值拷贝、循环变量、闭包捕获、Go 1.22、slice、map

复习状态：已成稿

来源：https://lc100.pages.dev/go

## 问题

`for range` 遍历时修改元素为什么不生效？闭包里捕获循环变量为什么容易出错？

这道题本质考察：

```text
range 拿到的 value 通常是元素拷贝；
Go 1.22 之前循环变量复用；
遍历开始时部分信息已经确定，比如 slice 的遍历长度。
```

## 先讲人话

`for range` 里的第二个变量，很多时候不是原集合里的元素本体，而是“拿出来的一份复印件”。

比如：

```go
for _, v := range users {
	v.Name = "new"
}
```

你改的是 `v` 这份拷贝，不是 `users` 里的原元素。所以修改不生效。

如果要改原集合，通常要用下标：

```go
for i := range users {
	users[i].Name = "new"
}
```

## 前置概念

| 概念 | 小白解释 |
| --- | --- |
| range | Go 的遍历语法，可以遍历 slice、array、map、string、channel 等。 |
| 值拷贝 | 复制一份值。修改副本不会影响原值。 |
| 指针 | 保存变量地址。通过指针可以修改原对象。 |
| 闭包捕获 | 内部函数引用了外层变量。执行内部函数时可能读到外层变量当前值。 |
| Go 1.22 循环变量变化 | Go 1.22 起，循环变量按每次迭代创建新变量，修复了很多闭包捕获坑。 |

## 30 秒短答

`for range` 遍历 slice 或 array 时，第二个变量是元素的值拷贝，不是原元素引用。你修改这个变量，只是修改副本，不会影响原集合。要修改原集合，应使用下标访问，比如 `s[i] = ...` 或 `s[i].Field = ...`。

Go 1.22 之前，`for range` 的循环变量在整个循环中是同一个变量，每次迭代只是重新赋值，所以闭包或 goroutine 捕获这个变量时，容易全部读到最后一次的值。Go 1.22 起改成每次迭代创建独立变量，这个经典坑已经被修复，但面试里仍然要能解释历史原因和旧代码风险。

另外，range slice 时，遍历开始就确定了长度，循环中 append 新元素通常不会被本次 range 遍历到。

## 1-2 分钟标准回答

`for range` 的 value 变量通常是值拷贝。比如遍历 `[]User` 时，`for _, u := range users` 里的 `u` 是 `users` 中每个元素的副本。修改 `u.Name` 只是在改副本，不会写回原 slice。要修改原元素，应该用 index，例如 `for i := range users { users[i].Name = ... }`。

第二个坑是循环变量捕获。在 Go 1.22 之前，range 循环变量在整个循环中复用同一个变量，每轮只是重新赋值。如果在循环里创建闭包或 goroutine，并在里面引用这个变量，等闭包真正执行时，读到的可能都是最后一次迭代的值。经典修复方式是在循环内部重新声明一个变量，或者把变量作为参数传入闭包。

Go 1.22 改了语义，每次迭代都会创建新的循环变量，因此这个闭包捕获问题在新版本中被修复。但面试里仍然建议说明版本差异，因为很多旧项目或老面试题还会考 Go 1.21 之前的行为。

遍历 map 时，value 也是拷贝，而且 map 遍历顺序不保证稳定。遍历时删除 key 是安全的，但新增 key 是否会在本次遍历中出现不确定。

## 代码例子

### 修改 slice 元素不生效

```go
type User struct {
	Name string
}

users := []User{{Name: "old"}}

for _, u := range users {
	u.Name = "new" // 改的是副本
}

fmt.Println(users[0].Name) // old
```

正确写法：

```go
for i := range users {
	users[i].Name = "new"
}

fmt.Println(users[0].Name) // new
```

### Go 1.21 及之前的闭包捕获坑

```go
funcs := make([]func(), 3)

for i := 0; i < 3; i++ {
	funcs[i] = func() {
		fmt.Println(i)
	}
}

for _, f := range funcs {
	f()
}
```

Go 1.21 及之前常见输出：

```text
3
3
3
```

原因：闭包捕获的是同一个 `i` 变量，循环结束后 `i == 3`。

旧版本修复方式：

```go
for i := 0; i < 3; i++ {
	i := i // 每轮创建独立副本
	funcs[i] = func() {
		fmt.Println(i)
	}
}
```

Go 1.22 之后，循环变量默认按每次迭代创建新变量，很多场景不再需要手动 `i := i`。

## 原理拆解

### 1. 为什么修改 value 不生效

这段代码：

```go
for _, v := range s {
	v = 100
}
```

可以粗略理解成：

```text
每次从 s 里取一个元素，复制给变量 v；
你修改的是变量 v；
s 里的原元素没有被改。
```

所以要改原集合，应写：

```go
for i := range s {
	s[i] = 100
}
```

### 2. range slice 时长度什么时候确定

遍历 slice 时，range 会在开始时确定要遍历的长度。循环体里 append 新元素，一般不会增加本次遍历次数。

```go
s := []int{1, 2, 3}
for _, v := range s {
	s = append(s, v)
}

fmt.Println(s) // 长度变成 6，但循环只执行 3 次
```

### 3. range map 的特殊点

map 的 range 有几个特点：

```text
遍历顺序不保证。
value 是拷贝。
遍历时 delete 当前或未访问 key 是安全的。
遍历时新增 key 是否会被遍历到，不确定。
```

如果需要稳定顺序：

```go
keys := make([]string, 0, len(m))
for k := range m {
	keys = append(keys, k)
}
sort.Strings(keys)

for _, k := range keys {
	fmt.Println(k, m[k])
}
```

## 结合我的经历

待补充。

可以结合 Go 项目里的这些场景讲：

- 批量修改结构体 slice 时，优先用下标遍历，避免改副本。
- 循环里启动 goroutine 时，显式传参或确认 Go 版本，避免捕获错误变量。
- 遍历 map 输出配置或日志时，如果需要稳定顺序，先排序 key。

## 常见追问

| 追问 | 回答要点 |
| --- | --- |
| for range slice 时，循环里 append 新元素会遍历到吗？ | 通常不会。range 开始时已经确定遍历长度，append 不会增加本轮 range 次数。 |
| 遍历 `[]*User` 时，`u.Name = "new"` 会生效吗？ | 会。value 拷贝的是指针，指针副本仍指向同一个 User 对象。注意如果给 `u = &User{}` 重新赋值，不会改 slice 里的指针。 |
| Go 1.22 对循环变量做了什么改变？ | 每次迭代创建独立循环变量，闭包捕获时不再都指向同一个变量。旧版本仍需手动拷贝或传参。 |
| range map 时能删除或新增 key 吗？ | 删除安全；新增 key 是否会在本次遍历中出现不确定。不应依赖这个行为。 |

## 易错点

- 以为 `for _, v := range s` 的 `v` 是原元素。
- 遍历结构体 slice 时修改 `v.Field`，以为会写回原 slice。
- 循环里启动 goroutine，忘记循环变量捕获问题。
- 依赖 map 的遍历顺序。
- 以为 append 后 range 会自动遍历新增元素。

## 复习要点

- range 的 value 通常是值拷贝。
- 修改原 slice 元素，用 index。
- Go 1.22 修复了循环变量捕获坑，但旧题仍要会讲。
- range slice 的长度在开始时确定。
- map range 顺序不保证，新增 key 的遍历结果不确定。
