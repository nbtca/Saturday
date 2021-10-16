var express = require("express");
const router = express.Router();
const Login = require("../controller/login");

router.post("/", Login.login);
module.exports = router;
