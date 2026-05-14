# 哨兵系统 — 开发日志

## 项目概述

- **项目**: Sentinel-pingan (哨兵：平安银行的守护神)
- **仓库**: https://github.com/magicman1324/Sentinel-pingan
- **PRD**: 金融级服务器状态监控面板与可观测性体系建设
- **分支**: `main`

---

## 2026-05-14 — 技术方案 & 项目骨架

### 阶段 0: 技术方案输出
- 阅读 PRD，输出浓缩版技术方案和选型
- 核心链路: `Go Agent (Cgroups v2) → Kafka → Flink (+Redis规则) → Alertmanager → 钉钉/邮件`
- 关键决策: Cgroups v2 而非 eBPF（金融求稳）、Flink At-Most-Once（延迟优先）、TiDB HTAP（消除Lambda架构）

### 阶段 1: 基础骨架搭建
- **Agent** (`agent/`): Go 采集器骨架，gopsutil采集，Sarama Kafka producer，环境变量配置
- **Backend** (`backend/`): Gin REST + gRPC 双栈，sqlx + TiDB，go-redis，规则CRUD骨架
- **Flink Job** (`flink-job/`): Java 17 + Flink 1.20，Kafka source → RuleEvaluator → Alert sink，Redis轮询热加载规则
- **Deploy** (`deploy/`): docker-compose (ZK+Kafka+Redis+TiDB+Backend+Flink)，K8s DaemonSet + Deployment
- **Schema** (`schema/tidb.sql`): rules/alerts/channels DDL
- **Configs** (`configs/`): Alertmanager 路由/抑制/钉钉+邮件通道配置

### 阶段 2: Go代码Review & 修复 (19处)
**Agent修复 (6处)**:
1. Collector接口 `interface{}` → `*model.Metrics` 消除boxing
2. CPU: `runtime.GOMAXPROCS(0)` 替代 `NumCPU()` (Go 1.25容器感知)
3. CPU: 首次 `PercentWithContext` non-blocking
4. Memory: 补齐 `OOMCount` 读取cgroup v2 `memory.events`
5. Network: 两次采样差值计算真实bytes/sec速率
6. Sender: `WaitForAll` + `Idempotent=true` 金融级可靠性

**Backend修复 (8处)**:
7. model/rule.go 新增 `Severity` 字段 (与Flink端对齐)
8. model/alert.go `AlertStatus` 类型常量替代裸string
9. repository: `ListAll` + `ListEnabled` 双方法
10. service: 补全 `UpdateRule`/`DeleteRule`/`ResolveAlert`/`GetEnabledRules`
11. service: `SyncEnabledRulesToRedis` + mutex 防并发覆盖
12. handler: strconv错误不再吞掉，id=0误删bug修复
13. handler: 校验 `name`/`metric` 必填字段
14. Shutdown + 10s timeout 优雅关机

**基础设施**:
- schema: rules表新增 `severity` 列
- agent/go.mod: 清理 viper/yaml 等未用依赖

### 阶段 3: P0 — Agent Cgroups v2 原生采集
- **删除 gopsutil** — 零外部依赖 (仅 sarama)
- **新增 `internal/cgroup/detect.go`** — 自动探测 cgroup v1/v2，统一路径读取
- **CPU采集器重写** — `/proc/stat` 差值计算 `percent_used`
- **Memory采集器重写** — cgroup v2 `memory.current/max/events`，fallback `/proc/meminfo`
- **Disk采集器重写** — `/proc/mounts` + `syscall.Statfs`
- **Network采集器重写** — `/proc/net/dev` 累计值差值计算速率
- **Hostname** — `/proc/sys/kernel/hostname` 优先，`os.Hostname()` fallback
- **编译验证**: `GOOS=linux go vet` 通过，`go build` 通过

### 阶段 4: P1 — Backend 完整规则引擎 API
- **gRPC proto**: `api/proto/monitor.proto` (RuleService + AlertService)
- **gRPC server**: 手写stub，AlertService 接入 business logic
- **Service层完善**: Create/Update/Delete + `publishRuleUpdate()` (Pipeline: SET + PUBLISH)
- **Redis pub/sub**: 规则变更 → `monitor:rule-updated` channel → Flink订阅
- **ValidateRule middleware**: name/metric必填 + operator白名单校验
- **分页总览**: alers接口返回 `{data, total, page, size}`
- **健康检查升级**: ping TiDB + Redis，返回 `{status, tidb, redis}`
- **审计日志**: `audit_logs` 表 + model + repository
- **编译验证**: `go vet` 通过

---

## 当前代码统计

| 模块 | 文件数 | Go 行数 | 状态 |
|------|--------|---------|------|
| Agent | 11 | 594 | ✅ vet clean |
| Backend | 14 | 755 | ✅ vet clean |
| Flink Job | 8 | ~300 | ⏳ 骨架待完善 |
| Configs | 2 | - | ✅ |
| Deploy | 3 | - | ✅ |
| Schema | 1 | 69 SQL | ✅ |
| **总计** | **39** | **1,349 Go** | |

## Git 提交历史

```
d75a503 feat(backend): validation middleware, pagination totals, health check, audit log
2e6b611 feat(backend): complete service CRUD + gRPC proto + Redis pub/sub
bf55870 refactor(agent): replace gopsutil with pure cgroup v2 /proc reads
6f8d5d9 docs: update README with tech stack and project overview
717ab64 Merge branch 'main' (initial repo README)
a3c6246 feat: initialize monitoring platform skeleton
```

## 待办

- [ ] P2: Flink Job — pub/sub改造 + 持续时间判断 + 告警去重 + KafkaSink
- [ ] P3: Alertmanager — 钉钉卡片模板 + 邮件模板 + 抑制规则 + 静默API + HA
- [ ] P4: Dashboard — React SPA 前端面板
- [ ] P5: DevOps — CI/CD + Helm Chart + 压测
