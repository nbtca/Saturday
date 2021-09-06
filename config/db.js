const Sequelize = require("sequelize");

const { dbConfig } = require("./config");
// sequelize-automate -t js -h rm-uf6s9l8ep4131lzt9go.mysql.rds.aliyuncs.com -d repairteam_build -u high_admin -p ***REMOVED*** -P 3306  -e mysql -o models

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
