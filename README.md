### 分布式微服务系统

微服务架构要解决的问题：

1. cmd: 命令行工具
   1. protoc-gen-go-http: protoc 插件
   1. protoc-gen-go-errors: protoc 插件
2. codec: 编码/解码/序列化
3. config: 配置中心定义及实现  
   通过实现`Source`接口来定义多个配置源，其中相同的配置项会因加载顺序的先后而被覆盖；另外，通过实现`Watcher`接口而实现配置修改的热更新。
4. endpoint: 服务端点定义及实现
5. errors: 统一错误定义
6. example: 测试案例
7. gateway: API网关
8. logger: 统一日志定义及实现
9.  metadata: 元信息/数据定义
10. metrics: 接口监控定义
11. midware: 中间件
    1. auth: 权限认证
    2. breaker: 熔断器
    3. logging: 日志
    4. metadata: 元数据
    5. metrics: 性能监测
    6. ratelimit: 限流器
    7. recovery: 异常恢复
    8. selector: 路由选择器/拦截器
    9. tracing: 链路追踪
    10. validate: 参数校验
12. pkg: 底层库包
13. registry: 服务注册与发现
14. selector: 路由与负载均衡
15. tool: 基础工具
16. transport: 网络传输协议


### 代码质量提升
```sh
# struct优化
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
# 查看哪些结构体有待优化空间
fieldalignment ./...
# --fix: 直接优化结构体（此时字段注释会被覆盖）
fieldalignment --fix ./...

# 漏洞检测
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

### 23-08-4 task
- 理清kratos运行机制
- 搬运kratos常用脚手架功能
- 建立一个配置中心——所有配置均可从配置中心获得

### 想法
参照[go-orb](https://github.com/go-orb/go-orb)框架的设计，可以很便利的通过yaml配置来加载微服务架构中的各个组件（或称插件，Plugins），所以事先将各个组件设计成插件化。
