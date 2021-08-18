const { respond } = require("../../utils/utils");
const admin = require("../../models/admin");
class Admin {
  constructor() {}
  async getAll(req, res, next) {
    try {
      let data = await admin.get();
      respond(res, 0, "Success", data);
    } catch (error) {
      next(error);
    }
  }
  async get(req, res, next) {
    try {
      let data = await admin.get(req.params.rid);
      //TODO error code
      data
        ? respond(res, 0, "Success", data)
        : respond(res, 123, "no such admin");
    } catch (error) {
      next(error);
    }
  }
  async create(req, res, next) {
    try {
      await admin.create(req.body.rid);
      respond(res, 0);
    } catch (error) {
      next(error);
    }
  }
  async delete(req, res, next) {
    try {
      await element.delete({
        rid: req.body.rid,
      });
      respond(res, 0);
    } catch (err) {
      next(err);
    }
  }
}
module.exports = new Admin();
