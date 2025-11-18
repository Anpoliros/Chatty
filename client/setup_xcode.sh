#!/bin/bash

# Chatty iOS App - Xcode Setup Script
# This script helps organize the project structure for Xcode

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

echo "🚀 Setting up Chatty iOS Project..."
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Xcode is installed
if ! command -v xcodebuild &> /dev/null; then
    echo -e "${RED}❌ Xcode is not installed. Please install Xcode from the App Store.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Xcode is installed${NC}"
echo ""

# Check if all source files exist
echo "📁 Checking source files..."

required_files=(
    "ChatApp/ChatApp.swift"
    "ChatApp/Info.plist"
    "ChatApp/Models/User.swift"
    "ChatApp/Models/Message.swift"
    "ChatApp/Models/Conversation.swift"
    "ChatApp/Views/Auth/LoginView.swift"
    "ChatApp/Views/Auth/RegisterView.swift"
    "ChatApp/Views/Conversation/ConversationListView.swift"
    "ChatApp/Views/Chat/ChatView.swift"
    "ChatApp/ViewModels/AuthViewModel.swift"
    "ChatApp/ViewModels/ConversationViewModel.swift"
    "ChatApp/ViewModels/ChatViewModel.swift"
    "ChatApp/Services/Network/APIClient.swift"
    "ChatApp/Services/Storage/UserDefaultsManager.swift"
    "ChatApp/Services/WebSocket/WebSocketManager.swift"
    "ChatApp/Utils/Constants.swift"
    "ChatApp/Utils/Extensions.swift"
)

missing_files=()
for file in "${required_files[@]}"; do
    if [ ! -f "$file" ]; then
        missing_files+=("$file")
    fi
done

if [ ${#missing_files[@]} -ne 0 ]; then
    echo -e "${RED}❌ Missing files:${NC}"
    for file in "${missing_files[@]}"; do
        echo "  - $file"
    done
    exit 1
fi

echo -e "${GREEN}✅ All source files found${NC}"
echo ""

# Display instructions
echo "========================================="
echo "  Next Steps to Open in Xcode"
echo "========================================="
echo ""
echo -e "${BLUE}Option 1: Create New Project (Recommended)${NC}"
echo ""
echo "1. Open Xcode"
echo "2. File → New → Project"
echo "3. Choose iOS → App"
echo "4. Project Name: ChatApp"
echo "5. Interface: SwiftUI"
echo "6. Language: Swift"
echo "7. Save in: $SCRIPT_DIR"
echo ""
echo "8. Delete default files:"
echo "   - ContentView.swift"
echo "   - ChatAppApp.swift (Xcode generated)"
echo ""
echo "9. Add source files:"
echo "   - Right-click ChatApp group"
echo "   - Add Files to 'ChatApp'"
echo "   - Select ChatApp folder"
echo "   - Check: 'Copy items if needed'"
echo "   - Check: 'Create groups'"
echo "   - Add to targets: ChatApp"
echo ""
echo "10. Add Starscream package:"
echo "    - Project → Package Dependencies"
echo "    - Click '+'"
echo "    - URL: https://github.com/daltoniam/Starscream"
echo ""
echo -e "${BLUE}Option 2: Use Xcode to Open Directory${NC}"
echo ""
echo "1. Open Xcode"
echo "2. File → Open"
echo "3. Select: $SCRIPT_DIR/ChatApp"
echo "4. Xcode will detect Swift files"
echo "5. Follow steps 10 from Option 1 to add Starscream"
echo ""
echo "========================================="
echo ""
echo -e "${YELLOW}📝 Before Running:${NC}"
echo ""
echo "1. Make sure server is running:"
echo "   cd ../server && make run"
echo ""
echo "2. Configure server address in ChatApp/Utils/Constants.swift"
echo "   - Simulator: use localhost"
echo "   - Real device: use your Mac's IP address"
echo ""
echo "3. Trust your developer account on the device (for real device testing)"
echo ""
echo "========================================="
echo ""
echo -e "${GREEN}✨ Setup Complete!${NC}"
echo ""
echo "For detailed instructions, see: SETUP.md"
echo "For project overview, see: README.md"
echo ""
echo "Happy coding! 🎉"
echo ""
