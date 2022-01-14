const express = require("express");
const router = express.Router();
const { auth } = require("../middleware/auth");

const userRouter = require("./user");
const elementsRouter = require("./elements");
const eventsRouter = require("./events");
const settingRouter = require("./setting");


// const Bot = require("../utils/bot");
// router.get("/mirai", (request, response) => {
//   Bot.msgTest();
//   response.send("mirai");
// })



router.use("/user", userRouter);
router.use("/elements", elementsRouter);
router.use("/events", auth, eventsRouter);
router.use("/setting", settingRouter);

module.exports = router;
