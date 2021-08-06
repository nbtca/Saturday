var express = require("express");
const router = express.Router();
const { respond } = require("../utils");
var host = "rm-uf6s9l8ep4131lzt9go.mysql.rds.aliyuncs.com";
const jwt = require("jsonwebtoken");

const mysql = require("serverless-mysql")({
  config: {
    host: host,
    user: "high_admin",
    password: "***REMOVED***",
    database: "repairteam_build",
  },
});

var elementsRouter = require("./routes/elements");
var eventsRouter = require("./routes/events");
var login = require("./routes/user");





module.exports = router;
