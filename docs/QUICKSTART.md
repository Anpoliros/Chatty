# Chatty 快速开始指南

## 5分钟快速体验

### 前置要求

确保已安装:
- Go 1.21+ ([安装指南](https://go.dev/doc/install))
- Docker & Docker Compose ([安装指南](https://docs.docker.com/get-docker/))

### 第一步: 克隆项目

```bash
git clone https://github.com/Anpoliros/Chatty.git
cd Chatty/server
```

### 第二步: 启动数据库服务

```bash
# 使用Docker Compose启动PostgreSQL和Redis
docker-compose up -d

# 验证服务状态
docker-compose ps
```

你应该看到:
```
NAME                IMAGE               STATUS
chatty-postgres     postgres:15-alpine  Up
chatty-redis        redis:7-alpine      Up
```

### 第三步: 配置环境

```bash
# 复制配置文件(可选,使用默认配置即可)
cp .env.example .env
```

### 第四步: 启动服务器

```bash
# 方式1: 使用Makefile
make run

# 方式2: 直接运行
go run cmd/server/main.go
```

服务器启动成功后,你会看到:
```
[INFO] Starting Chatty Server...
[INFO] Server is running on :8080
```

### 第五步: 测试API

打开新的终端窗口:

#### 1. 健康检查

```bash
curl http://localhost:8080/health
```

响应:
```json
{
  "status": "ok",
  "time": 1234567890
}
```

#### 2. 注册用户

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "password123",
    "nickname": "Alice"
  }'
```

#### 3. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "password123"
  }'
```

保存返回的`token`,后续请求需要使用。

#### 4. 测试WebSocket连接

你可以使用浏览器的开发者工具或WebSocket客户端:

```javascript
// 浏览器控制台
const ws = new WebSocket('ws://localhost:8080/api/v1/ws?user_id=1');

ws.onopen = () => {
  console.log('Connected!');

  // 发送消息
  ws.send(JSON.stringify({
    type: 'private',
    receiver_id: 2,
    conversation_id: 1,
    content: 'Hello!'
  }));
};

ws.onmessage = (event) => {
  console.log('Received:', JSON.parse(event.data));
};
```

## 常用命令

### 使用Makefile

```bash
# 查看所有可用命令
make help

# 编译项目
make build

# 运行项目
make run

# 开发模式(热重载)
make dev

# 运行测试
make test

# 启动Docker服务
make docker-up

# 停止Docker服务
make docker-down

# 查看Docker日志
make docker-logs

# 格式化代码
make fmt

# 清理构建文件
make clean
```

### 直接使用Go命令

```bash
# 运行
go run cmd/server/main.go

# 编译
go build -o bin/chatty-server cmd/server/main.go

# 测试
go test ./...

# 下载依赖
go mod download
go mod tidy
```

## 开发工具推荐

### 1. 安装Air(热重载工具)

```bash
go install github.com/air-verse/air@latest

# 使用Air运行
make dev
# 或
air
```

### 2. API测试工具

推荐使用以下工具之一:
- [Postman](https://www.postman.com/)
- [Insomnia](https://insomnia.rest/)
- [HTTPie](https://httpie.io/)
- curl

### 3. WebSocket测试工具

- [Postman WebSocket](https://www.postman.com/)
- [WebSocket King](https://websocketking.com/)
- 浏览器开发者工具

## 目录导航

```
server/
├── cmd/server/main.go          # 从这里开始
├── internal/
│   ├── api/                    # HTTP处理器
│   ├── service/                # 业务逻辑
│   ├── repository/             # 数据访问
│   ├── model/                  # 数据模型
│   ├── websocket/              # WebSocket
│   └── middleware/             # 中间件
├── config/config.yaml          # 配置文件
└── Makefile                    # 构建脚本
```

## 开发流程

### 1. 修改代码

编辑你需要修改的文件,例如:
- `internal/api/*.go` - 添加新的API端点
- `internal/service/*.go` - 添加业务逻辑
- `internal/model/*.go` - 添加数据模型

### 2. 测试修改

```bash
# 运行测试
go test ./...

# 或使用热重载
make dev
```

### 3. 格式化代码

```bash
make fmt
```

## 常见问题

### Q: 数据库连接失败

**A:** 确保Docker服务正在运行:
```bash
docker-compose ps
# 如果没有运行,启动它:
docker-compose up -d
```

### Q: 端口8080已被占用

**A:** 修改配置文件或使用环境变量:
```bash
PORT=8081 go run cmd/server/main.go
```

或修改 `config/config.yaml`:
```yaml
server:
  port: "8081"
```

### Q: go mod下载很慢

**A:** 使用Go模块代理:
```bash
export GOPROXY=https://goproxy.cn,direct
go mod download
```

### Q: 如何查看数据库数据?

**A:** 连接PostgreSQL:
```bash
# 使用docker
docker exec -it chatty-postgres psql -U chatty -d chatty_db

# 或使用GUI工具如DBeaver, pgAdmin
```

### Q: 如何重置数据库?

**A:** 删除并重新创建:
```bash
docker-compose down -v
docker-compose up -d
```

## 下一步

现在你已经成功运行了Chatty服务器!

接下来可以:
1. 📖 阅读 [API文档](API.md)
2. 🏗️ 查看 [项目状态](PROJECT_STATUS.md)
3. 💻 开始开发iOS客户端
4. 🤝 为项目做贡献

## 获取帮助

- 📝 查看文档: `docs/`目录
- 🐛 报告问题: [GitHub Issues](https://github.com/Anpoliros/Chatty/issues)
- 💬 讨论: [GitHub Discussions](https://github.com/Anpoliros/Chatty/discussions)

祝你使用愉快! 🎉
