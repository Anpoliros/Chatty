# Chatty iOS App - Xcode Setup Guide

本指南将帮助你在Xcode中快速设置和运行Chatty iOS应用。

## 方法一: 自动设置（推荐）⚡️

### 1. 运行setup脚本

```bash
cd client
chmod +x setup_xcode.sh
./setup_xcode.sh
```

脚本会自动：
- 创建Xcode项目
- 组织源文件
- 配置项目设置
- 添加必要的框架

### 2. 打开项目

```bash
open ChatApp/ChatApp.xcodeproj
```

### 3. 运行应用

- 选择模拟器或真机
- 点击运行按钮 (⌘R)

## 方法二: 手动设置 🛠

### 步骤1: 创建新项目

1. 打开Xcode
2. 选择 `File` → `New` → `Project`
3. 选择 `iOS` → `App`
4. 填写项目信息：
   - Product Name: `ChatApp`
   - Team: 选择你的开发团队
   - Organization Identifier: `com.yourname` (或你的标识符)
   - Interface: `SwiftUI`
   - Language: `Swift`
   - Storage: `None`
   - 取消勾选 `Include Tests`
5. 选择保存位置为 `client` 目录
6. 点击 `Create`

### 步骤2: 删除默认文件

删除Xcode自动生成的以下文件(保留项目结构)：
- `ContentView.swift`
- `ChatAppApp.swift` (Xcode生成的)

### 步骤3: 导入源文件

在Xcode中：

1. **右键点击 `ChatApp` 组** → `Add Files to "ChatApp"`

2. **选择以下目录**（确保勾选 "Create groups"）：
   - `ChatApp/Models`
   - `ChatApp/Views`
   - `ChatApp/ViewModels`
   - `ChatApp/Services`
   - `ChatApp/Utils`
   - `ChatApp/ChatApp.swift`

3. **确保勾选**：
   - ✅ Copy items if needed
   - ✅ Create groups
   - ✅ Add to targets: ChatApp

### 步骤4: 配置Info.plist

1. 在项目导航器中选择 `Info.plist`
2. 删除Xcode自动生成的内容
3. 右键点击 → `Open As` → `Source Code`
4. 将 `client/ChatApp/Info.plist` 的内容复制粘贴进去

或者直接替换：
```bash
cp ChatApp/Info.plist ChatApp/ChatApp/Info.plist
```

### 步骤5: 添加依赖包

Chatty需要Starscream库用于WebSocket通信：

1. 在Xcode中选择项目 (最顶部的蓝色图标)
2. 选择 `ChatApp` target
3. 点击 `Package Dependencies` 标签
4. 点击 `+` 按钮
5. 输入: `https://github.com/daltoniam/Starscream`
6. 点击 `Add Package`
7. 确认添加 `Starscream` 到你的app target

### 步骤6: 配置项目设置

1. 选择项目 → ChatApp target → `Signing & Capabilities`
2. 选择你的 `Team`
3. 确认 `Bundle Identifier` (如: `com.yourname.ChatApp`)

### 步骤7: 配置服务器地址

打开 `Utils/Constants.swift`，修改服务器地址：

```swift
enum API {
    // 如果使用模拟器，使用localhost
    static let baseURL = "http://localhost:8080/api/v1"
    static let wsURL = "ws://localhost:8080/api/v1/ws"

    // 如果使用真机，使用Mac的局域网IP
    // static let baseURL = "http://192.168.1.100:8080/api/v1"
    // static let wsURL = "ws://192.168.1.100:8080/api/v1/ws"
}
```

### 步骤8: 构建和运行

1. 选择模拟器（如 iPhone 15 Pro）
2. 按 `⌘R` 或点击运行按钮
3. 等待构建完成

## 项目结构验证

打开项目后，你应该看到如下结构：

```
ChatApp
├── ChatApp.swift                 # App入口
├── Info.plist                   # 配置文件
├── Models/                      # 数据模型
│   ├── User.swift
│   ├── Message.swift
│   └── Conversation.swift
├── Views/                       # UI视图
│   ├── Auth/
│   │   ├── LoginView.swift
│   │   └── RegisterView.swift
│   ├── Conversation/
│   │   └── ConversationListView.swift
│   └── Chat/
│       └── ChatView.swift
├── ViewModels/                  # 视图模型
│   ├── AuthViewModel.swift
│   ├── ConversationViewModel.swift
│   └── ChatViewModel.swift
├── Services/                    # 服务层
│   ├── Network/
│   │   └── APIClient.swift
│   ├── Storage/
│   │   └── UserDefaultsManager.swift
│   └── WebSocket/
│       └── WebSocketManager.swift
└── Utils/                       # 工具类
    ├── Constants.swift
    └── Extensions.swift
```

## 运行前检查清单

- [ ] 服务器正在运行 (`cd server && make run`)
- [ ] 数据库已启动 (`docker-compose up -d`)
- [ ] 已添加Starscream依赖包
- [ ] 已配置正确的服务器地址
- [ ] 已选择合适的模拟器或设备
- [ ] 项目可以成功构建 (⌘B)

## 真机测试

如果要在真机上测试：

### 1. 获取Mac的局域网IP

```bash
ifconfig | grep "inet " | grep -v 127.0.0.1
```

### 2. 修改Constants.swift

将 `localhost` 替换为你的Mac IP地址：

```swift
static let baseURL = "http://192.168.1.100:8080/api/v1"  // 使用你的IP
static let wsURL = "ws://192.168.1.100:8080/api/v1/ws"
```

### 3. 确保设备和Mac在同一网络

- iPhone连接到与Mac相同的WiFi网络
- 关闭Mac的防火墙或允许端口8080

### 4. 配置签名

- 在Xcode中选择你的Apple ID作为开发团队
- 选择你的设备
- 点击运行

## 常见问题

### Q: 编译错误 "Cannot find type 'Starscream' in scope"

**A:** 确保已添加Starscream包依赖。重启Xcode后重试。

### Q: 运行时错误 "Connection refused"

**A:**
1. 确认服务器正在运行
2. 检查服务器地址配置是否正确
3. 如果使用真机,检查IP地址和网络连接

### Q: "App Transport Security" 错误

**A:** 已在Info.plist中配置允许HTTP连接。如果仍有问题,检查Info.plist是否正确导入。

### Q: WebSocket连接失败

**A:** WebSocketManager需要Starscream库。确保：
1. 已添加Starscream依赖
2. 服务器WebSocket服务正常
3. 在WebSocketManager.swift中取消WebSocket实现的注释

### Q: 找不到某些Swift文件

**A:** 确保在添加文件时勾选了：
- ✅ Copy items if needed
- ✅ Add to targets: ChatApp

## 开发技巧

### 实时预览

SwiftUI支持实时预览。在任何View文件中按 `⌘⌥↩` 启动预览。

### 调试

- 设置断点: 点击代码行号
- 查看变量: 鼠标悬停在变量上
- 查看控制台: `⌘⇧Y`
- 查看网络请求: 在console中查看日志

### 快捷键

- `⌘R` - 运行
- `⌘B` - 构建
- `⌘.` - 停止
- `⌘⇧K` - 清理构建
- `⌘⇧Y` - 显示/隐藏控制台

## 下一步

1. ✅ 运行服务器
2. ✅ 启动iOS应用
3. 📱 注册账号
4. 💬 开始聊天!

## 获取帮助

- 查看 `README.md` 了解项目概述
- 查看 `docs/API.md` 了解API文档
- 查看 `docs/IMPLEMENTATION_SUMMARY.md` 了解实现细节

## 截图

运行成功后你将看到：
1. 登录界面 - 输入用户名和密码
2. 会话列表 - 查看所有聊天
3. 聊天界面 - 发送和接收消息

祝你使用愉快! 🎉
