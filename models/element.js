const { mysql } = require("../config/config");

exports.get = async (rid) => {
  let dbResults;
  try {
    if (rid == -1) {
      // roll
      dbResults = await mysql.query(
        "SELECT rid,ralias,name,class,gmt_create,gmt_modified,rprofile,event_count FROM repairelements WHERE rid !=0000000000 ORDER BY RAND() LIMIT 1;"
      );
      dbResults = dbResults[0];
    } else if (rid == null) {
      dbResults = await mysql.query(
        "SELECT rid,ralias,name,class,gmt_create,gmt_modified,rprofile,event_count FROM repairelements"
      );
    } else {
      dbResults = await mysql.query(
        "SELECT rid,ralias,name,class,gmt_create,gmt_modified,rprofile,event_count FROM repairelements WHERE rid=?",
        [rid]
      );
      dbResults = dbResults[0];
    }
  } catch (err) {
    return err;
  }
  await mysql.end();
  return dbResults;
};

exports.checkPassword = async (rid, password) => {
  let dbResults;
  try {
    dbResults = await mysql.query(
      "SELECT rpassword FROM repairelements WHERE rid=?",
      [rid]
    );
  } catch (error) {
    return err;
  }
  await mysql.end();
  let ans = dbResults[0].rpassword == password ? true : false;
  return ans;
};

exports.create = async (element) => {
  try {
    await mysql.query(
      "insert INTO repairelements (rid,ralias,rpassword,name,class,gmt_create,gmt_modified,rprofile,event_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
      [
        element.rid,
        element.alias,
        element.password,
        element.name,
        element.class,
        SYSDATE(),
        SYSDATE(),
        element.rprofile,
        0,
      ]
    );
  } catch (error) {
    return error;
  }
};

// exports.update = async (info) => {
//   await mysql.query(
//     "update repairelements set rpassword = ?,ralias = ?,name=?,class=?,gmt_modified = sysdate() where rid = ?;",
//     [info.password, info.alias, info.name, info.class, info.rid]
//   );
// };
