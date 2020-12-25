
# 用于快速构建项目的基础库

---

# 配置文件说明

> 一般使用toml作为配置文件, 可以使用命令行-c支持多配置文件<br>
> 配置来源优先级 命令行 > WithViper > WithConfig > WithFiles(Apollo分片优先级最高) > WithApollo > 默认配置文件<br>
> 注意: 多个配置文件如果存在同配置分片会智能合并, 同分片中完全相同的配置节点以最后的文件为准, 从apollo拉取的配置会覆盖相同的文件配置节点

+ 框架配置示例
```toml
[Frame]
Debug = true
FreeMemoryInterval = 2000
WaitServiceRunTime = 500
ContinueWaitServiceRunTime = 120000
```

+ 服务配置示例
```toml
[Services.ApiService]
Bind = ":8080"
IPWithNginxForwarded = false
IPWithNginxReal = false
ShowDetailedErrorOfProduction = false
```

+ 组件配置示例
```toml
[Components.Cache.default]
CacheDB = "memory"
Codec = "msgpack"
DirectReturnOnCacheFault = true
MemoryCacheDB.CleanupInterval = 300000
```

+ 更多配置参考 [core.Config](./core/config.go)
