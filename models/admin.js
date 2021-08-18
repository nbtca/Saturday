const { mysql } = require("../config/config");

exports.get = async (rid) => {
  let dbResults;
  if (rid) {
    dbResults = await mysql.query("SELECT * FROM admin WHERE rid=?", [rid]);
    dbResults = dbResults[0];
  } else {
    dbResults = await mysql.query("SELECT * FROM admin");
  }

  await mysql.end();
  return dbResults;
};
exports.create = async (rid) => {
  await mysql
    .query("SELECT * FROM admin WHERE rid=?", [rid])
    .then(async (res) => {
      if (res == null) {
        await mysql.query(
          "insert into admin (rid,gmt_create,gmt_modified) values (?,?,?)",
          [rid, SYSDATE(), SYSDATE()]
        );
      } else {
        throw new Error("rid already exists");
      }
    });
};
exports.delete = async (rid) => {
  try {
    await mysql.query("DELETE FROM admin where rid = ?", [rid]);
  } catch (error) {
    return error;
  }
};
