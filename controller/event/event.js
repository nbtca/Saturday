const { jsonPush, respond, dateToStr } = require("../../utils/utils");
const { actionSheet } = require("../../config/config");
const event = require("../../models/event");
const ElementModel = require("../../models/ElementModel");
const EventModel = require("../../models/EventModel");
const { appendLog } = require("./action");
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
        // if (item.rid) {
        //   item.alias = await ElementModel.findByFilter({ ralias }, { rid: req.params.rid });
        // }
        item.time = item.time.substring(0, 10) + " " + item.time.substring(11, 19);
        item.icon = actionSheet[item.type].icon;
        item.title = actionSheet[item.type].title;
      }
      data.event_log = temp;
      data.repair_description = JSON.parse(data.repair_description);

      respond(res, 0, "Success", data);
    } catch (error) {
      console.log(error);
    }
  }
  async getAll(req, res, next) {
    try {
      await EventModel.findAll(["eid", "user_description", "status", "model", "rid", "gmt_create"]).then(result => {
        for (let item of result) {
          item.dataValues.gmt_create = dateToStr(item.gmt_create, "time");
        }
        respond(res, 0, "Success", result);
      });
    } catch (error) {
      console.log(error);
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
      await EventModel.create({
        uid: req.body.uid,
        model: req.body.model,
        eqq: req.body.qq,
        ephone: req.body.phone,
        preference: req.body.preference,
        user_description: req.body.description,
        event_log: eventLog,
        gmt_create: new Date(),
        gmt_modified: new Date(),
      });
      respond(res, 0);
    } catch (error) {
      console.log(error);
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
        await EventModel.update(
          {
            model: req.body.model,
            ephone: req.body.phone,
            eqq: req.body.qq,
            user_description: req.body.description,
            preference: req.body.preference,
            event_log: eventLog,
            eid: req.body.eid,
          },
          { eid: thisEvent.eid }
        );
        respond(res, 0);
      } else {
        respond(res, 220, "Event has been accepted or deleted");
      }
    } catch (error) {
      next(error);
    }
  }

  async delete(req, res, next) {
    let thisEvent = req.event;
    try {
      appendLog("delete", thisEvent);
      await EventModel.update(thisEvent, { eid: thisEvent.eid });
      respond(res, 0);
    } catch (error) {
      next(error);
    }
  }

  async accept(req, res, next) {
    try {
      let thisEvent = req.event;
      appendLog("accept", thisEvent, { rid: res.locals.data.rid });
      console.log(thisEvent);
      await EventModel.update(thisEvent, { eid: thisEvent.eid });
      respond(res, 0);
    } catch (error) {
      console.log(error);
    }
  }
  //submit repair
  async submit(req, res, next) {
    try {
      console.log(req.body);
      let thisEvent = req.event;
      let repair_description = {
        time: new Date(),
        rid: res.locals.data.rid,
        description: req.body.description,
      };
      repair_description = jsonPush(thisEvent.repair_description, repair_description);
      appendLog("submit", thisEvent, {
        repair_description: repair_description,
        description: req.body.description,
      });
      await EventModel.update(thisEvent, { eid: thisEvent.eid });
      respond(res, 0);
    } catch (err) {
      next(err);
    }
  }

  async drop(req, res, next) {
    try {
      let thisEvent = req.event;
      appendLog("drop", thisEvent, { rid: res.locals.data.rid });
      thisEvent.rid = null;
      await EventModel.update(thisEvent, { eid: thisEvent.eid });
      respond(res, 0);
    } catch (err) {
      next(err);
    }
  }

  async close(req, res, next) {
    let aid = res.locals.data.aid;
    let thisEvent = req.event;
    try {
      appendLog(req.body.accept ? "close" : "reject", thisEvent, { aid: aid });
      await EventModel.update(thisEvent, { eid: thisEvent.eid });
      respond(res, 0);
    } catch (err) {
      next(err);
    }
  }

  async assign(req, res, next) {
    try {
      let thisEvent = req.event;
      appendLog("assign", thisEvent, {
        aid: res.locals.data.aid,
        rid: req.body.rid,
      });
      await EventModel.update(thisEvent, { eid: thisEvent.eid });
      respond(res, 0);
    } catch (err) {
      next(err);
    }
  }
}

module.exports = new Event();
