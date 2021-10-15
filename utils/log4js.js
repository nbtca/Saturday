const log4js = require("log4js");

let date = new Date();
let logName = date.getFullYear() + "-" + date.getMonth() + "-" + date.getDate() + ".log";
log4js.configure({
  appenders: {
    fileout: { type: "file", filename: "./logs/" + logName },
    consoleout: { type: "console" },
  },
  categories: {
    default: { appenders: ["fileout", "consoleout"], level: "debug" },
    anything: { appenders: ["consoleout"], level: "debug" },
  },
});

module.exports = log4js;
