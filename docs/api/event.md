---
layout: default
title: 事件
parent: API
nav_order: 2
---

# 事件

## 目录
- [事件状态表(status)](#事件状态表status)
- [事件行为表(action)](#事件行为表action)
- [获取指定事件](#获取指定事件)
- [获取全部事件](#获取全部事件)
- [获取认证成员接受的指定事件](#获取认证成员接受的指定事件)
- [获取认证成员接受的全部事件](#获取认证成员接受的全部事件)
- [认证成员接受事件](#认证成员接受事件)
- [认证成员提交事件](#认证成员提交事件)
- [认证成员修改事件提交](#认证成员修改事件提交)
- [认证成员放弃事件](#认证成员放弃事件)
- [管理员退回成员事件提交](#管理员退回成员事件提交)
- [管理员关闭事件](#管理员关闭事件)
- [报修人员创建事件](#报修人员创建事件)
- [报修人员更改事件](#报修人员更改事件)
- [报修人员取消事件](#报修人员取消事件)
- [报修人员获取指定事件](#报修人员获取指定事件)
- [报修人员获取全部事件](#报修人员获取全部事件)
- [获取报修人员全部事件](#获取报修人员全部事件)

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
GET /events/{eventId}
```

### 参数

| 名称     | 类型   | in   | 描述   |
| -------- | ------ | ---- | ------ |
| eventId | string | path | 事件ID |

### 示例

#### 请求

```
GET /members/eventId
```

#### 响应

```
{
  "eventId": 1,
  "clientId": 1,
  "model": "7590",
  "problem": "hackintosh",
  "member": {
    "memberId": "2333333333",
    "alias": "滑稽",
    "role": "member",
    "profile": "relaxing",
    "avatar": "",
    "createdBy": "0000000000",
    "gmtCreate": "2022-04-23 15:49:59",
    "gmtModified": "2022-04-30 17:29:46"
  },
  "closedBy": {},
  "status": "accepted",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    },
    {
      "logId": 2,
      "description": "",
      "memberId": "2333333333",
      "action": "accept",
      "gmtCreate": "2022-05-10 11:03:18"
    }
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
}
```

### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 获取全部事件

```
GET /events
```

### 参数

| 名称   | 类型    | in    | 描述      |
| ------ | ------- | ----- | --------- |
| offset | integer | query |           |
| limit  | integer | query | 默认为 30 |

### 示例

#### 请求

```
GET /events
```

#### 响应

```
{
  "eventId": 1,
  "clientId": 1,
  "model": "7590",
  "problem": "hackintosh",
  "member": {
    "memberId": "2333333333",
    "alias": "滑稽",
    "role": "member",
    "profile": "relaxing",
    "avatar": "",
    "createdBy": "0000000000",
    "gmtCreate": "2022-04-23 15:49:59",
    "gmtModified": "2022-04-30 17:29:46"
  },
  "closedBy": {},
  "status": "accepted",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    },
    {
      "logId": 2,
      "description": "",
      "memberId": "2333333333",
      "action": "accept",
      "gmtCreate": "2022-05-10 11:03:18"
    }
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
},
{
  "eventId": 2,
  "clientId": 2,
  "model": "",
  "problem": "下电影",
  "member": {},
  "closedBy": {},
  "status": "open",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    },
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
}
```

### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |


## 获取认证成员接受的指定事件

可以获取到事件中 client 的联系方式

```
GET /member/events/{eventId}
```

### 参数

| 名称          | 类型   | in     | 描述 |
| ------------- | ------ | ------ | ---- |
| Authorization | string | header |      |
| eventId       | number | path   |      |

### 示例

#### 请求
```
GET /member/events/1
```

#### 响应
``` 
{
  "eventId": 1,
  "clientId": 1,
  "model": "7590",
  "phone":"13333333333",
  "qq":"123456789",
  "contact_perference":"qq",
  "problem": "hackintosh",
  "member": {
    "memberId": "2333333333",
    "alias": "滑稽",
    "role": "member",
    "profile": "relaxing",
    "avatar": "",
    "createdBy": "0000000000",
    "gmtCreate": "2022-04-23 15:49:59",
    "gmtModified": "2022-04-30 17:29:46"
  },
  "closedBy": {},
  "status": "accepted",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    },
    {
      "logId": 2,
      "description": "",
      "memberId": "2333333333",
      "action": "accept",
      "gmtCreate": "2022-05-10 11:03:18"
    }
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
}
```

## 获取认证成员接受的全部事件

```
GET /member/events
```
### 参数

| 名称          | 类型   | in     | 描述 |
| ------------- | ------ | ------ | ---- |
| Authorization | string | header |      |

### 示例

#### 请求
```
GET /member/events
```

#### 响应

```
{
  "eventId": 1,
  "clientId": 1,
  "model": "7590",
  "phone":"13333333333",
  "qq":"123456789",
  "contact_perference":"qq",
  "problem": "hackintosh",
  "member": {
    "memberId": "2333333333",
    "alias": "滑稽",
    "role": "member",
    "profile": "relaxing",
    "avatar": "",
    "createdBy": "0000000000",
    "gmtCreate": "2022-04-23 15:49:59",
    "gmtModified": "2022-04-30 17:29:46"
  },
  "closedBy": {},
  "status": "accepted",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    },
    {
      "logId": 2,
      "description": "",
      "memberId": "2333333333",
      "action": "accept",
      "gmtCreate": "2022-05-10 11:03:18"
    }
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
},
{
  "eventId": 2,
  "clientId": 2,
  "model": "",
  "phone":"13333333333",
  "qq":"123456789",
    "contact_perference":"qq",
  "problem": "下电影",
  "member": {},
  "closedBy": {},
  "status": "open",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    },
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
}
```

### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 认证成员接受事件

+ 事件状态变更为`accepted`
+ 事件memberId变为成员id

```
POST /member/events/{eventId}/accept
```

### 参数

| 名称           | 类型   | in     | 描述   |
| -------------- | ------ | ------ | ------ |
| Authorizeation | string | header |        |
| eventId       | string | path   | 事件ID |

### 示例

#### 请求

```
POST /member/events/1/accept
```

#### 响应

```
{
  "eventId": 1,
  "clientId": 1,
  "model": "7590",
  "phone":"13333333333",
  "qq":"123456789",
  "contact_perference":"qq",
  "problem": "hackintosh",
  "member": {
    "memberId": "2333333333",
    "alias": "滑稽",
    "role": "member",
    "profile": "relaxing",
    "avatar": "",
    "createdBy": "0000000000",
    "gmtCreate": "2022-04-23 15:49:59",
    "gmtModified": "2022-04-30 17:29:46"
  },
  "closedBy": {},
  "status": "accepted",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    },
    {
      "logId": 2,
      "description": "",
      "memberId": "2333333333",
      "action": "accept",
      "gmtCreate": "2022-05-10 11:03:18"
    }
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
}
```

### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 认证成员提交事件

+ 事件状态变更为`committed`
+ 提醒管理员审核

```
POST /member/events/{eventId}/commit
```

### 参数

| 名称           | 类型   | in     | 描述     |
| -------------- | ------ | ------ | -------- |
| Authorizeation | string | header |          |
| eventId       | string | path   | 事件ID   |
| content        | string | body   | 维修描述 |

### 示例

#### 请求

```
POST /member/events/1/commit

{
	"content":"重装系统"
}
```

#### 响应

```
{
  "eventId": 1,
  "clientId": 1,
  "model": "7590",
  "phone":"13333333333",
  "qq":"123456789",
  "contact_perference":"qq",
  "problem": "hackintosh",
  "memberId": "2333333333",
  "closedBy": "",
  "status": "committed",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    },
    {
      "logId": 2,
      "description": "",
      "memberId": "2333333333",
      "action": "accept",
      "gmtCreate": "2022-05-10 11:03:18"
    },
    {
      "logId": 3,
      "description": "重装系统",
      "memberId": "2333333333",
      "action": "commit",
      "gmtCreate": "2022-05-10 11:03:18"
    },
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
}
```

### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 认证成员修改事件提交

```
PATCH /member/events/{eventId}/commit
```

### 参数

| 名称           | 类型   | in     | 描述   |
| -------------- | ------ | ------ | ------ |
| Authorizeation | string | header |        |
| eventId       | string | path   | 事件ID |

### 示例

#### 请求

```
PATCH /member/events/1/commit

{
	"content":"重装系统(ghost)"
}
```

#### 响应

```
{
  "eventId": 1,
  "clientId": 1,
  "model": "7590",
  "phone":"13333333333",
  "qq":"123456789",
  "contact_perference":"qq",
  "problem": "hackintosh",
  "memberId": "2333333333",
  "closedBy": "",
  "status": "committed",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    },
    {
      "logId": 2,
      "description": "",
      "memberId": "2333333333",
      "action": "accept",
      "gmtCreate": "2022-05-10 11:03:18"
    },
    {
      "logId": 3,
      "description": "重装系统",
      "memberId": "2333333333",
      "action": "commit",
      "gmtCreate": "2022-05-10 11:05:18"
    },
       {
      "logId": 4,
      "description": "重装系统(ghost)",
      "memberId": "2333333333",
      "action": "alterCommit",
      "gmtCreate": "2022-05-10 12:03:18"
    }
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
}
```

### Http 状态码

| HTTP Status Code | 描述 |
| ---------------- | ---- |
| 200              | OK   |
| 422              |      |



## 认证成员放弃事件

+ 事件状态变更为`open`
+ 清空事件memberId字段

```
DELETE /member/events/{eventId}/accept
```

### 参数

| 名称           | 类型   | in     | 描述   |
| -------------- | ------ | ------ | ------ |
| Authorizeation | string | header |        |
| eventId       | string | path   | 事件ID |

### 示例

#### 请求

```
DELETE /member/events/1/accept
```

#### 响应

```
{
  "eventId": 1,
  "clientId": 1,
  "model": "7590",
  "phone":"13333333333",
  "qq":"123456789",
  "contact_perference":"qq",
  "problem": "hackintosh",
  "memberId": "",
  "closedBy": "",
  "status": "open",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    },
    {
      "logId": 2,
      "description": "",
      "memberId": "2333333333",
      "action": "accept",
      "gmtCreate": "2022-05-10 11:03:18"
    },
    {
      "logId": 3,
      "description": "",
      "memberId": "2333333333",
      "action": "drop",
      "gmtCreate": "2022-05-10 11:03:18"
    },
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
}
```

### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 管理员退回成员事件提交

+ 事件状态变更为`accepted`

```
DELETE /events/{eventId}/commit
```

### 参数

| 名称           | 类型   | in     | 描述   |
| -------------- | ------ | ------ | ------ |
| Authorizeation | string | header |        |
| eventId       | string | path   | 事件ID |

### 示例

#### 请求

```
DELETE /events/events/1/commit
```

#### 响应

```
{
  "eventId": 1,
  "clientId": 1,
  "model": "7590",
  "phone":"13333333333",
  "qq":"123456789",
  "contact_perference":"qq",
  "problem": "hackintosh",
  "memberId": "2333333333",
  "closedBy": "",
  "status": "accepted",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    },
    {
      "logId": 2,
      "description": "",
      "memberId": "2333333333",
      "action": "accept",
      "gmtCreate": "2022-05-10 11:03:18"
    },
    {
      "logId": 3,
      "description": "重装系统",
      "memberId": "2333333333",
      "action": "commit",
      "gmtCreate": "2022-05-10 11:03:18"
    },
    {
      "logId": 4,
      "description": "",
      "memberId": "0000000000",
      "action": "reject",
      "gmtCreate": "2022-05-10 11:03:18"
    }
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
}
```

### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 管理员关闭事件

+ 事件状态变更为`closed`
+ 事件closedBy字段变更为管理员id

```
POST /events/{eventId}/close
```

### 参数

| 名称           | 类型   | in     | 描述   |
| -------------- | ------ | ------ | ------ |
| Authorizeation | string | header |        |
| eventId       | string | path   | 事件ID |

### 示例

#### 请求

```
POST /events/events/1/close
```

#### 响应

```
{
  "eventId": 1,
  "clientId": 1,
  "model": "7590",
  "phone":"13333333333",
  "qq":"123456789",
  "contact_perference":"qq",
  "problem": "hackintosh",
  "memberId": {
    "memberId": "2333333333",
    "alias": "滑稽",
    "role": "member",
    "profile": "relaxing",
    "avatar": "",
    "createdBy": "0000000000",
    "gmtCreate": "2022-04-23 15:49:59",
    "gmtModified": "2022-04-30 17:29:46"
  },
  "closedBy": {
    "memberId": "0000000000",
    "alias": "管理",
    "role": "admin",
    "profile": "",
    "avatar": "",
    "createdBy": "",
    "gmtCreate": "2022-04-30 17:28:42",
    "gmtModified": "2022-04-30 17:28:44"
  },
  "status": "closed",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    },
    {
      "logId": 2,
      "description": "",
      "memberId": "2333333333",
      "action": "accept",
      "gmtCreate": "2022-05-10 11:03:18"
    },
    {
      "logId": 3,
      "description": "重装系统",
      "memberId": "2333333333",
      "action": "commit",
      "gmtCreate": "2022-05-10 11:03:18"
    },
    {
      "logId": 4,
      "description": "",
      "memberId": "0000000000",
      "action": "close",
      "gmtCreate": "2022-05-10 11:03:18"
    }
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
}
```

### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 报修人员创建事件

```
POST /client/event
```
### 参数

| 名称               | 类型   | in     | 描述 |
| ------------------ | ------ | ------ | ---- |
| Authorizeation     | string | header |      |
| phone              | number | body   |      |
| qq                 | number | body   |      |
| contact_perference | string | body   |      |
| problem            | string | body   |      |


### 示例

#### 请求

```
POST /client/event
{
  "phone": "13333333333",
  "qq": "123456789",
  “contact_perference": "phone",
  "problem": "装轮子"
}
```

#### 响应

```
{
  "eventId": 1,
  "clientId": 1,
  "model": "7590",
  "phone": "13333333333",
  "qq": "123456789",
  "contact_perference":"qq",
  "problem": "装轮子",
  "memberId": "",
  "closedBy": "",
  "status": "open",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    }
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
}
```

### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |


## 报修人员更改事件

```
PATCH /client/events/{eventId}
```
### 参数

| 名称           | 类型   | in     | 描述   |
| -------------- | ------ | ------ | ------ |
| Authorizeation | string | header |        |
| eventId       | string | path   | 事件ID |
| phone          | number | body   |        |
| qq             | number | body   |        |
| problem        | string | body   |        |


### 示例

#### 请求

```
POST /client/events/1/
{
  "phone": "13333333333",
  "qq": "123456789",
  "problem": "装轮子"
}
```

#### 响应

```
{
  "eventId": 1,
  "clientId": 1,
  "model": "7590",
  "phone": "13333333333",
  "qq": "123456789",
  "contact_perference":"qq",
  "problem": "装轮子",
  "memberId": "",
  "closedBy": "",
  "status": "open",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    },
    {
      "logId": 2,
      "description": "",
      "memberId": "2333333333",
      "action": "update",
      "gmtCreate": "2022-05-10 11:03:18"
    }
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
}
```

### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |





## 报修人员取消事件

```
DELETE /client/events/{eventId}
```


### 示例

#### 请求

```
DELETE /client/events/1
```

#### 响应
```
{
  "eventId": 1,
  "clientId": 1,
  "model": "7590",
  "phone": "13333333333",
  "qq": "123456789",
  "contact_perference":"qq",
  "problem": "装轮子",
  "memberId": "",
  "closedBy": "",
  "status": "canceled",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    },
    {
      "logId": 2,
      "description": "",
      "memberId": "",
      "action": "cancel",
      "gmtCreate": "2022-05-10 11:03:18"
    }
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
}
```

### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |




## 报修人员获取指定事件

```
GET /client/events/{eventId}
```

| 名称           | 类型   | in     | 描述 |
| -------------- | ------ | ------ | ---- |
| Authorizeation | string | header |      |
| eventId        | number | path   |      |

### 示例

#### 请求

```
GET /client/events/1
```

#### 响应
```
{
  "eventId": 1,
  "clientId": 1,
  "model": "7590",
  "phone": "13333333333",
  "qq": "123456789",
  "contact_perference":"qq",
  "problem": "装轮子",
  "memberId": "",
  "closedBy": "",
  "status": "open",
  "logs": [
    {
      "logId": 1,
      "description": "",
      "memberId": "",
      "action": "create",
      "gmtCreate": "2022-05-10 11:00:26"
    }
  ],
  "gmtCreate": "2022-05-10 10:23:54",
  "gmtModified": "2022-05-12 23:22:44"
}
```


### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 报修人员获取全部事件

```
GET /client/events
```

### 参数

| 名称           | 类型    | in     | 描述      |
| -------------- | ------- | ------ | --------- |
| Authorizeation | string  | header |           |
| offset         | integer | query  |           |
| limit          | integer | query  | 默认为 30 |

### 示例

#### 请求

```
GET /client/events
```

#### 响应
```
[
  {
    "eventId": 1,
    "clientId": 1,
    "model": "7590",
    "phone": "13333333333",
    "qq": "123456789",
    "contact_perference":"qq",
    "problem": "装轮子",
    "memberId": "",
    "closedBy": "",
    "status": "open",
    "logs": [
      {
        "logId": 1,
        "description": "",
        "memberId": "",
        "action": "create",
        "gmtCreate": "2022-05-10 11:00:26"
      }
    ],
    "gmtCreate": "2022-05-10 10:23:54",
    "gmtModified": "2022-05-12 23:22:44"
  }
]
```


### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



## 获取报修人员全部事件

成员获取到一个用户的全部事件

```
GET /clients/{clientId)/events
```

### 参数

| 名称           | 类型    | in     | 描述      |
| -------------- | ------- | ------ | --------- |
| Authorizeation | string  | header |           |
| offset         | integer | query  |           |
| limit          | integer | query  | 默认为 30 |

### 示例

#### 请求

```
GET /clients/1/events
```

#### 响应
```
[
  {
    "eventId": 1,
    "clientId": 1,
    "model": "7590",
    "problem": "装轮子",
    "memberId": "",
    "closedBy": "",
    "status": "open",
    "logs": [
      {
        "logId": 1,
        "description": "",
        "memberId": "",
        "action": "create",
        "gmtCreate": "2022-05-10 11:00:26"
      }
    ],
    "gmtCreate": "2022-05-10 10:23:54",
    "gmtModified": "2022-05-12 23:22:44"
  }
]
```


### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |


