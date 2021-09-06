const jwt = require("jsonwebtoken");
const { cert } = require("../config/config");
const { respond } = require("../utils/utils");
exports.auth = (req, res, next) => {
  let decoded;
  try {
    decoded = jwt.verify(req.headers.authorization.slice(7), cert);
  } catch (err) {
    next(err);
  }
  if (decoded) {
    // 
    res.locals.data = decoded.data;
    req.role = decoded.data.role;
    next();
  } else {
    // respond(res, 11, "Token authentication expired");
    respond(res, 12, "wrong token");
  }
};

exports.isAdmin = (req, res, next) => {
  console.log(req.role);
  req.role == "admin" ? next() : respond(res, 250, "No admin permission");
};
