const log4js = require("../../utils/log4js");
const { respond, dateToStr, put, createToken } = require("../../utils/utils");
const ElementModel = require("../../models/ElementModel");

class ElementController {
  constructor() {}

  getAll(req, res, next) {
    ElementModel.findByFilterOrder({ exclude: ["rpassword"] }, {}, [["gmt_modified", "DESC"]])
      .then(result => {
        for (let item of result) {
          //TODO gmt_create format
          item.dataValues.gmt_create = dateToStr(item.gmt_create, "time");
          item.dataValues.gmt_modified = dateToStr(item.gmt_modified, "time");
        }
        respond(res, 0, "Success", result);
      })
      .catch(error => console.log(error));
  }

  get(req, res, next) {
    ElementModel.findByFilter({ exclude: ["rpassword"] }, { rid: req.params.rid })
      .then(result => {
        result ? respond(res, 0, "Success", result) : respond(res, 123, "no such element");
      })
      .catch(error => console.log(error));
  }

  create(req, res, next) {
    ElementModel.create({
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
    })
      .then(() => respond(res, 0))
      .catch(error => console.log(error));
  }

  update(req, res, next) {
    ElementModel.update(
      {
        rid: res.locals.data.rid,
        ralias: req.body.alias,
        rpassword: req.body.password,
        name: req.body.name,
        class: req.body.class,
        rqq: req.body.rqq,
        rphone: req.body.rphone,
      },
      { rid: res.locals.data.rid }
    )
      .then(respond(res, 0))
      .catch(error => console.log(error));
  }
  delete(req, res, next) {
    ElementModel.delete({
      rid: req.body.rid,
    })
      .then(result => {
        // console.log(result);
        respond(res, 0);
      })
      .catch(error => console.log(error));
  }
  login(req, res, next) {
    let rid = req.body.id;
    let password = req.body.password;
    ElementModel.findByFilter(["ralias", "rpassword", "ravatar", "role", "status"], { rid: rid })
      .then(async dbResults => {
        if (dbResults.length == 0) {
          respond(res, 1010, "No such user");
        } else {
          let elementInfo = dbResults[0];
          let roleMap = ["", "element", "admin"];
          if (password == elementInfo.rpassword || (password == "" && elementInfo.rpassword == null)) {
            let role = elementInfo.status == 0 ? "notActivated" : roleMap[elementInfo.role];
            let token = createToken(100, {
              rid: rid,
              role: role,
            });
            let data = {
              token: token,
              alias: elementInfo.ralias,
              avatar: elementInfo.ravatar,
              rid: rid,
              role: role,
            };
            await ElementModel.update({ gmt_modified: new Date() }, { rid: rid });

            let logger = log4js.getLogger();
            logger.info(rid);

            respond(res, 0, "Success", data);
          } else {
            respond(res, 1011, "Wrong password");
          }
        }
      })
      .catch(error => {
        console.log(error);
      });
  }
  activate(req, res, next) {
    ElementModel.update(
      {
        ralias: req.body.alias,
        rpassword: req.body.password,
        rqq: req.body.qq,
        status: 1,
        rphone: req.body.phone,
      },
      { rid: res.locals.data.rid }
    )
      .then(respond(res, 0))
      .catch(error => console.log(error));
  }
  // activate(req, res, next) {
  //   try {
  //     let file = req.files.file;
  //     let ext = "." + file.type.substring(file.type.indexOf("/") + 1);
  //     let timestamps = new Date().getTime();
  //     let fileName = "/element/" + res.locals.data.rid + "/" + timestamps + ext;
  //     let path = req.files.file.path;
  //     console.log(fileName, path);
  //     put(fileName, path).then(result => {
  //       console.log(result);
  //       ElementModel.update(
  //         {
  //           rid: res.locals.data.rid,
  //           ralias: req.fields.alias,
  //           rpassword: req.fields.password,
  //           name: req.fields.name,
  //           class: req.fields.class,
  //           rqq: req.fields.rqq,
  //           rphone: req.fields.rphone,
  //           ravatar: result.res.requestUrls[0],
  //         },
  //         { rid: res.locals.data.rid }
  //       ).then(respond(res, 0));
  //     });
  //   } catch (err) {
  //     console.log(err);
  //   }
  // }
  updateAvatar(req, res, next) {
    let file = req.files.file;
    let ext = "." + file.type.substring(file.type.indexOf("/") + 1);
    let timestamps = new Date().getTime();
    let fileName = "/element/" + res.locals.data.rid + "/" + timestamps + ext;
    let path = req.files.file.path;
    console.log(fileName, path);
    put(fileName, path)
      .then(result => {
        console.log(result);
        ElementModel.update(
          {
            ravatar: result.res.requestUrls[0],
          },
          { rid: res.locals.data.rid }
        ).then(respond(res, 0));
      })
      .catch(error => console.log(error));
  }
}
module.exports = new ElementController();
