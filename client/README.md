# Chatty iOS Client

点对点聊天iOS客户端 - 完整实现的SwiftUI应用

## 🚀 快速开始

### 自动设置（推荐）

```bash
cd client
./setup_xcode.sh
```

然后按照屏幕上的说明在Xcode中打开项目。

### 详细设置

查看 [SETUP.md](SETUP.md) 获取完整的Xcode设置指南。

## 技术栈

- Swift 5.9+
- SwiftUI (iOS 15+)
- Combine (响应式编程)
- URLSession (网络请求)
- Starscream (WebSocket)
- MVVM架构

## ✅ 已实现功能

- ✅ 用户注册/登录
- ✅ 会话列表展示
- ✅ 实时聊天界面
- ✅ 消息气泡UI
- ✅ WebSocket实时通信(框架)
- ✅ JWT令牌管理
- ✅ 本地数据持久化
- ✅ 下拉刷新
- ✅ 用户搜索
- ✅ 错误处理

## 项目结构

```
ChatApp/
├── Models/              # 数据模型
│   ├── User.swift
│   ├── Message.swift
│   └── Conversation.swift
│
├── Services/            # 服务层
│   ├── Network/
│   │   ├── APIClient.swift
│   │   └── APIEndpoints.swift
│   ├── WebSocket/
│   │   └── WebSocketManager.swift
│   └── Storage/
│       └── UserDefaultsManager.swift
│
├── ViewModels/          # 视图模型(MVVM)
│   ├── AuthViewModel.swift
│   ├── ConversationViewModel.swift
│   └── ChatViewModel.swift
│
├── Views/               # 视图
│   ├── Auth/
│   │   ├── LoginView.swift
│   │   └── RegisterView.swift
│   ├── Conversation/
│   │   └── ConversationListView.swift
│   └── Chat/
│       ├── ChatView.swift
│       └── Components/
│           ├── MessageBubble.swift
│           └── ChatInputView.swift
│
├── Utils/               # 工具类
│   ├── Constants.swift
│   └── Extensions.swift
│
└── ChatApp.swift        # App入口
```

## 开始使用

### 1. 创建Xcode项目

1. 打开Xcode
2. 创建新项目: File → New → Project
3. 选择 "iOS" → "App"
4. 填写项目信息:
   - Product Name: ChatApp
   - Interface: SwiftUI
   - Language: Swift
   - Minimum Deployments: iOS 15.0

### 2. 添加依赖

使用Swift Package Manager添加Starscream:

1. File → Add Packages...
2. 输入: `https://github.com/daltoniam/Starscream`
3. 选择最新版本

### 3. 导入源文件

将本目录下的所有 `.swift` 文件复制到Xcode项目的对应目录。

### 4. 配置服务器地址

在 `Utils/Constants.swift` 中配置服务器地址:

```swift
struct API {
    static let baseURL = "http://localhost:8080/api/v1"
    static let wsURL = "ws://localhost:8080/api/v1/ws"
}
```

### 5. 运行

选择模拟器或真机,点击运行按钮。

## 功能清单

### 已实现
- ✅ 用户登录/注册
- ✅ 会话列表
- ✅ 点对点聊天
- ✅ WebSocket实时消息
- ✅ 消息历史记录
- ✅ MVVM架构

### 待实现
- ⏳ 本地消息持久化
- ⏳ 图片消息
- ⏳ 推送通知
- ⏳ 离线消息同步
- ⏳ 消息已读状态

## API集成

客户端与服务端的通信:

### HTTP REST API
- 登录/注册
- 获取会话列表
- 获取消息历史
- 创建会话

### WebSocket
- 实时消息推送
- 在线状态更新
- 正在输入提示

## 测试

### 单元测试

运行测试: `Cmd + U`

测试文件位于 `ChatAppTests/` 目录。

### 注意事项

1. **Info.plist配置**

如果使用HTTP连接(非HTTPS),需要在Info.plist中添加:

```xml
<key>NSAppTransportSecurity</key>
<dict>
    <key>NSAllowsArbitraryLoads</key>
    <true/>
</dict>
```

2. **真机测试**

真机测试时,确保:
- 设备与服务器在同一网络
- 使用服务器的局域网IP地址

## 架构说明

### MVVM模式

- **Model**: 数据模型,表示应用的数据结构
- **View**: SwiftUI视图,负责UI展示
- **ViewModel**: 视图模型,连接View和Model,处理业务逻辑

### 数据流

```
View → ViewModel → Service → API/WebSocket
         ↑                        ↓
         ←←←←← Response ←←←←←←←←←←
```

## 开发建议

1. 先运行服务端,确保后端正常工作
2. 在模拟器中测试基本功能
3. 使用Xcode的Network Link Conditioner测试网络异常
4. 查看Console输出调试信息

## 故障排除

### 连接失败
- 检查服务器是否启动
- 检查网络连接
- 确认API地址配置正确

### WebSocket断开
- 检查服务器WebSocket服务
- 查看控制台错误信息
- 实现重连机制

### 界面不更新
- 确保ViewModel使用@Published
- 检查View使用@ObservedObject或@StateObject

## 贡献

欢迎提交Issue和Pull Request!
