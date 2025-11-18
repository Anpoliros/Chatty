import Foundation
import Combine

/// 聊天ViewModel
class ChatViewModel: ObservableObject {
    @Published var messages: [Message] = []
    @Published var isLoading = false
    @Published var errorMessage: String?
    @Published var isSending = false

    private var cancellables = Set<AnyCancellable>()
    private let apiClient = APIClient.shared
    private let storage = UserDefaultsManager.shared
    private let wsManager = WebSocketManager.shared

    let conversation: Conversation

    init(conversation: Conversation) {
        self.conversation = conversation
        loadMessages()
        subscribeToWebSocket()
    }

    // MARK: - Load Messages

    /// 加载消息列表
    func loadMessages(page: Int = 1) {
        isLoading = true
        errorMessage = nil

        let endpoint = Constants.API.Endpoint.messages(conversationId: conversation.id)
        let parameters = [
            "page": "\(page)",
            "page_size": "\(Constants.Pagination.defaultPageSize)"
        ]

        apiClient.get(endpoint, parameters: parameters)
            .receive(on: DispatchQueue.main)
            .sink { [weak self] completion in
                self?.isLoading = false
                if case .failure(let error) = completion {
                    self?.errorMessage = error.localizedDescription
                }
            } receiveValue: { [weak self] (response: MessagesResponse) in
                // 消息按时间倒序返回,需要反转
                self?.messages = response.messages.reversed()
            }
            .store(in: &cancellables)
    }

    // MARK: - Send Message

    /// 发送文本消息
    func sendMessage(content: String) {
        guard !content.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
            return
        }

        isSending = true

        let request = SendMessageRequest(
            conversationId: conversation.id,
            type: .text,
            content: content,
            mediaUrl: nil
        )

        apiClient.post(Constants.API.Endpoint.sendMessage, body: request)
            .receive(on: DispatchQueue.main)
            .sink { [weak self] completion in
                self?.isSending = false
                if case .failure(let error) = completion {
                    self?.errorMessage = error.localizedDescription
                }
            } receiveValue: { [weak self] (response: MessageResponse) in
                // 添加到消息列表
                self?.messages.append(response.data)

                // TODO: 通过WebSocket发送
                // 暂时通过HTTP已经发送成功
            }
            .store(in: &cancellables)
    }

    // MARK: - Mark as Read

    /// 标记消息为已读
    func markAsRead() {
        struct MarkAsReadRequest: Codable {
            let conversationId: UInt
            let messageId: UInt?

            enum CodingKeys: String, CodingKey {
                case conversationId = "conversation_id"
                case messageId = "message_id"
            }
        }

        let request = MarkAsReadRequest(
            conversationId: conversation.id,
            messageId: messages.last?.id
        )

        apiClient.post(Constants.API.Endpoint.markAsRead, body: request)
            .receive(on: DispatchQueue.main)
            .sink { _ in } receiveValue: { (_: [String: String]) in }
            .store(in: &cancellables)
    }

    // MARK: - WebSocket

    /// 订阅WebSocket消息
    private func subscribeToWebSocket() {
        wsManager.messagePublisher
            .receive(on: DispatchQueue.main)
            .sink { [weak self] wsMessage in
                guard let self = self else { return }

                // 如果是当前会话的消息,添加到列表
                if wsMessage.conversationId == self.conversation.id {
                    // TODO: 将WSMessage转换为Message并添加
                    print("Received WebSocket message for conversation \(self.conversation.id)")
                }
            }
            .store(in: &cancellables)
    }

    // MARK: - Helper Methods

    /// 获取当前用户ID
    var currentUserId: UInt {
        storage.currentUser?.id ?? 0
    }
}
