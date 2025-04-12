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

   ```
   DB_URL=<USER>:<PASSWORD>@(<ADDRESS>:<PORT>)/<DATABASE>

   ACCESS_KEY_ID=<YOUR_ACCESS_KEY_ID>
   ACCESS_KEY_SECRET=<YOUR_ACCESS_KEY_SECRET>
  
   MAIL_HOST=<YOUR_MAIL_HOST>
   MAIL_PORT=<YOUR_MAIL_PORT>
   MAIL_USERNAME=<YOUR_MAIL_USERNAME>
   MAIL_PASSWORD=<YOUR_MAIL_PASSWORD>

   PORT=<PORT_TO_LISTEN>
   # 以下为可选配置
   NSQ_HOST=<YOUR_NSQD_HOST>:<YOUR_NSQD_TCP_PORT (4150 IN COMMON)>
   NSQ_SECRET=<YOUR_NSQ_SECRET>
   EVENT_TOPIC=<YOUR_NSQ_EVENT_TOPIC_NAME>
   LOG_TOPIC=<YOUR_NSQ_LOG_TOPIC_NAME>

   RPC_PORT=<YOUR_RPC_PORT>

   LOGTO_APPID=<LOGTO_APPID>
   LOGTO_APP_SECRET=<LOGTO_APP_SECRET>
   LOGTO_ENDPOINT=<LOGTO_ENDPOINT>

   GITHUB_OWNER=<Github_Repo_Owner>
   GITHUB_REPO=<Github_Repo_Name>
   GITHUB_TOKEN=<Github_Token>
   ```

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
| Accept | Assign member to github issue  | Not implemented |
| Commit | Add comment in github issue  | |
| Drop | Add comment in github issue  | |
| Approve | Add comment in github issue and close issue | |

### From Github to Saturday

| Event Action | Github Action | Description |
| --- | --- | ---|
