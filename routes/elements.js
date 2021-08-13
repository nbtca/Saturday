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
    //TODO error code
    data
      ? respond(res, 0, "Success", data)
      : respond(res, 123, "no such element");
  } catch (error) {
    next(error);
  }
});

router.put("/", async (req, res, next) => {
  try {
    element.update({
      password: req.body.password,
      alias: req.body.alias,
      name: req.body.name,
      class: req.body.class,
      rid: res.locals.data.rid,
    });
    respond(0);
  } catch (err) {
    next(err);
  }
});

module.exports = router;
