const express = require("express");
const router = express.Router();
const { isAdmin } = require("../middleware/auth");
const { isEidVaild, isCurrentElement } = require("../middleware/event");

const Event = require("../controller/event/event");

router.post("/", Event.create);
router.get("/", Event.getAll);
router.get("/:eid", isEidVaild, Event.get);
router.use("/", isEidVaild);
router.put("/", Event.update);
router.post("/delete", Event.delete);
router.post("/accept", Event.accept);
router.post("/submit", isCurrentElement, Event.submit);
router.post("/drop", isCurrentElement, Event.drop);
router.post("/assign", isAdmin, Event.assign);
router.post("/close", isAdmin, Event.close);

module.exports = router;
