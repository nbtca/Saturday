const log4js = require("../../utils/log4js");
const { respond, dateToStr, put, createToken } = require("../../utils/utils");
const SettingModel = require("../../models/SettingModel");
class SettingController {
  constructor() {}
  async get(req, res, next) {
    try {
      let setting = await SettingModel.findAll();
      let settingInJson = JSON.parse(setting[0].setting);
      respond(res, 0, "Success", settingInJson);
    } catch (e) {
      console.log(e);
    }
  }
  async update(req, res, next) {
    try {
      let jsonString = JSON.stringify(req.body.setting);
      await SettingModel.update({
        setting: jsonString,
      });
      respond(res, 0, "Success", req.body.setting);
    } catch (e) {
      console.log(e);
    }
  }
}
module.exports = new SettingController();
