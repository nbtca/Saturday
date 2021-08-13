const jwt = require("jsonwebtoken");
const { cert } = require("../config/config");
const { respond } = require("../utils");
exports.auth = (req, res, next) => {
  let decoded;
  try {
    decoded = jwt.verify(req.headers.authorization.slice(7), cert);
  } catch (err) {
    next(err);
  }
  if (decoded.data) {
    res.locals.data = decoded.data;
    res.role = decoded.data.role;
    next();
  } else {
    respond(res, 11, "Token authentication expired");
  }
};

exports.isAdmin = (req, res, next) => {
  req.role == "admin" ? next() : respond(res, 250, "No admin permission");
};
