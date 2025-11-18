import Foundation
import Combine

/// 认证ViewModel
class AuthViewModel: ObservableObject {
    @Published var isLoggedIn = false
    @Published var isLoading = false
    @Published var errorMessage: String?
    @Published var currentUser: User?

    private var cancellables = Set<AnyCancellable>()
    private let apiClient = APIClient.shared
    private let storage = UserDefaultsManager.shared

    init() {
        // 检查登录状态
        isLoggedIn = storage.isLoggedIn
        currentUser = storage.currentUser
    }

    // MARK: - Login

    /// 用户登录
    func login(username: String, password: String) {
        guard !username.isEmpty, !password.isEmpty else {
            errorMessage = "Please enter username and password"
            return
        }

        isLoading = true
        errorMessage = nil

        let request = UserLoginRequest(username: username, password: password)

        apiClient.post(Constants.API.Endpoint.login, body: request)
            .receive(on: DispatchQueue.main)
            .sink { [weak self] completion in
                self?.isLoading = false
                if case .failure(let error) = completion {
                    self?.errorMessage = error.localizedDescription
                }
            } receiveValue: { [weak self] (response: AuthResponse) in
                self?.handleAuthSuccess(response)
            }
            .store(in: &cancellables)
    }

    // MARK: - Register

    /// 用户注册
    func register(username: String, email: String, password: String, nickname: String?) {
        guard username.isValidUsername else {
            errorMessage = "Username must be 3-20 characters, letters, numbers and underscore only"
            return
        }

        guard email.isValidEmail else {
            errorMessage = "Invalid email format"
            return
        }

        guard password.isValidPassword else {
            errorMessage = "Password must be at least 6 characters"
            return
        }

        isLoading = true
        errorMessage = nil

        let request = UserRegisterRequest(
            username: username,
            email: email,
            password: password,
            nickname: nickname
        )

        apiClient.post(Constants.API.Endpoint.register, body: request)
            .receive(on: DispatchQueue.main)
            .sink { [weak self] completion in
                self?.isLoading = false
                if case .failure(let error) = completion {
                    self?.errorMessage = error.localizedDescription
                }
            } receiveValue: { [weak self] (response: AuthResponse) in
                self?.handleAuthSuccess(response)
            }
            .store(in: &cancellables)
    }

    // MARK: - Logout

    /// 用户登出
    func logout() {
        // 调用API登出
        apiClient.post(Constants.API.Endpoint.logout, body: ["": ""])
            .receive(on: DispatchQueue.main)
            .sink { _ in } receiveValue: { (_: [String: String]) in }
            .store(in: &cancellables)

        // 清除本地数据
        storage.clearAll()
        isLoggedIn = false
        currentUser = nil

        // 断开WebSocket
        WebSocketManager.shared.disconnect()
    }

    // MARK: - Private Methods

    private func handleAuthSuccess(_ response: AuthResponse) {
        // 保存token和用户信息
        storage.authToken = response.token
        storage.currentUser = response.user
        storage.isLoggedIn = true

        // 更新状态
        isLoggedIn = true
        currentUser = response.user

        // 连接WebSocket
        WebSocketManager.shared.connect(userId: response.user.id)
    }
}
