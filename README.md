# vaultpull

> Sync secrets from HashiCorp Vault into local `.env` files with namespace filtering.

---

## Installation

```bash
go install github.com/youruser/vaultpull@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/vaultpull.git
cd vaultpull
go build -o vaultpull .
```

---

## Usage

Set your Vault address and token, then run `vaultpull` with a namespace and output path:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

vaultpull --namespace secret/myapp --output .env
```

This will pull all secrets under `secret/myapp` and write them to a local `.env` file:

```env
DATABASE_URL=postgres://user:pass@localhost:5432/mydb
API_KEY=abc123
REDIS_URL=redis://localhost:6379
```

### Flags

| Flag | Description | Default |
|-------------|--------------------------------------|---------|
| `--namespace` | Vault secret path / namespace | _(required)_ |
| `--output` | Output file path | `.env` |
| `--overwrite` | Overwrite existing file if it exists | `false` |
| `--addr` | Vault server address | `$VAULT_ADDR` |
| `--token` | Vault token | `$VAULT_TOKEN` |

---

## Requirements

- Go 1.21+
- HashiCorp Vault with a valid token and read permissions on the target namespace

---

## License

[MIT](LICENSE) © 2024 youruser