var express = require("express");
const router = express.Router();
const User = require("../controller/user/user");


router.get("/:uid", User.get);
router.post("/uid", User.getUid);
router.post("/",User.create)

module.exports = router;
