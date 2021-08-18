var express = require("express");
const router = express.Router();
const { auth } = require("../middleware/auth");

var loginRouter = require("./login");
var userRouter = require("./user");
var elementsRouter = require("./elements");
var adminRouter = require("./admin");
var eventsRouter = require("./events");

router.use("/login", loginRouter);
router.use("/user", userRouter);
router.use("/elements", auth, elementsRouter);
router.use("/admin", auth, adminRouter);
router.use("/events", auth, eventsRouter);

module.exports = router;
