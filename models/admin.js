const { mysql } = require("../config/config");

exports.get = async (rid) => {
  let dbResults;
  try {
    dbResults = await mysql.query("SELECT * FROM admin WHERE rid=?", [rid]);
  } catch (err) {
    return err;
  }
  await mysql.end();
  return dbResults[0];
};
