var express = require("express");
var router = express.Router();
const { jsonPush, respond } = require("../utils");
const { cert } = require("../config/config");
const user = require("../models/user");

router.post("/login", async (req, res, next) => {
  let rid = req.body.id;
  let password = req.body.password;
  try {
    dbResults = user.get(rid);
    if (dbResults == null) {
      respond(res, 101, "no such user");
    } else {
      if (dbResults.rpassword == password) {
        let isPasswordEmpty = !password ? true : false;
        let isAdmin = user.isAdmin(rid);
        let role = isAdmin ? "admin" : "element";
        let data = {
          rid: rid,
          role: role,
          aid: isAdmin.aid,
        };
        let token = jwt.sign(
          {
            exp: Math.floor(Date.now() / 1000) + 24 * 60 * 60,
            data: data,
          },
          cert
        );
        let data = {
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
        respond(res, 0, "Success", data);
      } else {
        respond(res, 102, "Wrong password");
      }
    }
  } catch (err) {
    next(err);
  }
  await mysql.end();
});

//TODO complele and test
router.put("/", async (req, res, next) => {
  try {
    //   await mysql.query(
    //     "UPDATE repairelements SET rpassword = ?,ralias = ?,name=?,class=?,gmt_modified = SYSDATE() WHERE rid = ?;",
    //     [
    //       req.body.password,
    //       req.body.alias,
    //       req.body.name,
    //       req.body.class,
    //       res.locals.data.rid,
    //     ]
    //   );
    user.update({
      password: req.body.password,
      alias: req.body.alias,
      name: req.body.name,
      class: req.body.class,
      rid: res.locals.data.rid,
    });
    respond(0);
  } catch (err) {
    next(err);
  }
  res.send(returnObj);
});

module.exports = router;
