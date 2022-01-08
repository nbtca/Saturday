var express = require("express");
const router = express.Router();
const { auth } = require("../middleware/auth");

var loginRouter = require("./login");
var userRouter = require("./user");
var elementsRouter = require("./elements");
var eventsRouter = require("./events");
var settingRouter = require("./setting");


const Bot = require("../utils/bot");
router.get("/mirai", (request, response) => {
  Bot.msgTest();
  response.send("mirai");
})

router.use("/login", loginRouter);
router.use("/user", userRouter);
router.use("/elements", auth, elementsRouter);
router.use("/events", auth, eventsRouter);
router.use("/setting", settingRouter);

module.exports = router;
