const express = require("express");
const router = express.Router();
const User = require("../controller/user/user");

//TODO complele and test
router.post("/wxlogin", User.wxLogin);
router.get("/:uid", User.get);
router.post("/uid", User.getUid);
router.post("/",User.create)

module.exports = router;
