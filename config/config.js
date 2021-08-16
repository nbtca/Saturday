exports.cert = "02000163";

var host = "rm-uf6s9l8ep4131lzt9go.mysql.rds.aliyuncs.com";
exports.mysql = require("serverless-mysql")({
  config: {
    host: host,
    user: "high_admin",
    password: "02000163",
    database: "repairteam_build",
  },
});
exports.actionSheet = {
  create: {
    type: "create",
    title: "提交",
    icon: "add_circle",
    auth: {
      role: ["user"],
    },
    targetStatus: 0,
    alterItem: [],
    logItem: [],
  },
  delete: {
    type: "delete",
    title: "取消",
    icon: "remove_circle",
    auth: {
      role: ["currentUser"],
      formerStatus: [0, 1],
    },
    targetStatus: -1,
    alterItem: [],
    logItem: [],
  },
  close: {
    type: "close",
    title: "完成",
    icon: "check_circle",
    auth: {
      role: ["admin"],
      formerStatus: [2],
    },
    targetStatus: 3,
    alterItem: ["aid"],
    logItem: ["aid"],
  },
  reject: {
    type: "reject",
    title: "完成",
    icon: "sentiment_very_dissatisfied",
    auth: {
      role: ["admin"],
      formerStatus: [2],
    },
    targetStatus: 1,
    alterItem: ["aid"],
    logItem: ["aid"],
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
    auth: {
      role: ["element"],
      formerStatus: [0],
    },
    targetStatus: 1,
    alterItem: ["rid"],
    logItem: ["rid"],
  },
  drop: {
    type: "drop",
    title: "放弃",
    icon: "sentiment_very_dissatisfied",
    auth: {
      role: ["currentElement"],
      formerStatus: [1],
    },
    targetStatus: 0,
    alterItem: [],
    logItem: ["rid"],
  },
  assign: {
    type: "assign",
    title: "指派",
    icon: "accept_circle",
    auth: {
      role: ["admin"],
      formerStatus: [0],
    },
    targetStatus: 1,
    alterItem: ["rid"],
    logItem: ["aid", "rid"],
  },
  submit: {
    type: "submit",
    title: "提交维修",
    //TODO change the icon
    icon: "sentiment_very_dissatisfied",
    auth: {
      role: ["currentElement"],
      formerStatus: [1],
    },
    targetStatus: 2,
    alterItem: ["repair_description"],
    logItem: ["description"],
  },
};
