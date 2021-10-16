const jwt = require("jsonwebtoken");
const log4js = require("../utils/log4js");
const { respond } = require("../utils/utils");
const { cert } = require("../config/config");
const ElementModel = require("../models/ElementModel");
class Login {
  constructor() {}
  async login(req, res, next) {
    let rid = req.body.id;
    let password = req.body.password;
    try {
      let dbResults = await ElementModel.findByFilter(["ralias", "rpassword", "ravatar", "role", "status"], { rid: rid });
      if (dbResults.length == 0) {
        respond(res, 1010, "No such user");
      } else {
        let elementInfo = dbResults[0];
        let roleMap = ["", "element", "admin"];
        if (password == elementInfo.rpassword || (password == "" && elementInfo.rpassword == null)) {
          let role = roleMap[elementInfo.role];
          if (elementInfo.status == 0) role = "notActivated";
          let info = {
            rid: rid,
            role: role,
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
            alias: elementInfo.ralias,
            avatar: elementInfo.ravatar,
            rid: rid,
            role: role,
          };
          await ElementModel.update(
            {
              gmt_modified: new Date(),
            },
            { rid: rid }
          );
          // TODO auto
          let logger = log4js.getLogger();
          logger.info(rid);
          respond(res, 0, "Success", data);
        } else {
          respond(res, 1011, "Wrong password");
        }
      }
    } catch (err) {
      next(err);
    }
  }
}
module.exports = new Login();
