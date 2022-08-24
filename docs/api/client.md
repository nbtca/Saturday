---
layout: default
title: 报修人员
parent: API
nav_order: 3
---

# 报修人员
{: .no_toc}

## 目录
- [获取指定报修人员](#获取指定报修人员)
- [获取全部报修人员](#获取全部报修人员)
- [创建报修人员](#创建报修人员)
- [创建报修人员Token](#创建报修人员token)

## 获取指定报修人员

```
GET /clients/{clientId}
```



## 获取全部报修人员

```
GET /clients/{clientId}
```



## 创建报修人员

```
POST /client
```



## 通过微信创建报修人员Token

该接口接受微信小程序内获取的wx.login.code，并通过code，从微信提供的接口获取用户的openId，若code缺少或微信接口返回错误，则返回相应的错误。在获取到openId后，若已存对应openId的用户，则返回对应的报修人员信息以及令牌，若不存在拥有对应openId的用户，则创建一个相应openId的用户。

```
POST /clients/token/wechat
```


### 参数

| 名称 | 类型   | in   | 描述          |
| ---- | ------ | ---- | ------------- |
| code | String | path | wx.login.code |


### 示例

#### 请求

// TODO {clientId} unknown

```
POST /clients/token/wechat
```

#### 响应

```
{
  "uid": "6f1f9702-66d6-447b-8e55-5f0c647c8d3a",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTM0NDMxNjIsImRhdGEiOnsidWlkIjoiNmYxZjk3MDItNjZkNi00NDdiLThlNTUtNWYwYzY0N2M4ZDNhIiwicm9sZSI6InVzZXIifSwiaWF0IjoxNjUzMzU2NzYyfQ.ocAxJGhw6Xt2vt7bwGcMeRPLOQOmaspznyu9aI7G670"
}
```

### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |

