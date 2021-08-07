var express = require("express");
const router = express.Router();
const { respond } = require("../utils");
var host = "rm-uf6s9l8ep4131lzt9go.mysql.rds.aliyuncs.com";

const mysql = require("serverless-mysql")({
  config: {
    host: host,
    user: "high_admin",
    password: "***REMOVED***",
    database: "repairteam_build",
  },
});

//TODO user or user/elemet/admin ?
var userRouter = require("./user");

var elementsRouter = require("./elements");
var eventsRouter = require("./events");

router.post("/login", async (req, res, next) => {
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
        let cert = "***REMOVED***";
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

router.post("/register", async (req, res, next) => {
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

router.use("/events", eventsRouter);
router.use("/elements", elementsRouter);

module.exports = router;
