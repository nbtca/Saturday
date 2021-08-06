const { NotExtended } = require("http-errors");
const jwt = require("jsonwebtoken");
const { cert } = require("../config/config");
const { respond } = require("../utils");
exports.auth = (req, res, NotExtended) => {
  let decoded;
  let returnObj = {
    resultCode: null,
    resultMsg: null,
  };
  try {
    decoded = jwt.verify(req.headers.authorization, cert);
  } catch (err) {
    next(err);
  }
  if (decoded.data) {
    res.locals.data = decoded.data;
    next();
  } else {
    respond(res, 11, "Token authentication expired");
  }
};
