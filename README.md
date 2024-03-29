## 分布式微服务系统

微服务架构要解决的问题：

1. cmd: 命令行工具
   1. protoc-gen-go-http: protoc 插件
   2. protoc-gen-go-errors: protoc 插件
2. codec: 编码/解码/序列化（原encoding）
3. config: 配置中心定义及实现  
   通过实现`Source`接口来定义多个配置源，其中相同的配置项会因加载顺序的先后而被覆盖；另外，通过实现`Watcher`接口而实现配置修改的热更新。
4. constant: 常量
5. endpoint: 服务端点定义及实现
6. errors: 统一错误定义
7. example: 测试案例
8. logger: 统一日志定义及实现
9. metadata: 元信息/数据定义
10. metrics: 接口监控定义
11. midware: 中间件
    1. auth: 权限认证
    2. breaker: 熔断器（原circuitbreaker）
    3. logging: 日志
    4. metadata: 元数据
    5. metrics: 性能监测
    6. filter: 路由过滤器（原selector）
    7. ratelimit: 限流器
    8. recovery: 异常恢复
    9. tracing: 链路追踪
    10. validate: 参数校验
12. pkg: 底层库包
13. proto: proto文件
14. registry: 服务注册与发现
15. selector: 路由与负载均衡
16. transport: 网络传输协议
17. util: 基础工具

### 代码评审

```sh
# struct优化
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
# 查看哪些结构体有待优化空间（--test=false: 忽略测试文件）
fieldalignment --test=false ./...
# --fix: 直接优化结构体（此时字段注释会被覆盖）
fieldalignment --fix ./...

# 漏洞检测
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

### 安装protoc插件

```sh
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/carmel/go-micro/cmd/protoc-gen-go-http@latest
go install github.com/carmel/go-micro/cmd/protoc-gen-go-errors@latest
```

### 部署架构

![architecture](go-micro.png)
