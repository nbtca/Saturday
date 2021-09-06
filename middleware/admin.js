const { respond } = require("../utils/utils");
exports.isAidValid = async (req, res, next) => {
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
