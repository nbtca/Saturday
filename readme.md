# Saturday
> still relaxing

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
   RPC_PORT=<YOUR_RPC_PORT>
   ```
5. 启动服务
   在项目根目录下运行
   ``` sh
   go run main.go
   ```
6. 服务运行在`8080`端口

## 测试
1. 安装 `Docker`
2. 运行测试

   ```sh
   go test <floder>
   ```
