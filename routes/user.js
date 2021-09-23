var express = require("express");
const router = express.Router();
const { respond } = require("../utils/utils");
const User = require("../controller/user/user");

//TODO complele and test

router.get("/:uid", User.get);
router.post("/uid", User.getUid);
router.post("/",User.create)

module.exports = router;
