const event = require("../models/event");

exports.isEidVaild = async (req, res, next) => {
  try {
    eid = req.params.eid || req.body.eid;
    req.event = await event.get(eid);
    if (req.event) {
      next();
    } else {
      res.status(210).send({ error: "Eid does not exist" });
      // respond(res, 210, "Eid does not exist");
    }
  } catch (error) {
    next(error);
  }
};
