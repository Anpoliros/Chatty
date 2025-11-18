import SwiftUI

/// 聊天视图
struct ChatView: View {
    @StateObject private var viewModel: ChatViewModel
    @State private var messageText = ""
    @State private var scrollProxy: ScrollViewProxy?

    init(conversation: Conversation) {
        _viewModel = StateObject(wrappedValue: ChatViewModel(conversation: conversation))
    }

    var body: some View {
        VStack(spacing: 0) {
            // Messages List
            ScrollViewReader { proxy in
                ScrollView {
                    LazyVStack(spacing: 8) {
                        if viewModel.isLoading {
                            ProgressView()
                                .padding()
                        }

                        ForEach(viewModel.messages) { message in
                            MessageBubble(
                                message: message,
                                isFromCurrentUser: message.isMe(userId: viewModel.currentUserId)
                            )
                            .id(message.id)
                        }
                    }
                    .padding()
                }
                .onAppear {
                    scrollProxy = proxy
                    scrollToBottom()
                }
                .onChange(of: viewModel.messages.count) { _ in
                    scrollToBottom()
                }
            }

            Divider()

            // Input Area
            ChatInputView(
                text: $messageText,
                isSending: viewModel.isSending,
                onSend: sendMessage
            )
        }
        .navigationTitle(viewModel.conversation.displayName(currentUserId: viewModel.currentUserId))
        .navigationBarTitleDisplayMode(.inline)
        .onAppear {
            viewModel.markAsRead()
        }
    }

    private func sendMessage() {
        guard !messageText.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else {
            return
        }

        let text = messageText
        messageText = ""
        hideKeyboard()

        viewModel.sendMessage(content: text)
    }

    private func scrollToBottom() {
        guard let lastMessage = viewModel.messages.last else { return }
        DispatchQueue.main.asyncAfter(deadline: .now() + 0.1) {
            withAnimation {
                scrollProxy?.scrollTo(lastMessage.id, anchor: .bottom)
            }
        }
    }
}

// MARK: - Message Bubble

struct MessageBubble: View {
    let message: Message
    let isFromCurrentUser: Bool

    var body: some View {
        HStack {
            if isFromCurrentUser {
                Spacer()
            }

            VStack(alignment: isFromCurrentUser ? .trailing : .leading, spacing: 4) {
                // Sender name (for received messages)
                if !isFromCurrentUser, let sender = message.sender {
                    Text(sender.displayName)
                        .font(.caption)
                        .foregroundColor(.gray)
                        .padding(.leading, 12)
                }

                // Message content
                Text(message.content)
                    .padding(Constants.UI.messageBubblePadding)
                    .background(isFromCurrentUser ? Color.blue : Color(.systemGray5))
                    .foregroundColor(isFromCurrentUser ? .white : .primary)
                    .cornerRadius(Constants.UI.cornerRadius)

                // Time
                Text(message.timeString)
                    .font(.caption2)
                    .foregroundColor(.gray)
                    .padding(.horizontal, 12)
            }
            .frame(maxWidth: Constants.UI.messageBubbleMaxWidth, alignment: isFromCurrentUser ? .trailing : .leading)

            if !isFromCurrentUser {
                Spacer()
            }
        }
    }
}

// MARK: - Chat Input View

struct ChatInputView: View {
    @Binding var text: String
    let isSending: Bool
    let onSend: () -> Void

    var body: some View {
        HStack(spacing: 12) {
            // Text Field
            TextField("Message", text: $text)
                .textFieldStyle(PlainTextFieldStyle())
                .padding(10)
                .background(Color(.systemGray6))
                .cornerRadius(20)
                .disabled(isSending)

            // Send Button
            Button(action: onSend) {
                if isSending {
                    ProgressView()
                        .progressViewStyle(CircularProgressViewStyle())
                } else {
                    Image(systemName: "arrow.up.circle.fill")
                        .font(.system(size: 32))
                        .foregroundColor(text.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty ? .gray : .blue)
                }
            }
            .disabled(text.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty || isSending)
        }
        .padding(.horizontal)
        .padding(.vertical, 8)
        .background(Color(.systemBackground))
    }
}

struct ChatView_Previews: PreviewProvider {
    static var previews: some View {
        NavigationView {
            ChatView(conversation: Conversation(
                id: 1,
                type: .private,
                name: "Test User",
                avatar: nil,
                lastMessageId: nil,
                createdAt: Date(),
                updatedAt: Date(),
                lastMessage: nil,
                participants: nil
            ))
        }
    }
}
