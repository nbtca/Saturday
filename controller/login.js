const jwt = require("jsonwebtoken");
const { respond } = require("../utils/utils");
const { mysql, cert } = require("../config/config");
const ElementModel = require("../models/ElementModel");
class Login {
  constructor() {}
  async login(req, res, next) {
    let rid = req.body.id;
    let password = req.body.password;
    try {
      let dbResults = await ElementModel.findByFilter(["ralias", "ravatar", "role", "rpassword"], { rid: rid });
      console.log(dbResults[0]);
      if (dbResults == null) {
        respond(res, 101, "no such user");
      } else {
        if (1) {
          console.log(dbResults[0].rpassword);
          if (password == dbResults[0].rpassword || (password == "" && dbResults[0].rpassword == null)) {
            let isActivated = password != "" ? true : false;
            let role = dbResults[0].role == 2 ? "admin" : "element";
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
              alias: dbResults[0].ralias,
              avatar: dbResults[0].ravatar,
              rid: rid,
              role: role,
              isActivated: isActivated,
            };
            await ElementModel.update(
              {
                gmt_modified: new Date(),
              },
              { rid: rid }
            );
            // TODO auto
            respond(res, 0, "Success", data);
          } else {
            respond(res, 102, "Wrong password");
          }
        }
      }
    } catch (err) {
      next(err);
    }
    await mysql.end();
  }
}
module.exports = new Login();
