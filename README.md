# Passman CLI

A secure, cross-platform command-line password manager that keeps your credentials safe with end-to-end encryption.

### [🛡️Visit Product Page](https://passman.subratdwivedi.dev)

![Platform Support](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## ✨ Features

- **End-to-End Encryption** — Your passwords are encrypted locally before being sent to the server. Only you can decrypt them.
- **Zero-Knowledge Architecture** — Your master password never leaves your device. The server only stores encrypted data.
- **Secure Key Management** — Uses a background agent to keep your encryption key in memory, auto-locking after inactivity.
- **Cross-Platform** — Works on Windows, macOS, and Linux (including WSL2).
- **Interactive TUI** — Beautiful terminal user interface for browsing, searching, editing, and deleting passwords.
- **Clipboard Integration** — Copy passwords to clipboard with automatic clearing after 60 seconds.
- **Offline Fallback Storage** — Automatic encrypted file-based storage when system keyring is unavailable.
- **Self-Update** — Update to the latest version with a single command (`pman update`).

## 🔐 How It Works

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   You       │────▶│  Passman    │────▶│   Server    │
│             │     │    CLI      │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
       │                   │                   │
       │                   │                   │
  Master Password    Local Encryption    Encrypted Data
  (never sent)       with your key       (cannot decrypt)
```

### Security Architecture

1. **Registration/Login**: Your master password is used to derive an encryption key using Argon2id (memory-hard KDF).
2. **Key Storage**: The derived key is held in a background agent process, never written to disk.
3. **Encryption**: All passwords are encrypted with AES-256-GCM before leaving your device.
4. **Auto-Lock**: The agent automatically wipes the key from memory after 10 minutes of inactivity.
5. **Re-authentication**: When the agent locks, you'll be prompted for your master password to continue.

## 📥 Installation

### Windows

1. Download the latest `pman-windows-amd64.exe` from [Releases](https://github.com/subrat-dwi/passman-cli/releases)
2. Rename to `pman.exe` and move to a directory in your PATH (e.g., `C:\Users\YourName\bin\`)
3. Or add the download location to your PATH environment variable

**Using PowerShell:**
```powershell
# Create bin directory if it doesn't exist
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\bin"

# Download latest release
Invoke-WebRequest -Uri "https://github.com/subrat-dwi/passman-cli/releases/latest/download/pman-windows-amd64.exe" -OutFile "$env:USERPROFILE\bin\pman.exe"

# Add to PATH (run as Administrator, or add manually via System Properties)
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:USERPROFILE\bin", "User")
```

### macOS

1. Download the latest release for your architecture:
   - Apple Silicon (M1/M2/M3): `pman-darwin-arm64`
   - Intel: `pman-darwin-amd64`

**Using Terminal:**
```bash
# For Apple Silicon
curl -L https://github.com/subrat-dwi/passman-cli/releases/latest/download/pman-darwin-arm64 -o /usr/local/bin/pman

# For Intel
curl -L https://github.com/subrat-dwi/passman-cli/releases/latest/download/pman-darwin-amd64 -o /usr/local/bin/pman

# Make executable
chmod +x /usr/local/bin/pman
```

**Note**: On first run, macOS may block the app. Go to **System Preferences → Security & Privacy → General** and click "Allow Anyway".

### Linux

```bash
# Download
sudo curl -L https://github.com/subrat-dwi/passman-cli/releases/latest/download/pman-linux-amd64 -o /usr/local/bin/pman

# Make executable
sudo chmod +x /usr/local/bin/pman

# Verify installation
pman --help

```

### Build from Source

Requires Go 1.21 or later.

```bash
git clone https://github.com/subrat-dwi/passman-cli.git
cd passman-cli
go build -o pman cmd/pman/main.go

# Move to PATH
mv pman /usr/local/bin/  # Linux/macOS
# Or on Windows: move pman.exe to a directory in your PATH
```

## 🚀 Quick Start

### 1. Create an Account

```bash
pman auth register
```

You'll be prompted for:
- **Email**: Your account email
- **Master Password**: Choose a strong password (min 8 chars, requires uppercase, lowercase, number, and special character)

> ⚠️ **Important**: Remember your master password! It cannot be recovered if lost.

### 2. Login

```bash
pman auth login
```

### 3. Add a Password

```bash
pman create
```

Enter the service name (e.g., "Gmail"), username, and password.

### 4. View Your Passwords

```bash
pman list
```

This opens an interactive list where you can:
- **↑/↓** — Navigate
- **/** — Filter/search
- **Enter** — View password details
- **u** — Edit password
- **d** — Delete password
- **c** — Copy password to clipboard (in detail view)
- **q** — Quit

## 📖 Command Reference

### Authentication

| Command | Description |
|---------|-------------|
| `pman auth register` | Create a new account |
| `pman auth login` | Login to your account |
| `pman auth logout` | Logout and clear local credentials |

### Password Management

| Command | Description |
|---------|-------------|
| `pman list` | List all passwords (interactive) |
| `pman create` | Add a new password |

### Agent Management

| Command | Description |
|---------|-------------|
| `pman agent status` | Check if agent is locked/unlocked |
| `pman agent lock` | Manually lock the agent (wipe key from memory) |

### Other

| Command | Description |
|---------|-------------|
| `pman version` | Show version information |
| `pman update` | Update pman to the latest version |
| `pman update --check` | Check for updates without installing |
| `pman --help` | Show help |

## 🔒 Security Best Practices

### Choosing a Master Password

Your master password is the key to all your passwords. Choose wisely:

- ✅ Use at least 12 characters
- ✅ Include uppercase, lowercase, numbers, and symbols
- ✅ Use a passphrase (e.g., "correct-horse-battery-staple!")
- ❌ Don't reuse passwords from other services
- ❌ Don't use personal information (birthdays, names)

### Operational Security

- **Lock when away**: Run `pman agent lock` before leaving your computer
- **Don't share your master password**: It's the only password you need to remember
- **Use on trusted devices only**: Avoid public or shared computers

## 💾 Data Storage

### Credentials Storage

| Platform | Primary Storage | Fallback |
|----------|----------------|----------|
| Windows | Windows Credential Manager | Encrypted file |
| macOS | Keychain | Encrypted file |
| Linux/WSL | Encrypted file | — |

Fallback storage location: `~/.passman/vault.enc`

### What's Stored Locally

- Access token (for API authentication)
- Salt (for key derivation)
- Key verifier (encrypted token to verify master password)

**Never stored**: Your master password or encryption key

## 🔧 Troubleshooting

### "You're not logged in"

Run `pman auth login` to authenticate.

### "Vault locked. Enter master password"

The agent has auto-locked due to inactivity. Enter your master password to continue.

### "Incorrect master password"

The password you entered doesn't match. Try again with the correct master password.

### "Cannot connect to password agent"

The background agent isn't running. It should start automatically. Try:
1. Close all pman processes
2. Run any pman command again

### "Account setup incomplete"

Your local credentials are missing. Run `pman auth login` to restore them.

### Linux: Segmentation Fault

If you're using an older binary, download the latest release. The Linux version now uses file-based storage to avoid keyring issues.

## 🛠️ Configuration

Configuration file location: `~/.passman/config.yaml`

```yaml
api_base_url: https://api.passman.example.com
```

## 🔄 Updating

### Self-Update (Recommended)

Pman can update itself to the latest version:

```bash
# Check for updates
pman update --check

# Update to the latest version
pman update
```

**Note**: On Windows, you may need to run as Administrator. On Linux/macOS, you may need `sudo` depending on where pman is installed.

### Manual Update

Alternatively, download the latest release from the [Releases page](https://github.com/subrat-dwi/passman-cli/releases) and replace your existing binary.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — Terminal UI framework
- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [go-keyring](https://github.com/zalando/go-keyring) — Cross-platform keyring access
- [Argon2](https://github.com/P-H-C/phc-winner-argon2) — Password hashing
- [go-selfupdate](https://github.com/creativeprojects/go-selfupdate) — Self-update mechanism

---

**Developed by [Subrat Dwivedi](https://www.subratdwivedi.dev)**
