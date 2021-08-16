var express = require("express");
const router = express.Router();
const Login = require("../controller/login");
// TODO test
router.post("/", Login.login);
module.exports = router;
