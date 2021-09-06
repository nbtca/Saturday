const BaseModel = require("./BaseModel");
const RepairElementsModel = require("../config/sequelize/repairelements");

class ElementModel extends BaseModel {
  constructor() {
    super(RepairElementsModel());
  }
}
module.exports = new ElementModel();
