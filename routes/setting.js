const express = require("express");
const router = express.Router();
const { isAdmin } = require("../middleware/auth");
const Setting = require("../controller/setting/setting");

router.get("/", Setting.get);
router.put("/", isAdmin, Setting.update);

module.exports = router;
