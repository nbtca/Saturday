const { jsonPush, respond } = require("../utils");
const event = require("../models/event");
const element = require("../models/element");

// A:admin U:user E:element CE:current element
// delete (1-3)->0 U
// accept 1->2 E
// cancel   2->1 CE
// submit 2->3 CE,A
// assign 1->2 A
// reject 3->(1,2) A
// close  ?->4 A

class Event {
  constructor() {}
  async get(req, res, next) {
    let data = req.event;
    try {
      let temp = JSON.parse(data.event_log);
      for (let i = 0; i < temp.length; i++) {
        if (temp[i].rid) {
          temp[i].alias = element.get(temp[i].rid).ralias;
        }
        temp[i].time =
          temp[i].time.substring(0, 10) + " " + temp[i].time.substring(11, 19);
        //TODO use action sheet
        if (temp[i].type == "create") {
          temp[i].title = "提交";
          temp[i].icon = "add_circle";
        } else if (temp[i].type == "delete") {
          temp[i].title = "取消";
          temp[i].icon = "remove_circle";
        } else if (temp[i].type == "close") {
          temp[i].title = "完成";
          temp[i].icon = "check_circle";
        } else if (temp[i].type == "update") {
          temp[i].title = "更新";
          temp[i].icon = "update_circle";
        } else if (temp[i].type == "accept") {
          temp[i].title = "接受";
          temp[i].icon = "accept_circle";
        } else if (temp[i].type == "cancel") {
          temp[i].title = "放弃";
          temp[i].icon = "sentiment_very_dissatisfied";
        } else if (temp[i].type == "reject") {
          temp[i].title = "退回";
          temp[i].icon = "sentiment_very_dissatisfied";
        } else if (temp[i].type == "assign") {
          temp[i].title = "指派";
          temp[i].icon = "accept_circle";
        } else if (temp[i].type == "submit") {
          temp[i].title = "提交维修";
          temp[i].icon = "sentiment_very_dissatisfied";
        }
      }
      data.event_log = temp;
      data.repair_description = JSON.parse(data.repair_description);
      respond(res, 0, "Success", data);
      return;
    } catch (error) {
      next(error);
    }
  }
  async getAll(req, res, next) {
    try {
      let data = await event.get();
      respond(res, 0, "Success", data);
    } catch (error) {
      next(error);
    }
  }
  async creat(req, res, next) {
    let eventLog = JSON.stringify([
      {
        type: "create",
        time: new Date(),
      },
    ]);
    try {
      await event.creat({
        uid: req.body.uid,
        model: req.body.model,
        qq: req.body.qq,
        phone: req.body.phone,
        preference: req.body.preference,
        description: req.body.description,
        eventLog: eventLog,
      });
      respond(res, 0);
    } catch (error) {
      next(error);
    }
  }
  async update(req, res, next) {
    try {
      let thisEvent = req.event;
      if (thisEvent.status <= 1) {
        let addEventLog = {
          type: "update",
          time: new Date(),
        };
        eventLog = jsonPush(thisEvent.event_log, addEventLog);
        await event.update({
          model: req.body.model,
          phone: req.body.phone,
          qq: req.body.qq,
          description: req.body.description,
          preference: req.body.preference,
          eventLog: eventLog,
          eid: req.body.eid,
        });
        respond(res, 0);
      } else {
        // TODO error code
        respond(res, 220, "Event has been accepted or deleted");
      }
    } catch (error) {
      next(error);
    }
  }
  async accept(req, res, next) {
    let eid = req.body.eid;
    try {
      let thisEvent = req.event;
      if (thisEvent.rid == null && thisEvent.status == 0) {
        let rid = res.locals.data.rid;
        let addEventLog = {
          type: "accept",
          time: new Date(),
          rid: rid,
        };
        let eventLog = jsonPush(thisEvent.event_log, addEventLog);
        await event.accept(rid, eventLog, eid);
        respond(res, 0);
      } else {
        respond(res, 220, "Event has been accepted or deleted");
      }
    } catch (error) {
      next(error);
    }
  }
  async submit(req, res, next) {
    try {
      let thisEvent = req.event;
      let eventLog = {
        type: "submit",
        time: new Date(),
        rid: res.locals.data.rid,
        description: req.body.description,
      };
      let description = {
        time: new Date(),
        rid: res.locals.data.rid,
        description: req.body.description,
      };
      eventLog = jsonPush(thisEvent.event_log, eventLog);
      description = jsonPush(thisEvent.repair_description, description);
      await event.submit(eventLog, description, req);
      respond(res, 0);
    } catch (err) {
      next(err);
    }
  }
  async cancel(req, res, next) {
    let eid = req.body.eid;
    try {
      let thisEvent = req.event;
      if (thisEvent.status == 1) {
        let addeventLog = {
          type: "cancel",
          time: new Date(),
          rid: res.locals.data.rid,
        };
        let eventLog = jsonPush(thisEvent.event_log, addeventLog);
        await event.cancel(eventLog, eid);
        respond(res, 0);
      } else {
        respond(res, 220, "Event has been cancel or closed");
      }
    } catch (err) {
      next(err);
    }
  }
  async close(req, res, next) {
    let eid = req.body.eid;
    let aid = res.locals.data.aid;
    let status;
    try {
      let thisEvent = req.event;
      if (thisEvent.status == 2) {
        let addeventLog = {
          type: "",
          time: new Date(),
          aid: aid,
        };
        if (req.body.accept) {
          addeventLog.type = "close";
          status = 3;
        } else {
          addeventLog.type = "reject";
          status = 0;
        }
        let eventLog = jsonPush(thisEvent.event_log, addeventLog);
        await event.close(aid, eventLog, status, eid);
        respond(res, 0);
      } else {
        respond(res, 251, "Event status error");
      }
    } catch (err) {
      next(err);
    }
  }
  async assign(req, res, next) {
    let eid = req.body.eid;
    let aid = res.locals.data.aid;
    let rid = req.body.rid;
    try {
      let thisEvent = req.event;
      let addeventLog = {
        type: "assign",
        time: new Date(),
        aid: aid,
        rid: rid,
      };
      let eventLog = jsonPush(thisEvent.event_log, addeventLog);
      await event.assign(rid, eventLog, eid);
      respond(res, 0);
    } catch (err) {
      next(err);
    }
  }
}
module.exports = new Event();
