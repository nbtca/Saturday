const { mysql } = require("../config/config");

exports.get = async (eid) => {
  let dbResults;
  try {
    if (eid) {
      dbResults = await mysql.query("SELECT * FROM `event` WHERE eid=?", [eid]);
      dbResults = dbResults[0];
    } else {
      dbResults = await mysql.query(
        "SELECT eid,user_description,status,rid,gmt_create,gmt_modified FROM `event` ORDER BY gmt_modified DESC"
      );
    }
  } catch (err) {
    return err;
  }
  await mysql.end();
  return dbResults;
};

exports.accept = async (rid, eventLog, eid) => {
  try {
    await mysql.query(
      "UPDATE `event` SET rid=?,event_log=?,status=1 WHERE eid=?",
      [rid, eventLog, eid]
    );
  } catch (err) {
    return err;
  }
  await mysql.end();
};
