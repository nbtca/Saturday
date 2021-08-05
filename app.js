var createError = require("http-errors");
var express = require("express");
var path = require("path");
var cookieParser = require("cookie-parser");
var logger = require("morgan");
var host = "rm-uf6s9l8ep4131lzt9go.mysql.rds.aliyuncs.com";
const jwt = require("jsonwebtoken");

const mysql = require("serverless-mysql")({
  config: {
    host: host,
    user: "high_admin",
    password: "02000163",
    database: "repairteam_build",
  },
});

var indexRouter = require("./routes/admin");
var usersRouter = require("./routes/user");
var login = require("./routes/login");

var app = express();

// view engine setup
app.set("views", path.join(__dirname, "views"));
app.set("view engine", "jade");

app.use(logger("dev"));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, "public")));

app.post("/login", async (req, res, next) => {
  let rid = req.body.id;
  let password = req.body.password;
  let returnObj = {
    resultCode: null,
    resultMsg: null,
    data: null,
  };
  try {
    let dbResults = await mysql.query(
      "SELECT rpassword,ralias FROM repairelements WHERE rid=?",
      [rid]
    );
    if (dbResults.length == 0) {
      returnObj.resultCode = 101;
      returnObj.resultMsg = "no such user";
    } else {
      if (dbResults[0].rpassword == password) {
        let isPasswordEmpty = !password ? true : false;
        let isAdmin = await mysql.query("SELECT aid FROM admin WHERE rid=?", [
          rid,
        ]);
        let role = isAdmin ? "admin" : "element";
        let data = {
          rid: rid,
          role: role,
          aid: isAdmin[0].aid,
        };
        let cert = "02000163";
        let token = jwt.sign(
          {
            exp: Math.floor(Date.now() / 1000) + 24 * 60 * 60,
            data: data,
          },
          cert
        );
        returnObj.data = {
          token: token,
          alias: dbResults[0].ralias,
          rid: rid,
          role: role,
          isPasswordEmpty: isPasswordEmpty,
        };
        await mysql.query(
          "UPDATE repairelements SET  gmt_modified=SYSDATE() WHERE rid=?",
          [rid]
        );
        returnObj.resultCode = 0;
        returnObj.resultMsg = "Success";
      } else {
        returnObj.resultCode = 102;
        returnObj.resultMsg = "Wrong password";
      }
    }
  } catch (err) {
    next(err);
  }
  await mysql.end();
  res.send(returnObj);
});

app.use((req, res, next) => {
  // let decoded;
  // let cert = "02000163";
  // let returnObj = {
  //   resultCode: null,
  //   resultMsg: null,
  // };
  // try {
  //   decoded = jwt.verify(req.headers.authorization, cert);
  // } catch (err) {
  //   next(err);
  // }
  // if (decoded.data) {
  //   res.locals.data = decoded.data;
  //   next();
  // } else {
  //   returnObj.resultCode = 11;
  //   returnObj.resultMsg = "Token authentication expired";
  //   res.send(returnObj);
  // }
  next();
});

app.post("/register", async (req, res, next) => {
  let returnObj = {
    resultCode: null,
    resultMsg: null,
    data: null,
  };
  let dbResults;
  try {
    dbResults = await mysql.query(
      "SELECT rpassword FROM repairelements WHERE rid = ?;",
      [res.locals.data.rid]
    );
    if (dbResults[0].rpassword == null || dbResults[0].rpassword == "") {
      await mysql.query(
        "UPDATE repairelements SET rpassword = ?,ralias = ?,name=?,class=?,gmt_modified = SYSDATE() WHERE rid = ?;",
        [
          req.body.password,
          req.body.alias,
          req.body.name,
          req.body.class,
          res.locals.data.rid,
        ]
      );
      returnObj.resultCode = 0;
      returnObj.resultMsg = "Success";
      returnObj.data = { alias: req.body.alias };
    } else {
      returnObj.resultCode = 10;
      returnObj.resultMsg = "Already registered";
    }
  } catch (err) {
    next(err);
  }
  await mysql.end();
  res.send(returnObj);
});

app.post("/events", async (req, res, next) => {
  let returnObj = {
    resultCode: null,
    resultMsg: null,
    data: null,
  };
  let dbResults;
  try {
    dbResults = await mysql.query(
      "SELECT eid,user_description,status,rid,gmt_create,gmt_modified FROM `event` ORDER BY gmt_modified DESC"
    );
    returnObj.resultCode = 0;
    returnObj.resultMsg = "Success";
    returnObj.data = dbResults;
  } catch (err) {
    next(err);
  }
  await mysql.end();
  res.send(returnObj);
});

//有问题
app.use("/event", async (req, res, next) => {
  let returnObj = {
    resultCode: null,
    resultMsg: null,
  };
  let dbResults;
  try {
    dbResults = await mysql.query("SELECT eid FROM `event` WHERE eid=?", [
      req.body.eid,
    ]);
  } catch (err) {
    next(err);
  }
  await mysql.end();
  if (dbResults.length != 0) {
    console.log("经过了event1");
    next();
  } else {
    console.log("经过了event2");
    returnObj.resultCode = 210;
    returnObj.resultMsg = "Eid does not exist";
    res.send(returnObj);
  }
});

app.use("/event/manage", async (req, res, next) => {
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

app.post("/event/getEventDetail", async (req, res, next) => {
  let returnObj = {
    resultCode: null,
    resultMsg: null,
    data: null,
  };
  let dbResults;
  try {
    dbResults = await mysql.query("SELECT * FROM `event` WHERE eid=?", [
      req.body.eid,
    ]);
    returnObj.resultCode = 0;
    returnObj.resultMsg = "Success";
    returnObj.data = dbResults[0];
    let temp = JSON.parse(returnObj.data.event_log);
    for (let i = 0; i < temp.length; i++) {
      if (temp[i].rid) {
        aliasResult = await mysql.query(
          "SELECT ralias FROM repairelements WHERE rid=?",
          [temp[i].rid]
        );
        temp[i].alias = aliasResult[0].ralias;
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

    returnObj.data.event_log = temp;
  } catch (err) {
    next(err);
  }
  await mysql.end();
  res.send(returnObj);
});

app.post("/event/handleEvent", async (req, res, next) => {
  let returnObj = {
    resultCode: null,
    resultMsg: null,
  };
  let dbResults;
  try {
    dbResults = await mysql.query(
      "SELECT rid,status,event_log  FROM `event` WHERE eid=?",
      [req.body.eid]
    );
    if (dbResults[0].rid == null && dbResults[0].status == 0) {
      let addeventLog = {
        type: "accept",
        time: new Date(),
        rid: res.locals.data.rid,
      };
      eventLog = jsonPush(dbResults[0].event_log, addeventLog);
      await mysql.query(
        "UPDATE `event` SET rid=?,event_log=?,status=1 WHERE eid=?",
        [res.locals.data.rid, eventLog, req.body.eid]
      );
      returnObj.resultCode = 0;
      returnObj.resultMsg = "Success";
    } else {
      returnObj.resultCode = 220;
      returnObj.resultMsg = "Event has been accepted or deleted";
    }
  } catch (err) {
    next(err);
  }
  await mysql.end();
  res.send(returnObj);
});

app.use("/event/edit", async (req, res, next) => {
  let returnObj = {
    resultCode: null,
    resultMsg: null,
  };
  let dbResults;
  try {
    dbResults = await mysql.query("SELECT rid FROM `event` WHERE eid=?", [
      req.body.eid,
    ]);
  } catch (err) {
    next(err);
  }
  await mysql.end();
  if (dbResults[0].rid == res.locals.data.rid) {
    console.log("经过了edit1");
    next();
  } else {
    console.log("经过了edit1");
    returnObj.resultCode = 230;
    returnObj.resultMsg = "No edit permission";
    res.send(returnObj);
  }
});

app.post("/event/edit/submitRepair", async (req, res, next) => {
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

app.post("/event/edit/cancelEvent", async (req, res, next) => {
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

app.post("/event/manage/checkEvent", async (req, res, next) => {
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

app.post("/manage/returnRepairelements", async (req, res, next) => {
  let returnObj = {
    resultCode: null,
    resultMsg: null,
    data: null,
  };
  let dbResults;
  try {
    if (req.body.random) {
      dbResults = await mysql.query(
        "SELECT rid,ralias,name,class,gmt_create,gmt_modified,rprofile,event_count FROM repairelements WHERE rid !=0000000000 ORDER BY RAND() LIMIT 1;"
      );
    } else {
      dbResults = await mysql.query(
        "SELECT rid,ralias,name,class,gmt_create,gmt_modified,rprofile,event_count FROM repairelements"
      );
    }
    returnObj.resultCode = 0;
    returnObj.resultMsg = "Success";
    returnObj.data = dbResults;
    console.log(dbResults);
  } catch (err) {
    next(err);
  }
  await mysql.end();
  res.send(returnObj);
});

app.post("/event/manage/assignEvent", async (req, res, next) => {
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

app.use("/", indexRouter);
app.use("/users", usersRouter);

// catch 404 and forward to error handler
app.use(function (req, res, next) {
  next(createError(404));
});

// error handler
app.use(function (err, req, res, next) {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get("env") === "development" ? err : {};

  // render the error page
  res.status(err.status || 500);
  res.render("error");
});

function jsonPush(str, data) {
  if (!str) {
    str = "[]";
  }
  let temp = JSON.parse(str);
  temp.push(data);
  return JSON.stringify(temp);
}

module.exports = app;
