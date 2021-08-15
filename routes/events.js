const express = require("express");
const router = express.Router();
const { isAdmin } = require("../middleware/auth");
const { isEidVaild, isCurrentUser } = require("../middleware/event");

const Event = require("../controller/event");

router
  .get("/:eid", isEidVaild, Event.get)
  .use("/", isEidVaild)
  .get("/", Event.getAll)
  .post("/", Event.creat)
  .put("/", Event.update)
  .post("/accept", Event.accept)
  .post("/submit", isCurrentUser, Event.submit)
  .post("/cancel", isCurrentUser, Event.cancel)
  .post("/close", isAdmin, Event.close)
  .post("/assign", isAdmin, Event.assign);

module.exports = router;
