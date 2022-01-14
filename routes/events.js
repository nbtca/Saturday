const express = require("express");
const router = express.Router();
const { isAdmin, auth } = require("../middleware/auth");
const { isEidValid: isEidValid, isCurrentElement } = require("../middleware/event");

const Event = require("../controller/event/event");

router.post("/", Event.create);
router.get("/", Event.getAll);
router.get("/:eid", isEidValid, Event.get);
router.use("/", isEidValid);
router.put("/", Event.update);
router.post("/delete", Event.delete);
router.post("/accept", Event.accept);
router.post("/submit", isCurrentElement, Event.submit);
router.post("/alterSubmit", isCurrentElement, Event.alterSubmit);
router.post("/drop", isCurrentElement, Event.drop);
router.use(isAdmin);
router.post("/close", Event.close);
router.post("/reject", Event.reject);
router.post("/assign", Event.assign);

module.exports = router;
