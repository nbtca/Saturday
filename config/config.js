exports.cert = "***REMOVED***";

const config = {
  host: "rm-uf6s9l8ep4131lzt9go.mysql.rds.aliyuncs.com",
  user: "high_admin",
  password: "***REMOVED***",
  database: "repairteam_build2",
};
exports.dbConfig = config;
exports.mysql = require("serverless-mysql")({ config: config });
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
    alterItem: ["closed_by"],
    logItem: ["closed_by"],
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
    alterItem: [],
    logItem: ["rejected_by"],
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
      formerStatus: [1, 2],
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
    logItem: ["assigned_by", "rid"],
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
  alterSubmit: {
    type: "alterSubmit",
    title: "修改维修描述",
    icon: "123",
    auth: {
      role: ["currentElement"],
      formerStatus: [2],
    },
    targetStatus: 2,
    alterItem: ["repair_description"],
    logItem: ["description"],
  },
};



exports.ossConfig = {
  region: "oss-cn-hangzhou",
  accessKeyId: "LTAI5tCyeZFdHskUvpTRCyPp",
  accessKeySecret: "r79n1DQaL5Y0lpremWGguBoHFA3aky",
  bucket: "sunday-res",
}