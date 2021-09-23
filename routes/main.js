var express = require("express");
const router = express.Router();
const { auth } = require("../middleware/auth");
const formidableMiddleware = require("express-formidable");

const { put } = require("../test");

var loginRouter = require("./login");
var userRouter = require("./user");
var elementsRouter = require("./elements");
var eventsRouter = require("./events");

router.post("/test", formidableMiddleware(), async (req, res, next) => {
  console.log("pass test");
  console.log(req.headers);
  console.log(req.fields);
  console.log(req.files.file.path);
  let file = req.files.file;
  let ext = "." + file.type.substring(file.type.indexOf("/") + 1);
  let fileName = "test" + ext;
  let path = req.files.file.path;
  console.log(fileName, path);

  await put("/element/rid/avatar.png", path);
  // res.send(req.body);
  res.send();
});
// const upload = multer({ dest: 'C:/Users/Administrator/Desktop/image' })
// router.post('/test', upload.single('file'), function(req, res, next){
// 	var file = req.file;
// 	res.json(200, {
//   		msg : 'success',
//   		imgs_url : 'http://example.com/image/' + file.filename //返回图片URL
// 	});
// });
router.use("/login", loginRouter);
router.use("/user", userRouter);
router.use("/elements", auth, elementsRouter);
router.use("/events", auth, eventsRouter);

module.exports = router;
