const { jsonPush, respond } = require("../utils/utils");
const { actionSheet } = require("../config/config");
const event = require("../models/event");
const element = require("../models/element");
const { Action } = require("../utils/action");
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
      for (let item of temp) {
        if (item.rid) {
          item.alias = element.get(item.rid).ralias;
        }
        item.time =
          item.time.substring(0, 10) + " " + item.time.substring(11, 19);
        item.icon = actionSheet[item.type].icon;
        item.title = actionSheet[item.type].title;
      }
      data.event_log = temp;
      data.repair_description = JSON.parse(data.repair_description);
      respond(res, 0, "Success", data);
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
  async create(req, res, next) {
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

  async delete(req, res, next) {
    let thisEvent = req.event;
    try {
      let del = Action("delete");
      del.perform(thisEvent);
      await event.update(thisEvent);
      respond(res, 0);
      // if (thisEvent.status <= 1) {
      //   let addEventLog = {
      //     type: "delete",
      //     time: new Date(),
      //   };
      //   let eventLog = jsonPush(thisEvent.event_log, addEventLog);
      //   event.delete(eventLog, req.body.eid);
      //   respond(res, 0);
      // } else {
      //   respond(res, 123, "no permission");
      // }
    } catch (error) {
      next(error);
    }
  }

  async accept(req, res, next) {
    try {
      let thisEvent = req.event;
      let accept = Action("accept");
      thisEvent = accept.perform(thisEvent, {
        rid: res.locals.data.rid,
      });
      await event.update(thisEvent);
      respond(res, 0);
    } catch (error) {
      next(error);
    }
  }

  async submit(req, res, next) {
    try {
      let thisEvent = req.event;
      let submit = Action("submit");
      let repair_description = {
        time: new Date(),
        rid: res.locals.data.rid,
        description: req.body.description,
      };
      repair_description = jsonPush(
        thisEvent.repair_description,
        repair_description
      );
      thisEvent = submit.perform(thisEvent, {
        repair_description: repair_description,
        description: req.body.description,
      });
      await event.update();
      respond(res, 0);
    } catch (err) {
      next(err);
    }
  }

  async drop(req, res, next) {
    try {
      let thisEvent = req.event;
      let drop = Action("drop");
      thisEvent = drop.perform(thisEvent, {
        rid: res.locals.data.rid,
      });
      thisEvent.rid = null;
      await event.update(thisEvent);
      respond(res, 0);
    } catch (err) {
      next(err);
    }
  }

  async close(req, res, next) {
    let aid = res.locals.data.aid;
    let thisEvent = req.event;
    try {
      if (req.body.accept) {
        let close = Action("close");
        thisEvent = close.perform(thisEvent, {
          aid: aid,
        });
      } else {
        let reject = Action("reject");
        thisEvent = reject.perform(thisEvent, {
          aid: aid,
        });
      }
      await event.update(thisEvent);
      respond(res, 0);
    } catch (err) {
      next(err);
    }
  }

  async assign(req, res, next) {
    try {
      let thisEvent = req.event;
      let assign = Action("assign");
      thisEvent = assign.perform(thisEvent, {
        aid: res.locals.data.aid,
        rid: req.body.rid,
      });
      await event.update(thisEvent);
      respond(res, 0);
    } catch (err) {
      next(err);
    }
  }
}

module.exports = new Event();
