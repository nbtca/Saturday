var express = require("express");
const router = express.Router();
const Admin = require("../controller/admin/admin");

router.get("/", Admin.getAll);
router.get("/:aid", Admin.get);
router.post("/", Admin.create);
router.delete("/", Admin.delete);
module.exports = router;
