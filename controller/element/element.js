const { respond, dateToStr, put } = require("../../utils/utils");
const ElementModel = require("../../models/ElementModel");
class ElementController {
  constructor() {
    // ElementModel = new ElementModel();
  }
  getAll(req, res, next) {
    try {
      ElementModel.findByFilterOrder({ exclude: ["rpassword"] }, {}, [
        ["gmt_modified", "DESC"],
      ]).then(result => {
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
  get(req, res, next) {
    try {
      ElementModel.findByFilter(
        { exclude: ["rpassword"] },
        { rid: req.params.rid }
      ).then(result => {
        result
          ? respond(res, 0, "Success", result)
          : respond(res, 123, "no such element");
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
        created_by: res.locals.data.rid,
        gmt_create: new Date(),
        gmt_modified: new Date(),
        // TODO:gmt_expire
      });
      respond(res, 0);
    } catch (error) {
      next(error);
    }
  }
  activate(req, res, next) {
    try {
      let file = req.files.file;
      let ext = "." + file.type.substring(file.type.indexOf("/") + 1);
      let timestamps = new Date().getTime();
      let fileName = "/element/" + res.locals.data.rid + "/" + timestamps + ext;
      let path = req.files.file.path;
      console.log(fileName, path);
      put(fileName, path).then(result => {
        console.log(result);
        ElementModel.update(
          {
            rid: res.locals.data.rid,
            ralias: req.fields.alias,
            password: req.fields.password,
            name: req.fields.name,
            class: req.fields.class,
            rqq: req.fields.rqq,
            rphone: req.fields.rphone,
            ravatar: result.res.requestUrls[0],
          },
          { rid: res.locals.data.rid }
        ).then(respond(res, 0));
      });
    } catch (err) {
      console.log(err);
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
