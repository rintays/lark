# lark

A Golang CLI for Feishu/Lark inspired by gog.

## Usage

Configure credentials (writes `~/.config/lark/config.json` by default):

```bash
lark auth login --app-id <APP_ID> --app-secret <APP_SECRET>
```

Optionally override the API base URL:

```bash
lark auth login --app-id <APP_ID> --app-secret <APP_SECRET> --base-url https://open.feishu.cn
```

Fetch a tenant access token (cached in config):

```bash
lark auth
```

Get tenant info:

```bash
lark whoami
```

Send a message:

```bash
lark msg --chat-id <CHAT_ID> --text "hello"
```

### Global flags

- `--config` override the config path.
- `--json` output JSON.
