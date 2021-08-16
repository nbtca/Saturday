const event = require("../models/event");
const { respond } = require("../utils/utils");
exports.isEidVaild = async (req, res, next) => {
  try {
    let eid = req.params.eid || req.body.eid;
    if (eid) {
      req.event = await event.get(eid);
      if (req.event != null) {
        next();
      } else {
        respond(res, 210, "Eid does not exist");
      }
    }
  } catch (error) {
    next(error);
  }
};

exports.isCurrentElement = async (req, res, next) => {
  try {
    eid = req.params.eid || req.body.eid;
    rid = req.event.rid;
    if (rid == res.locals.data.rid) {
      next();
    } else {
      respond(res, 230, "No edit permission");
    }
  } catch (error) {
    next(error);
  }
};
