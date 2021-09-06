const { respond } = require("../../utils/utils");
const admin = require("../../models/admin");
const AdminModel = require("../../models/AdminModel");
class Admin {
  constructor() {
    this.Admin = AdminModel;
  }
  async getAll(req, res, next) {
    try {
      this.Admin.findAll().then((result) => respond(res, 0, "Success", result));
      // let data = await admin.get();
      // respond(res, 0, "Success", data);
    } catch (error) {
      next(error);
    }
  }
  async get(req, res, next) {
    try {
      await this.Admin.findByFilter({}, { rid: req.params.rid }).then(
        (result) => {
          result
            ? respond(res, 0, "Success", result)
            : respond(res, 123, "no such admin");
        }
      );
    } catch (error) {
      next(error);
    }
  }
  async create(req, res, next) {
    try {
      await admin.create({
        rid: req.body.rid,
        gmt_create: new Date(),
        gmt_modified: new Date(),
      });
      respond(res, 0);
    } catch (error) {
      next(error);
    }
  }

  async delete(req, res, next) {
    try {
      await this.Admin.delete({
        rid: req.body.rid,
      }).then((result) => {
        // console.log(result);
        respond(res, 0);
      });
    } catch (err) {
      next(err);
    }
  }
}
module.exports = new Admin();
