import Foundation
import Combine

/// 会话列表ViewModel
class ConversationViewModel: ObservableObject {
    @Published var conversations: [Conversation] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private var cancellables = Set<AnyCancellable>()
    private let apiClient = APIClient.shared
    private let storage = UserDefaultsManager.shared

    init() {
        loadConversations()
    }

    // MARK: - Load Conversations

    /// 加载会话列表
    func loadConversations() {
        guard storage.isLoggedIn else { return }

        isLoading = true
        errorMessage = nil

        apiClient.get(Constants.API.Endpoint.conversations)
            .receive(on: DispatchQueue.main)
            .sink { [weak self] completion in
                self?.isLoading = false
                if case .failure(let error) = completion {
                    self?.errorMessage = error.localizedDescription
                }
            } receiveValue: { [weak self] (response: ConversationsResponse) in
                self?.conversations = response.conversations
            }
            .store(in: &cancellables)
    }

    // MARK: - Create Conversation

    /// 创建私聊会话
    func createPrivateConversation(withUser userId: UInt, completion: @escaping (Conversation?) -> Void) {
        let request = CreateConversationRequest.privateConversation(participantId: userId)

        apiClient.post(Constants.API.Endpoint.conversations, body: request)
            .receive(on: DispatchQueue.main)
            .sink { [weak self] completionResult in
                if case .failure(let error) = completionResult {
                    self?.errorMessage = error.localizedDescription
                    completion(nil)
                }
            } receiveValue: { [weak self] (response: ConversationResponse) in
                // 添加到列表
                self?.conversations.insert(response.conversation, at: 0)
                completion(response.conversation)
            }
            .store(in: &cancellables)
    }

    // MARK: - Helper Methods

    /// 获取当前用户ID
    var currentUserId: UInt {
        storage.currentUser?.id ?? 0
    }
}
