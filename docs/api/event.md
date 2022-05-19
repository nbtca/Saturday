---
layout: default
title: 事件
parent: API
nav_order: 2
---

# 事件

## 基础

![](https://clas-bucket.oss-cn-hangzhou.aliyuncs.com/uPic/GhhXcd.png)

## 事件状态表(status)

| 状态名 |           | 描述                                 |
| ------ | --------- | ------------------------------------ |
| 待处理 | open      | 维修事件未被成员接受                 |
| 取消   | cancelled | 维修事件被用户取消，不需要再进行处理 |
| 受理   | accepted  | 维修事件已被成员接受                 |
| 待审核 | committed | 成员提交了维修描述，管理员尚未审核   |
| 关闭   | closed    | 维修事件已解决，不能再更改该事件     |

## 事件行为表(action)

| 操作名   |             | 操作权限       | 事件状态变更           | 描述                                         |
| -------- | ----------- | -------------- | ---------------------- | -------------------------------------------- |
| 创建     | create      | client         | nil => open            | 用户创建了维修事件                           |
| 受理     | accept      | member         | open => accepted       | 成员接受了维修事件                           |
| 取消     | cancel      | current client | open => canceled       | 用户取消了自己创建的维修事件                 |
| 放弃     | drop        | current member | accept => open         | 成员放弃了自己接受的维修事件                 |
| 提交     | commit      | current member | accept => committed    | 成员维修完成，添加维修描述后提交给管理员审核 |
| 修改提交 | alterCommit | current member | committed => committed | 成员修改 未被审核的维修提交                  |
| 拒绝提交 | reject      | admin          | committed => accepted  | 管理员拒绝提交                               |
| 关闭     | close       | admin          | committed => closed    | 管理员通过提交                               |

## 获取指定事件

```
GET /events/{event_id}
```

## 参数

| 名称     | 类型   | in   | 描述 |
| -------- | ------ | ---- | ---- |
| event_id | String | path | 学号 |

## 示例

### 请求

```
GET /members/event_id
```

### 响应

```
{
  "event_id": 1,
  "client_id": 1,
  "model": "7590",
  "problem": "hackintosh",
  "member_id": "",
  "closed_by": "",
  "status": "open",
  "logs": [
    {
      "log_id": 1,
      "description": "",
      "member_id": "",
      "action": "create",
      "gmt_create": "2022-05-10 11:00:26"
    },
    {
      "log_id": 2,
      "description": "",
      "member_id": "2333333333",
      "action": "accept",
      "gmt_create": "2022-05-10 11:03:18"
    }
  ],
  "gmt_create": "2022-05-10 10:23:54",
  "gmt_modified": "2022-05-12 23:22:44"
}
```

## Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 获取全部事件

```
GET /events
```



## 获取认证成员接受的指定事件

```
GET /member/events/{event_id}
```



## 获取认证成员接受的全部事件

```
GET /member/events
```



## 认证成员接受事件

+ 事件状态变更为`accepted`
+ 事件member_id变为成员id

```
POST /member/events/{event_id}/accept
```

## 参数

| 名称           | 类型   | in     | 描述 |
| -------------- | ------ | ------ | ---- |
| Authorizeation | string | header |      |
| event_id       | String | path   | 学号 |

## 示例

### 请求

```
POST /member/events/1/accept
```

### 响应

```
{
  "event_id": 1,
  "client_id": 1,
  "model": "7590",
  "problem": "hackintosh",
  "member_id": "2333333333",
  "closed_by": "",
  "status": "accepted",
  "logs": [
    {
      "log_id": 1,
      "description": "",
      "member_id": "",
      "action": "create",
      "gmt_create": "2022-05-10 11:00:26"
    },
    {
      "log_id": 2,
      "description": "",
      "member_id": "2333333333",
      "action": "accept",
      "gmt_create": "2022-05-10 11:03:18"
    }
  ],
  "gmt_create": "2022-05-10 10:23:54",
  "gmt_modified": "2022-05-12 23:22:44"
}
```

## Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 认证成员提交事件

+ 事件状态变更为`committed`
+ 提醒管理员审核

```
POST /member/events/{event_id}/commit
```

## 参数

| 名称           | 类型   | in     | 描述     |
| -------------- | ------ | ------ | -------- |
| Authorizeation | string | header |          |
| event_id       | String | path   | 学号     |
| content        | string | body   | 维修描述 |

## 示例

### 请求

```
POST /member/events/1/commit

{
	"content":"重装系统"
}
```

### 响应

```
{
  "event_id": 1,
  "client_id": 1,
  "model": "7590",
  "problem": "hackintosh",
  "member_id": "2333333333",
  "closed_by": "",
  "status": "committed",
  "logs": [
    {
      "log_id": 1,
      "description": "",
      "member_id": "",
      "action": "create",
      "gmt_create": "2022-05-10 11:00:26"
    },
    {
      "log_id": 2,
      "description": "",
      "member_id": "2333333333",
      "action": "accept",
      "gmt_create": "2022-05-10 11:03:18"
    },
    {
      "log_id": 3,
      "description": "重装系统",
      "member_id": "2333333333",
      "action": "commit",
      "gmt_create": "2022-05-10 11:03:18"
    },
  ],
  "gmt_create": "2022-05-10 10:23:54",
  "gmt_modified": "2022-05-12 23:22:44"
}
```

## Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 认证成员修改事件提交

```
PATCH /member/events/{event_id}/commit
```

## 参数

| 名称           | 类型   | in     | 描述 |
| -------------- | ------ | ------ | ---- |
| Authorizeation | string | header |      |
| event_id       | String | path   | 学号 |

## 示例

### 请求

```
PATCH /member/events/1/commit

{
	"content":"重装系统(ghost)"
}
```

### 响应

```
{
  "event_id": 1,
  "client_id": 1,
  "model": "7590",
  "problem": "hackintosh",
  "member_id": "2333333333",
  "closed_by": "",
  "status": "committed",
  "logs": [
    {
      "log_id": 1,
      "description": "",
      "member_id": "",
      "action": "create",
      "gmt_create": "2022-05-10 11:00:26"
    },
    {
      "log_id": 2,
      "description": "",
      "member_id": "2333333333",
      "action": "accept",
      "gmt_create": "2022-05-10 11:03:18"
    },
    {
      "log_id": 3,
      "description": "重装系统",
      "member_id": "2333333333",
      "action": "commit",
      "gmt_create": "2022-05-10 11:05:18"
    },
       {
      "log_id": 4,
      "description": "重装系统(ghost)",
      "member_id": "2333333333",
      "action": "alterCommit",
      "gmt_create": "2022-05-10 12:03:18"
    }
  ],
  "gmt_create": "2022-05-10 10:23:54",
  "gmt_modified": "2022-05-12 23:22:44"
}
```

## Http 状态码

| HTTP Status Code | 描述 |
| ---------------- | ---- |
| 200              | OK   |
| 422              |      |



## 认证成员放弃事件

+ 事件状态变更为`open`
+ 清空事件member_id字段

```
DELETE /member/events/{event_id}/accept
```

## 参数

| 名称           | 类型   | in     | 描述 |
| -------------- | ------ | ------ | ---- |
| Authorizeation | string | header |      |
| event_id       | String | path   | 学号 |

## 示例

### 请求

```
DELETE /member/events/1/accept
```

### 响应

```
{
  "event_id": 1,
  "client_id": 1,
  "model": "7590",
  "problem": "hackintosh",
  "member_id": "",
  "closed_by": "",
  "status": "open",
  "logs": [
    {
      "log_id": 1,
      "description": "",
      "member_id": "",
      "action": "create",
      "gmt_create": "2022-05-10 11:00:26"
    },
    {
      "log_id": 2,
      "description": "",
      "member_id": "2333333333",
      "action": "accept",
      "gmt_create": "2022-05-10 11:03:18"
    },
    {
      "log_id": 3,
      "description": "",
      "member_id": "2333333333",
      "action": "drop",
      "gmt_create": "2022-05-10 11:03:18"
    },
  ],
  "gmt_create": "2022-05-10 10:23:54",
  "gmt_modified": "2022-05-12 23:22:44"
}
```

## Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 管理员退回成员事件提交

+ 事件状态变更为`accepted`

```
DELETE /events/{event_id}/commit
```

## 参数

| 名称           | 类型   | in     | 描述 |
| -------------- | ------ | ------ | ---- |
| Authorizeation | string | header |      |
| event_id       | String | path   | 学号 |

## 示例

### 请求

```
DELETE /events/events/1/commit
```

### 响应

```
{
  "event_id": 1,
  "client_id": 1,
  "model": "7590",
  "problem": "hackintosh",
  "member_id": "2333333333",
  "closed_by": "",
  "status": "accepted",
  "logs": [
    {
      "log_id": 1,
      "description": "",
      "member_id": "",
      "action": "create",
      "gmt_create": "2022-05-10 11:00:26"
    },
    {
      "log_id": 2,
      "description": "",
      "member_id": "2333333333",
      "action": "accept",
      "gmt_create": "2022-05-10 11:03:18"
    },
    {
      "log_id": 3,
      "description": "重装系统",
      "member_id": "2333333333",
      "action": "commit",
      "gmt_create": "2022-05-10 11:03:18"
    },
    {
      "log_id": 4,
      "description": "",
      "member_id": "0000000000",
      "action": "reject",
      "gmt_create": "2022-05-10 11:03:18"
    }
  ],
  "gmt_create": "2022-05-10 10:23:54",
  "gmt_modified": "2022-05-12 23:22:44"
}
```

## Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 管理员关闭事件

+ 事件状态变更为`closed`
+ 事件closed_by字段变更为管理员id

```
POST /events/{event_id}/close
```

## 参数

| 名称           | 类型   | in     | 描述 |
| -------------- | ------ | ------ | ---- |
| Authorizeation | string | header |      |
| event_id       | String | path   | 学号 |

## 示例

### 请求

```
POST /events/events/1/close
```

### 响应

```
{
  "event_id": 1,
  "client_id": 1,
  "model": "7590",
  "problem": "hackintosh",
  "member_id": "2333333333",
  "closed_by": "0000000000",
  "status": "closed",
  "logs": [
    {
      "log_id": 1,
      "description": "",
      "member_id": "",
      "action": "create",
      "gmt_create": "2022-05-10 11:00:26"
    },
    {
      "log_id": 2,
      "description": "",
      "member_id": "2333333333",
      "action": "accept",
      "gmt_create": "2022-05-10 11:03:18"
    },
    {
      "log_id": 3,
      "description": "重装系统",
      "member_id": "2333333333",
      "action": "commit",
      "gmt_create": "2022-05-10 11:03:18"
    },
    {
      "log_id": 4,
      "description": "",
      "member_id": "0000000000",
      "action": "close",
      "gmt_create": "2022-05-10 11:03:18"
    }
  ],
  "gmt_create": "2022-05-10 10:23:54",
  "gmt_modified": "2022-05-12 23:22:44"
}
```

## Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 报修人员创建事件

```
POST /clients/event
```



## 报修人员更改事件

```
POST /clients/events/{event_id}
```



## 报修人员取消事件

```
POST /clients/events/{event_id}
```



## 获取报修人员事件

```
GET /client/events/{event_id}
```



## 获取报修人员全部事件

```
GET /client/events
```