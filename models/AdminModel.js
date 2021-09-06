const BaseModel = require("./BaseModel");
const Admin = require("../config/sequelize/admin");

class AdminModel extends BaseModel {
  constructor() {
    super(Admin());
  }
}
module.exports = new AdminModel();
