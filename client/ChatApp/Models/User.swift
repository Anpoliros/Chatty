import Foundation

/// 用户模型
struct User: Codable, Identifiable, Equatable {
    let id: UInt
    let username: String
    let email: String
    let nickname: String
    let avatar: String?
    let status: String  // online, offline, away
    let lastSeenAt: Date
    let createdAt: Date

    enum CodingKeys: String, CodingKey {
        case id
        case username
        case email
        case nickname
        case avatar
        case status
        case lastSeenAt = "last_seen_at"
        case createdAt = "created_at"
    }

    /// 用户显示名称
    var displayName: String {
        nickname.isEmpty ? username : nickname
    }

    /// 是否在线
    var isOnline: Bool {
        status == "online"
    }
}

/// 用户注册请求
struct UserRegisterRequest: Codable {
    let username: String
    let email: String
    let password: String
    let nickname: String?
}

/// 用户登录请求
struct UserLoginRequest: Codable {
    let username: String
    let password: String
}

/// 认证响应
struct AuthResponse: Codable {
    let message: String
    let user: User
    let token: String
}
