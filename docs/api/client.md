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
GET /clients/{client_id}
```



## 获取全部报修人员

```
GET /clients/{client_id}
```



## 创建报修人员

```
POST /client
```



## 创建报修人员Token

```
POST /clients/token/wechat
```


### 参数

| 名称 | 类型   | in   | 描述          |
| ---- | ------ | ---- | ------------- |
| code | String | path | wx.login.code |


### 示例

#### 请求

// TODO {client_id} unknown

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

