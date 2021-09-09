const { respond, dateToStr } = require("../../utils/utils");
const ElementModel = require("../../models/ElementModel");
class ElementController {
  constructor() {
    // ElementModel = new ElementModel();
  }
  async getAll(req, res, next) {
    try {
      await ElementModel.findAll({ exclude: ["rpassword"] }).then(result => {
        for (let item of result) {
          //TODO gmt_create format
          item.dataValues.gmt_create = dateToStr(item.gmt_create, "time");
          item.dataValues.gmt_modified = dateToStr(item.gmt_modified, "time");
        }
        respond(res, 0, "Success", result);
      });
    } catch (error) {
      console.error(error);
    }
  }
  async get(req, res, next) {
    try {
      await ElementModel.findByFilter({ exclude: ["rpassword"] }, { rid: req.params.rid }).then(result => {
        result ? respond(res, 0, "Success", result) : respond(res, 123, "no such element");
      });
    } catch (error) {
      next(error);
    }
  }
  async create(req, res, next) {
    try {
      await ElementModel.create({
        rid: req.body.rid,
        alias: req.body.alias,
        password: req.body.password,
        name: req.body.name,
        class: req.body.class,
        rqq: req.body.rqq,
        rphone: req.body.rphone,
        ravatar: req.body.ravatar,
        profile: req.body.profile,
        gmt_create: new Date(),
        gmt_modified: new Date(),
      });
      respond(res, 0);
    } catch (error) {
      next(error);
    }
  }
  async update(req, res, next) {
    try {
      await ElementModel.update(
        {
          rid: res.locals.data.rid,
          ralias: req.body.alias,
          password: req.body.password,
          name: req.body.name,
          class: req.body.class,
          rqq: req.body.rqq,
          rphone: req.body.rphone,
          ravatar: req.body.ravatar,
        },
        { rid: res.locals.data.rid }
      ).then(respond(res, 0));
    } catch (err) {
      next(err);
    }
  }
  async delete(req, res, next) {
    try {
      await ElementModel.delete({
        rid: req.body.rid,
      }).then(result => {
        // console.log(result);
        respond(res, 0);
      });
    } catch (err) {
      next(err);
    }
  }
}
module.exports = new ElementController();
