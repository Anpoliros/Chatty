import SwiftUI

/// 注册视图
struct RegisterView: View {
    @StateObject private var authViewModel = AuthViewModel()
    @Environment(\.presentationMode) var presentationMode

    @State private var username = ""
    @State private var email = ""
    @State private var password = ""
    @State private var confirmPassword = ""
    @State private var nickname = ""
    @State private var localError = ""

    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    // Logo
                    Image(systemName: "person.crop.circle.badge.plus")
                        .resizable()
                        .frame(width: 80, height: 80)
                        .foregroundColor(.blue)
                        .padding(.top, 40)
                        .padding(.bottom, 20)

                    Text("Create Account")
                        .font(.title)
                        .fontWeight(.bold)

                    // Username Field
                    VStack(alignment: .leading, spacing: 5) {
                        Text("Username")
                            .font(.caption)
                            .foregroundColor(.gray)
                        TextField("Username", text: $username)
                            .textFieldStyle(RoundedBorderTextFieldStyle())
                            .autocapitalization(.none)
                            .disableAutocorrection(true)
                        Text("3-20 characters, letters, numbers and underscore only")
                            .font(.caption2)
                            .foregroundColor(.gray)
                    }
                    .padding(.horizontal)

                    // Email Field
                    VStack(alignment: .leading, spacing: 5) {
                        Text("Email")
                            .font(.caption)
                            .foregroundColor(.gray)
                        TextField("Email", text: $email)
                            .textFieldStyle(RoundedBorderTextFieldStyle())
                            .autocapitalization(.none)
                            .keyboardType(.emailAddress)
                            .disableAutocorrection(true)
                    }
                    .padding(.horizontal)

                    // Nickname Field (Optional)
                    VStack(alignment: .leading, spacing: 5) {
                        Text("Nickname (Optional)")
                            .font(.caption)
                            .foregroundColor(.gray)
                        TextField("Nickname", text: $nickname)
                            .textFieldStyle(RoundedBorderTextFieldStyle())
                    }
                    .padding(.horizontal)

                    // Password Field
                    VStack(alignment: .leading, spacing: 5) {
                        Text("Password")
                            .font(.caption)
                            .foregroundColor(.gray)
                        SecureField("Password", text: $password)
                            .textFieldStyle(RoundedBorderTextFieldStyle())
                        Text("At least 6 characters")
                            .font(.caption2)
                            .foregroundColor(.gray)
                    }
                    .padding(.horizontal)

                    // Confirm Password Field
                    VStack(alignment: .leading, spacing: 5) {
                        Text("Confirm Password")
                            .font(.caption)
                            .foregroundColor(.gray)
                        SecureField("Confirm Password", text: $confirmPassword)
                            .textFieldStyle(RoundedBorderTextFieldStyle())
                    }
                    .padding(.horizontal)

                    // Error Message
                    if !localError.isEmpty {
                        Text(localError)
                            .foregroundColor(.red)
                            .font(.caption)
                            .padding(.horizontal)
                    }

                    if let errorMessage = authViewModel.errorMessage {
                        Text(errorMessage)
                            .foregroundColor(.red)
                            .font(.caption)
                            .padding(.horizontal)
                    }

                    // Register Button
                    Button(action: register) {
                        if authViewModel.isLoading {
                            ProgressView()
                                .progressViewStyle(CircularProgressViewStyle(tint: .white))
                        } else {
                            Text("Register")
                                .fontWeight(.semibold)
                        }
                    }
                    .frame(maxWidth: .infinity)
                    .frame(height: 50)
                    .background(Color.blue)
                    .foregroundColor(.white)
                    .cornerRadius(Constants.UI.cornerRadius)
                    .padding(.horizontal)
                    .disabled(authViewModel.isLoading)
                    .padding(.top, 10)

                    Spacer()
                }
            }
            .navigationBarTitle("Register", displayMode: .inline)
            .navigationBarItems(leading: Button("Cancel") {
                presentationMode.wrappedValue.dismiss()
            })
        }
    }

    private func register() {
        // Validate inputs
        localError = ""

        guard !username.isEmpty, !email.isEmpty, !password.isEmpty else {
            localError = "Please fill in all required fields"
            return
        }

        guard password == confirmPassword else {
            localError = "Passwords do not match"
            return
        }

        hideKeyboard()

        // Call register
        authViewModel.register(
            username: username,
            email: email,
            password: password,
            nickname: nickname.isEmpty ? nil : nickname
        )

        // Observe login state
        // If successful, dismiss the view
        DispatchQueue.main.asyncAfter(deadline: .now() + 0.5) {
            if authViewModel.isLoggedIn {
                presentationMode.wrappedValue.dismiss()
            }
        }
    }
}

struct RegisterView_Previews: PreviewProvider {
    static var previews: some View {
        RegisterView()
    }
}
