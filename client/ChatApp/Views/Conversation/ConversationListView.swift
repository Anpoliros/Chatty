import SwiftUI

/// 会话列表视图
struct ConversationListView: View {
    @StateObject private var viewModel = ConversationViewModel()
    @StateObject private var authViewModel = AuthViewModel()
    @State private var showingNewChat = false

    var body: some View {
        NavigationView {
            ZStack {
                if viewModel.isLoading && viewModel.conversations.isEmpty {
                    ProgressView("Loading...")
                } else if viewModel.conversations.isEmpty {
                    emptyView
                } else {
                    conversationList
                }
            }
            .navigationTitle("Chats")
            .navigationBarItems(
                leading: Button(action: {
                    authViewModel.logout()
                }) {
                    Image(systemName: "arrow.left.square")
                },
                trailing: Button(action: {
                    showingNewChat = true
                }) {
                    Image(systemName: "square.and.pencil")
                }
            )
            .sheet(isPresented: $showingNewChat) {
                NewChatView { conversation in
                    // 跳转到新创建的会话
                    showingNewChat = false
                }
            }
            .refreshable {
                viewModel.loadConversations()
            }
        }
    }

    // MARK: - Conversation List

    private var conversationList: some View {
        List(viewModel.conversations) { conversation in
            NavigationLink(destination: ChatView(conversation: conversation)) {
                ConversationRow(
                    conversation: conversation,
                    currentUserId: viewModel.currentUserId
                )
            }
        }
        .listStyle(PlainListStyle())
    }

    // MARK: - Empty View

    private var emptyView: some View {
        VStack(spacing: 20) {
            Image(systemName: "message")
                .resizable()
                .frame(width: 80, height: 80)
                .foregroundColor(.gray.opacity(0.5))

            Text("No conversations yet")
                .font(.headline)
                .foregroundColor(.gray)

            Text("Start a new chat to begin messaging")
                .font(.subheadline)
                .foregroundColor(.gray)
                .multilineTextAlignment(.center)

            Button(action: { showingNewChat = true }) {
                Label("New Chat", systemImage: "square.and.pencil")
                    .padding()
                    .background(Color.blue)
                    .foregroundColor(.white)
                    .cornerRadius(10)
            }
        }
        .padding()
    }
}

// MARK: - Conversation Row

struct ConversationRow: View {
    let conversation: Conversation
    let currentUserId: UInt

    var body: some View {
        HStack(spacing: 12) {
            // Avatar
            Circle()
                .fill(Color.blue.opacity(0.2))
                .frame(width: Constants.UI.avatarSize, height: Constants.UI.avatarSize)
                .overlay(
                    Text(conversation.displayName(currentUserId: currentUserId).prefix(1).uppercased())
                        .foregroundColor(.blue)
                        .fontWeight(.semibold)
                )

            // Content
            VStack(alignment: .leading, spacing: 4) {
                HStack {
                    Text(conversation.displayName(currentUserId: currentUserId))
                        .font(.headline)

                    Spacer()

                    Text(conversation.lastMessageTime)
                        .font(.caption)
                        .foregroundColor(.gray)
                }

                HStack {
                    Text(conversation.lastMessagePreview)
                        .font(.subheadline)
                        .foregroundColor(.gray)
                        .lineLimit(1)

                    Spacer()

                    // Unread badge
                    let unread = conversation.unreadCount(currentUserId: currentUserId)
                    if unread > 0 {
                        Text("\(unread)")
                            .font(.caption)
                            .foregroundColor(.white)
                            .padding(.horizontal, 8)
                            .padding(.vertical, 4)
                            .background(Color.blue)
                            .clipShape(Capsule())
                    }
                }
            }
        }
        .padding(.vertical, 4)
    }
}

// MARK: - New Chat View

struct NewChatView: View {
    @Environment(\.presentationMode) var presentationMode
    @StateObject private var conversationViewModel = ConversationViewModel()
    @State private var searchText = ""
    @State private var searchResults: [User] = []
    @State private var isSearching = false

    let onCreate: (Conversation) -> Void

    var body: some View {
        NavigationView {
            VStack {
                // Search Bar
                SearchBar(text: $searchText, onSearch: searchUsers)
                    .padding()

                if isSearching {
                    ProgressView("Searching...")
                } else if searchResults.isEmpty {
                    VStack(spacing: 10) {
                        Image(systemName: "magnifyingglass")
                            .font(.largeTitle)
                            .foregroundColor(.gray)
                        Text("Search for users to start chatting")
                            .foregroundColor(.gray)
                    }
                    .padding()
                } else {
                    List(searchResults) { user in
                        Button(action: {
                            createConversation(with: user)
                        }) {
                            HStack {
                                Circle()
                                    .fill(Color.blue.opacity(0.2))
                                    .frame(width: 40, height: 40)
                                    .overlay(
                                        Text(user.displayName.prefix(1).uppercased())
                                            .foregroundColor(.blue)
                                    )

                                VStack(alignment: .leading) {
                                    Text(user.displayName)
                                        .font(.headline)
                                    Text("@\(user.username)")
                                        .font(.caption)
                                        .foregroundColor(.gray)
                                }

                                Spacer()

                                Image(systemName: "chevron.right")
                                    .foregroundColor(.gray)
                            }
                        }
                    }
                }
            }
            .navigationTitle("New Chat")
            .navigationBarItems(trailing: Button("Cancel") {
                presentationMode.wrappedValue.dismiss()
            })
        }
    }

    private func searchUsers() {
        guard !searchText.isEmpty else { return }

        isSearching = true

        APIClient.shared.get(
            Constants.API.Endpoint.searchUsers,
            parameters: ["keyword": searchText]
        )
        .receive(on: DispatchQueue.main)
        .sink { completion in
            isSearching = false
        } receiveValue: { (response: [String: [User]]) in
            searchResults = response["users"] ?? []
        }
        .store(in: &conversationViewModel.cancellables)
    }

    private func createConversation(with user: User) {
        conversationViewModel.createPrivateConversation(withUser: user.id) { conversation in
            if let conversation = conversation {
                onCreate(conversation)
            }
        }
    }
}

// MARK: - Search Bar

struct SearchBar: View {
    @Binding var text: String
    let onSearch: () -> Void

    var body: some View {
        HStack {
            Image(systemName: "magnifyingglass")
                .foregroundColor(.gray)

            TextField("Search users...", text: $text, onCommit: onSearch)
                .textFieldStyle(PlainTextFieldStyle())

            if !text.isEmpty {
                Button(action: { text = "" }) {
                    Image(systemName: "xmark.circle.fill")
                        .foregroundColor(.gray)
                }
            }
        }
        .padding(8)
        .background(Color(.systemGray6))
        .cornerRadius(10)
    }
}

struct ConversationListView_Previews: PreviewProvider {
    static var previews: some View {
        ConversationListView()
    }
}
