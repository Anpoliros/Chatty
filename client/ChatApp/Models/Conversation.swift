import Foundation

/// 会话类型
enum ConversationType: String, Codable {
    case `private`
    case group
}

/// 会话成员
struct ConversationMember: Codable, Identifiable {
    let id: UInt
    let userId: UInt
    let unreadCount: Int
    let joinedAt: Date
    let user: User?

    enum CodingKeys: String, CodingKey {
        case id
        case userId = "user_id"
        case unreadCount = "unread_count"
        case joinedAt = "joined_at"
        case user
    }
}

/// 会话模型
struct Conversation: Codable, Identifiable {
    let id: UInt
    let type: ConversationType
    let name: String?
    let avatar: String?
    let lastMessageId: UInt?
    let createdAt: Date
    let updatedAt: Date
    let lastMessage: Message?
    let participants: [ConversationMember]?

    enum CodingKeys: String, CodingKey {
        case id
        case type
        case name
        case avatar
        case lastMessageId = "last_message_id"
        case createdAt = "created_at"
        case updatedAt = "updated_at"
        case lastMessage = "last_message"
        case participants
    }

    /// 获取会话显示名称
    func displayName(currentUserId: UInt) -> String {
        if type == .group {
            return name ?? "Group Chat"
        }

        // 私聊:显示对方的名称
        if let otherUser = participants?.first(where: { $0.userId != currentUserId })?.user {
            return otherUser.displayName
        }

        return "Unknown"
    }

    /// 获取会话头像
    func displayAvatar(currentUserId: UInt) -> String? {
        if type == .group {
            return avatar
        }

        // 私聊:显示对方的头像
        if let otherUser = participants?.first(where: { $0.userId != currentUserId })?.user {
            return otherUser.avatar
        }

        return nil
    }

    /// 获取未读消息数
    func unreadCount(currentUserId: UInt) -> Int {
        if let member = participants?.first(where: { $0.userId == currentUserId }) {
            return member.unreadCount
        }
        return 0
    }

    /// 最后一条消息预览
    var lastMessagePreview: String {
        guard let lastMsg = lastMessage else {
            return "No messages yet"
        }

        switch lastMsg.type {
        case .text:
            return lastMsg.content
        case .image:
            return "[Image]"
        case .file:
            return "[File]"
        case .audio:
            return "[Audio]"
        case .video:
            return "[Video]"
        }
    }

    /// 最后消息时间
    var lastMessageTime: String {
        guard let lastMsg = lastMessage else {
            return ""
        }
        return lastMsg.timeString
    }
}

/// 创建会话请求
struct CreateConversationRequest: Codable {
    let type: ConversationType
    let participantId: UInt?
    let name: String?
    let memberIds: [UInt]?

    enum CodingKeys: String, CodingKey {
        case type
        case participantId = "participant_id"
        case name
        case memberIds = "member_ids"
    }

    /// 创建私聊会话
    static func privateConversation(participantId: UInt) -> CreateConversationRequest {
        return CreateConversationRequest(
            type: .private,
            participantId: participantId,
            name: nil,
            memberIds: nil
        )
    }

    /// 创建群聊会话
    static func groupConversation(name: String, memberIds: [UInt]) -> CreateConversationRequest {
        return CreateConversationRequest(
            type: .group,
            participantId: nil,
            name: name,
            memberIds: memberIds
        )
    }
}

/// 会话响应
struct ConversationResponse: Codable {
    let message: String
    let conversation: Conversation
}

/// 会话列表响应
struct ConversationsResponse: Codable {
    let conversations: [Conversation]
}
