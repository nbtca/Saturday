const event = require("../models/event");

exports.isEidVaild = async (req, res, next) => {
  try {
    eid = req.params.eid || req.body.eid;
    req.event = await event.get(eid);
    if (req.event) {
      next();
    } else {
      respond(res, 210, "Eid does not exist");
    }
  } catch (error) {
    next(error);
  }
};

exports.isCurrentUser = async (req, res, next) => {
  try {
    eid = req.params.eid || req.body.eid;
    rid = await event.get(eid).rid;
    if (rid == res.locals.data.rid) {
      next();
    } else {
      respond(res, 230, "No edit permission");
    }
  } catch (error) {
    next(error);
  }
};
