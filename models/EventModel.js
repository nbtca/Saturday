const BaseModel = require("./BaseModel");
const Event = require("../config/sequelize/event");

class EventModel extends BaseModel {
  constructor() {
    super(Event());
  }
}
module.exports = new EventModel();
