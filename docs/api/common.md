---
layout: default
title: 通用 
parent: API
nav_order: 1
---

# 通用（Common）

## 目录
- [文件上传](#文件上传)

## 文件上传
请求类型需为 multipart/form-data
```
PATCH /member/avatar
```
#### 参数

| 名称          | 类型   | in     | 描述                      |
| ------------- | ------ | ------ | ------------------------- |
| Authorization | string | header |                           |
| file          | File   | form   | 需上传的文件（小于 10mb） |

#### 示例

##### 请求
```
PATCH /upload
```

##### 响应

```
{
  url:"https://sunday-res.oss-cn-hangzhou.aliyuncs.com/weekend/1662184635.jpg"
}
```

#### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |
