import Foundation

/// 消息类型
enum MessageType: String, Codable {
    case text
    case image
    case file
    case audio
    case video
}

/// 消息模型
struct Message: Codable, Identifiable, Equatable {
    let id: UInt
    let conversationId: UInt
    let senderId: UInt
    let type: MessageType
    let content: String
    let mediaUrl: String?
    let isRead: Bool
    let createdAt: Date
    let sender: User?

    enum CodingKeys: String, CodingKey {
        case id
        case conversationId = "conversation_id"
        case senderId = "sender_id"
        case type
        case content
        case mediaUrl = "media_url"
        case isRead = "is_read"
        case createdAt = "created_at"
        case sender
    }

    /// 是否是自己发送的消息
    func isMe(userId: UInt) -> Bool {
        senderId == userId
    }

    /// 格式化时间显示
    var timeString: String {
        let formatter = DateFormatter()
        let calendar = Calendar.current

        if calendar.isDateInToday(createdAt) {
            formatter.dateFormat = "HH:mm"
        } else if calendar.isDateInYesterday(createdAt) {
            return "Yesterday"
        } else {
            formatter.dateFormat = "MM/dd"
        }

        return formatter.string(from: createdAt)
    }
}

/// 发送消息请求
struct SendMessageRequest: Codable {
    let conversationId: UInt
    let type: MessageType
    let content: String
    let mediaUrl: String?

    enum CodingKeys: String, CodingKey {
        case conversationId = "conversation_id"
        case type
        case content
        case mediaUrl = "media_url"
    }
}

/// 消息响应
struct MessageResponse: Codable {
    let message: String
    let data: Message
}

/// 获取消息列表响应
struct MessagesResponse: Codable {
    let messages: [Message]
    let page: Int
    let pageSize: Int

    enum CodingKeys: String, CodingKey {
        case messages
        case page
        case pageSize = "page_size"
    }
}
