# envchain-cli

> Secure per-project environment variable manager backed by OS keychain with team sync support

---

## Installation

```bash
go install github.com/yourorg/envchain-cli@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourorg/envchain-cli/releases).

---

## Usage

**Store a secret for a project:**

```bash
envchain set myproject AWS_SECRET_ACCESS_KEY s3cr3tvalue
```

**Run a command with injected environment variables:**

```bash
envchain run myproject -- aws s3 ls
```

**List stored keys for a project:**

```bash
envchain list myproject
```

**Export variables for team sync (encrypted):**

```bash
envchain export myproject --out myproject.env.enc
envchain import myproject --in myproject.env.enc
```

Secrets are stored in your OS keychain (macOS Keychain, GNOME Keyring, or Windows Credential Manager) and never written to disk in plaintext.

---

## How It Works

- Each project has its own isolated namespace in the OS keychain
- Team sync exports use age encryption — only authorized keys can decrypt
- No plaintext `.env` files are ever created or committed

---

## Requirements

- Go 1.21+
- macOS, Linux (with `libsecret`), or Windows

---

## License

MIT © 2024 yourorg