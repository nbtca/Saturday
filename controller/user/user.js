const UserModel = require("../../models/UserModel");
const { respond, uuid } = require("../../utils/utils");
class UserController {
  constructor() {}
  getUid(req, res, next) {
    try {
      UserModel.findByFilter( ["uid"] , { uopenid: req.body.openid }).then(
        result => {
          if (result.length == 0) {
            respond(res,122,"No such user");
          } else {
            respond(res, 0, "Success", result[0])
          }
        }
      );
    } catch (error) {
      console.error(error);
    }
  }
  get(req, res, next) {
    try {
      UserModel.findByFilter({}, { uid: req.params.uid }).then(result => {
        result ? respond(res, 0, "Success", result) : respond(res, 123, "user");
      });
    } catch (error) {
      next(error);
    }
  }
  create(req, res, next) {
    try {
      let id = uuid();
      UserModel.create({
        uid: id,
        uopenid: req.body.openid,
        gmt_create: new Date(),
        gmt_modified: new Date(),
      }).then(result => {
        respond(res, 0, "Success", { uid: id });
      });
    } catch (error) {
      console.log(error);
    }
  }
}
module.exports = new UserController();
