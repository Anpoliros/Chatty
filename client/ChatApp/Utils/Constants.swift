import Foundation

/// 全局常量
enum Constants {
    /// API配置
    enum API {
        // 服务器地址 - 根据实际情况修改
        static let baseURL = "http://localhost:8080/api/v1"
        static let wsURL = "ws://localhost:8080/api/v1/ws"

        // 超时设置
        static let timeout: TimeInterval = 30

        // API端点
        enum Endpoint {
            static let register = "/auth/register"
            static let login = "/auth/login"
            static let logout = "/auth/logout"
            static let me = "/users/me"
            static let searchUsers = "/users/search"
            static let updateProfile = "/users/profile"
            static let conversations = "/conversations"
            static func messages(conversationId: UInt) -> String {
                "/messages/\(conversationId)"
            }
            static let sendMessage = "/messages"
            static let markAsRead = "/messages/read"
        }
    }

    /// UserDefaults键
    enum UserDefaultsKey {
        static let authToken = "authToken"
        static let currentUser = "currentUser"
        static let isLoggedIn = "isLoggedIn"
    }

    /// UI配置
    enum UI {
        static let cornerRadius: CGFloat = 12
        static let padding: CGFloat = 16
        static let animationDuration: Double = 0.3

        // 消息气泡
        static let messageBubbleMaxWidth: CGFloat = 260
        static let messageBubblePadding: CGFloat = 12

        // 头像
        static let avatarSize: CGFloat = 40
        static let avatarCornerRadius: CGFloat = 20
    }

    /// 颜色配置
    enum Colors {
        // 使用系统颜色名称,在SwiftUI中可以通过Color直接访问
        static let primaryColor = "blue"
        static let secondaryColor = "gray"

        // 消息气泡颜色
        static let myMessageBubble = "blue"
        static let otherMessageBubble = "gray"
    }

    /// 分页配置
    enum Pagination {
        static let defaultPageSize = 50
        static let maxPageSize = 100
    }
}
