## 2026-05-14 (cont.) — P2: Flink Job 完善

### 改造内容
1. **死信处理**: MetricDeserializer catch JSON异常→log+return null (Flink drop + metrics counter)
2. **Redis pub/sub**: 后台 daemon 线程 subscribe `monitor:rule-updated` 替代 5s 轮询，Backend 规则变更后秒级生效
3. **持续时间判断**: `MapState<ruleId, firstViolationMs>` + `TimerService`，metric 持续超阈值 N 秒才触发告警
4. **告警去重**: `MapState<ruleId, lastFireMs>`，60s 窗口内同 host+rule 不重复告警
5. **KafkaSink**: 替换 `PrintSinkFunction`，`AT_MOST_ONCE` 低延迟投递到 `alerts` topic
6. **Checkpointing**: 60s 周期，保障 ValueState/MapState 容错恢复

### 本次提交
`ac74aac` feat(flink): pub/sub hot-reload, duration tracking, alert dedup, KafkaSink

---

## 2026-05-14 — 提交历史总结

```
ac74aac feat(flink): pub/sub hot-reload, duration tracking, alert dedup, KafkaSink
041d5e9 docs: add development diary with full history and code review
d75a503 feat(backend): validation middleware, pagination totals, health check, audit log
2e6b611 feat(backend): complete service CRUD + gRPC proto + Redis pub/sub
bf55870 refactor(agent): replace gopsutil with pure cgroup v2 /proc reads
6f8d5d9 docs: update README with tech stack and project overview
717ab64 Merge branch 'main'
a3c6246 feat: initialize monitoring platform skeleton
```

### 进度

| Phase | 状态 | 提交 |
|-------|------|------|
| P0 Agent | ✅ | bf55870 |
| P1 Backend | ✅ | d75a503 |
| P2 Flink | ✅ | ac74aac |
| P3 Alertmanager | ⏳ | — |
| P4 Dashboard | ⏳ | — |
| P5 DevOps | ⏳ | — |
