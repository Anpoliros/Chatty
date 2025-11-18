# Chatty Server

Chatty聊天服务的后端服务。

## 目录结构说明

```
server/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口,路由配置
│
├── internal/                     # 内部包(不对外暴露)
│   ├── api/                     # HTTP API处理器
│   │   ├── auth.go              # 认证相关API
│   │   ├── user.go              # 用户相关API
│   │   └── message.go           # 消息相关API
│   │
│   ├── websocket/               # WebSocket处理
│   │   ├── hub.go               # 连接管理中心
│   │   ├── client.go            # 客户端连接
│   │   ├── message.go           # 消息定义
│   │   └── handler.go           # WebSocket处理器
│   │
│   ├── service/                 # 业务逻辑层
│   │   ├── user_service.go      # 用户业务逻辑
│   │   └── message_service.go   # 消息业务逻辑
│   │
│   ├── repository/              # 数据访问层
│   │   ├── database.go          # 数据库连接
│   │   ├── redis.go             # Redis连接
│   │   ├── user_repository.go   # 用户数据访问
│   │   └── message_repository.go # 消息数据访问
│   │
│   ├── model/                   # 数据模型
│   │   ├── user.go              # 用户模型
│   │   └── message.go           # 消息模型
│   │
│   └── middleware/              # 中间件
│       ├── auth.go              # JWT认证中间件
│       └── cors.go              # CORS中间件
│
├── pkg/                         # 公共库(可对外暴露)
│   ├── logger/                  # 日志工具
│   │   └── logger.go
│   └── utils/                   # 工具函数
│       ├── config.go            # 配置管理
│       └── jwt.go               # JWT工具
│
├── config/                      # 配置文件
│   └── config.yaml              # 配置文件示例
│
├── migrations/                  # 数据库迁移脚本
│
├── .env.example                 # 环境变量示例
├── docker-compose.yaml          # Docker Compose配置
├── go.mod                       # Go模块定义
└── README.md                    # 本文件
```

## 开发指南

### 添加新功能

#### 1. 定义数据模型

在 `internal/model/` 中定义:

```go
// internal/model/feature.go
package model

type Feature struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### 2. 实现数据访问层

在 `internal/repository/` 中实现:

```go
// internal/repository/feature_repository.go
package repository

type FeatureRepository struct {
    db *gorm.DB
}

func NewFeatureRepository(db *gorm.DB) *FeatureRepository {
    return &FeatureRepository{db: db}
}

func (r *FeatureRepository) Create(feature *model.Feature) error {
    return r.db.Create(feature).Error
}
```

#### 3. 实现业务逻辑层

在 `internal/service/` 中实现:

```go
// internal/service/feature_service.go
package service

type FeatureService struct {
    featureRepo *repository.FeatureRepository
}

func NewFeatureService(featureRepo *repository.FeatureRepository) *FeatureService {
    return &FeatureService{featureRepo: featureRepo}
}

func (s *FeatureService) CreateFeature(name string) (*model.Feature, error) {
    feature := &model.Feature{Name: name}
    err := s.featureRepo.Create(feature)
    return feature, err
}
```

#### 4. 实现API处理器

在 `internal/api/` 中实现:

```go
// internal/api/feature.go
package api

func CreateFeature(c *gin.Context) {
    var req struct {
        Name string `json:"name" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 调用service层
    // ...

    c.JSON(200, gin.H{"message": "success"})
}
```

#### 5. 注册路由

在 `cmd/server/main.go` 中注册:

```go
v1.POST("/features", api.CreateFeature)
```

### 数据库操作

#### 自动迁移

项目使用GORM的AutoMigrate功能,在应用启动时自动创建/更新表结构:

```go
db.AutoMigrate(&model.User{}, &model.Message{})
```

#### 手动迁移

如果需要更复杂的迁移,可以在 `migrations/` 目录下创建SQL文件。

### WebSocket开发

#### 消息类型

在 `internal/websocket/message.go` 中定义新的消息类型:

```go
const (
    MessageTypeCustom MessageType = "custom"
)
```

#### 处理消息

在 `internal/websocket/hub.go` 的 `handleBroadcast` 方法中添加处理逻辑:

```go
case MessageTypeCustom:
    // 处理自定义消息
```

### 配置管理

配置优先级: 环境变量 > 配置文件

#### 添加新配置项

1. 在 `pkg/utils/config.go` 中添加字段:

```go
type Config struct {
    NewField string `yaml:"new_field"`
}
```

2. 在 `config/config.yaml` 中添加默认值:

```yaml
new_field: "default_value"
```

3. 支持环境变量覆盖:

```go
if val := os.Getenv("NEW_FIELD"); val != "" {
    config.NewField = val
}
```

### 测试

#### 单元测试

```bash
go test ./...
```

#### 测试特定包

```bash
go test ./internal/service
```

#### 测试覆盖率

```bash
go test -cover ./...
```

### 性能优化

#### 数据库连接池

在 `internal/repository/database.go` 中配置:

```go
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

#### Redis缓存

使用Redis缓存热点数据,减少数据库查询。

#### 消息队列

对于高并发场景,可以使用Redis Pub/Sub或NATS作为消息队列。

## 常见问题

### 1. 数据库连接失败

检查配置文件中的数据库连接信息,确保PostgreSQL服务正在运行。

```bash
docker-compose ps
```

### 2. WebSocket连接失败

检查防火墙设置,确保端口8080可访问。在生产环境建议使用wss(WebSocket over TLS)。

### 3. JWT token过期

在配置文件中调整token过期时间:

```yaml
jwt:
  expire_time: 168  # 小时
```

### 4. 跨域问题

在 `internal/middleware/cors.go` 中配置允许的源:

```go
AllowOrigins: []string{"https://your-domain.com"}
```

## 性能指标

### 预期性能

- 并发连接: 10,000+
- 消息延迟: < 100ms
- 吞吐量: 10,000+ msg/s

### 压力测试

使用工具如Apache Bench、wrk或自定义脚本进行压力测试。

## 监控

建议集成以下监控工具:

- Prometheus: 指标收集
- Grafana: 可视化
- ELK Stack: 日志分析

## 安全建议

1. 使用HTTPS/WSS
2. 定期更新依赖
3. 使用强密码策略
4. 限制API请求频率
5. 输入验证和过滤
6. SQL注入防护(GORM已内置)
7. XSS防护
