# Saturday
> still relaxing

## 简介
维修队后端

## 如何运行
1. 安装`golang`,`mysql`
2. 安装项目依赖

   在项目根目录下运行
   ``` sh
   go get
   ```
3. 导入数据库

   在`mysql`中新建数据库，并将`assets/saturday.sql`导入
4. 添加配置文件

   在项目根目录下新建`.env`文件，添加`DB_URL`
   ```sh
   touch .env
   echo DB_URL="<账号>:<密码>@<地址>/<数据库名称>"
   ```
5. 启动服务
   在项目根目录下运行
   ``` sh
   go run main.go
   ```
