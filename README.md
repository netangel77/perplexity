Perplexity Chat GUI
A minimal Go app that runs a local web server, opens your default browser automatically, and provides a ChatGPT-style chat UI connected to the Perplexity /chat/completions API.

✨ Features
✅ Zero external dependencies (pure stdlib + //go:embed)

✅ ChatGPT-style UI with message history

✅ Model selection (sonar-small-online, sonar-pro, etc.)

✅ Citations display from Perplexity search results

✅ Single binary (go build → ship anywhere)

✅ Cross-platform (Windows/Linux/macOS)

🚀 Quick Start
Prerequisites
Go 1.18+ (for //go:embed)

Perplexity API key from Settings > API

1. Setup
bash
mkdir perplexity-gui && cd perplexity-gui
go mod init perplexity-gui
2. Create Files
Save the code from below as main.go.

3. Run
bash
go run .
Your browser opens automatically to http://127.0.0.1:XXXX:

Paste your pplx-... API key

Select model (sonar-small-online recommended for testing)

Start chatting!

📦 Build Single Executable
bash
go build -ldflags="-s -w" -o perplexity-chat.exe
Distribute perplexity-chat.exe (4MB) → runs anywhere, no install needed.

💰 Business Use (Your AI Agent Strategy)
Perfect for client delivery:

text
YourLogo-Perplexity-Research.exe  → $99 one-time
DevOps-Troubleshooter.exe         → $299 setup + $99/mo
Verticals ready:

DevOps troubleshooting agent

Market research tool

Lead gen research bot

Competitor analysis desktop

🛠️ Code (main.go)
go
// [Complete code from previous response - the working version]
🔧 Customization
Feature	Edit
Add Dark Mode	index.html CSS
Custom Domain	Change sonar-small-online → your hosted model
File Upload	Add /api/upload endpoint
Enterprise Logo	Replace header in index.html
Tray Icon	Add systray dependency
📱 Screenshots
text
[Expected UI: API Key input → Model dropdown → Chat area → Send button]
Paste key → "AI trends 2026" → Citations appear below
🐛 Troubleshooting
Issue	Fix
Port busy	Kills server, uses next free port
No browser opens	Copy URL from terminal
401 API Error	Check pplx- key + credits in dashboard
CORS	Localhost only, no CORS issues
📄 License
MIT - Free for commercial use.
