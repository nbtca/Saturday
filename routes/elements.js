var express = require("express");
const router = express.Router();
const { respond } = require("../utils");
const element = require("../models/element");
/* GET home page. */

router.get("/", async (req, res, next) => {
  try {
    let data = await element.get();
    respond(res, 0, "Success", data);
  } catch (error) {
    next(error);
  }
});

router.get("/:rid", async (req, res, next) => {
  try {
    let data = await element.get(req.params.rid);
    respond(res, 0, "Success", data);
  } catch (error) {
    next(error);
  }
});
module.exports = router;
