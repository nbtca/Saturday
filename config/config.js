exports.cert = "***REMOVED***";

var host = "rm-uf6s9l8ep4131lzt9go.mysql.rds.aliyuncs.com";
exports.mysql = require("serverless-mysql")({
  config: {
    host: host,
    user: "high_admin",
    password: "***REMOVED***",
    database: "repairteam_build",
  },
});
exports.actionSheet = {
  create: {
    type: "create",
    title: "提交",
    icon: "add_circle",
    role: ["user"],
  },
  delete: {
    type: "delete",
    title: "取消",
    icon: "remove_circle",
    role: ["user"],
  },
  close: {
    type: "close",
    title: "完成",
    icon: "check_circle",
    role: ["admin"],
  },
  update: {
    type: "update",
    title: "更新",
    icon: "update_circle",
    role: ["user"],
  },
  accept: {
    type: "accept",
    title: "接受",
    icon: "accept_circle",
    role: ["element"],
  },
  cancel: {
    type: "cancel",
    title: "放弃",
    icon: "sentiment_very_dissatisfied",
    role: ["currentElement"],
  },
  assign: {
    type: "assign",
    title: "指派",
    icon: "accept_circle",
    role: ["admin"],
  },
  submit: {
    type: "submit",
    title: "提交维修",
    //TODO change the icon
    icon: "sentiment_very_dissatisfied",
    role: ["currentElement"],
  },
};
