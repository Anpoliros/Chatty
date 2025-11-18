import Foundation
import Combine

/// API错误类型
enum APIError: Error, LocalizedError {
    case invalidURL
    case requestFailed(Error)
    case invalidResponse
    case decodingFailed(Error)
    case unauthorized
    case serverError(String)

    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Invalid URL"
        case .requestFailed(let error):
            return "Request failed: \(error.localizedDescription)"
        case .invalidResponse:
            return "Invalid response from server"
        case .decodingFailed(let error):
            return "Failed to decode response: \(error.localizedDescription)"
        case .unauthorized:
            return "Unauthorized. Please login again."
        case .serverError(let message):
            return message
        }
    }
}

/// API客户端
class APIClient {
    static let shared = APIClient()

    private let baseURL: String
    private let session: URLSession
    private let decoder: JSONDecoder
    private let encoder: JSONEncoder

    private init() {
        self.baseURL = Constants.API.baseURL

        let config = URLSessionConfiguration.default
        config.timeoutIntervalForRequest = Constants.API.timeout
        self.session = URLSession(configuration: config)

        self.decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        decoder.dateDecodingStrategy = .custom { decoder in
            let container = try decoder.singleValueContainer()
            let dateString = try container.decode(String.self)

            if let date = Date.iso8601Formatter.date(from: dateString) {
                return date
            }

            let formatter = DateFormatter()
            formatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ss.SSSZ"
            if let date = formatter.date(from: dateString) {
                return date
            }

            throw DecodingError.dataCorruptedError(in: container, debugDescription: "Cannot decode date string \(dateString)")
        }

        self.encoder = JSONEncoder()
        encoder.keyEncodingStrategy = .convertToSnakeCase
    }

    // MARK: - Request Methods

    /// GET请求
    func get<T: Decodable>(_ endpoint: String, parameters: [String: String]? = nil) -> AnyPublisher<T, APIError> {
        return request(endpoint, method: "GET", parameters: parameters)
    }

    /// POST请求
    func post<T: Encodable, R: Decodable>(_ endpoint: String, body: T) -> AnyPublisher<R, APIError> {
        return request(endpoint, method: "POST", body: body)
    }

    /// PUT请求
    func put<T: Encodable, R: Decodable>(_ endpoint: String, body: T) -> AnyPublisher<R, APIError> {
        return request(endpoint, method: "PUT", body: body)
    }

    /// DELETE请求
    func delete<R: Decodable>(_ endpoint: String) -> AnyPublisher<R, APIError> {
        return request(endpoint, method: "DELETE")
    }

    // MARK: - Private Methods

    private func request<R: Decodable>(_ endpoint: String, method: String, parameters: [String: String]? = nil, body: Encodable? = nil) -> AnyPublisher<R, APIError> {
        guard var urlComponents = URLComponents(string: baseURL + endpoint) else {
            return Fail(error: APIError.invalidURL).eraseToAnyPublisher()
        }

        // 添加查询参数
        if let parameters = parameters {
            urlComponents.queryItems = parameters.map { URLQueryItem(name: $0.key, value: $0.value) }
        }

        guard let url = urlComponents.url else {
            return Fail(error: APIError.invalidURL).eraseToAnyPublisher()
        }

        var request = URLRequest(url: url)
        request.httpMethod = method
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        // 添加认证token
        if let token = UserDefaultsManager.shared.authToken {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        // 添加请求体
        if let body = body {
            do {
                request.httpBody = try encoder.encode(body)
            } catch {
                return Fail(error: APIError.requestFailed(error)).eraseToAnyPublisher()
            }
        }

        return session.dataTaskPublisher(for: request)
            .tryMap { data, response in
                guard let httpResponse = response as? HTTPURLResponse else {
                    throw APIError.invalidResponse
                }

                // 检查状态码
                switch httpResponse.statusCode {
                case 200...299:
                    return data
                case 401:
                    throw APIError.unauthorized
                default:
                    // 尝试解析错误消息
                    if let errorResponse = try? self.decoder.decode([String: String].self, from: data),
                       let errorMessage = errorResponse["error"] {
                        throw APIError.serverError(errorMessage)
                    }
                    throw APIError.serverError("Server error: \(httpResponse.statusCode)")
                }
            }
            .decode(type: R.self, decoder: decoder)
            .mapError { error in
                if let apiError = error as? APIError {
                    return apiError
                } else if error is DecodingError {
                    return APIError.decodingFailed(error)
                } else {
                    return APIError.requestFailed(error)
                }
            }
            .eraseToAnyPublisher()
    }
}
