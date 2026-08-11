# Marketing Engine 架构梳理

> 可视化图解：[Marketing Engine 架构图解](./marketing-engine-architecture.html)
>
> 所属项目索引：[Marketing 项目材料](./README.md)

## 专项补充文档

如果只想背深挖点，可以直接看这三篇：

- [Marketing Engine 核心知识深挖](./marketing-engine-key-knowledge-deep-dive.md)：Plan、Trigger、Processor、Handler、RuleCenter 和服务边界。
- [Marketing Processor / Handler / RuleCenter 深挖](./marketing-processor-handler-rulecenter-deep-dive.md)：筛选链、执行动作、RuleCenter 边界、重复触达防守。
- [Marketing 稳定性和排障深挖](./marketing-reliability-troubleshooting-deep-dive.md)：OOMKilled、大批量人群、Kafka、Handler 失败和配置风险。
- [Marketing PlanDispatch 生命周期深挖](./marketing-plan-dispatch-lifecycle-deep-dive.md)：PlanRawMessage、PlanDispatch、Trigger 分支、DoExecuteHandler 和异常记录。
- [Marketing User Group Processor 链深挖](./marketing-user-group-processor-chain-deep-dive.md)：Processor 注册、三类 Processor、ConditionType、Delay 和筛选优化。
- [Marketing 用户名单 Redis / Roaring Bitmap 深挖](./marketing-user-list-redis-roaring-bitmap-deep-dive.md)：用户名单缓存、Redis 二级 Roaring Bitmap、LocalCache、RemoteLocalCache 和缓存一致性。
- [Marketing Handler 执行和风险控制深挖](./marketing-handler-risk-control-deep-dive.md)：Handler 接口、Notify、Voucher、PNAR、HandlerRules 和资损/骚扰风险。
- [Marketing Consumer / Batch History 可靠性深挖](./marketing-consumer-batch-history-reliability-deep-dive.md)：consumer 注册、topic 并发限制、batch history 对账和异常状态。
- [Marketing 技术优化点和深挖清单](./marketing-technical-optimization-points.md)：配置化、并发、批次对账、防重复、可观测性和配置风险。

## 一、Engine 是什么？一句话定位

**Marketing Engine 是一个可配置的运营自动化引擎**——通过 Admin 页面配置 Plan（计划），定义"什么时候、对哪些用户、做什么动作"，系统自动执行，无需开发介入。

面试开场：

> 我负责的 Marketing Engine 本质上是一个 **事件驱动 + 定时调度的规则引擎**，核心抽象是 Plan。每个 Plan 定义了触发条件（Trigger）、目标用户筛选（User Group Processor）、和执行动作（Handler），形成一条完整的 `触发 → 筛选 → 执行` 流水线。

---

## 项目路径速查

| 项目 | 本机路径 | 这里看什么 |
|---|---|---|
| `insurance-marketing` | `/Users/si.chen/GolandProjects/insurance-marketing` | Marketing 主服务，Engine、RuleCenter、Group Center、Notify/Reward/Repo 都在这里。 |
| `insurance-marketing-data` | `/Users/si.chen/GolandProjects/insurance-marketing-data` | 定时/离线场景下批量拉用户、查 ES/S3/CSV、写离线执行历史、推 Kafka。 |
| `insurance-marketing-api` | `/Users/si.chen/GolandProjects/insurance-marketing-api` | Marketing 主服务 API/proto。 |
| `insurance-marketing-data-api` | `/Users/si.chen/GolandProjects/insurance-marketing-data-api` | Marketing-Data API/proto。 |
| `insurance-operator-bff` | `/Users/si.chen/GolandProjects/insurance-operator-bff` | O-BFF / Admin 聚合入口；涉及后台页面、审批、批量、proxy 时优先确认这里。 |

常用代码入口：

| 想看什么 | 推荐路径 |
|---|---|
| Engine 统一入口 `PlanDispatch` | `/Users/si.chen/GolandProjects/insurance-marketing/src/engine/internal/manager/impl/plan_manager_impl.go` |
| PlanManager 接口 | `/Users/si.chen/GolandProjects/insurance-marketing/src/engine/internal/manager/plan_manager.go` |
| User Group Processor 链式匹配 | `/Users/si.chen/GolandProjects/insurance-marketing/src/engine/internal/biz/impl/user_group_biz_impl.go` |
| Processor 三类接口 | `/Users/si.chen/GolandProjects/insurance-marketing/src/engine/internal/processor/user_group/` |
| Handler 插件 | `/Users/si.chen/GolandProjects/insurance-marketing/src/engine/internal/handler/` |
| RuleCenter | `/Users/si.chen/GolandProjects/insurance-marketing/src/basic/rule-center/` |
| Kafka consumer | `/Users/si.chen/GolandProjects/insurance-marketing/src/engine/consumer/inner/default_consumer.go` |
| Marketing-Data 离线拉数编排 | `/Users/si.chen/GolandProjects/insurance-marketing-data/src/manager/impl/offline_data_fetch_manager_impl.go` |

---

## 二、核心模型：Plan

Plan 是引擎的最小执行单元，数据库里一行记录，Admin 上一个可视化配置项：

```
Plan
├── TriggerType      → 怎么触发（定时=1 / MQ=2 / 实时事件=3）
├── TopicName        → Kafka topic（MQ/Event 模式下的消息来源）
├── SchedulerCron    → Cron 表达式（定时模式下的调度频率）
├── UserGroupId      → 关联的用户分组
│   └── UserGroupProcessorRelation[]  → 一组 Processor 链
│       ├── ProcessorCode      → 处理器代码
│       ├── ConditionType      → AND/OR/Delay 组合逻辑
│       └── ProcessorParams    → JSON 参数（Admin 页面上填的那个）
├── HandlerName      → 执行什么动作（发通知/发券/数据清洗...）
├── HandlerParams    → 动作参数
├── HandlerRules     → 前置正则规则校验
└── TriggerRuleParam → 触发层的 RuleCenter 校验
```

---

## 三、完整数据流：三阶段流水线

```
┌─────────────┐    ┌─────────────────────┐    ┌──────────────┐
│  Stage 1     │    │  Stage 2             │    │  Stage 3      │
│  触发 Trigger │───▶│  筛选 User Group     │───▶│  执行 Handler  │
└─────────────┘    └─────────────────────┘    └──────────────┘
```

### Stage 1: 触发（三种入口殊途同归）

```
定时触发 (TriggerType=1)
    Scheduler → Marketing-Data 批量拉用户 → 推送到 Kafka topic
    → Consumer 消费 → PlanManager.PlanDispatch()

MQ 触发 (TriggerType=2)
    Marketing-Data 推数据到 Kafka topic
    → Consumer 消费 → PlanManager.PlanDispatch()

实时事件 (TriggerType=3)
    C-BFF ReportEvent RPC → PublishUserEvent → 内部 Kafka topic
    → Consumer 消费 → PlanManager.PlanDispatch()
```

**核心设计**：三种触发方式最终都汇聚到 `PlanManager.PlanDispatch()`，这是引擎的统一入口。

### Stage 2: 用户筛选（Processor 链式执行）

User Group 下挂载一组 **Processor**，按 ConditionType 组合：

```
UserGroupProcessorRelation[]
    │
    ├── ConditionType=1 (UserList)     → SchedulePullUserListProcessorI
    │       定义「从哪里拉用户」，返回 ES 查询条件
    │       例：SchedulePullUserListByPolicyInfo（按 partner_id + 保单状态拉用户）
    │
    ├── ConditionType=2 (RealtimeAnd)  → EventCheckProcessorI / AdditionalCheckProcessorI
    │       实时校验，所有条件必须全部通过
    │       例：CheckUserWithPolicyStatusChangeInfo（保单状态变更事件匹配）
    │
    ├── ConditionType=3 (RealtimeOr)   → 同上
    │       至少一个条件通过即可
    │
    ├── ConditionType=4 (DelayAnd)     → 延迟匹配
    │       发布延迟事件，等时间窗口到了再重新匹配
    │
    └── ConditionType=5 (DelayOr)      → 同上
```

**Processor 的三大族**：

| 族 | 接口 | 职责 | 典型实现 |
|---|---|---|---|
| **EventCheck** | `EventCheckProcessorI` | 实时事件匹配（MQ/Event 触发时用） | 保单状态变更匹配、页面行为匹配 |
| **SchedulePullUserList** | `SchedulePullUserListProcessorI` | 定义 ES 查询条件拉用户（定时触发时用） | 按保单属性拉、按用户营销信息拉、按车辆信息拉 |
| **AdditionalCheck** | `AdditionalCheckProcessorI` | 额外过滤条件（任何触发模式都可叠加） | 有没有买过产品、是不是在某个圈人组里、产品折扣状态 |

**执行顺序**：

```
MatchUserGroupByProcessors(message)
    1. Realtime AND  → 全部通过？
    2. Realtime OR   → 至少一个通过？
    3. Delay AND/OR  → 发延迟事件，当前返回 false
    4. UserList      → 填充 PlaceholderFieldMap（定时模式下从 ES 结果取字段）
```

### Stage 3: 执行 Handler

筛选通过后，进入 Handler：

```
DoExecuteHandler(handlerParams)
    │
    ├── 1. PreHandle()    → 校验参数、构建请求对象
    ├── 2. ruleCheck()    → HandlerRules 正则匹配（前置条件）
    ├── 3. fillParam()    → 变量替换（占位符 → 真实值）
    ├── 4. Handle()       → 执行业务逻辑
    └── 5. Record         → 记录执行历史、上报监控指标
```

**Handler 种类**：

| 类别 | 代表 | 做什么 |
|---|---|---|
| **Notify** | `NotifyHandler` | 发推送/WhatsApp/AI Call |
| **Voucher** | `DispatchVoucherHandler` | 发优惠券 |
| **PNAR** | `ExpiryPnArHandler` | 保单到期提醒推送 |
| **Chat** | `WordingChatHandler` | 商城内弹窗/文案 |
| **DataClean** | `MetricsHandler` | 数据指标更新 |
| **Mix** | `RewardWithNotifyHandler` | 发奖励+发通知（组合动作） |

---

## 四、RuleCenter：通用条件判断中心

**RuleCenter 是独立于 Engine 的通用条件判断中心**，提供统一的 `RuleCheck(ruleName, ruleParams) → bool` 接口。

```go
type RuleCenterService interface {
    RuleCheck(ctx, param *RuleCenterCheckParam) bool
}
```

被三个层面调用：

| 调用方 | 场景 |
|---|---|
| AdditionalCheck Processors | 用户筛选时（"买没买过产品"） |
| Plan.TriggerRuleParam | 触发层前置校验（"这个 plan 要不要执行"） |
| Handler 内部 | 执行前二次校验（"到期日有没有变"） |

已注册的 8 个 Rule：

| Rule | 用途 |
|---|---|
| `CheckHasPolicies` | 用户有没有某些状态的保单（支持 partner_id + product_ids） |
| `CheckHasRecentPolicy` | N 分钟内有没有新保单 |
| `CheckIsReferrer` | 是不是推荐人 |
| `CheckIsUbbCheckFail` | 车辆 UBB 检查有没有失败 |
| `CheckMotorExpiredLatest` | 车险到期日是否最新 |
| `CheckProductDiscount` | 产品有没有折扣 |
| `CheckReferralActivityProcessor` | 推荐活动是否在有效期 |
| `CheckOrderRewardInfoProcessor` | 订单有没有特定类型的奖励 |

---

## 五、PlanManager 编排全景

```
PlanDispatch(message)
    ├── PlanCheck(planId)           // 校验 Plan 状态、有效期
    ├── FillingUpBasicParam(msg)    // 填充 AccountId、ProductId 等
    │
    ├── if MQ/Event:
    │       mqPlanExecute()
    │         ├── MatchKafkaUserGroupListByMessage()   // 匹配 UserGroup
    │         └── GetPlanIdsByUserGroupIds()            // 找到关联的 Plan
    │
    ├── if Scheduled:
    │       scheduledPlanExecute()
    │         ├── QueryPlanDetailByPlanIdFromCache()
    │         └── MatchUserGroupByProcessors()          // Processor 链式校验
    │
    └── ExecutePlan(message)
            └── for each planId:
                    MatchHandlerParam() → AsyncExecuteHandler()
                        └── DoExecuteHandler()
                                ├── PreHandle
                                ├── ruleCheck
                                ├── fillParam
                                ├── Handle
                                └── RecordPlanHistory
```

---

## 六、目录结构

主服务路径：`/Users/si.chen/GolandProjects/insurance-marketing`

```
src/engine/
├── boot/                          # 引擎启动引导
├── consumer/inner/                # Kafka 消费者
├── internal/
│   ├── biz/                       # 业务逻辑层（PlanBiz, UserGroupBiz）
│   ├── handler/                   # 执行器（notify, voucher, pnar, chat, mix...）
│   ├── manager/                   # 编排层（PlanManager）
│   ├── processor/user_group/      # 用户筛选处理器
│   │   ├── event_check/           # 实时事件匹配
│   │   ├── schedule_pull_user_list/ # 定时拉用户
│   │   └── additional_check/      # 额外过滤条件
│   └── task_handler/              # 任务条件/奖励处理
├── service/                       # 对外服务接口
└── task/                          # 定时任务
```

---

## 七、设计哲学（面试加分点）

### 1. 策略模式 + 自注册

所有 Processor、Handler、Rule 都通过 `init()` 自注册到全局 map，运行时通过 Code/Name 查找。新增能力 = 新增一个文件 + `init()` 注册，零侵入。

### 2. 三种触发殊途同归

无论定时、MQ、实时事件，最终都汇聚到 `PlanDispatch` → 同一套筛选和执行逻辑。统一入口消除了特殊分支。

### 3. 配置化 > 代码化

Processor Params、Handler Params 都是 JSON，运营在 Admin 页面配置即可。把"做什么"编码在系统里，把"对谁做、什么时候做"交给配置。

### 4. RuleCenter 解耦

业务条件判断抽成独立服务层，Processor / Handler / Trigger 都可以复用同一条规则，避免逻辑散落在各处。

---

## 八、面试叙事顺序

1. **定位**："我负责的是保险运营自动化引擎"
2. **核心模型**："核心抽象是 Plan，定义触发 → 筛选 → 执行三阶段"
3. **触发机制**："支持三种触发方式，最终汇聚到统一入口"
4. **筛选能力**："User Group Processor 链式执行，支持 AND/OR/延迟组合"
5. **执行能力**："Handler 插件化，覆盖推送、发券、数据清洗等场景"
6. **RuleCenter**："通用条件判断中心，被筛选层、触发层、执行层复用"
7. **设计亮点**："策略模式 + 自注册、配置化驱动、三种触发殊途同归"

## 九、可讲排障案例

- [Marketing 容器 OOMKilled 排障](./oom-killed-troubleshooting-star.md)：适合回答稳定性、线上/测试环境排障、容器内存、Go 服务内存问题、如何区分业务 panic 和系统级退出。
- [OOMKilled 排障图解](./oom-killed-troubleshooting-star.html)：适合快速复习证据链和口述顺序。
