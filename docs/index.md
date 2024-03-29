---
layout: default
title: 
nav_order: 1
has_children: false 
---

# 简介 & 规划

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
GET /events/{eventId} // get an event's public info
// delete 
// following requires Authorization in header
PUT /events/{eventId}/accept // accept event
    // 要求事件状态为未接受
  // following requires to be current member (memberId==event.memberId)
GET /member/evnets // get the private info of all events that is accepted by member
GET /member/evnets/{member_id} // get the private info of the event
POST /member/events/{evemt_id}/commit // commit event for admin approval (event status: accepted=>committed)
PUT /member/events/{member_id}/commit // alter commit (event status: committed)
DELETE /member/events/{eventId}/accept // drop event (event status: accepted,committed)

	// following requires role to be admin
DELETE /events/{eventId}/commit // reject commit (event status: committed=>accepted)
PUT /events/{eventId}/close // close event (event status: accepted=>closed)
PUT /events/{eventId}/{member_id} // assign event to member (event status: created => accepted(by assigned member))

//---client(报修人员)---
GET /clinets
GET /clients/{clientId}
POST /client
POST /client/token
```

