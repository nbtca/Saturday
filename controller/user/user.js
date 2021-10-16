const UserModel = require("../../models/UserModel");
const { respond, uuid, createToken } = require("../../utils/utils");
const axios = require("axios");
const log4js = require("../../utils/log4js");
class UserController {
  constructor() {}
  getUid(req, res, next) {
    try {
      UserModel.findByFilter(["uid"], { uopenid: req.body.openid }).then(
        result => {
          if (result.length == 0) {
            respond(res, 122, "No such user");
          } else {
            respond(res, 0, "Success", result[0]);
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
  async wxLogin(req, res, next) {
    const wxAuthUrl = "https://api.weixin.qq.com/sns/jscode2session";
    axios
      .get(wxAuthUrl, {
        params: {
          appid: "***REMOVED***",
          secret: "***REMOVED***",
          js_code: req.body.code,
          grant_type: "authorization_code",
        },
      })
      .then(wxData => {
        if (wxData.openid == null) {
          respond(res, 555, "wrong code");
          throw new Error("wrong code");
        }
        let openid = wxData.data.openid;
        return UserModel.findByFilter(["uid"], { uopenid: openid });
      })
      .then(async result => {
        let uid;
        if (result.length == 0) {
          uid = uuid();
          await UserModel.create({
            uid: uid,
            uopenid: openid,
            gmt_create: new Date(),
            gmt_modified: new Date(),
          });
        } else {
          uid = result[0].uid;
        }
        //生成token
        let token = createToken(1, { uid: uid });
        let data = {
          uid: uid,
          token: token,
        };
        respond(res, 0, "Success", data);
      })
      .catch(error => {
        console.log(error);
      });
  }
}
module.exports = new UserController();
