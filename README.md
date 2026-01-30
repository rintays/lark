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
lark msg send --chat-id <CHAT_ID> --text "hello"
```

Send a message to a user by email:

```bash
lark msg send --receive-id-type email --receive-id user@example.com --text "hello"
```

List recent chats:

```bash
lark chats list --limit 10
```

Search users:

```bash
lark users search --email user@example.com
lark users search --mobile "+1-555-0100"
lark users search --name "Ada"
lark users search --name "Ada" --department-id 0
```

### Global flags

- `--config` override the config path.
- `--json` output JSON.
