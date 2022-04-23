# API文档
## 成员（Member）

### 获取指定成员

```
GET /members/{member_id}
```

#### 参数

| 名称      | 类型    | in   | 描述 |
| --------- | ------- | ---- | ---- |
| member_id | integer | path |      |

#### 示例

#### 请求

```
GET /members/2333333333
```

#### 响应

```
{
  "member_id": "2333333333",
  "alias": "滑稽",
  "name": "滑稽",
  "section": "计算机233",
  "profile": "",
  "phone": "",
  "qq": "",
  "avatar": "",
  "created_by": "",
  "gmt_create": "2022-04-17T19:35:55.000Z",
  "gmt_modified": "2022-04-17T19:35:55.000Z",
  "role": "member"
}
```

#### Http状态码

| HTTP Status Code | 描述      |
| ---------------- | --------- |
| **200**          | OK        |
| 404              | Resource not found|
### 获取指定成员

```
GET /members/{member_id}
```

#### 参数

| 名称      | 类型    | in   | 描述 |
| --------- | ------- | ---- | ---- |
| member_id | integer | path |      |

#### 示例

#### 请求

```
GET /members
```

#### 响应

```
[
  {
    "member_id": "2333333333",
    "alias": "滑稽",
    "name": "滑稽",
    "section": "计算机233",
    "profile": "",
    "phone": "",
    "qq": "",
    "avatar": "",
    "created_by": "",
    "gmt_create": "2022-04-17T19:35:55.000Z",
    "gmt_modified": "2022-04-17T19:35:55.000Z",
    "role": "member"
  }
]
```

#### Http状态码

| HTTP Status Code | 描述      |
| ---------------- | --------- |
| **200**          | OK        |
| 404              | Resource not found|



### 创建成员
+ 需要身份为管理员
+ member_id（学号）需为唯一

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
| alias `可选`   | string  | body   | 昵称 |
| phone `可选`   | string  | body   |
| qq    `可选`   | string  | body   |
| avata `可选`   | string  | body   | 头像 |

#### 示例 

#### 请求

```
POST /members/1234567890

{
  "member_id": "2333333333",
  "alias": "滑稽",
  "name": "滑稽",
  "section": "计算机233",
  "profile": "relaxing",
  "phone": "12356839487",
  "qq": "123456",
}
```

#### 响应

```
{
  "member_id": "2333333333",
  "alias": "滑稽",
  "name": "滑稽",
  "section": "计算机233",
  "profile": "relaxing",
  "phone": "12356839487",
  "qq": "123456",
  "avatar": "",
  "created_by": "",
  "gmt_create": "2022-04-17T19:35:55.000Z",
  "gmt_modified": "2022-04-17T19:35:55.000Z",
  "role": "member"
}
```



#### Http状态码

| HTTP Status Code | 描述                 |
| ---------------- | -------------------- |
| 200              | OK                   |
| 404              | Resource not found           |
| 422              | Unprocessable Entity |



