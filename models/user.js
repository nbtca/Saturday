const { mysql } = require("../config/config");

exports.get = async (rid) => {
  let dbResults;
  try {
    dbResults = await mysql.query("SELECT * FROM repairelements WHERE rid=?", [
      rid,
    ]);
  } catch (err) {
    return err;
  }
  await mysql.end();
  return dbResults[0];
};

exports.update = async (info) => {
  await mysql.query(
    "UPDATE repairelements SET rpassword = ?,ralias = ?,name=?,class=?,gmt_modified = SYSDATE() WHERE rid = ?;",
    [info.password, info.alias, info.name, info.class, info.rid]
  );
};

exports.isAdmin = async (rid) => {
  let dbResults;
  try {
    dbResults = await mysql.query("SELECT aid FROM admin WHERE rid=?", [rid]);
  } catch (error) {
    return err;
  }
  await mysql.end();
  return dbResults[0];
};
