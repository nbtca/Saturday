const BaseModel = require("./BaseModel");
const User = require("../config/user");

class UserModel extends BaseModel {
  constructor() {
    super(User());
  }
}
module.exports = new UserModel();
