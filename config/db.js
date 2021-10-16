const Sequelize = require("sequelize");

const { dbConfig } = require("./config");

const db = new Sequelize(dbConfig.database, dbConfig.user, dbConfig.password, {
  host: dbConfig.host,
  dialect: "mysql",
});

db.authenticate()
  .then(() => {
    console.log("数据库连接成功...");
  })
  .catch((err) => {
    console.error("数据库连接失败...", err);
  });
module.exports = db;
