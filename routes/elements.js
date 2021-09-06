var express = require("express");
const router = express.Router();
const { test } = require("../test");
var Element = require("../controller/element/element");
router.get("/", (...e) => Element.getAll(...e));
// router.get("/test", (e) => test(e));
router.get("/:rid", (...e) => Element.get(...e));
router.post("/", (...e) => Element.create(...e));
router.put("/", (...e) => Element.update(...e));
router.delete("/", (...e) => Element.delete(...e));

module.exports = router;
