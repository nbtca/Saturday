---
layout: default
title: 成员
parent: API
nav_order: 1
---

## 成员（Member）

### 获取指定成员 已完成

```
GET /members/{member_id}
```

#### 参数

| 名称      | 类型    | in   | 描述 |
| --------- | ------- | ---- | ---- |
| member_id | integer | path |      |

#### 示例

##### 请求

```
GET /members/2333333333
```

##### 响应

```
{
  "member_id": "2333333333",
  "alias": "滑稽",
  "role": "member",
  "profile": "relaxing",
  "avatar": "",
  "created_by": "0000000000",
  "gmt_create": "2022-04-23 15:49:59",
  "gmt_modified": "2022-04-30 17:29:46"
}
```

#### Http 状态码

| HTTP Status Code | 描述               |
| ---------------- | ------------------ |
| **200**          | OK                 |
| 404              | Resource not found |

### 获取全部成员 已完成

```
GET /members
```

#### 参数

| 名称   | 类型    | in    | 描述      |
| ------ | ------- | ----- | --------- |
| offset | integer | query |           |
| limit  | integer | query | 默认为 30 |

#### 示例

##### 请求

```
GET /members
```

##### 响应

```
[
  {
    "member_id": "0000000000",
    "alias": "管理",
    "role": "admin",
    "profile": "",
    "avatar": "",
    "created_by": "",
    "gmt_create": "2022-04-30 17:28:42",
    "gmt_modified": "2022-04-30 17:28:44"
  },
  {
    "member_id": "2333333333",
    "alias": "滑稽",
    "role": "member",
    "profile": "relaxing",
    "avatar": "",
    "created_by": "0000000000",
    "gmt_create": "2022-04-23 15:49:59",
    "gmt_modified": "2022-04-30 17:29:46"
  }
]
```

#### Http 状态码

| HTTP Status Code | 描述               |
| ---------------- | ------------------ |
| **200**          | OK                 |
| 404              | Resource not found |

### 创建用户 Token 已完成

返回认证用户信息以及 token

```
POST /members/{member_id}/token
```

#### 参数

| 名称      | 类型   | in   | 描述 |
| --------- | ------ | ---- | ---- |
| member_id | string | path | 学号 |
| password  | string | body | 姓名 |

#### 示例

##### 请求

```
POST /members/2333333333

{
  "password": "123456"
}
```

##### 响应

```
{
  "member_id": "2333333333",
  "alias": "滑稽",
  "name": "滑稽",
  "section": "计算机233",
  "role": "member",
  "profile": "relaxing",
  "phone": "12356839487",
  "qq": "123456",
  "avatar": "",
  "created_by": "0000000000",
  "gmt_create": "2022-04-23 15:49:59",
  "gmt_modified": "2022-04-30 17:29:46",
  "token": "not implemented"
}
```

#### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 404              | Resource not found   |
| 422              | Unprocessable Entity |

### 获取认证用户信息 已完成

```
GET /member
```

#### 参数

| 名称           | 类型   | in     | 描述 |
| -------------- | ------ | ------ | ---- |
| Authorizeation | string | header |      |

#### 示例

##### 请求

```
GET /member
```

##### 响应

```
{
  "member_id": "2333333333",
  "alias": "滑稽",
  "name": "滑稽",
  "section": "计算机233",
  "role": "member",
  "profile": "relaxing",
  "phone": "12356839487",
  "qq": "123456",
  "avatar": "",
  "created_by": "0000000000",
  "gmt_create": "2022-04-23 15:49:59",
  "gmt_modified": "2022-04-30 17:29:46"
}
```

#### Http 状态码

| HTTP Status Code | 描述               |
| ---------------- | ------------------ |
| **200**          | OK                 |
| 404              | Resource not found |

### 用户激活 未完成

- 用户在初次设定密码后激活
- member_inactive=>member
- admin_inactive=>admin

```
PUT /member/active
```

#### 参数

| 名称           | 类型    | in     | 描述 |
| -------------- | ------- | ------ | ---- |
| Authorizeation | string  | header |      |
| member_id      | integer | path   | 学号 |
| password       | string  | body   | 密码 |
| alias `可选`   | string  | body   | 昵称 |
| phone `可选`   | string  | body   |      |
| qq `可选`      | string  | body   |      |
| avatar `可选`  | string  | body   | 头像 |
| profile `可选` | string  | body   | 简介 |

#### 示例

##### 请求

```
PUT /member

{
  "alias": "滑稽",
  "phone": "12356839487",
  "qq": "123456",
  "password":"123456"
}
```

##### 响应

```
{
  "member_id": "2333333333",
  "alias": "滑da稽",
  "name": "滑稽",
  "section": "计算机233",
  "role": "member",
  "profile": "want to relax",
  "phone": "12356839487",
  "qq": "123456",
  "avatar": "",
  "created_by": "0000000000",
  "gmt_create": "2022-04-23 15:49:59",
  "gmt_modified": "2022-04-30 17:29:46"
}
```



### 用户更新信息 未完成

```
PUT /member
```

#### 参数

| 名称            | 类型    | in     | 描述 |
| --------------- | ------- | ------ | ---- |
| Authorizeation  | string  | header |      |
| alias `可选`    | string  | body   | 昵称 |
| member_id       | integer | path   | 学号 |
| phone `可选`    | string  | body   |      |
| qq `可选`       | string  | body   |      |
| avatar `可选`   | string  | body   | 头像 |
| profile `可选`  | string  | body   | 简介 |
| password `可选` | string  | body   | 密码 |

#### 示例

##### 请求

```
PUT /member

{
  "member_id": "2333333333",
  "alias": "滑da稽",
  "name": "滑稽",
  "profile": "want to relax",
  "phone": "12356839487",
  "qq": "123456"
}
```

##### 响应

```
{
  "member_id": "2333333333",
  "alias": "滑da稽",
  "name": "滑稽",
  "section": "计算机233",
  "role": "member",
  "profile": "want to relax",
  "phone": "12356839487",
  "qq": "123456",
  "avatar": "",
  "created_by": "0000000000",
  "gmt_create": "2022-04-23 15:49:59",
  "gmt_modified": "2022-04-30 17:29:46"
}
```

### 用户修改头像 未完成

```
PUT /member/avater
```

### 创建成员 已完成

- 需要身份为管理员
- member_id（学号）需为唯一

```
POST /members/{member_id}
```

#### 参数

| 名称           | 类型    | in     | 描述 |
| -------------- | ------- | ------ | ---- |
| Authorizeation | string  | header |      |
| member_id      | integer | path   | 学号 |
| name           | string  | body   | 姓名 |
| section        | string  | body   | 班级 |
| role           | string  | body   | 权限 |
| alias `可选`   | string  | body   | 昵称 |
| phone `可选`   | string  | body   |
| qq `可选`      | string  | body   |
| avatar `可选`  | string  | body   | 头像 |

#### 示例

##### 请求

```
POST /members/3000000000

{
  "alias": "小稽",
  "name": "滑小稽",
  "section": "计算机233",
  "profile": "。。。",
  "role": "member_inactive",
  "phone": "12352439487",
  "qq": "123456"
}
```

##### 响应

```
{
  "member_id": "3000000000",
  "alias": "小稽",
  "name": "滑小稽",
  "section": "计算机233",
  "role": "member_inactive",
  "profile": "。。。",
  "phone": "",
  "qq": "123456",
  "avatar": "",
  "created_by": "2333333333",
  "gmt_create": "2022-04-30 23:06:44",
  "gmt_modified": "2022-04-30 23:06:44"
}
```

#### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 404              | Resource not found   |
| 422              | Unprocessable Entity |

### 创建多个用户 未完成

```
POST /members
```

// TODO

### 修改用户基本信息 已完成

```
PUT /members/{member_id}
```

#### 参数

| 名称           | 类型    | in     | 描述 |
| -------------- | ------- | ------ | ---- |
| Authorizeation | string  | header |      |
| member_id      | integer | path   | 学号 |
| name           | string  | body   | 姓名 |
| section        | string  | body   | 班级 |
| role           | string  | body   | 权限 |

#### 示例

##### 请求

```
PATCH /members/2333333333

{
  "name": "滑稽",
  "section": "计算机322",
  "role":"admin"
}
```

##### 响应

```
{
  "member_id": "2333333333",
  "alias": "滑稽",
  "name": "滑稽",
  "section": "计算机322",
  "profile": "relaxing",
  "phone": "12356839487",
  "qq": "123456",
  "avatar": "",
  "created_by": "",
  "gmt_create": "2022-04-17T19:35:55.000Z",
  "gmt_modified": "2022-04-17T19:35:55.000Z",
  "role": "admin"
}
```

#### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 404              | Resource not found   |
| 422              | Unprocessable Entity |
