const { respond } = require("../utils/utils");
exports.isAidValid = async (req, res, next) => {
  req.role == "admin" ? next() : respond(res, 250, "No admin permission");
};
