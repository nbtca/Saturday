const actionSheet = require("../../config/actionSheet");
const { jsonPush } = require("../../utils/utils");
class Action {
  constructor(type, role) {
    this.ref = actionSheet[type];
    // this.isRoleValid(role);
  }
  isRoleValid(role) {
    let flag = false;
    for (let r of role) {
      console.log(r);
      if (this.ref.auth.role.indexOf(r) != -1) {
        console.log("valid");
        flag = true;
        break;
      }
    }
    if (!flag) {
      throw new Error("invalid role");
    }
    return flag;
  }
  isFormerStatusValid(event) {
    let status = event.status;
    if (this.ref.auth.formerStatus.indexOf(status) != -1) {
      return true;
    } else {
      throw new Error("invalid former status");
    }
  }
  constructActionLog(data) {
    let log = {
      type: this.ref.type,
      time: new Date(),
    };
    for (let item of this.ref.logItem) {
      if (data[item] != null) {
        log[item] = data[item];
      } else {
        throw new Error("missing item:" + item);
      }
    }
    return log;
  }
  perform(event, data) {
    this.isFormerStatusValid(event);
    event.event_log = jsonPush(event.event_log, this.constructActionLog(data));
    event.status = this.ref.targetStatus;
    for (let item of this.ref.alterItem) {
      if (data[item] != null) {
        event[item] = data[item];
      } else {
        throw new Error("missing item:" + item);
      }
    }
    return event;
  }
}

exports.appendLog = (type, event, data) => {
  let action = new Action(type);
  return action.perform(event, data);
};

// let role = ["currentUser", "element"];
// let role = "element";
// let action = new Action("accept", role);
// let event = {
//   status: 0,
// };
// try {
//   let nevent = action.perform(event, {
//     rid: "123",
//   });
//   console.log(nevent);
// } catch (error) {
//   console.error(error);
// }
