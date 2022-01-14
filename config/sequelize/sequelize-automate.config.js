const { dbConfig } = require("../config");

module.exports = {
  dbOptions: {
    database: dbConfig.database,
    username: dbConfig.user,
    password: dbConfig.password,
    dialect: "mysql",
    host: dbConfig.host,
    port: 3306,
    logging: false,
  },
  options: {
    type: "js",
    dir: "./config/sequelize",
  }
};
