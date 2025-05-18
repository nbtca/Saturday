# Saturday

## 简介

+ 使用 Golang，gin 搭建的维修队后端
+ [API文档](https://nbtca.github.io/Saturday/api)
+ 使用此后端服务的项目
  + [Sunday](https://github.com/nbtca/Sunday) (管理系统)
  + [Hawaii](https://github.com/nbtca/Hawaii) (维修小程序)

## 如何运行

1. 安装 `Golang` , `Mysql`
2. 安装项目依赖

   在项目根目录下运行

   ``` sh
   go get
   ```

3. 导入数据库

   在`Mysql`中新建数据库，并将`assets/saturday.sql`导入
4. 添加配置文件

   在项目根目录下新建`.env`文件，添加配置

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

   在项目根目录下运行

   ``` sh
   go run main.go
   ```

5. 服务运行在`8080`端口

## 测试

1. 安装 `Docker`
2. 运行测试

   ```sh
   go test <floder>
   ```

## Syncing with Github Issue

The aim is to achieve a two-way sync between Saturday and Github Issues. The following table outlines the actions taken in Saturday and their corresponding actions in Github.

### From Saturday to Github

| Event Action | Github Action | Comment |
| --- | --- | ---|
| Create | Create Github Issue | |
| Cancel | Close Github Issue as not planned  | |
| Commit | Add comment in github issue  | |
| Drop | Add comment in github issue  | |
| close | Add comment in github issue and close issue | |

### From Github to Saturday

| Event Action | Github Action | Description |
| --- | --- | ---|
| accept | @nbtca-bot accept | |
| drop | @nbtca-bot drop | |
| commit | @nbtca-bot commit | |
| alterAccept |  | edit previous comment |
| reject | @nbtca-bot reject | |
| close | @nbtca-bot close | |
