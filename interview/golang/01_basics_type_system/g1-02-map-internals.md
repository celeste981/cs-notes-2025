# G1-02 map 底层实现：哈希冲突如何解决？扩容是怎么做的？

分类：基础与类型系统

材料类型：interview question / knowledge topic

难度：L2/L3

优先级：P0

关键词：map、hmap、bmap、bucket、tophash、overflow bucket、渐进式扩容、并发安全

复习状态：已成稿

来源：https://lc100.pages.dev/go

## 问题

Go 的 `map` 底层是怎么实现的？哈希冲突如何解决？扩容是怎么做的？

这道题通常还会追问：

```text
map 为什么遍历无序？
map 为什么不能并发读写？
sync.Map 和 map + RWMutex 怎么选？
```

## 先讲人话

`map` 可以理解成很多个“桶”。每个 key 先算出一个 hash，再根据 hash 找到应该放进哪个桶。

一个桶里可以放多个 key/value。如果几个 key 算出来落到同一个桶，这就是哈希冲突。Go 的处理方式是：

```text
桶里先放最多 8 个键值对；
如果还放不下，就挂 overflow bucket；
桶太满或 overflow 太多时，再扩容整理。
```

扩容不是一次性搬完所有桶，因为那样某一次写入会突然很慢。Go 会在后续 map 操作中逐步搬迁旧桶，这叫渐进式扩容。

## 前置概念

| 概念 | 小白解释 |
| --- | --- |
| hash | 把 key 算成一个数字，方便快速定位。 |
| bucket | map 底层的桶，一个桶最多放 8 个 key/value。 |
| tophash | hash 的高 8 位，用来快速过滤不可能匹配的 key。 |
| overflow bucket | 当前桶放满后额外挂上的溢出桶。 |
| 负载因子 | 元素数量和桶数量的比例。太高说明桶太挤，查询效率会下降。 |
| 渐进式扩容 | 扩容时不一次性搬完，而是在后续读写删过程中分批迁移。 |

## 30 秒短答

Go 的 `map` 底层核心结构是 `runtime.hmap`，里面记录元素数量、桶数量指数 `B`、当前桶数组 `buckets`、扩容时的旧桶数组 `oldbuckets` 等信息。桶的数量是 `2^B`。

每个桶 `bmap` 最多存 8 个 key/value，并保存 8 个 `tophash` 用于快速判断。哈希冲突主要通过桶内 8 个槽位加 overflow bucket 解决。

扩容有两类：负载因子过高会触发翻倍扩容；overflow bucket 太多但元素并不多时，会触发等量扩容，用来整理稀疏桶和回收 overflow。扩容是渐进式的，后续 map 操作会顺手搬迁旧桶。

另外，普通 `map` 不是并发安全的，并发读写会触发 fatal error，不能靠 `recover` 兜底。

## 1-2 分钟标准回答

Go 的 `map` 是哈希表实现。底层有一个 `hmap` 结构，记录当前元素数量 `count`、桶数量的指数 `B`、桶数组指针 `buckets`，以及扩容过程中的旧桶数组 `oldbuckets`。桶数量是 `2^B`。

每个桶可以放 8 个 key/value。桶里会先保存 8 个 `tophash`，也就是 hash 值的高 8 位，用来快速过滤。真正比较时，先看 `tophash` 是否匹配，再比较 key 本身。key 和 value 在桶里是分开连续存放的，这样可以减少内存 padding。

哈希冲突通过桶内槽位和 overflow bucket 解决。如果多个 key 落到同一个桶，先放到这个桶的 8 个槽位里；如果放满了，就链接额外的 overflow bucket。实际情况下，每个桶有 8 个槽位，所以冲突链通常不会太长。

map 扩容有两种触发条件。第一种是负载因子太高，比如元素数量相对桶数量太多，会触发翻倍扩容。第二种是 overflow bucket 太多，但元素数量不一定很多，这时会触发等量扩容，目的不是增加容量，而是整理数据，减少 overflow 链。

扩容不是一次性搬迁所有桶，而是渐进式的。扩容期间 `oldbuckets` 指向旧桶数组，后续每次读、写、删 map 时，runtime 会顺带搬迁一部分旧桶。这样避免某一次操作因为全量 rehash 出现明显卡顿。

## 原理拆解

### 1. hmap 里最关键的字段

面试不需要逐字背结构体，但可以记住这几个字段：

| 字段 | 含义 |
| --- | --- |
| `count` | 当前 map 中 key/value 的数量。 |
| `B` | 桶数量的指数，桶数是 `2^B`。 |
| `buckets` | 当前桶数组。 |
| `oldbuckets` | 扩容时的旧桶数组，扩容完成后为 nil。 |

### 2. 每个 bucket 放 8 个键值对

一个 bucket 大致包含：

```text
tophash[8]
keys[8]
values[8]
overflow pointer
```

为什么不是 `key/value/key/value` 交替存？

```text
把 key 连续放、value 连续放，可以减少内存对齐导致的 padding 浪费。
```

### 3. 查询流程

```text
1. 对 key 计算 hash。
2. 用 hash 的低位选择 bucket。
3. 在 bucket 的 8 个 tophash 中快速过滤。
4. tophash 匹配后，再比较真实 key。
5. 当前 bucket 没找到，就沿 overflow bucket 继续找。
```

### 4. 扩容的两种情况

| 触发原因 | 扩容方式 | 目的 |
| --- | --- | --- |
| 负载因子过高 | 翻倍扩容 | 增加 bucket 数量，降低冲突概率。 |
| overflow bucket 太多 | 等量扩容 | 整理桶，减少溢出链，提高查询效率。 |

负载因子可以简单理解为：

```text
元素数量 / 桶数量
```

Go runtime 里常见阈值约为 6.5。面试时说“大约 6.5”即可，不需要把源码常量背成绝对规则。

### 5. 渐进式扩容

如果 map 有很多元素，一次性搬迁所有桶会导致某一次操作非常慢。

Go 的做法是：

```text
扩容时先保留 oldbuckets；
新写入和查询逐步使用新 buckets；
后续 map 操作顺手搬迁 1 到 2 个旧 bucket；
全部搬完后 oldbuckets 置空。
```

这能把一次大成本拆成很多次小成本。

## 扩展：常见的 map / 字典结构实现方式

这里要先区分两个概念：

```text
Go 里的 map：特指 Go runtime 实现的哈希表。
更泛化的 map / dictionary / associative array：指“根据 key 找 value”的抽象数据结构。
```

所以你说的数组、跳表、红黑树不算错，但它们不是 Go `map` 的底层实现，而是“字典结构”可以选择的不同实现方式。

| 实现方式 | 查找复杂度 | 是否有序 | 典型场景 |
| --- | --- | --- | --- |
| 哈希表 Hash Table | 平均 `O(1)`，极端冲突可能退化 | 无序 | Go `map`、Java `HashMap`、大多数内存 KV 查询。 |
| 有序数组 / 排序数组 | 二分查找 `O(log n)`，插入删除 `O(n)` | 有序 | 数据量小、构建后主要查询、不频繁修改。 |
| 直接寻址数组 Direct Address Table | `O(1)` | 按下标天然有序 | key 范围很小且连续，比如状态码、ASCII 字符表。空间可能浪费。 |
| 红黑树 / 平衡二叉搜索树 | `O(log n)` | 有序 | Java `TreeMap`、C++ `std::map`，需要按 key 排序、范围查询。 |
| 跳表 Skip List | 平均 `O(log n)` | 有序 | Redis Sorted Set 的一种核心结构、部分 KV 存储的内存索引。实现相对简单，适合范围查询。 |
| B 树 / B+ 树 | `O(log n)` | 有序 | 数据库和文件系统索引，适合磁盘/页式存储，减少随机 I/O。 |
| Trie / 前缀树 | 和 key 长度相关 | 按字符路径组织 | 字符串前缀匹配、路由匹配、词典检索。 |

面试里可以这样总结：

> 如果只是问 Go 的 `map`，答案就是哈希表，核心是 bucket、tophash、overflow bucket 和渐进式扩容。  
> 如果问“字典结构有哪些实现”，那就可以扩展说：哈希表适合无序快速查找；红黑树、跳表适合有序遍历和范围查询；数组适合 key 范围小或数据构建后少修改；B+ 树更偏数据库/磁盘索引场景。

容易说错的一点：

```text
不要说 Go map 是红黑树或跳表实现。
Go map 是哈希表；红黑树和跳表是其他有序 map 的常见实现。
```

## 结合我的经历

Marketing 项目里两种方案都有用到，选择逻辑比较清楚。

### 1. `sync.Map`：用于读多写少、按 key 独立更新的缓存

Marketing 里自己写的 `sync.Map` 主要有三处：

| 位置 | 用法 | 为什么适合 `sync.Map` |
| --- | --- | --- |
| `src/common/config/config.go` | `MarketingConfig.config *sync.Map`，`RefMarketingConfig` 按 `cid` 读取多区域业务配置；配置变更 listener 会重新初始化。 | 配置读取频繁，写入主要发生在启动和配置刷新；每个 `cid` 的 value 独立，读写逻辑简单。 |
| `src/common/biz/impl/product_biz_impl.go` | `FreeProductMap sync.Map`，`CheckFreeProduct` 先按 `productId` 查本地缓存，miss 时调用 Product RPC，再 `Store`。 | 典型读多写少缓存；key 之间没有复杂一致性约束，只需要单 key 的 `Load/Store`。 |
| `src/basic/group-center/internal/biz/impl/group_biz_impl.go` | `accessTimeCache sync.Map`，记录 `groupId -> time.Time`，对 group last access time 做 5 分钟写入节流。 | 用户组校验链路并发高；每个 group 独立更新，适合用 `sync.Map` 避免手写锁。 |

这类场景的共同点是：

```text
key 相对独立；
主要操作是 Load / Store；
不需要在锁内维护多个 key 的一致性；
读多写少或者写入分散。
```

面试可以这样讲：

> 在 Marketing 里，我会把 `sync.Map` 用在读多写少、key 独立的缓存场景。比如产品是否 free 的缓存 `FreeProductMap`，就是按 productId 先 Load，miss 后查 Product RPC 再 Store；用户组 last access time 的节流缓存也是按 groupId 独立更新。这种场景不需要复杂事务语义，用 `sync.Map` 可以减少手写锁代码，也避免普通 map 并发读写的风险。

### 2. `map + RWMutex`：用于注册表、需要遍历或组合操作的 map

Marketing 里 `map + sync.RWMutex` 主要用于 handler 注册表：

| 位置 | 用法 | 为什么适合 `map + RWMutex` |
| --- | --- | --- |
| `src/engine/internal/handler/handler.go` | 全局 `handlerMap map[string]Handler`，`InjectHandler` 写入，`RefHandler` 和 `GetHandlerList` 读取。 | 需要保护普通 map，同时 `GetHandlerList` 要遍历整个 map；用 RWMutex 可以让查询走读锁、注册走写锁。 |
| `src/engine/internal/handler/pnar/pn_ar_handler.go` | 全局 `pnArHandlerMap map[string]PnArHandler`，注入和读取 PN/AR handler。 | 同样是注册表模式，写入集中在初始化，运行期主要读取；RWMutex 表达更直接。 |

这类场景的共同点是：

```text
map 是一个整体注册表；
需要遍历、返回列表或维护整体一致性；
写入通常发生在初始化或插件注册阶段；
希望类型保持明确，不想频繁做 interface{} 类型断言。
```

面试可以这样讲：

> 对于 handler 注册表，我更倾向 `map + RWMutex`。Marketing 的 `handlerMap` 和 `pnArHandlerMap` 都是这个模式：注册时用写锁，查 handler 和列出 handler 时用读锁。这里不只是单 key 缓存，还会遍历整个 map 生成 Admin 展示列表，所以用普通 map 加读写锁更自然，类型也更清晰。

### 3. 另一个相关点：只需要串行化临界区时用 `Mutex`

`src/common/template/placeholder_handler.go` 里有一个 `sync.Mutex`，不是为了保护 map，而是为了保护共享的 `baseTemp.Parse/Execute` 流程。注释里也写了 `baseTemp.Parse` 并发场景需要优化。

这说明锁的选择不是“有 map 就选某个固定方案”，而是看临界区到底保护什么：

```text
单 key 并发缓存：sync.Map。
注册表 / 需要遍历 / 类型明确：map + RWMutex。
共享对象的串行操作：Mutex。
整份配置快照读多写少：项目里更多用 LocalAtomicValueHolder 做不可变快照替换。
```

### 4. 项目里更常见的缓存模式：不可变快照 + atomic holder

Marketing 很多业务缓存没有用 `sync.Map` 或 `map + RWMutex`，而是用 `LocalAtomicValueHolder` 保存一整份不可变缓存对象，例如：

- 用户组本地 Roaring Bitmap 缓存：`memoryGroupListCache`
- Reward 配置缓存：`rewardCache`
- Cross Anti-Harassment 配置缓存：`crossAntiHarassmentCache`
- Plan、Task、Asset、Activity 等配置缓存

这种模式适合“定时整批 reload，运行期高频只读”的场景：

```text
reload 时构建一份新 map；
构建完成后一次性替换缓存指针；
读请求只读当前快照，不在原 map 上并发写。
```

面试可以补一句：

> 如果是整份配置定时刷新，我不会优先用 `sync.Map` 做逐个 key 更新。Marketing 里更多是构建 immutable cache item，然后用 atomic holder 一次性替换。这样读路径更轻，也避免读到半更新状态。

## 常见追问

| 追问 | 回答要点 |
| --- | --- |
| map 遍历为什么无序？ | Go 故意让 map range 的遍历顺序不稳定，避免开发者依赖顺序。需要有序就取出 key 排序。 |
| map 并发读写会怎样？ | 普通 map 并发读写会触发 fatal error，不是普通 panic，不能依赖 recover。需要加锁、用 channel 串行化，或使用 `sync.Map`。 |
| sync.Map 适合什么场景？ | 适合读多写少、key 相对稳定、缓存类场景。写多、需要复杂复合操作、需要稳定遍历时，`map + sync.RWMutex` 通常更直接。 |
| 删除 map 中不存在的 key 会怎样？ | 安全，无效果。读取不存在的 key 会返回 value 类型零值，可以用 `v, ok := m[k]` 区分。 |

## 易错点

- 把 `map` 当成并发安全结构使用。
- 以为 `recover` 可以兜住并发读写 map 的 fatal error。
- 以为 map 遍历顺序只是“碰巧无序”，实际 Go 设计上就不保证顺序。
- 只知道 map 会扩容，但说不清 overflow bucket 太多也会触发等量扩容。

## 复习要点

- map 是哈希表，桶数量是 `2^B`。
- 一个 bucket 放 8 个 key/value，冲突多了挂 overflow bucket。
- 扩容分翻倍扩容和等量扩容。
- 扩容是渐进式搬迁。
- 普通 map 不能并发读写。
