const express = require("express");
const router = express.Router();
const { jsonPush, respond } = require("../utils");
const { isAdmin } = require("../middleware/auth");
const { isEidVaild, isCurrentUser } = require("../middleware/event");
const event = require("../models/event");
const element = require("../models/element");

router.use("/", isEidVaild);

router.get("/", async (req, res, next) => {
  try {
    let data = await event.get();
    respond(res, 0, "Success", data);
  } catch (error) {
    next(error);
  }
});

router.use("/:eid", isEidVaild);

router.get("/:eid", async (req, res, next) => {
  let data = req.event;
  let temp = JSON.parse(data.event_log);
  for (let i = 0; i < temp.length; i++) {
    if (temp[i].rid) {
      temp[i].alias = element.get(temp[i].rid).ralias;
    }
    temp[i].time =
      temp[i].time.substring(0, 10) + " " + temp[i].time.substring(11, 19);
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
  respond(res, 0, "Success", data);
});

// A:admin U:user E:element CE:current element
// delete (1-3)->0 U
// accept 1->2 E
// cancel   2->1 CE
// submit 2->3 CE,A
// assign 1->2 A
// reject 3->(1,2) A
// close  ?->4 A

// router.put("/:eid", async (req, res, next) => {
//   status = req.body.status;
//   if (status == 0) {
//   } else if (status == 1) {
//   } else if (status == 1) {
//   } else if (status == 1) {
//   } else if (status == 1) {
//   }
// });

router.post("/accept", async (req, res, next) => {
  let eid = req.body.eid;
  try {
    let data = await event.get(eid);
    if (data[0].rid == null && data[0].status == 0) {
      let rid = res.locals.data.rid;
      let addeventLog = {
        type: "accept",
        time: new Date(),
        rid: rid,
      };
      eventLog = jsonPush(data[0].event_log, addeventLog);
      await event.accept(rid, eventLog, eid);
      respond(res, 0);
    } else {
      respond(res, 220, "Event has been accepted or deleted");
    }
  } catch (error) {
    next(err);
  }
});

router.post("/submit", isCurrentUser, async (req, res, next) => {
  try {
    let dbResults = event.get(req.body.eid);
    let addeventLog = {
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
    let eventLog = jsonPush(dbResults.event_log, addeventLog);
    description = jsonPush(dbResults.repair_description, description);
    await event.submit(eventLog, description, req);
    respond(res, 0);
  } catch (err) {
    next(err);
  }
});

router.post("/cancel", isCurrentUser, async (req, res, next) => {
  let eid = req.body.eid;
  try {
    let dbResults = event.get(eid);
    if (dbResults.status == 1) {
      let addeventLog = {
        type: "cancel",
        time: new Date(),
        rid: res.locals.data.rid,
      };
      let eventLog = jsonPush(dbResults[0].event_log, addeventLog);
      await event.cancel(eventLog, eid);
      respond(res, 0);
    } else {
      respond(res, 220, "Event has been cancel or closed");
    }
  } catch (err) {
    next(err);
  }
});

router.post("/close", isAdmin, async (req, res, next) => {
  let eid = req.body.eid;
  let aid = res.locals.data.aid;
  let status;
  try {
    let dbResults = event.get(eid);
    if (dbResults.status == 2) {
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
      let eventLog = jsonPush(dbResults.event_log, addeventLog);
      await event.close(aid, eventLog, status, eid);
      respond(res, 0);
    } else {
      respond(res, 251, "Event status error");
    }
  } catch (err) {
    next(err);
  }
});

router.post("/assign", isAdmin, async (req, res, next) => {
  let eid = req.body.eid;
  let aid = res.locals.data.aid;
  let rid = req.body.rid;
  try {
    let dbResults = event.get(eid);
    let addeventLog = {
      type: "assign",
      time: new Date(),
      aid: aid,
      rid: rid,
    };
    let eventLog = jsonPush(dbResults.event_log, addeventLog);
    await assign(rid, eventLog, eid);
    respond(res, 0);
  } catch (err) {
    next(err);
  }
});

module.exports = router;
