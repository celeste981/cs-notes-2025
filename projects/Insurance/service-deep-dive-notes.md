# Insurance 服务深挖笔记

## 使用方式

这份笔记不是背诵稿，而是面试追问时的“证据库”。先记住每个服务的边界，再挑 1-2 个你最能讲清楚的链路展开。

## BFF 层

### O-BFF / Operator BFF

面试专项：[O-BFF 面试专项笔记](o-bff-interview-notes.md)

定位：

> Insurance Admin 后台聚合层。Admin 页面统一调用 O-BFF，O-BFF 再编排 Product、Order、Policy、Marketing、Promotion、Risk Control 等下游服务。

关键能力：

- 后台 API 聚合：Service -> Manager/Biz -> Integrate。
- Interceptor 治理：权限、公共参数、脱敏、审批信息、操作日志、重复请求、安全日志、RPC proxy / orchestration。
- Approval Center：高风险操作审批。
- Batch Operate / Task Center：批量导入、异步任务、结果下载。
- Data Fix Center：数据修复需要审批和可追踪执行。
- Operator Executer / Processor：把后台动作标准化为可执行单元。
- ES 查询：复杂后台列表和报表。

代码证据：

- `/Users/si.chen/GolandProjects/insurance-operator-bff/src/service/`
- `/Users/si.chen/GolandProjects/insurance-operator-bff/src/manager/`
- `/Users/si.chen/GolandProjects/insurance-operator-bff/src/integrate/`
- `/Users/si.chen/GolandProjects/insurance-operator-bff/src/approval_center/`
- `/Users/si.chen/GolandProjects/insurance-operator-bff/src/batch_operate_center/`
- `/Users/si.chen/GolandProjects/insurance-operator-bff/src/data_fix_center/`

追问回答：

> O-BFF 不是单纯转发，它要统一权限、审批、脱敏、批量任务、数据修复和下游聚合。后台操作风险比 C 端查询更高，所以很多动作不能直接改库或直调下游，需要走审批和任务中心。

面试扩展：

- 为什么需要 O-BFF：Admin 操作跨多个下游服务，且需要统一审批、审计、脱敏和批量任务。
- 怎么避免大泥球：O-BFF 只做后台编排和治理，领域事实仍属于 Product、Order、Policy 等 owner 服务。
- 批量导入导出怎么讲：上传/查询 -> 创建 batch/task -> 校验/切批 -> 可选审批 -> Task Center 异步执行 -> 记录成功/失败和错误明细。
- 监控怎么讲：Batch Operate、Task Center、RPC Proxy、ES Adapter 有指标，新增指标要避免高基数和敏感 label。

### B-BFF / Business BFF

定位：

> 对接 MP、DP、CL 等 B 端/合作渠道的保险聚合层，核心是插件化工作流。

关键能力：

- Product：产品查询、详情、授权、展示。
- Pricing：保费计算、折扣、价格测试。
- Order：订单创建、更新、状态同步。
- Policy：保单查询、状态监听。
- Plugin Workflow：Rollout、Auth、DataFilling、PresaleCheck、PremiumsCalc、Marketing、Priority、PriceTesting、OrderCreation、OrderUpdate。
- Consumer：MP/DP 订单状态、保单状态、退款、理赔通知。

代码证据：

- `/Users/si.chen/GolandProjects/insurance-business-bff/src/manager/workflow/`
- `/Users/si.chen/GolandProjects/insurance-business-bff/src/manager/plugins/`
- `/Users/si.chen/GolandProjects/insurance-business-bff/src/manager/plugin_selector/`
- `/Users/si.chen/GolandProjects/insurance-business-bff/src/consumer/`
- `/Users/si.chen/GolandProjects/insurance-business-bff/src/service/insurance_business_partner_service.go`
- `/Users/si.chen/GolandProjects/insurance-business-bff/src/service/insurance_calculate_service.go`

追问回答：

> 插件化的原因是渠道和产品差异非常多。如果每个市场、产品、渠道都写 if-else，后面不可维护。插件选择器根据配置和灰度策略选择具体插件，工作流负责把插件串成完整业务流程。

### C-BFF / Consumer BFF

定位：

> C 端用户保险购买和管理入口，聚合产品展示、报价、下单、支付、投保、保单、理赔等能力。

关键能力：

- `motor`：车险报价、车型、车辆信息、地区处理器、投保。
- `accident_health`：意外险/健康险产品、保费、购买预检、保单。
- `sales`：销售订单、支付回调、投保触发。
- 多地区：ID、TH、MY、PH、SG 等地区差异通过 region handler 或策略隔离。

代码证据：

- `/Users/si.chen/GolandProjects/insurance-consumer-bff/src/motor/service/`
- `/Users/si.chen/GolandProjects/insurance-consumer-bff/src/motor/manager/`
- `/Users/si.chen/GolandProjects/insurance-consumer-bff/src/motor/manager/impl/region_handler/`
- `/Users/si.chen/GolandProjects/insurance-consumer-bff/src/accident_health/service/`
- `/Users/si.chen/GolandProjects/insurance-consumer-bff/src/accident_health/manager/`

追问回答：

> C-BFF 的复杂度主要来自产品线和地区差异。比如车险要处理车辆、车主、车型库、报价和地区规则；意外健康险又更关注多人投保、资格预检、优惠和保单展示。

### L-BFF / Logistics BFF

定位：

> 物流保险业务编排层，主要由 OFG/OMS/OA/ODO 消息驱动，而不是同步用户请求驱动。

关键能力：

- Consumer：消费物流状态、店铺授权、商品变更、理赔审核等消息。
- Adaptor Chain：补全 Account、Order、Shop、Item、Category、OMS/OFG 数据。
- Filter Chain：店铺授权、类目、金额、市场、物流方式、时间窗口。
- State Machine：RSF、DLV、ORS、SLS+ 等场景状态。
- Processor：CommitOrder、ConfirmOrder、SurrenderOrder、Claim Report。

本地状态：

- 当前本地没有找到 `insurance-logistics-bff` clone，主要来源是 Confluence `LBFF功能分析`。

追问回答：

> L-BFF 是典型事件驱动系统。物流状态可能乱序、重复、延迟，所以要靠可靠消费、状态机恢复、幂等和过滤器短路机制保证流程稳定。

## 核心服务层

### Product

定位：

> 产品域服务，维护产品、计划、保障、价格规则、费率和部分风控/可见性规则。

关键模块：

- `product/service`：产品、产品扩展、保费、plan/layer 等服务。
- `product/engine`：保费计算器和可见性过滤。
- `product/biz`：product、plan、benefit、pricing、category、broker/insurer、finance config。
- `product/integrate`：Policy、Risk、UEA、motor quotation 等下游/外部依赖。

代码证据：

- `/Users/si.chen/GolandProjects/insurance-product/src/product/engine/premium_calculator.go`
- `/Users/si.chen/GolandProjects/insurance-product/src/product/service/product_service.go`
- `/Users/si.chen/GolandProjects/insurance-product/src/product/biz/pricing_biz.go`
- `/Users/si.chen/GolandProjects/insurance-product/src/product/biz/product_biz.go`

面试说法：

> Product 是售前链路的核心事实源之一，BFF 查询产品和计算保费时都会依赖它。它的难点是产品、plan、benefit、费率规则和地区差异比较多。

### Order

定位：

> 保险订单、支付、财务账单、invoice 等订单域能力。

关键模块：

- `order`：订单创建、状态机、订单历史、policy status 变更处理。
- `payment`：支付单、支付回调、transaction、SPM/CSS。
- `financial`：账单、清结算相关。
- `invoice`：发票、refund note、self-billed/consolidated 等。

代码证据：

- `/Users/si.chen/GolandProjects/insurance-order/src/order/statemachine/order_status_machine.go`
- `/Users/si.chen/GolandProjects/insurance-order/src/order/event/handle_payment_result_event.go`
- `/Users/si.chen/GolandProjects/insurance-order/src/payment/service/insurance_payment_order_service.go`
- `/Users/si.chen/GolandProjects/insurance-order/src/order/integrate/policy_client.go`

面试说法：

> Order 的关键是订单状态机、支付结果幂等和与 Policy 的状态联动。支付成功不等于整个流程结束，还要触发投保、保单状态回传和后续补偿。

### Policy

定位：

> 保单生命周期服务，覆盖投保、保单查询、理赔、批改/取消、续保。

关键模块：

- `policy_center`：保单中心、投保、查询、状态机。
- `policy_claim`：理赔提交、查询、审核、赔付状态。
- `policy_endorse`：取消/批改/退保。
- `policy_renewal`：续保窗口和续保逻辑。

代码证据：

- `/Users/si.chen/GolandProjects/insurance-policy/src/policy_center/state_machine/policy_fsm.go`
- `/Users/si.chen/GolandProjects/insurance-policy/src/policy_center/state_machine_v2/`
- `/Users/si.chen/GolandProjects/insurance-policy/src/policy_claim/state_machine/claim_fsm.go`
- `/Users/si.chen/GolandProjects/insurance-policy/src/policy_endorse/state_machine/cancellation_fsm.go`

面试说法：

> Policy 是出单和保单状态的事实源。它要对接 Product、Order、Account、EA、Document、Risk 等，同时用状态机控制保单、理赔和取消流程。

### Promotion / Voucher

定位：

> 促销、优惠券、库存和价格中心能力，为购买链路和 Marketing 发券提供支持。

关键模块：

- `promotion`：活动、产品活动、库存、税费规则。
- `voucher`：券查询、发放、核销、VSS 集成。
- `repertory`：库存/预算。
- `price_center_v2`：原价、促销价、券、cashback/coinback/direct discount 的价格计算。
- `promotion_risk_check`：活动/券/库存风险校验。

代码证据：

- `/Users/si.chen/GolandProjects/insurance-promotion/src/promotion/`
- `/Users/si.chen/GolandProjects/insurance-promotion/src/voucher/`
- `/Users/si.chen/GolandProjects/insurance-promotion/src/price_center_v2/pricing_engine/`
- `/Users/si.chen/GolandProjects/insurance-promotion/src/promotion_risk_check/`

面试说法：

> Promotion 不是简单发券，它还要处理活动有效性、库存/预算、价格计算、券优先级和风控。Marketing 执行动作时通常会调用它做发券或活动权益。

### Account

定位：

> 统一账户域，维护 SPIN 内部 account 与外部平台用户信息映射。

关键模块：

- `account_manager`：账户查询/生成。
- `auth_manager`、`login_manager`：登录和授权相关。
- `kyc_biz`：KYC 信息。
- `generate_user_consumer`：消费生成用户信息消息。
- `ea_client`：向外部账户/用户系统取数。

代码证据：

- `/Users/si.chen/GolandProjects/insurance-account/src/service/insurance_account_manager_service.go`
- `/Users/si.chen/GolandProjects/insurance-account/src/biz/account_biz.go`
- `/Users/si.chen/GolandProjects/insurance-account/src/biz/kyc_biz.go`
- `/Users/si.chen/GolandProjects/insurance-account/src/consumer/generate_user_consumer.go`

面试说法：

> Account 是跨服务的基础维度。订单、保单、营销、人群和风控都需要稳定的 accountId，否则只靠外部 userId 很难统一查询和分片。

### EA / External Aggregator

定位：

> 外部 broker / insurer 适配层，封装不同保司接口、签名、字段转换和调用记录。

关键模块：

- broker property / vehicle / cargo / claim / insure service。
- 多 broker client：Igloo、Fuse、Qoala、Pasarpolis、Bolttech、ZA、Allianz、MSIG、PolicyStreet、SeaInsure 等。
- converter：把 SPIN 内部模型转换成保司接口模型。
- broker_v2 template：用模板化方式复用 auth、insure、query、claim、policy status 等外部调用。

代码证据：

- `/Users/si.chen/GolandProjects/insurance-external-aggregator/src/service/`
- `/Users/si.chen/GolandProjects/insurance-external-aggregator/src/integrate/broker/`
- `/Users/si.chen/GolandProjects/insurance-external-aggregator/src/integrate/broker_v2/`
- `/Users/si.chen/GolandProjects/insurance-external-aggregator/src/converter/`

面试说法：

> EA 的价值是防腐层。Policy/Order/Product 不应该直接适配每家保司的签名、字段和异常，所以 EA 把外部差异收口，内部服务保持统一模型。

## 数据支撑层

### Marketing-Data

面试专项：[Marketing-Data 面试专项笔记](marketing-data-interview-notes.md)

定位：

> Marketing 的离线数据服务，负责批量人群拉取、ES 查询、文件数据处理、批次历史和 Kafka 推送。

关键模块：

- Internal API：`OfflineDataFetchAndPush`、`QueryOfflinePlanHistory`。
- Admin API：`QueryEsList`。
- Manager：batch 幂等、history 初始化、异步拉数、结果更新、监控。
- Handler：按 DataSource 注册 ES/File/Hive/Insight 等实现。
- Repo/Monitor：`offline_plan_history_tab` 记录执行状态，`plan_batch_extract_total_insurance` 上报批次结果和耗时。

代码证据：

- `/Users/si.chen/GolandProjects/insurance-marketing-data/src/service/marketing_data_internal_service.go`
- `/Users/si.chen/GolandProjects/insurance-marketing-data/src/manager/impl/offline_data_fetch_manager_impl.go`
- `/Users/si.chen/GolandProjects/insurance-marketing-data/src/handler/offline_data_fetch_handler.go`

追问回答：

> batch_id 幂等很重要，因为离线任务失败重试时，重复拉取和重复推送都可能造成重复触达用户。

面试扩展：

- 为什么拆服务：Marketing 负责计划执行，Marketing-Data 负责重数据抽取，失败模式和资源消耗不同。
- 主链路：Marketing 传 `plan_id`、`batch_id`、`data_source`、`group_condition`、`dest_topic` -> Marketing-Data 校验和幂等 -> 写 history -> 异步 handler 拉数 -> 推 Kafka -> 更新结果。
- DataSource 扩展：实现 `OfflineDataFetchHandler`，注册 `data_source` 初始化函数，主流程不用改。
- `QueryEsList`：用于 Admin 预览 ES 查询结果，`size <= 100`，不承担大批量抽取。

### Canal / ES

定位：

> MySQL 到 ES 的近实时同步链路，服务于 Admin 查询、宽表、分表聚合和报表分析。

关键模块：

- Canal Server / Instance：解析 MySQL binlog。
- Kafka：承接 Canal server 到 adapter 的消息。
- Adapter：按 sync script 写 ES。
- ES mapping：字段类型和索引结构。
- DLQ / 补数：处理异常消息阻塞和最终失败。

代码证据：

- `/Users/si.chen/GolandProjects/canal-adapter/README.md`
- `/Users/si.chen/GolandProjects/canal-adapter/es7x/`
- `/Users/si.chen/GolandProjects/canal-adapter/es6x/`
- `/Users/si.chen/GolandProjects/canal-adapter/docs/id_vn_rds_expansion_canal_td.md`

追问回答：

> ES 查询问题要分层定位：业务库有没有数据，binlog 有没有被 Canal 消费，Kafka/Adapter 有没有阻塞，ES mapping 是否正确，最后才看查询 DSL。

## 面试防守边界

| 可以强讲 | 谨慎说法 |
| --- | --- |
| Marketing、O-BFF、Marketing-Data、Canal 是我们组相关链路，我能讲架构和排障方式。 | 不说“我个人独立 owner 所有服务”。 |
| B-BFF/C-BFF/L-BFF/Product/Order/Policy/Promotion/Account/EA 是上下游核心服务，我理解边界和调用契约。 | 不把这些都说成个人核心项目，除非能补 PR/Jira 证据。 |
| ES/Canal 是查询支撑链路，MySQL owner 服务是事实源。 | 不说 ES 是强一致事实源。 |
