const { respond } = require("../../utils/utils");
const ElementModel = require("../../models/ElementModel");
class ElementController {
  constructor() {
    this.Element = new ElementModel();
  }
  async getAll(req, res, next) {
    try {
      await this.Element.findAll().then((result) => {
        respond(res, 0, "Success", result);
      });
    } catch (error) {
      console.error(error);
    }
  }
  async get(req, res, next) {
    try {
      await this.Element.findByFilter({}, { rid: req.params.rid }).then(
        (result) => {
          result
            ? respond(res, 0, "Success", result)
            : respond(res, 123, "no such element");
        }
      );
    } catch (error) {
      next(error);
    }
  }
  async create(req, res, next) {
    try {
      await this.Element.create({
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
  async test(req, res, next) {
    console.log(req);
  }
  async update(req, res, next) {
    try {
      await this.Element.update(
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
      await this.Element.delete({
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
module.exports = new ElementController();
