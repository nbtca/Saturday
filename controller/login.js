const jwt = require("jsonwebtoken");
const { respond } = require("../utils/utils");
const { mysql, cert } = require("../config/config");
const ElementModel = require("../models/ElementModel");
const admin = require("../models/admin");
class Login {
  constructor() {}
  async login(req, res, next) {
    let rid = req.body.id;
    let password = req.body.password;
    try {
      let Element = new ElementModel();
      let dbResults = await Element.findByFilter({}, { rid: rid });
      if (dbResults == null) {
        respond(res, 101, "no such user");
      } else {
        if (1) {
          // if (await element.checkPassword(rid, password)) {
          let isActivated = !password ? true : false;
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
            avatar: dbResults.ravatar,
            rid: rid,
            role: role,
            isActivated: isActivated,
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
  }
}
module.exports = new Login();
