# Insurance 业务流程地图

## 先讲人话

Insurance 的核心不是“一个接口”，而是几条端到端链路：

- C 端买保险：用户看产品、算保费、下单、支付、投保、查保单、理赔。
- B 端嵌入式保险：MP/DP/CL 等渠道把保险嵌在自己的业务流程里，B-BFF 做聚合和插件编排。
- 物流保险：不是用户主动下单，而是 OFG/OMS 等物流消息驱动状态机。
- Admin 后台：运营/客服/Dev 在 Admin 操作，O-BFF 聚合下游服务，并通过审批、批量、任务和数据修复控风险。
- Marketing：运营配置计划后，系统按定时/MQ/RPC 触发筛人和执行通知、发券、PNAR。
- 数据同步：Canal 把 MySQL 变更同步到 ES，支撑后台查询和报表。

## 1. C 端购买闭环

```text
App / H5
  -> C-BFF
  -> Product: 产品/计划/保障/价格基础信息
  -> Promotion / Voucher: 活动、优惠券、折扣
  -> Risk Control: 投保前风控
  -> Order / Payment: 创建订单、支付
  -> Policy: 支付成功后投保出单
  -> EA / UEA: 对接保司或 broker
  -> Document / Notification / Marketing: 电子保单、通知、后续营销
```

面试怎么说：

> C-BFF 的价值是把用户侧复杂流程聚合成前端好用的 API。前端不需要知道 Product、Promotion、Order、Policy、EA 怎么拆，只关心产品展示、保费计算、下单支付和保单查询。系统内部用 manager/biz/integrate 分层，把不同产品线和地区差异隔离起来。

追问点：

- 为什么需要 BFF：屏蔽多服务调用、适配前端视图、承接地区/产品差异。
- 风险点：支付回调幂等、投保失败补偿、金额一致性校验、状态机合法流转。

## 2. B 端 MP/DP/CL 嵌入式保险

### MP 商品险流程

```text
MP PDP / Checkout
  -> B-BFF
  -> Product: category whitelist / 产品查询 / 风控规则 / 保费计算
  -> Promotion: itembass / 营销展示
  -> ODO Kafka: 创单和订单状态消息
  -> B-BFF consumer
  -> Account / EA: 用户和 MP 订单数据补全
  -> Order: 创建/更新保险订单
  -> Policy: 完成订单后投保出单
  -> ODP: 查询保单列表和详情
```

关键事实：

- PDP 和 Checkout 会先判断 category whitelist，只有命中后才调用保险产品和保费接口。
- 创单与投保更新可以拆不同 topic，避免大促创单高吞吐影响低吞吐的保司投保链路。
- B-BFF 通过插件工作流串联 Rollout、Auth、DataFilling、PresaleCheck、PremiumsCalc、Marketing、Priority、OrderCreation/Update。

### CL 贷款险流程

```text
CashLoan
  -> B-BFF calculate_premiums
  -> Product: 按 loan tenure / amount / interest 计算保费
  -> B-BFF create_policy_order
  -> Order: 创建保险订单
  -> CashLoan confirm / cancel
  -> B-BFF update
  -> Policy: 发放成功投保，发放失败取消
```

### DP 场景流程

```text
DP 支付成功
  -> Kafka order_status_dp
  -> B-BFF consumer
  -> Order / Policy
  -> 直接投保或查询保单详情
```

面试抓手：

> B-BFF 的难点不是单个 RPC，而是把不同渠道、不同产品、不同市场抽象成插件化工作流。新增渠道或策略时，通常新增 plugin 或配置 selector，而不是把所有逻辑堆在一个 if-else 里。

## 3. 物流保险链路

```text
OFG / OMS / OA / ODO Kafka
  -> L-BFF Consumer
  -> Adaptor Chain: 补全 Account / Order / Shop / Item / Category / OMS / OFG 数据
  -> Filter Chain: 店铺授权、类目、金额、市场、物流方式、时间窗口
  -> State Machine: RSF / DLV / ORS / SLS+
  -> Processor: CommitOrder / ConfirmOrder / SurrenderOrder
  -> Insurance Core / Order / Policy
  -> Claim Report Center
```

典型场景：

- RSF：退货申请、退货审核通过、退货完成，驱动退货运费保险订单状态。
- DLV：订单创建、发货、签收，驱动配送保险确认。
- ORS / SLS+：海外退货或跨境物流状态变化驱动确认、理赔报案。

面试抓手：

> L-BFF 的核心是消息驱动。物流链路很多状态来自外部系统，不能靠用户同步请求推进，所以用 Consumer + Adaptor + Filter + State Machine + Processor 的模型，既能补齐数据，又能把业务规则和状态流转控制住。

## 4. Admin 后台运营链路

```text
Insurance Admin
  -> O-BFF Service / Interceptor
  -> 权限 / 脱敏 / 参数校验 / 审批判断
  -> Manager / Biz 编排
  -> Product / Order / Policy / Promotion / Marketing / Risk Control
  -> DB / ES / Config Center / Kafka / Reliable Event
```

常见后台能力：

- 配置查询和更新：产品、活动、券、营销计划、灰度、开关。
- 审批：高风险配置先创建审批实例，审批通过后执行。
- 批量操作：CSV 上传、批量校验、审批、异步执行、结果下载。
- 数据修复：人工修复必须可追踪、可审批、可回滚或可审计。
- ES 查询：后台列表和报表走 ES，避免扫业务库。

面试抓手：

> O-BFF 是后台聚合入口。它把 Admin 页面操作转换成多个下游服务调用，并在入口处统一做审批、批量、任务、脱敏和审计，降低前端对接复杂度，也降低后台误操作风险。

## 5. Marketing 运营计划链路

```text
Admin 配置 Plan / User Group
  -> Marketing Plan Engine
  -> Trigger: GoScheduler / Kafka / RPC Event
  -> User Group Processor
  -> RuleCenter / AdditionalCheck
  -> Handler: PN / AR / Voucher / Chat / Event / Recommendation
  -> UEA / Promotion / Voucher / Policy / Order / C-BFF
  -> 执行结果、batch history、Prometheus
```

和 Marketing-Data 的关系：

```text
Marketing 定时计划
  -> Marketing-Data OfflineDataFetchAndPush
  -> ES / S3 / CSV / Hive / Insight
  -> offline_plan_history
  -> Kafka
  -> Marketing 消费并执行 Handler
```

面试抓手：

> Marketing 是执行引擎，Marketing-Data 是离线数据助手。拆分后，大查询和文件处理不会拖慢主服务，也能减少上线主服务时打断离线任务的风险。

## 6. ES / Canal 查询链路

```text
业务服务 MySQL
  -> binlog
  -> Canal Server / Instance
  -> Kafka
  -> Canal Adapter
  -> ES mapping / sync script
  -> Elasticsearch index
  -> O-BFF / Admin 查询
```

关键点：

- MySQL owner 服务仍然是事实源，ES 是查询和分析视图。
- Admin 的复杂列表、宽表和分表聚合适合走 ES。
- Mapping 要先定义再同步字段；字段类型错了通常要重建索引或补数据。
- Canal 同步异常要看 instance、Kafka、Adapter、ES bulk response、DLQ 或补数窗口。

面试抓手：

> Canal/ES 的价值是让后台查询快，但它不是强一致主链路。排障时要先判断是 MySQL 没变、binlog 没消费、Kafka/Adapter 阻塞、ES mapping 错，还是查询 DSL 不对。

## 7. 一分钟总回答

> 我们 Insurance 后端有几条主链路：用户侧通过 C-BFF 完成产品展示、保费计算、下单支付、投保和保单/理赔查询；B 端嵌入式场景通过 B-BFF 对接 MP、DP、CL，用插件化工作流做授权、数据填充、售前校验、保费计算、营销和订单创建；物流保险通过 L-BFF 消费 OFG/OMS 消息，用适配器链补数据、过滤器链校验规则、状态机推进订单；后台运营通过 O-BFF 聚合 Product、Order、Policy、Promotion、Marketing 等服务，并承接审批、批量、任务和数据修复；营销侧通过 Marketing Plan Engine 执行通知、发券、PNAR，离线大数据由 Marketing-Data 拉取；查询和报表依赖 Canal 把 MySQL 同步到 ES。我的重点是后台运营和营销域，但我会把这些上下游链路一起讲清楚。
