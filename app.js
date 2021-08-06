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

var elementsRouter = require("./routes/elements");
var eventsRouter = require("./routes/events");
var login = require("./routes/user");

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

app.use("/events", eventsRouter);
app.use("/elements", elementsRouter);

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

module.exports = app;
