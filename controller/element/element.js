const { respond } = require("../../utils/utils");
const element = require("../../models/element");
class Element {
  constructor() {}
  async getAll(req, res, next) {
    try {
      let data = await element.get();
      respond(res, 0, "Success", data);
    } catch (error) {
      next(error);
    }
  }
  async get(req, res, next) {
    try {
      let data = await element.get(req.params.rid);
      //TODO error code
      data
        ? respond(res, 0, "Success", data)
        : respond(res, 123, "no such element");
    } catch (error) {
      next(error);
    }
  }
  async create(req, res, next) {
    console.log("pass create");
    try {
      await element.create({
        rid: req.body.rid,
        alias: req.body.alias,
        password: req.body.password,
        name: req.body.name,
        class: req.body.class,
        rqq: req.body.rqq,
        rphone: req.body.rphone,
        ravatar: req.body.ravatar,
        profile: req.body.profile,
      });
      respond(res, 0);
    } catch (error) {
      next(error);
    }
  }
  async update(req, res, next) {
    try {
      await element.update({
        rid: res.locals.data.rid,
        alias: req.body.alias,
        password: req.body.password,
        name: req.body.name,
        class: req.body.class,
        rqq: req.body.rqq,
        rphone: req.body.rphone,
        ravatar: req.body.ravatar,
      });
      respond(res, 0);
    } catch (err) {
      next(err);
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
module.exports = new Element();
