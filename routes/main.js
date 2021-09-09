var express = require("express");
const router = express.Router();
const { auth } = require("../middleware/auth");
const formidableMiddleware = require("express-formidable");

const { putBuffer } = require("../test");

var loginRouter = require("./login");
var userRouter = require("./user");
var elementsRouter = require("./elements");
var adminRouter = require("./admin");
var eventsRouter = require("./events");

router.post("/test", formidableMiddleware(), async (req, res, next) => {
  console.log("pass test");
  console.log(req.headers);
  console.log(req.body);
  console.log(req.files);
  let buffer = new Buffer(req.files.file);

  await putBuffer(buffer);

  res.send(req.body);
});
router.use("/login", loginRouter);
router.use("/user", userRouter);
router.use("/elements", auth, elementsRouter);
router.use("/admin", auth, adminRouter);
router.use("/events", auth, eventsRouter);

module.exports = router;
