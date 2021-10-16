const BaseModel = require("./BaseModel");
const User = require("../config/sequelize/user");

class UserModel extends BaseModel {
  constructor() {
    super(User());
  }
}
module.exports = new UserModel();
