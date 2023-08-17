### 分布式微服务系统

微服务架构要解决的问题：

1. 服务注册与发现 registry
2. 路由与负载均衡 selector
3. 配置中心 config
4. 传输协议 transport
5. 服务端点 endpoint
6. 日志管理 log
7. 中间件 midware
   1. 认证 auth
   2. 熔断器 circuitbreaker
   3. 限流器 ratelimit
   4. 链路追踪 tracing
   5. 性能监控 metrics
8. API 网关 gateway
9. 工具类 util


### 结构体字段排列优化
```sh
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
# 查看哪些结构体有待优化空间
fieldalignment ./...
# --fix: 直接优化结构体（此时字段注释会被覆盖）
fieldalignment --fix ./...
```
