# Chatty API文档

## 基本信息

- Base URL: `http://localhost:8080/api/v1`
- 认证方式: JWT Bearer Token
- 内容类型: `application/json`

## 认证

### 注册

创建新用户账号。

**请求**

```
POST /auth/register
Content-Type: application/json
```

**请求体**

```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "nickname": "John Doe"
}
```

**参数说明**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名,3-20字符 |
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码,最少6字符 |
| nickname | string | 否 | 昵称,不填则使用username |

**响应**

```json
{
  "message": "Registration successful",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "nickname": "John Doe"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 登录

用户登录获取访问令牌。

**请求**

```
POST /auth/login
Content-Type: application/json
```

**请求体**

```json
{
  "username": "johndoe",
  "password": "password123"
}
```

**响应**

```json
{
  "message": "Login successful",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "nickname": "John Doe",
    "status": "online"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 登出

用户登出。

**请求**

```
POST /auth/logout
Authorization: Bearer <token>
```

**响应**

```json
{
  "message": "Logout successful"
}
```

## 用户

### 获取当前用户信息

**请求**

```
GET /users/me
Authorization: Bearer <token>
```

**响应**

```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "nickname": "John Doe",
  "avatar": "https://example.com/avatar.jpg",
  "status": "online",
  "last_seen_at": "2024-01-01T12:00:00Z",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### 搜索用户

**请求**

```
GET /users/search?keyword=john
Authorization: Bearer <token>
```

**参数说明**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 是 | 搜索关键词 |

**响应**

```json
{
  "users": [
    {
      "id": 1,
      "username": "johndoe",
      "nickname": "John Doe",
      "avatar": "https://example.com/avatar.jpg",
      "status": "online"
    }
  ]
}
```

### 更新用户资料

**请求**

```
PUT /users/profile
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**

```json
{
  "nickname": "New Nickname",
  "avatar": "https://example.com/new-avatar.jpg"
}
```

**响应**

```json
{
  "message": "Profile updated",
  "user": {
    "id": 1,
    "username": "johndoe",
    "nickname": "New Nickname",
    "avatar": "https://example.com/new-avatar.jpg"
  }
}
```

## 会话

### 获取会话列表

获取当前用户的所有会话。

**请求**

```
GET /conversations
Authorization: Bearer <token>
```

**响应**

```json
{
  "conversations": [
    {
      "id": 1,
      "type": "private",
      "name": "Jane Doe",
      "avatar": "https://example.com/avatar2.jpg",
      "last_message": {
        "id": 100,
        "content": "Hello!",
        "created_at": "2024-01-01T12:00:00Z",
        "sender": {
          "id": 2,
          "username": "janedoe",
          "nickname": "Jane Doe"
        }
      },
      "unread_count": 3,
      "updated_at": "2024-01-01T12:00:00Z"
    },
    {
      "id": 2,
      "type": "group",
      "name": "Project Team",
      "avatar": "https://example.com/group.jpg",
      "participants": [
        {
          "id": 1,
          "username": "johndoe"
        },
        {
          "id": 2,
          "username": "janedoe"
        }
      ],
      "last_message": {
        "id": 101,
        "content": "Meeting at 3pm",
        "created_at": "2024-01-01T13:00:00Z"
      },
      "unread_count": 0,
      "updated_at": "2024-01-01T13:00:00Z"
    }
  ]
}
```

### 创建会话

创建私聊或群聊会话。

**请求(私聊)**

```
POST /conversations
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体(私聊)**

```json
{
  "type": "private",
  "participant_id": 2
}
```

**请求体(群聊)**

```json
{
  "type": "group",
  "name": "Project Team",
  "member_ids": [2, 3, 4]
}
```

**响应**

```json
{
  "message": "Conversation created",
  "conversation": {
    "id": 3,
    "type": "private",
    "name": "Jane Doe",
    "participants": [
      {
        "id": 1,
        "user_id": 1
      },
      {
        "id": 2,
        "user_id": 2
      }
    ],
    "created_at": "2024-01-01T14:00:00Z"
  }
}
```

## 消息

### 获取消息列表

获取指定会话的消息历史。

**请求**

```
GET /messages/:conversationId?page=1&page_size=50
Authorization: Bearer <token>
```

**参数说明**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| conversationId | uint | 是 | - | 会话ID(路径参数) |
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 50 | 每页数量 |

**响应**

```json
{
  "messages": [
    {
      "id": 100,
      "conversation_id": 1,
      "sender_id": 2,
      "type": "text",
      "content": "Hello!",
      "is_read": true,
      "created_at": "2024-01-01T12:00:00Z",
      "sender": {
        "id": 2,
        "username": "janedoe",
        "nickname": "Jane Doe",
        "avatar": "https://example.com/avatar2.jpg"
      }
    },
    {
      "id": 101,
      "conversation_id": 1,
      "sender_id": 1,
      "type": "image",
      "content": "",
      "media_url": "https://example.com/image.jpg",
      "is_read": false,
      "created_at": "2024-01-01T12:01:00Z",
      "sender": {
        "id": 1,
        "username": "johndoe",
        "nickname": "John Doe"
      }
    }
  ],
  "page": 1,
  "page_size": 50,
  "total": 150
}
```

### 发送消息

通过HTTP发送消息(也可通过WebSocket)。

**请求**

```
POST /messages
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**

```json
{
  "conversation_id": 1,
  "type": "text",
  "content": "Hello, how are you?"
}
```

**消息类型**

- `text`: 文本消息
- `image`: 图片消息
- `file`: 文件消息
- `audio`: 语音消息
- `video`: 视频消息

**响应**

```json
{
  "message": "Message sent",
  "data": {
    "id": 102,
    "conversation_id": 1,
    "sender_id": 1,
    "type": "text",
    "content": "Hello, how are you?",
    "is_read": false,
    "created_at": "2024-01-01T12:30:00Z"
  }
}
```

### 标记已读

标记会话中的消息为已读。

**请求**

```
POST /messages/read
Authorization: Bearer <token>
Content-Type: application/json
```

**请求体**

```json
{
  "conversation_id": 1,
  "message_id": 102
}
```

**响应**

```json
{
  "message": "Marked as read"
}
```

### 删除消息

删除消息(软删除)。

**请求**

```
DELETE /messages/:id
Authorization: Bearer <token>
```

**响应**

```json
{
  "message": "Message deleted",
  "id": "102"
}
```

## WebSocket

### 连接

建立WebSocket连接。

**请求**

```
WS /ws?user_id=<user_id>
Upgrade: websocket
```

在生产环境中,应使用JWT token进行认证:

```
WS /ws
Authorization: Bearer <token>
```

### 消息格式

#### 发送消息

```json
{
  "type": "private",
  "receiver_id": 2,
  "conversation_id": 1,
  "content": "Hello via WebSocket!",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

#### 接收消息

```json
{
  "type": "private",
  "sender_id": 2,
  "conversation_id": 1,
  "content": "Hi there!",
  "message_id": 103,
  "timestamp": "2024-01-01T12:01:00Z"
}
```

#### 正在输入

```json
{
  "type": "typing",
  "sender_id": 2,
  "conversation_id": 1,
  "extra": {
    "is_typing": true
  }
}
```

#### 已读回执

```json
{
  "type": "read",
  "sender_id": 2,
  "conversation_id": 1,
  "extra": {
    "message_id": 100
  }
}
```

#### 系统消息

```json
{
  "type": "system",
  "content": "User joined the conversation",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## 错误响应

### 格式

```json
{
  "error": "Error message description"
}
```

### 常见错误码

| 状态码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未认证或token无效 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突(如用户名已存在) |
| 500 | 服务器内部错误 |

### 错误示例

```json
{
  "error": "Invalid username or password"
}
```

```json
{
  "error": "Authorization header is required"
}
```

```json
{
  "error": "Conversation not found"
}
```

## 速率限制

为防止滥用,API实施了速率限制:

- 认证接口: 5次/分钟
- 其他接口: 100次/分钟
- WebSocket消息: 50条/秒

超过限制将返回429状态码。

## 分页

支持分页的接口使用以下参数:

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码 |
| page_size | int | 50 | 每页数量 |

响应包含分页信息:

```json
{
  "data": [...],
  "page": 1,
  "page_size": 50,
  "total": 150
}
```

## 文件上传

### 上传头像

```
POST /upload/avatar
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <binary data>
```

### 上传聊天图片

```
POST /upload/image
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <binary data>
```

**响应**

```json
{
  "url": "https://example.com/uploads/image.jpg"
}
```
