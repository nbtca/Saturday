const BaseModel = require("./BaseModel");
const Setting = require("../config/sequelize/setting");

class SettingModel extends BaseModel {
  constructor() {
    super(Setting());
  }
}
module.exports = new SettingModel();
