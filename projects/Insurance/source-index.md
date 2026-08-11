# Insurance 资料来源索引

这个文件只记录“从哪里来的信息”。后面准备面试回答时，优先用这些来源做证据，避免凭感觉讲系统。

## Confluence 已检索页面

| 页面 | URL | 可用信息 | 面试用途 |
| --- | --- | --- | --- |
| OBFF 功能分析 | https://confluence.shopee.io/pages/viewpage.action?pageId=2968425075 | O-BFF 是 Admin 后台聚合入口，包含 Product、Policy、Order、Claim、Risk、Finance、Approval、Batch、Task、Data Fix、Operator Executer、ES 查询等模块。 | 回答“Admin 后台为什么需要 BFF”“O-BFF 做什么”。 |
| B-BFF 功能分析 | https://confluence.shopee.io/pages/viewpage.action?pageId=2968424010 | 对接 MP/DP/CL 等渠道；核心是 Product/Pricing/Order/Policy 聚合和插件化工作流，包含 Rollout、Auth、DataFilling、PresaleCheck、PremiumsCalc、Marketing、Priority、PriceTesting、OrderCreation、OrderUpdate。 | 回答“B 端嵌入式保险售卖链路”“插件化架构”。 |
| CBFF 功能分析 | https://confluence.shopee.io/pages/viewpage.action?pageId=2968423411 | C 端 BFF，覆盖基础服务、motor、accident health、sales、car library、free insurance、push、admin；流程包括产品展示、报价、预检、创建订单、支付、投保、保单/理赔查询。 | 回答“C 端保险购买闭环”。 |
| LBFF 功能分析 | https://confluence.shopee.io/pages/viewpage.action?pageId=2968424743 | 物流保险 BFF，面向 OFG/OMS 消息；使用 Consumer、Adaptor Chain、Filter Chain、状态机、Processor 支撑 DLV、RSF、ORS、SLS+ 等场景。 | 回答“物流保险为什么适合消息驱动和状态机”。 |
| EC Buyer业务梳理 | https://confluence.shopee.io/pages/viewpage.action?pageId=2578909398 | MP/CL/DP 的完整购买流程、接口、topic、产品类型、市场差异；MP 有 category whitelist、PDP 产品查询、Checkout 保费试算、ODO 创单/状态消息、ODP 保单详情。 | 回答“MP/CL/DP 业务流程”。 |
| [Marketing] 运营计划架构 | https://confluence.shopee.io/pages/viewpage.action?pageId=1825819385 | Marketing 与 Marketing-Data 拆分；Marketing 是运营主系统，Marketing-Data 处理离线数据提取和实时 MQ 转发；能力包括定时离线人群、实时营销消息、RPC 活动查询、异步事件、推荐产品。 | 回答“Marketing 为什么拆 Marketing-Data”。 |
| insurance-marketing-data 技术设计 | https://confluence.shopee.io/pages/viewpage.action?pageId=1613567933 | 从 Marketing 拆出离线/实时数据提取；保留离线数据接口，新增 batch history、pre_filter、reliable event/consumer；上线切换 go-scheduler 和预处理开关。 | 回答“离线数据服务的边界和上线风险”。 |
| ES 使用分享 | https://confluence.shopee.io/pages/viewpage.action?pageId=2347394775 | Admin 使用 ES 做宽表查询、分表聚合和数据分析；Canal 通过 binlog、server、adapter、Kafka、ES mapping 把 MySQL 变更同步到 ES。 | 回答“Canal/ES 同步链路和排障”。 |
| 同步阻塞优化 | https://confluence.shopee.io/pages/viewpage.action?pageId=3045066306 | Canal/ES 同步遇到异常消息可能阻塞 partition；方案关注 DLQ、异常分类、指标、失败窗口和补数。 | 回答“ES 同步异常怎么处理”。 |

## 本地 Project KB

| 路径 | 可用信息 |
| --- | --- |
| `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/SPIN_Global_Service_Map.md` | Insurance 全局服务地图：入口层、核心服务、外部集成、Admin/UI、数据支撑。 |
| `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/System_Architecture.md` | O-BFF 分层、审批、批量、任务、数据修复、可靠事件、proxy/orchestration。 |
| `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/system/Admin_Platform_Component_Map.md` | O-BFF Admin 配置化组件地图：proxy、import、export、mask、approval、operation history 的配置入口和运行链路。 |
| `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/decisions/ADR-0002-admin-current-assembly-configuration.md` | 当前 Assembly 配置方案：proxy、batch import/export、mask、approval mapping 的推荐配置方式和禁用旧口径。 |
| `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/tech/Admin_Downstream_API_Contract.md` | O-BFF 接入下游 API 时的 Gateway、字段、operator、分页、批量导入/导出契约。 |
| `/Users/si.chen/GolandProjects/insurance-operator-bff/project-kb/tech/OBFF_Monitoring_Alerts.md` | O-BFF Batch Operate、Task Center、RPC Proxy、ES Adapter 等监控指标和告警查询口径。 |
| `/Users/si.chen/GolandProjects/insurance-marketing/project-kb/Cross_App_Overview.md` | Marketing 上游 C-BFF/H5/O-BFF/B-BFF/L-BFF/Policy events/Data report，下游 Account/Product/Promotion/Policy/Order/Marketing-Data/UEA 等。 |
| `/Users/si.chen/GolandProjects/insurance-marketing/project-kb/Tech_Architecture.md` | Marketing 技术栈、Engine/basic/common/model/repo 分层、Kafka、reliable-event、MySQL sharding、Redis、ES、S3、Config Center。 |

## 本地代码已扫路径

| 服务 | 当前路径 | 已确认代码痕迹 |
| --- | --- | --- |
| O-BFF | `/Users/si.chen/GolandProjects/insurance-operator-bff` | `src/service`、`src/manager`、`src/integrate`、`approval_center`、`batch_operate_center`、`task_center`、`data_fix_center`。 |
| Marketing | `/Users/si.chen/GolandProjects/insurance-marketing` | `src/engine`、`src/basic/rule-center`、`src/basic`、`src/common/integrate`。 |
| Marketing-Data | `/Users/si.chen/GolandProjects/insurance-marketing-data` | `OfflineDataFetchAndPush`、`QueryOfflinePlanHistory`、`QueryEsList`、DataSource handler 自注册、ES/File/Hive/Insight 实现、`offline_plan_history_tab`、`plan_batch_extract_total_insurance`。 |
| Canal Adapter | `/Users/si.chen/GolandProjects/canal-adapter` | REST launcher、ES6/ES7 adapter、sync config、RDS 扩容 TD、Canal materials。 |
| B-BFF | `/Users/si.chen/GolandProjects/insurance-business-bff` | `src/manager/workflow`、`src/manager/plugins`、`src/manager/plugin_selector`、MP/DP/CL service 和 consumer。 |
| C-BFF | `/Users/si.chen/GolandProjects/insurance-consumer-bff` | `src/motor`、`src/accident_health`、`src/sales` 相关 service/manager/biz/region_handler。 |
| Product | `/Users/si.chen/GolandProjects/insurance-product` | `src/product/service`、`engine/premium_calculator`、`biz/product/plan/benefit/pricing`、`integrate/unified_ea/risk/policy`。 |
| Order | `/Users/si.chen/GolandProjects/insurance-order` | `src/order`、`src/payment`、`src/financial`、`src/invoice`，含状态机、payment result event、policy status consumer。 |
| Policy | `/Users/si.chen/GolandProjects/insurance-policy` | `policy_center`、`policy_claim`、`policy_endorse`、`policy_renewal`，含 policy/claim/cancellation 状态机和外部 broker 集成。 |
| Promotion | `/Users/si.chen/GolandProjects/insurance-promotion` | `promotion`、`voucher`、`repertory`、`price_center_v2`、`promotion_risk_check`。 |
| Account | `/Users/si.chen/GolandProjects/insurance-account` | `account/auth/login/kyc` biz，`generate_user_consumer`，EA client。 |
| EA | `/Users/si.chen/GolandProjects/insurance-external-aggregator` | broker property/vehicle/cargo/claim/insure service，多 broker client/converter，HTTP client 和 broker_v2 template。 |

## 当前缺口

- 本地没有找到 `insurance-logistics-bff` 当前 clone；L-BFF 暂时以 Confluence 功能分析为主。
- Product/Order/Policy/Promotion 等不是重点个人项目，当前只整理边界和追问点，不写成“我个人独立负责”。
- 个人 PR、Jira、上线记录仍然是 `待补充`。
