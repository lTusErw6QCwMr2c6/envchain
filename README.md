# envchain

> Manage and chain environment variable profiles across projects with secret store integration.

## Installation

```bash
go install github.com/yourusername/envchain@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envchain.git && cd envchain && go build ./...
```

## Usage

Define a profile in `.envchain.yaml`:

```yaml
profiles:
  dev:
    AWS_REGION: us-east-1
    DB_HOST: localhost
    DB_PASS: secret://vault/myapp/db_password
  staging:
    extends: dev
    DB_HOST: staging.db.internal
```

Run a command with a loaded profile:

```bash
envchain run --profile dev -- go run ./cmd/server
```

Chain multiple profiles together:

```bash
envchain run --profile base,dev,overrides -- make deploy
```

Export resolved variables to your shell:

```bash
eval "$(envchain export --profile dev)"
```

List available profiles:

```bash
envchain list
```

## Secret Store Integration

`envchain` supports resolving secrets at runtime using the `secret://` URI scheme. Supported backends include HashiCorp Vault, AWS Secrets Manager, and local encrypted files.

```bash
envchain run --profile prod --secret-backend vault -- ./app
```

## License

MIT © [yourusername](https://github.com/yourusername)