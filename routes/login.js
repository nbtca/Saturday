var express = require("express");
const router = express.Router();
const jwt = require("jsonwebtoken");
const { respond } = require("../utils");
const { mysql, cert } = require("../config/config");
const element = require("../models/element");
const admin = require("../models/admin");

// TODO test
router.post("/", async (req, res, next) => {
  let rid = req.body.id;
  let password = req.body.password;
  try {
    dbResults = await element.get(rid);
    if (dbResults == null) {
      respond(res, 101, "no such user");
    } else {
      if (await element.checkPassword(rid, password)) {
        let isPasswordEmpty = !password ? true : false;
        let isAdmin = await admin.get(rid);
        let role = isAdmin ? "admin" : "element";
        let info = {
          rid: rid,
          role: role,
          aid: isAdmin.aid,
        };
        let token = jwt.sign(
          {
            exp: Math.floor(Date.now() / 1000) + 24 * 60 * 6000,
            data: info,
          },
          cert
        );
        let data = {
          token: token,
          alias: dbResults.ralias,
          rid: rid,
          role: role,
          isPasswordEmpty: isPasswordEmpty,
        };
        // TODO auto
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
module.exports = router;
