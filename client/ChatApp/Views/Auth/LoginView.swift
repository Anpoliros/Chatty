import SwiftUI

/// 登录视图
struct LoginView: View {
    @StateObject private var authViewModel = AuthViewModel()
    @State private var username = ""
    @State private var password = ""
    @State private var showRegister = false

    var body: some View {
        NavigationView {
            VStack(spacing: 20) {
                Spacer()

                // Logo
                Image(systemName: "message.fill")
                    .resizable()
                    .frame(width: 80, height: 80)
                    .foregroundColor(.blue)
                    .padding(.bottom, 40)

                // Title
                Text("Chatty")
                    .font(.largeTitle)
                    .fontWeight(.bold)

                Text("Welcome Back")
                    .font(.headline)
                    .foregroundColor(.gray)
                    .padding(.bottom, 30)

                // Username Field
                TextField("Username", text: $username)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
                    .autocapitalization(.none)
                    .disableAutocorrection(true)
                    .padding(.horizontal)

                // Password Field
                SecureField("Password", text: $password)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
                    .padding(.horizontal)

                // Error Message
                if let errorMessage = authViewModel.errorMessage {
                    Text(errorMessage)
                        .foregroundColor(.red)
                        .font(.caption)
                        .padding(.horizontal)
                }

                // Login Button
                Button(action: login) {
                    if authViewModel.isLoading {
                        ProgressView()
                            .progressViewStyle(CircularProgressViewStyle(tint: .white))
                    } else {
                        Text("Login")
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

                // Register Link
                Button(action: { showRegister = true }) {
                    Text("Don't have an account? Register")
                        .foregroundColor(.blue)
                }
                .padding(.top, 10)

                Spacer()
            }
            .navigationBarHidden(true)
            .sheet(isPresented: $showRegister) {
                RegisterView()
            }
        }
    }

    private func login() {
        hideKeyboard()
        authViewModel.login(username: username, password: password)
    }
}

struct LoginView_Previews: PreviewProvider {
    static var previews: some View {
        LoginView()
    }
}
