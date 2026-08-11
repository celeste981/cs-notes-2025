# 项目故事：Marketing 容器 OOMKilled 排障

> 可视化图解：[OOMKilled 排障图解](./oom-killed-troubleshooting-star.html)
>
> 所属项目索引：[Marketing 项目材料](./README.md)

## 一句话定位

这是一个 **容器内存超限导致服务实例被系统杀掉** 的排障案例。重点不是业务逻辑 bug，而是通过 CMDB、STDOUT/STDERR、termination signal 和内存 limit 证据，把 `exit code 137` 定性为 OOMKilled，并给出恢复和后续定位路径。

面试开场可以这样说：

> 我遇到过一次 Marketing 服务在 UAT 环境容器反复异常退出的问题。页面只显示 `Container exited with status 137`，我通过 CMDB 的 task 状态、STDOUT/STDERR 和 cgroup OOM 事件确认是内存超过 2GiB limit 后被系统 `SIGKILL`，然后把业务 panic、启动 smoke test 报错和真正 OOM 证据区分开，给出临时扩容和后续排查方向。

## STAR 故事

### Situation

UAT 环境的 Marketing SDU 出现容器异常退出，CMDB 上能看到类似：

- `OOMKilled - Container exited with status 137`
- `Running Instances: 0 of 1 Instances`
- task 状态在 `STARTING` / 重启过程中
- 部署内存为 `2048 MiB`

同时应用日志里还有其他错误，例如 SeaTalk webhook panic、启动 smoke test 的 gRPC check retry，这些日志容易误导判断。

### Task

我的目标是快速回答三个问题：

1. 容器是不是因为业务 panic 崩了？
2. `status 137` 是否确实是 OOMKilled？
3. 如果是 OOM，短期怎么恢复，后续应该从哪里查根因？

### Action

我按排障链路做了几步：

1. 先看 CMDB task 的状态和资源配置，确认 SDU 当前是 `0 of 1 Instances`，memory limit 是 `2048 MiB`。
2. 查 `View STDOUT / STDERR` 和 termination 相关日志，看到关键证据：
   - `exit code=137`
   - `caused by OOM kill`
   - `memory.events.oom_kill=1`
   - `usage=2150584320`
   - `limit=2147483648`
3. 对比 `usage` 和 `limit`，确认进程在使用约 2.15GB 时超过 2GiB limit，被系统 `SIGKILL`。
4. 排除干扰项：
   - SeaTalk webhook panic 是 reliable-event 里被 recover 后记录的错误，不等价于容器退出。
   - smoke test 的 gRPC check retry 后面有 passed，不是 OOM 的直接原因。
   - code-lens profiling report 只是 OOM 前最后在输出的日志，不能单独判定为根因。
5. 给出处理建议：
   - 短期把该 UAT SDU memory 从 2GiB 提到 3GiB 或 4GiB，让实例先恢复。
   - 回看 OOM 前 5 到 15 分钟的 Grafana 内存曲线和业务日志，判断是瞬时峰值还是持续上涨。
   - 重点检查大批量任务、consumer、scheduler、大查询、一次性构造大 slice/map、缓存增长和 goroutine 并发。

### Result

已确认这次退出的直接原因是 **OOMKilled**，不是普通业务 panic：

```text
exit code=137
memory.events.oom_kill=1
usage=2150584320
limit=2147483648
```

实际扩容结果、最终根因和是否有代码修复：`待补充`。

## 30 秒口述版

这次问题表面上是容器 `status 137`，同时日志里还有业务 panic 和 smoke test retry，容易混淆。我先从 CMDB task 和 STDOUT/STDERR 找直接退出证据，看到 `exit code=137`、`memory.events.oom_kill=1`，并且内存 `usage` 超过 2GiB `limit`，所以定性为 OOMKilled。然后我把 SeaTalk panic 和 smoke test retry 作为伴随现象排除，建议先临时扩容恢复，再回看 OOM 前的内存曲线、批任务和 consumer 负载定位根因。

## 1 分钟口述版

我排查过一次 Marketing UAT 容器异常退出。最开始页面上只看到 `Container exited with status 137`，同时日志里有 SeaTalk webhook panic 和启动 smoke test 的 gRPC check retry。我的判断方式是先不被业务错误带偏，而是找容器退出的直接证据。

我进入 CMDB 的 task 页面，确认当前实例数是 `0 of 1`，memory limit 是 2GiB。然后看 STDOUT/STDERR 和 termination 相关日志，里面有 `exit code=137`、`caused by OOM kill`、`memory.events.oom_kill=1`，并且 `usage=2150584320` 已经超过 `limit=2147483648`。这就能明确说明是进程超过容器内存上限后被系统 `SIGKILL`，不是普通 panic。

后续我把其他日志做了区分：SeaTalk panic 是 reliable-event recover 后记录的错误，smoke test retry 后面也通过了，都不是直接退出原因。处理建议是先把 UAT 这个 SDU 从 2GiB 临时扩到 3 或 4GiB 恢复服务，再查 OOM 前 5 到 15 分钟的内存曲线、consumer、scheduler、大查询和最近发版 diff，判断是瞬时峰值还是内存持续增长。

## 技术深挖点

- 架构设计：Marketing 是事件驱动和定时调度结合的服务，consumer、scheduler、reliable-event 都可能带来批量内存峰值。
- 难点取舍：排障时要先找容器退出的直接证据，不能只看到业务 panic 就下结论。
- 性能或稳定性：OOM 需要区分瞬时峰值、持续泄漏和部署 limit 过低。
- 数据一致性：本案例没有直接涉及数据一致性；如果 OOM 发生在批任务中，需要确认任务是否可重试、是否幂等、是否会重复发送或重复发券。
- 可观测性：关键证据来自 CMDB task、STDOUT/STDERR、termination log、Grafana 内存曲线和业务日志。

## 可量化结果

- 定位到容器 memory limit：`2147483648 bytes`，约 2GiB。
- 定位到 OOM 时内存 usage：`2150584320 bytes`，略超 2GiB limit。
- 找到强证据：`memory.events.oom_kill=1`。
- 最终恢复耗时、扩容后稳定时长、根因修复收益：`待补充`。

## 面试官可能追问

| 追问 | 回答要点 | 风险 |
| --- | --- | --- |
| 为什么 `137` 可以判断是 OOM？ | `137 = 128 + 9`，对应 `SIGKILL`；但单靠 137 不够，真正强证据是 `OOMKilled` / `memory.events.oom_kill=1` / `usage > limit`。 | 不要说所有 137 都一定是 OOM，SIGKILL 也可能来自人为 kill 或平台操作。 |
| panic 和 OOM 有什么区别？ | panic 是 Go 程序内部异常，可能被 recover；OOM 是进程超过内存限制后被系统杀掉。这个案例里 reliable-event panic 被 recover，不是直接退出原因。 | 如果 panic 高频刷屏，也可能间接推高内存或日志压力，不能完全忽略。 |
| 为什么不直接扩容就完了？ | 扩容是恢复手段，不是根因分析。后续还要看内存曲线、批任务、consumer、scheduler、大对象和最近 diff。 | 需要补充最终根因，否则只能说定位了直接原因。 |
| 你会怎么查 Go 服务内存问题？ | 看 pprof / heap profile、goroutine 数、Top alloc、对象保留路径；结合业务日志判断是否大 slice/map、缓存无界增长、JSON/CSV 全量构造或并发过高。 | 本次是否拿到 pprof：待补充，不能假装已经做过。 |
| 如果是批处理导致的，怎么优化？ | 分页、流式处理、限制 batch size 和并发、避免一次性加载全量、及时释放大对象、保证任务幂等和可重试。 | 需要结合具体 handler 或 task 代码，不能泛泛而谈。 |
| 如果发生在线上，处理顺序是什么？ | 先恢复实例和流量，必要时扩容或回滚；再保留日志和 profile 证据；最后做根因修复和回归。 | 要避免在 live 上做破坏性操作。 |

## 真实性检查

- 是否有事实依据：有。证据来自 CMDB / STDOUT / STDERR 中的 `exit code=137`、`memory.events.oom_kill=1`、`usage` 和 `limit`。
- 是否能解释细节：能解释 exit code 137、SIGKILL、OOMKilled、usage/limit 对比，以及为什么排除业务 panic 作为直接原因。
- 是否容易被追问穿：最终业务根因和恢复结果还缺失，面试时不要说成“我完成了根因修复”，更适合表述为“我完成了直接原因定位，并给出恢复和后续排查路径”。

## 复习要点

- `exit code 137`：通常表示进程收到 `SIGKILL`，OOMKilled 是常见原因。
- `OOMKilled` 强证据：平台状态、`memory.events.oom_kill=1`、`usage > limit`。
- 容器内存 limit：超过 cgroup limit 后，即使宿主机还有内存，容器进程也可能被杀。
- Go 排查方向：heap profile、goroutine、无界缓存、大对象、批量处理、JSON 序列化、日志放大。
- 面试表达边界：区分“直接退出原因已确认”和“业务根因已修复”。
