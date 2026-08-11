# G1-03 defer 的执行顺序和陷阱：defer、return、函数返回值的执行时序

分类：基础与类型系统

材料类型：interview question / knowledge topic

难度：L2

优先级：P0

关键词：defer、return、命名返回值、匿名返回值、LIFO、参数求值、闭包

复习状态：已成稿

来源：https://lc100.pages.dev/go

## 问题

`defer` 的执行顺序是什么？`defer`、`return` 和函数返回值之间是什么执行时序？

这道题常见考法是判断函数返回值，尤其是命名返回值和匿名返回值的区别。

## 先讲人话

`defer` 的意思是：这件事先登记下来，等函数快返回时再做。

但有两个细节非常重要：

```text
defer 的函数调用顺序：后登记的先执行。
return 不是一步完成：先给返回值赋值，再执行 defer，最后真正返回。
```

所以 `defer` 能不能改返回值，取决于它改的是不是“真正要返回的那个变量”。

## 前置概念

| 概念 | 小白解释 |
| --- | --- |
| defer | 延迟执行，函数返回前执行。常用于释放资源、解锁、关闭文件。 |
| LIFO | Last In First Out，后进先出。多个 defer 时，最后注册的最先执行。 |
| 命名返回值 | 函数签名里给返回值起了名字，比如 `func f() (result int)`。 |
| 匿名返回值 | 函数签名里只写类型，比如 `func f() int`。 |
| 闭包 | 函数引用了外层变量，执行时可能读到变量的最新值。 |

## 30 秒短答

`defer` 在函数返回前执行，多个 `defer` 按 LIFO 顺序执行。关键时序是：

```text
return 赋值 -> defer 执行 -> 函数真正返回
```

`return` 不是原子操作。对于命名返回值，`return` 会先把值赋给命名返回变量，然后执行 `defer`，所以 `defer` 可以修改最终返回值。对于匿名返回值，`return` 时通常已经把局部变量的值拷贝到返回位置，`defer` 再改局部变量，不会影响最终返回结果。

另外，`defer` 的参数在注册时立即求值；如果 `defer` 里是闭包引用外部变量，则执行时读到的是变量最新值。

## 1-2 分钟标准回答

`defer` 用来延迟执行函数，通常在资源释放场景使用，比如 `defer file.Close()`、`defer mu.Unlock()`。它的执行时机是当前函数即将返回之前。

多个 `defer` 的执行顺序是后进先出，也就是最后注册的最先执行。这一点类似栈。

最容易出错的是 `defer` 和 `return` 的关系。Go 里的 `return` 不是一步完成的，可以拆成三步理解：第一步，把 return 后面的表达式赋给返回值；第二步，执行 defer；第三步，函数真正返回。

如果函数使用命名返回值，返回变量本身就在函数作用域里，`defer` 闭包可以修改它，因此会影响最终返回结果。如果函数使用匿名返回值，`return result` 时会先把局部变量 `result` 的值拷贝到返回位置，之后 `defer` 修改局部变量，不会再影响返回值。

此外，`defer` 调用的参数是在注册 defer 时求值的，不是执行 defer 时才求值。但如果 defer 注册的是闭包，闭包内部引用外层变量，那执行时会读取变量的最新值。

## 代码例子

```go
func f() (result int) {
	defer func() {
		result++ // 修改命名返回值
	}()
	return 0 // 先赋值 result=0，再执行 defer
}
// f() 返回 1

func g() int {
	result := 0
	defer func() {
		result++ // 修改局部变量，不影响返回值
	}()
	return result // 返回值已拷贝
}
// g() 返回 0
```

解释：

```text
f 使用命名返回值，defer 修改的是最终要返回的 result。
g 使用匿名返回值，return result 时已经把 0 拷贝到返回位置，defer 改的是局部变量。
```

再看参数立即求值：

```go
func main() {
	i := 1
	defer fmt.Println(i) // 注册 defer 时，参数 i 已经求值为 1
	i = 2
}
// 输出 1
```

闭包则不同：

```go
func main() {
	i := 1
	defer func() {
		fmt.Println(i) // 执行 defer 时读取 i 的最新值
	}()
	i = 2
}
// 输出 2
```

## 原理拆解

### 1. 多个 defer 的顺序

```go
defer fmt.Println("A")
defer fmt.Println("B")
defer fmt.Println("C")
```

输出顺序：

```text
C
B
A
```

因为 defer 是栈结构，后注册先执行。

### 2. return 的三步

把：

```go
return x
```

理解成：

```text
1. 返回值 = x
2. 执行 defer
3. 真正返回
```

这就是命名返回值能被 defer 修改的原因。

### 3. defer 的常见用途

| 场景 | 示例 |
| --- | --- |
| 释放文件 | `defer f.Close()` |
| 解锁 | `mu.Lock(); defer mu.Unlock()` |
| 恢复 panic | `defer func() { recover() }()` |
| 记录耗时 | 函数入口记录 start，defer 里打耗时日志 |

### 4. defer 在循环里的坑

```go
for _, name := range files {
	f, _ := os.Open(name)
	defer f.Close()
}
```

问题：`defer` 要等整个函数返回才执行，不是每次循环结束执行。如果循环很多次，文件会一直不关闭，可能导致资源耗尽。

更稳的写法：

```go
for _, name := range files {
	func() {
		f, _ := os.Open(name)
		defer f.Close()
		// 使用 f
	}()
}
```

或者在循环体里手动关闭。

## 结合我的经历

待补充。

可以结合 Go 服务里的这些场景讲：

- 操作数据库事务时，使用 `defer` 做 rollback 兜底，但 commit 成功后要避免误回滚。
- 使用锁时，把 `Unlock` 紧跟在 `Lock` 后面 defer，降低忘记释放锁的风险。
- 文件、网络连接、HTTP response body 用完后 `defer Close()`，但在大循环里要谨慎。

## 常见追问

| 追问 | 回答要点 |
| --- | --- |
| defer 在 for 循环里有什么问题？ | defer 在函数返回时才执行，大量循环会堆积资源释放。可以用匿名函数缩小作用域，或手动释放。 |
| Go 1.14 对 defer 做了什么优化？ | 引入 open-coded defer，让部分 defer 在编译期展开，降低运行时分配和调用开销。面试重点是知道 defer 已经比早期版本快很多，但热点循环里仍要谨慎。 |
| recover 必须放在哪里才生效？ | 必须在 defer 函数里直接调用，才能捕获当前 goroutine 的 panic。不能跨 goroutine recover。 |

## 易错点

- 以为 `return` 是一步完成，忽略 defer 可以插在中间改命名返回值。
- 分不清 defer 参数立即求值和闭包延迟读取变量。
- 在循环中大量 defer，导致资源释放延迟。
- 以为 recover 能跨 goroutine 捕获 panic。

## 复习要点

- 多个 defer：后进先出。
- return 时序：赋值 -> defer -> 返回。
- 参数立即求值，闭包执行时取最新变量。
- 命名返回值可以被 defer 修改，匿名返回值通常不行。
