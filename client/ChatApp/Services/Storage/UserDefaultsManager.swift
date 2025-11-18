import Foundation

/// UserDefaults管理器
class UserDefaultsManager {
    static let shared = UserDefaultsManager()

    private let defaults = UserDefaults.standard
    private let encoder = JSONEncoder()
    private let decoder = JSONDecoder()

    private init() {}

    // MARK: - Auth Token

    var authToken: String? {
        get {
            defaults.string(forKey: Constants.UserDefaultsKey.authToken)
        }
        set {
            defaults.set(newValue, forKey: Constants.UserDefaultsKey.authToken)
        }
    }

    // MARK: - Current User

    var currentUser: User? {
        get {
            guard let data = defaults.data(forKey: Constants.UserDefaultsKey.currentUser) else {
                return nil
            }
            return try? decoder.decode(User.self, from: data)
        }
        set {
            if let user = newValue, let data = try? encoder.encode(user) {
                defaults.set(data, forKey: Constants.UserDefaultsKey.currentUser)
            } else {
                defaults.removeObject(forKey: Constants.UserDefaultsKey.currentUser)
            }
        }
    }

    // MARK: - Login Status

    var isLoggedIn: Bool {
        get {
            defaults.bool(forKey: Constants.UserDefaultsKey.isLoggedIn)
        }
        set {
            defaults.set(newValue, forKey: Constants.UserDefaultsKey.isLoggedIn)
        }
    }

    // MARK: - Clear All

    /// 清除所有数据(登出时调用)
    func clearAll() {
        authToken = nil
        currentUser = nil
        isLoggedIn = false
    }
}
