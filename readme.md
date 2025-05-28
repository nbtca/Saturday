# Saturday

A restful API for nbtca repair service.

Related projects:

- [Sunday](https://github.com/nbtca/Sunday) (Web frontend)
- [Hawaii](https://github.com/nbtca/Hawaii) (WeChat mini program)

## Configuration

| Key                        | Description                   |
| -------------------------- | ----------------------------- |
| `DB_DATASOURCE`            | PostgreSQL connection string  |
| `ALIYUN_ACCESS_KEY_ID`     | Aliyun access key ID          |
| `ALIYUN_ACCESS_KEY_SECRET` | Aliyun access key secret      |
| `WECHAT_APPID`             | WeChat app ID                 |
| `WECHAT_SECRET`            | WeChat secret                 |
| `MAIL_HOST`                | SMTP server host              |
| `MAIL_PORT`                | SMTP port (e.g., 465)         |
| `MAIL_USERNAME`            | SMTP username                 |
| `MAIL_PASSWORD`            | SMTP password                 |
| `TESTING_MAIL_RECEIVER_ADDRESS`    | Mail receiver used for testing (optional)         |
| `LOGTO_APPID`              | Logto app ID                  |
| `LOGTO_APP_SECRET`         | Logto app secret              |
| `LOGTO_ENDPOINT`           | Logto endpoint URL            |
| `TESTING_LOGTO_USER_ID`       | Logto test user ID (optional) |
| `GITHUB_OWNER`             | GitHub repo owner             |
| `GITHUB_REPO`              | GitHub repository name        |
| `GITHUB_TOKEN`             | GitHub personal access token  |
| `GITHUB_WEBHOOK_SECRET`    | GitHub webhook signing secret |
| `DIFY_API_ENDPOINT`        | Dify API base URL             |
| `DIFY_API_KEY`             | Dify API key                  |
| `NSQ_HOST`                 | NSQ daemon host and port      |
| `NSQ_SECRET`               | NSQ secret (optional)         |
| `NSQ_EVENT_TOPIC`          | NSQ event topic name          |
| `NSQ_LOG_TOPIC`            | NSQ log topic name            |
| `SERVER_PORT`              | API server listen port        |

### Consul

You can also use consul to manage the configuration. The configuration will be loaded from Consul if the `CONSUL_HTTP_ADDR` and `CONSUL_HTTP_TOKEN` are set. The key will be used as the prefix for the configuration keys.

| Key                        | Description                   |
| -------------------------- | ----------------------------- |
| `CONSUL_HTTP_ADDR`         | Consul address                |
| `CONSUL_HTTP_TOKEN`        | Consul token                  |
| `CONSUL_KEY`               | Consul key                    |