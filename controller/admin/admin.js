const { respond } = require("../../utils/utils");
const AdminModel = require("../../models/AdminModel");
class Admin {
  constructor() {
    // this.Admin = AdminModel;
  }
  async getAll(req, res, next) {
    try {
      AdminModel.findAll().then((result) => respond(res, 0, "Success", result));
      // let data = await admin.get();
      // respond(res, 0, "Success", data);
    } catch (error) {
      next(error);
    }
  }
  async get(req, res, next) {
    try {
      await AdminModel.findByFilter({}, { aid: req.params.aid }).then(
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
      await AdminModel.create({
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
      await AdminModel.delete({
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
