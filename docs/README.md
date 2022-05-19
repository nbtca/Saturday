---
layout: default
title: 
nav_order: 1
has_children: false 
---
# 简介 & 规划
> under construction

## TODOS

- [ ] 设计模式
- [ ] API设计
  - [x] 错误码
  - [ ] 写文档
  - [x] restful
- [x] 项目管理
  - [x] 工具选择
  - [ ] 进度/任务安排
- [ ] 测试
- [x] 参数校验
- [x] 日志
- [x] 字段权限
- [x] 结构
- [x] 分离配置
- [ ] 错误处理

## 进度安排

### 4.22前

+ 确认设计模式，数据库设计
+ 完成文档
+ 确定协作工具（工作流程）
+ 完成示例接口及测试
+ 确认业务逻辑

### 5.6前

+ 完成member接口
+ 完成测试
+ 回顾

### 5.20前

+ 完成余下接口



## 字典

| 资源                 | 名称                 | 描述 |
| -------------------- | -------------------- | ---- |
| 成员(一般)           | member               |      |
| 成员((一般)(未激活)  | member_not_activated |     |
| 成员(管理员)         | admin                |      |
| 成员(管理员)(未激活) | admin_not_activated  |  激活后成为管理员   |
| 成员(封存)           | member_archived      |  无写权限    |
| 维修事件             | event                |      |
| 报修人员             | client               |      |
| 维修事件创建         | create               |      |
| 维修事件放弃         | drop                 |      |
| 维修事件接受         | accept               |      |
| 维修事件提交审核     | commit               |      |
| 维修事件审核通过     | close                |      |
| 维修事件审核退回     | reject               |      |
| 维修事件指派         | assign               |      |
|                      |                      |      |

## API设计

```javascript
// private info => contacts,password,adress...
// ---Member(成员)---
GET /members // get all members' public info
GET /members/{member_id} // get a member's public info
POST /members/token // create token(login)

// following requires Authorization in header
PUT /member/activate // role member_inactive admin_inactive
    // change member status to activated
    // 要求成员先前role为未激活

    //  following  requires role not contains not_activated
GET /member // get a member's private info // role member
PUT /member // update member info // role member
PUT /member/avatar // change avatar // role member

	// following requires role to be admin
POST /members // bulk create // role admin
POST /members/{member_id} //create member  // role admin
PUT /members/{member_id} // update member info  // role admin

// ---Event(事件)---
GET /events // get all events' public info
GET /events/{event_id} // get an event's public info
// delete 
// following requires Authorization in header
PUT /events/{event_id}/accept // accept event
    // 要求事件状态为未接受
  // following requires to be current member (memberId==event.memberId)
GET /member/evnets // get the private info of all events that is accepted by member
GET /member/evnets/{member_id} // get the private info of the event
POST /member/events/{evemt_id}/commit // commit event for admin approval (event status: accepted=>committed)
PUT /member/events/{member_id}/commit // alter commit (event status: committed)
DELETE /member/events/{event_id}/accept // drop event (event status: accepted,committed)

	// following requires role to be admin
DELETE /events/{event_id}/commit // reject commit (event status: committed=>accepted)
PUT /events/{event_id}/close // close event (event status: accepted=>closed)
PUT /events/{event_id}/{member_id} // assign event to member (event status: created => accepted(by assigned member))

//---client(报修人员)---
GET /clinets
GET /clients/{client_id}
POST /client
POST /client/token
```
