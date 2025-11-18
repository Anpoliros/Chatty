# Chatty - 开源聊天服务框架

一个现代化的、可扩展的聊天服务框架,包含服务端和客户端实现。

## 项目特性

- 🚀 高性能的WebSocket实时通信
- 💬 支持单聊和群聊
- 🔐 JWT身份认证
- 📱 iOS客户端支持
- 🗄️ PostgreSQL + Redis架构
- 🐳 Docker支持,快速部署
- 📝 完整的消息历史记录
- ✅ 已读/未读状态追踪
- 🌐 RESTful API设计

## 技术栈

### 服务端
- **语言**: Go 1.21+
- **Web框架**: Gin
- **WebSocket**: gorilla/websocket
- **数据库**: PostgreSQL 15+
- **缓存**: Redis 7+
- **ORM**: GORM

### iOS客户端
- **语言**: Swift
- **UI框架**: SwiftUI / UIKit
- **WebSocket**: Starscream
- **本地存储**: GRDB / Realm
- **网络**: URLSession + Combine

## 项目结构

```
Chatty/
├── server/                 # 服务端
│   ├── cmd/
│   │   └── server/        # 主入口
│   ├── internal/
│   │   ├── api/           # HTTP API处理
│   │   ├── websocket/     # WebSocket处理
│   │   ├── service/       # 业务逻辑
│   │   ├── repository/    # 数据访问
│   │   ├── model/         # 数据模型
│   │   └── middleware/    # 中间件
│   ├── pkg/               # 公共库
│   │   ├── logger/        # 日志
│   │   └── utils/         # 工具函数
│   ├── config/            # 配置文件
│   └── migrations/        # 数据库迁移
│
├── client/                # iOS客户端
│   └── ChatApp/
│       ├── Models/        # 数据模型
│       ├── Views/         # UI视图
│       ├── ViewModels/    # 视图模型
│       ├── Services/      # 服务层
│       └── Utils/         # 工具类
│
└── docs/                  # 文档
```

## 快速开始

### 服务端

#### 1. 环境要求

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (可选)

#### 2. 使用Docker启动数据库

```bash
cd server
docker-compose up -d
```

这将启动PostgreSQL和Redis容器。

#### 3. 配置

复制配置文件示例:

```bash
cp .env.example .env
cp config/config.yaml config/config.local.yaml
```

根据需要修改配置文件。

#### 4. 安装依赖

```bash
go mod download
```

#### 5. 运行服务

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动。

### API文档

#### 认证接口

**注册**
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "nickname": "Test User"
}
```

**登录**
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

**登出**
```
POST /api/v1/auth/logout
Authorization: Bearer <token>
```

#### WebSocket连接

```
WS /api/v1/ws?user_id=<user_id>
Authorization: Bearer <token>
```

#### 消息接口

**获取会话列表**
```
GET /api/v1/conversations
Authorization: Bearer <token>
```

**获取消息列表**
```
GET /api/v1/messages/:conversationId?page=1&page_size=50
Authorization: Bearer <token>
```

**发送消息**
```
POST /api/v1/messages
Authorization: Bearer <token>
Content-Type: application/json

{
  "conversation_id": 1,
  "type": "text",
  "content": "Hello!"
}
```

## WebSocket消息协议

### 消息类型

- `private`: 私聊消息
- `group`: 群聊消息
- `typing`: 正在输入
- `read`: 已读回执
- `system`: 系统消息

### 消息格式

```json
{
  "type": "private",
  "sender_id": 1,
  "receiver_id": 2,
  "conversation_id": 1,
  "content": "Hello!",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 开发指南

### 数据库迁移

服务启动时会自动执行数据库迁移(Auto Migrate)。如果需要手动迁移,可以使用GORM的迁移功能。

### 日志

使用内置的logger包:

```go
import "github.com/Anpoliros/Chatty/server/pkg/logger"

log := logger.New()
log.Info("Server started")
log.Error("Error occurred:", err)
```

### 添加新的API端点

1. 在 `internal/model` 定义数据模型
2. 在 `internal/repository` 实现数据访问
3. 在 `internal/service` 实现业务逻辑
4. 在 `internal/api` 实现HTTP处理
5. 在 `cmd/server/main.go` 注册路由

## 部署

### Docker部署

```bash
# 构建镜像
docker build -t chatty-server .

# 运行容器
docker run -d -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-password \
  chatty-server
```

### 环境变量

- `PORT`: 服务端口(默认:8080)
- `DB_HOST`: 数据库主机
- `DB_PORT`: 数据库端口
- `DB_USER`: 数据库用户
- `DB_PASSWORD`: 数据库密码
- `DB_NAME`: 数据库名称
- `REDIS_HOST`: Redis主机
- `REDIS_PORT`: Redis端口
- `JWT_SECRET`: JWT密钥

## 路线图

### MVP版本 (v0.1)
- [x] 用户注册/登录
- [x] WebSocket连接管理
- [x] 单聊功能
- [x] 消息历史记录
- [ ] 在线状态同步
- [ ] 消息已读/未读
- [ ] 推送通知

### 进阶功能 (v0.2+)
- [ ] 群聊
- [ ] 语音/视频消息
- [ ] 文件传输
- [ ] 消息撤回/删除
- [ ] @提醒
- [ ] 端到端加密

## 贡献

欢迎贡献代码!请遵循以下步骤:

1. Fork本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

- 项目链接: https://github.com/Anpoliros/Chatty
- 问题反馈: https://github.com/Anpoliros/Chatty/issues

## 致谢

- [Gin](https://github.com/gin-gonic/gin) - Web框架
- [GORM](https://gorm.io/) - ORM库
- [gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket实现
