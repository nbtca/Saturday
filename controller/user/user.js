const UserModel = require("../../models/UserModel");
const { respond, uuid, createToken } = require("../../utils/utils");
const { weChatConfig } = require("../../config")
const axios = require("axios");
class UserController {
  constructor() { }
  getUid(req, res, next) {
    UserModel.findByFilter(["uid"], { uopenid: req.body.openid })
      .then(result => {
        if (result.length == 0) {
          respond(res, 122, "No such user");
        } else {
          respond(res, 0, "Success", result[0]);
        }
      })
      .catch(error => {
        console.log(error);
      });
  }
  get(req, res, next) {
    UserModel.findByFilter({}, { uid: req.params.uid })
      .then(result => {
        result ? respond(res, 0, "Success", result) : respond(res, 123, "user");
      })
      .catch(error => {
        console.log(error);
      });
  }
  create(req, res, next) {
    let id = uuid();
    UserModel.create({
      uid: id,
      uopenid: req.body.openid,
      gmt_create: new Date(),
      gmt_modified: new Date(),
    })
      .then(result => {
        respond(res, 0, "Success", { uid: id });
      })
      .catch(error => {
        console.log(error);
      });
  }
  wxLogin(req, res, next) {
    const wxAuthUrl = "https://api.weixin.qq.com/sns/jscode2session";
    let openid;
    axios
      .get(wxAuthUrl, {
        params: {
          appid: weChatConfig.appid,
          secret: weChatConfig.secret,
          js_code: req.body.code,
          grant_type: "authorization_code",
        },
      })
      .then(wxData => {
        if (wxData.data.openid == null) {
          respond(res, 555, "wrong code");
          throw new Error("wrong code");
        }
        openid = wxData.data.openid;
        console.log(openid);
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
        let token = createToken(1, { uid: uid, role: "user" });
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
