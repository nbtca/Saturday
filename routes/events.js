const express = require("express");
const router = express.Router();
const { jsonPush, respond } = require("../utils");
const { isEidVaild } = require("../middleware/event");
const { mysql } = require("../config/config");
const event = require("../models/event");
const element = require("../models/element");

router.use("/", isEidVaild);

router.get("/", async (req, res, next) => {
  let data = await event.get();
  respond(res, 0, "Success", data);
});

router.use("/:eid", isEidVaild);

router.get("/:eid", async (req, res, next) => {
  let data = req.event;
  let temp = JSON.parse(data.event_log);
  for (let i = 0; i < temp.length; i++) {
    if (temp[i].rid) {
      temp[i].alias = element.getAlias(temp[i].rid);
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



router.put("/:eid", async (req, res, next) => {
  status = req.body.status;
  if (status == 0) {
  } else if (status == 1) {
  } else if (status == 1) {
  } else if (status == 1) {
  } else if (status == 1) {
  }
});

router.post("/accept", async (req, res, next) => {
  let data = await event.get(req.body.eid);
  if (data[0].rid == null && data[0].status == 0) {
    let addeventLog = {
      type: "accept",
      time: new Date(),
      rid: res.locals.data.rid,
    };
    eventLog = jsonPush(data[0].event_log, addeventLog);
    await event.accept();
    respond(res, 0);
  } else {
    respond(res, 220, "Event has been accepted or deleted");
  }
});

//TODO move to middleware
// router.use("/event/edit", async (req, res, next) => {
//   let dbResults;
//   try {
//     dbResults = await mysql.query("SELECT rid FROM `event` WHERE eid=?", [
//       req.body.eid,
//     ]);
//   } catch (err) {
//     next(err);
//   }
//   await mysql.end();
//   if ((req.rid = res.locals.data.rid)) {
//     next();
//   } else {
//     respond(res, 230, "No edit permission");
//   }
// });

router.post("/edit/submit", async (req, res, next) => {
  let returnObj = {
    resultCode: null,
    resultMsg: null,
  };
  let dbResults;
  try {
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
    dbResults = await mysql.query(
      "SELECT event_log,repair_description FROM `event` WHERE eid=?",
      [req.body.eid]
    );
    let eventLog = jsonPush(dbResults[0].event_log, addeventLog);
    description = jsonPush(dbResults[0].repair_description, description);
    await mysql.query(
      "UPDATE `event` SET event_log=?,repair_description=?,status=2 WHERE eid=?",
      [eventLog, description, req.body.eid]
    );
    returnObj.resultCode = 0;
    returnObj.resultMsg = "Success";
  } catch (err) {
    next(err);
  }
  await mysql.end();
  res.send(returnObj);
});

// 1
router.post("/edit/cancelEvent", async (req, res, next) => {
  let returnObj = {
    resultCode: null,
    resultMsg: null,
  };
  let dbResults;
  try {
    dbResults = await mysql.query(
      "SELECT rid,status,event_log FROM `event` WHERE eid=?",
      [req.body.eid]
    );
    if (dbResults[0].rid == res.locals.data.rid && dbResults[0].status == 1) {
      let addeventLog = {
        type: "cancel",
        time: new Date(),
        rid: res.locals.data.rid,
      };
      let eventLog = jsonPush(dbResults[0].event_log, addeventLog);

      // eventLog = JSON.parse(eventLog[0].event_log);
      // eventLog.push(addeventLog);
      // eventLog = JSON.stringify(eventLog);
      await mysql.query(
        "UPDATE `event` SET rid=?,event_log=?,status=? WHERE eid=?",
        [null, eventLog, 0, req.body.eid]
      );
      returnObj.resultCode = 0;
      returnObj.resultMsg = "Success";
    } else {
      returnObj.resultCode = 220;
      returnObj.resultMsg = "Event has been cancel or closed";
    }
  } catch (err) {
    next(err);
  }
  await mysql.end();
  res.send(returnObj);
});

router.use("/manage", async (req, res, next) => {
  let returnObj = {
    resultCode: null,
    resultMsg: null,
  };
  let dbResults;
  try {
    dbResults = await mysql.query("SELECT aid FROM admin WHERE rid=?", [
      res.locals.data.rid,
    ]);
  } catch (err) {
    next(err);
  }
  await mysql.end();
  res.locals.data.aid = dbResults[0].aid;
  if (dbResults[0].aid) {
    console.log("经过了manage1");
    next();
  } else {
    console.log("经过了manage2");
    returnObj.resultCode = 250;
    returnObj.resultMsg = "No admin permission";
    res.send(returnObj);
  }
});

router.post("/manage/checkEvent", async (req, res, next) => {
  let returnObj = {
    resultCode: null,
    resultMsg: null,
  };
  let dbResults;
  let status;
  try {
    dbResults = await mysql.query(
      "SELECT status,event_log FROM `event` WHERE eid=?",
      [req.body.eid]
    );
    if (dbResults[0].status == 2) {
      let addeventLog = {
        type: "",
        time: new Date(),
        aid: res.locals.data.aid,
      };
      if (req.body.accept) {
        addeventLog.type = "close";
        status = 3;
      } else {
        addeventLog.type = "reject";
        status = 0;
      }
      let eventLog = jsonPush(dbResults[0].event_log, addeventLog);
      // eventLog = JSON.parse(dbResults[0].event_log);
      // eventLog.push(addeventLog);
      // eventLog = JSON.stringify(eventLog);
      await mysql.query(
        "UPDATE `event` SET aid=?,event_log=?,status=? WHERE eid=?",
        [res.locals.data.aid, eventLog, status, req.body.eid]
      );
      returnObj.resultCode = 0;
      returnObj.resultMsg = "Success";
    } else {
      returnObj.resultCode = 251;
      returnObj.resultMsg = "Event status error";
    }
  } catch (err) {
    next(err);
  }
  await mysql.end();
  res.send(returnObj);
});

router.post("/manage/assignEvent", async (req, res, next) => {
  let returnObj = {
    resultCode: null,
    resultMsg: null,
  };
  let dbResults;
  try {
    dbResults = await mysql.query("SELECT event_log FROM `event` WHERE eid=?", [
      req.body.eid,
    ]);
    let addeventLog = {
      type: "assign",
      time: new Date(),
      aid: res.locals.data.aid,
      rid: req.body.rid,
    };
    let eventLog = jsonPush(dbResults[0].event_log, addeventLog);
    // eventLog = JSON.parse(dbResults[0].event_log);
    // eventLog.push(addeventLog);
    // eventLog = JSON.stringify(eventLog);
    await mysql.query(
      "UPDATE `event` SET rid=?,event_log=?,status=? WHERE eid=?",
      [req.body.rid, eventLog, 1, req.body.eid]
    );
    returnObj.resultCode = 0;
    returnObj.resultMsg = "Success";
  } catch (err) {
    next(err);
  }
  await mysql.end();
  res.send(returnObj);
});

module.exports = router;
