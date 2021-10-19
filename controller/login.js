const log4js = require("../utils/log4js");
const { respond, createToken } = require("../utils/utils");
const ElementModel = require("../models/ElementModel");
class Login {
  constructor() {}
  async login(req, res, next) {
    let rid = req.body.id;
    let password = req.body.password;
     ElementModel.findByFilter(
      ["ralias", "rpassword", "ravatar", "role", "status"],
      { rid: rid }
    )
      .then(async dbResults => {
        if (dbResults.length == 0) {
          respond(res, 1010, "No such user");
        } else {
          let elementInfo = dbResults[0];
          let roleMap = ["", "element", "admin"];
          if (
            password == elementInfo.rpassword ||
            (password == "" && elementInfo.rpassword == null)
          ) {
            let role =
              elementInfo.status == 0
                ? "notActivated"
                : roleMap[elementInfo.role];
            let token = createToken(100, {
              rid: rid,
              role: role,
            });
            let data = {
              token: token,
              alias: elementInfo.ralias,
              avatar: elementInfo.ravatar,
              rid: rid,
              role: role,
            };
            await ElementModel.update(
              { gmt_modified: new Date() },
              { rid: rid }
            );

            let logger = log4js.getLogger();
            logger.info(rid);
            respond(res, 0, "Success", data);
          } else {
            respond(res, 1011, "Wrong password");
          }
        }
      })
      .catch(error => {
        console.log(error);
      });
  }
}
module.exports = new Login();
