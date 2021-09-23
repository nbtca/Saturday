const EventModel = require("../models/EventModel");

const { respond } = require("../utils/utils");
exports.isEidValid = (req, res, next) => {
  try {
    let eid = req.params.eid || req.body.eid;
    if (eid) {
      EventModel.findByFilter({}, { eid: eid }).then(result => {
        if (result.length != 0) {
          req.event = result[0].dataValues;
          if (req.event.rid == res.locals.data.rid) {
            req.role = "currentElement";
          }
          next();
        } else {
          respond(res, 2010, "Eid does not exist");
        }
      });
    }
  } catch (error) {
    console.log(error);
  }
};

exports.isCurrentElement = async (req, res, next) => {
  try {
    if (req.role == "currentElement") {
      next();
    } else {
      respond(res, 2020, "No edit permission");
    }
  } catch (error) {
    next(error);
  }
};
