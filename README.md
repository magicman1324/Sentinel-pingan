# Sentinel-pingan

**哨兵：平安银行的守护神** — 金融级服务器状态监控面板与可观测性平台

## 架构

```
Go Agent (Cgroups v2) → Kafka → Flink (+ Redis 规则引擎) → Alertmanager → 钉钉/邮件
                                            ↓
                                      TiDB (HTAP) ← Backend API (REST+gRPC)
```

## 技术栈

| 组件 | 选型 |
|------|------|
| 边缘Agent | Go 1.25+ (PGO优化, 容器感知, 单二进制 <5MB) |
| 消息队列 | Apache Kafka (削峰填谷, 防误报双重阈值) |
| 流计算 | Apache Flink (5~17万 QPS, <200ms延迟) |
| 规则引擎 | Redis 热加载 (<1s生效, 原子/复合规则) |
| 告警调度 | Prometheus Alertmanager (去重/分组/抑制/Gossip HA) |
| 存储 | TiDB HTAP (强一致 + 实时分析) |
| 后端 | Go (Gin REST + gRPC) |

## 项目结构

```
├── agent/          # Go 采集Agent (CPU/内存/磁盘/网络)
├── backend/        # 管控API (规则CRUD + 告警查询)
├── flink-job/      # Flink流处理Job (规则评估 + 告警生成)
├── configs/        # Alertmanager路由 + Agent配置
├── deploy/         # docker-compose + K8s manifests
└── schema/         # TiDB DDL
```

## 快速开始

```bash
make dev-up      # 启动全部依赖: Kafka + Redis + TiDB + Flink
make agent       # 编译 Agent    → agent/bin/monitor-agent
make backend     # 编译 Backend  → backend/bin/monitor-backend
make migrate     # 初始化 TiDB 表结构
```

## 性能基线

- Agent 冷启动: < 50ms | CPU: < 1% | 内存: < 20MB
- 端到端延迟: < 3s (网卡→钉钉)
- 系统可用性: 99.999%

## License

MIT
