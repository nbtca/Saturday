---
layout: default
title: 成员
parent: API
nav_order: 1
---

# 成员（Member）

## 目录

- [获取指定成员 已完成](#获取指定成员-已完成)
- [获取全部成员 已完成](#获取全部成员-已完成)
- [创建成员 Token 已完成](#创建成员-token-已完成)
- [获取认证成员信息 已完成](#获取认证成员信息-已完成)
- [成员激活](#成员激活)
- [成员更新信息](#成员更新信息)
- [成员修改头像 已完成](#成员修改头像-已完成)
- [创建成员 已完成](#创建成员-已完成)
- [创建多个成员 未完成](#创建多个成员-未完成)
- [修改成员基本信息 已完成](#修改成员基本信息-已完成)

---

## 获取指定成员 已完成
通过URI中指定的成员ID获取成员公开信息
```
GET /members/{memberId}
```

#### 参数


| 名称     | 类型    | in   | 描述 |
| -------- | ------- | ---- | ---- |
| memberId | integer | path |      |

#### 示例

##### 请求
```
GET /members/2333333333
```

##### 响应

```
{
  "memberId": "2333333333",
  "alias": "滑稽",
  "role": "member",
  "profile": "relaxing",
  "avatar": "",
  "createdBy": "0000000000",
  "gmtCreate": "2022-04-23 15:49:59",
  "gmtModified": "2022-04-30 17:29:46"
}
```

#### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |

## 获取全部成员 已完成

```
GET /members
```

#### 参数

| 名称   | 类型    | in    | 描述      |
| ------ | ------- | ----- | --------- |
| offset | integer | query |           |
| limit  | integer | query | 默认为 50 |

#### 示例

##### 请求

```
GET /members
```

##### 响应

```
[
  {
    "memberId": "0000000000",
    "alias": "管理",
    "role": "admin",
    "profile": "",
    "avatar": "",
    "createdBy": "",
    "gmtCreate": "2022-04-30 17:28:42",
    "gmtModified": "2022-04-30 17:28:44"
  },
  {
    "memberId": "2333333333",
    "alias": "滑稽",
    "role": "member",
    "profile": "relaxing",
    "avatar": "",
    "createdBy": "0000000000",
    "gmtCreate": "2022-04-23 15:49:59",
    "gmtModified": "2022-04-30 17:29:46"
  }
]
```

#### Http 状态码


| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |

## 创建成员 Token 已完成

接受维修人员 ID 及密码两个参数，若密码正确，则返回对应的维修人员信息以及令牌，若缺少参数，参数类型错误，密码错误或用户不存在则返回相应的错误信息。

```
POST /members/{memberId}/token
```

#### 参数

| 名称     | 类型   | in   | 描述 |
| -------- | ------ | ---- | ---- |
| memberId | string | path | 学号 |
| password | string | body | 密码 |

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
  "memberId": "2333333333",
  "alias": "滑稽",
  "name": "滑稽",
  "section": "计算机233",
  "role": "member",
  "profile": "relaxing",
  "phone": "12356839487",
  "qq": "123456",
  "avatar": "",
  "createdBy": "0000000000",
  "gmtCreate": "2022-04-23 15:49:59",
  "gmtModified": "2022-04-30 17:29:46",
  "token": "not implemented"
}
```

#### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |

## 获取认证成员信息 已完成
通过请求 Header 中 Authorization 字段获取成员信息，返回相应的成员详细信息，若认证失败或认证身份不为member，则返回相应的错误信息。

```
GET /member
```

#### 参数

| 名称          | 类型   | in     | 描述 |
| ------------- | ------ | ------ | ---- |
| Authorization | string | header |      |

#### 示例

##### 请求

```
GET /member
```

##### 响应

```
{
  "memberId": "2333333333",
  "alias": "滑稽",
  "name": "滑稽",
  "section": "计算机233",
  "role": "member",
  "profile": "relaxing",
  "phone": "12356839487",
  "qq": "123456",
  "avatar": "",
  "createdBy": "0000000000",
  "gmtCreate": "2022-04-23 15:49:59",
  "gmtModified": "2022-04-30 17:29:46"
}
```

#### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |

## 成员激活 

- 成员在初次设定密码后激活
- member_inactive=>member
- admin_inactive=>admin

```
PATCH /member/active
```

#### 参数

| 名称           | 类型   | in     | 描述 |
| -------------- | ------ | ------ | ---- |
| Authorization  | string | header |      |
| password       | string | body   | 密码 |
| alias `可选`   | string | body   | 昵称 |
| phone `可选`   | string | body   |      |
| qq `可选`      | string | body   |      |
| profile `可选` | string | body   | 简介 |

#### 示例

##### 请求

```
PATCH /member/active

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
  "memberId": "2333333333",
  "alias": "滑da稽",
  "name": "滑稽",
  "section": "计算机233",
  "role": "member",
  "profile": "want to relax",
  "phone": "12356839487",
  "qq": "123456",
  "avatar": "",
  "createdBy": "0000000000",
  "gmtCreate": "2022-04-23 15:49:59",
  "gmtModified": "2022-04-30 17:29:46"
}
```

#### http 状态码

| http status code | 描述                 |
| ---------------- | -------------------- |
| 200              | ok                   |
| 422              | unprocessable entity |

## 成员更新信息

通过请求 Header 中 Authorization 字段获取成员信息，根据请求body内字段更新成员信息，若认证失败或认证身份不为member，则返回相应的错误信息。
```
PUT /member
```

#### 参数

| 名称            | 类型    | in     | 描述 |
| --------------- | ------- | ------ | ---- |
| Authorization   | string  | header |      |
| alias `可选`    | string  | body   | 昵称 |
| memberId        | integer | path   | 学号 |
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
  "memberId": "2333333333",
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
  "memberId": "2333333333",
  "alias": "滑da稽",
  "name": "滑稽",
  "section": "计算机233",
  "role": "member",
  "profile": "want to relax",
  "phone": "12356839487",
  "qq": "123456",
  "avatar": "",
  "createdBy": "0000000000",
  "gmtCreate": "2022-04-23 15:49:59",
  "gmtModified": "2022-04-30 17:29:46"
}
```
#### http 状态码

| http status code | 描述                 |
| ---------------- | -------------------- |
| 200              | ok                   |
| 422              | unprocessable entity |

## 成员修改头像 已完成

```
PATCH /member/avatar
```
#### 参数

| 名称          | 类型   | in     | 描述 |
| ------------- | ------ | ------ | ---- |
| Authorization | string | header |      |
| url           | string | body   |      |

#### 示例

##### 请求

```
PATCH /member/avatar

{
  url:"https://sunday-res.oss-cn-hangzhou.aliyuncs.com/weekend/1662184635.jpg"
}
```

##### 响应

```
{
  "memberId": "3000000000",
  "alias": "小稽",
  "name": "滑小稽",
  "section": "计算机233",
  "role": "member_inactive",
  "profile": "。。。",
  "phone": "",
  "qq": "123456",
  "avatar": "https://sunday-res.oss-cn-hangzhou.aliyuncs.com/weekend/1662184635.jpg",
  "createdBy": "2333333333",
  "gmtCreate": "2022-04-30 23:06:44",
  "gmtModified": "2022-04-30 23:06:44"
}
```

#### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |


## 创建成员 已完成

需要身份为管理员,memberId（学号）需为唯一，若已存在，则返回错误信息。

```
POST /members/{memberId}
```

#### 参数

| 名称          | 类型    | in     | 描述 |
| ------------- | ------- | ------ | ---- |
| Authorization | string  | header |      |
| memberId      | integer | path   | 学号 |
| name          | string  | body   | 姓名 |
| section       | string  | body   | 班级 |
| role          | string  | body   | 权限 |
| alias `可选`  | string  | body   | 昵称 |
| phone `可选`  | string  | body   |
| qq `可选`     | string  | body   |
| avatar `可选` | string  | body   | 头像 |

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
  "memberId": "3000000000",
  "alias": "小稽",
  "name": "滑小稽",
  "section": "计算机233",
  "role": "member_inactive",
  "profile": "。。。",
  "phone": "",
  "qq": "123456",
  "avatar": "",
  "createdBy": "2333333333",
  "gmtCreate": "2022-04-30 23:06:44",
  "gmtModified": "2022-04-30 23:06:44"
}
```

#### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |

## 创建多个成员 未完成

```
POST /members
```

// TODO

## 修改成员基本信息 已完成

```
PUT /members/{memberId}
```

#### 参数

| 名称          | 类型    | in     | 描述 |
| ------------- | ------- | ------ | ---- |
| Authorization | string  | header |      |
| memberId      | integer | path   | 学号 |
| name          | string  | body   | 姓名 |
| section       | string  | body   | 班级 |
| role          | string  | body   | 权限 |

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
  "memberId": "2333333333",
  "alias": "滑稽",
  "name": "滑稽",
  "section": "计算机322",
  "profile": "relaxing",
  "phone": "12356839487",
  "qq": "123456",
  "avatar": "",
  "createdBy": "",
  "gmtCreate": "2022-04-17T19:35:55.000Z",
  "gmtModified": "2022-04-17T19:35:55.000Z",
  "role": "admin"
}
```

#### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |