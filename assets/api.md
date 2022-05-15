# API 文档



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



## 事件

### 基础

![](https://clas-bucket.oss-cn-hangzhou.aliyuncs.com/uPic/GhhXcd.png)

#### 事件状态表(status)

| 状态名 |           | 描述                                 |
| ------ | --------- | ------------------------------------ |
| 待处理 | open      | 维修事件未被成员接受                 |
| 取消   | cancelled | 维修事件被用户取消，不需要再进行处理 |
| 受理   | accepted  | 维修事件已被成员接受                 |
| 待审核 | committed | 成员提交了维修描述，管理员尚未审核   |
| 关闭   | closed    | 维修事件已解决，不能再更改该事件     |

#### 事件行为表(action)

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

### 获取指定事件

```
GET /events/{event_id}
```

#### 参数

| 名称     | 类型   | in   | 描述 |
| -------- | ------ | ---- | ---- |
| event_id | String | path | 学号 |

#### 示例

##### 请求

```
GET /members/event_id
```

##### 响应

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

#### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



### 获取全部事件

```
GET /events
```



### 认证成员接受事件

```
POST /member/events/{event_id}/accept
```

#### 参数

| 名称           | 类型   | in     | 描述 |
| -------------- | ------ | ------ | ---- |
| Authorizeation | string | header |      |
| event_id       | String | path   | 学号 |

#### 示例

##### 请求

```
GET /members/event_id
```

##### 响应

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

#### Http 状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 422              | Unprocessable Entity |



### 获取认证成员接受的指定事件

```
GET /member/events/{event_id}
```



### 获取认证成员接受的全部事件

```
GET /member/events
```



### 认证成员提交事件

```
POST /member/events/{event_id}/commit
```



### 认证成员修改事件提交

```
PATCH /member/events/{event_id}/commit
```



### 认证成员放弃事件

```
DELETE /member/events/{event_id}/accept
```



### 管理员退回成员事件提交

```
DELETE /events/{event_id}/commit
```



### 管理员关闭事件

```
POST /events/{event_id}/close
```



### 报修人员创建事件

```
POST /clients/event
```



### 报修人员更改事件

```
POST /clients/events/{event_id}
```



### 报修人员取消事件

```
POST /clients/events/{event_id}
```



### 获取报修人员事件

```
GET /client/events/{event_id}
```



### 获取报修人员全部事件

```
GET /client/events
```





## 报修人员



### 获取指定报修人员

```
GET /clients/{client_id}
```



### 获取全部报修人员

```
GET /clients/{client_id}
```



### 创建报修人员

```
POST /client
```



### 创建报修人员Token

```
POST /client/{client_id}/token
```
