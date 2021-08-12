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

exports.submit = async (eventLog, description, req) => {
  try {
    await mysql.query(
      "UPDATE `event` SET event_log=?,repair_description=?,status=2 WHERE eid=?",
      [eventLog, description, req.body.eid]
    );
  } catch (err) {
    return err;
  }
  await mysql.end();
};

exports.cancel = async (eventLog, eid) => {
  try {
    await mysql.query(
      "UPDATE `event` SET rid=?,event_log=?,status=? WHERE eid=?",
      [null, eventLog, 0, eid]
    );
  } catch (err) {
    return err;
  }
  await mysql.end();
};

exports.close = async (aid, eventLog, status, eid) => {
  try {
    await mysql.query(
      "UPDATE `event` SET aid=?,event_log=?,status=? WHERE eid=?",
      [aid, eventLog, status, eid]
    );
  } catch (err) {
    return err;
  }
  await mysql.end();
};
exports.assign = async (rid, eventLog, eid) => {
  try {
    await mysql.query(
      "UPDATE `event` SET rid=?,event_log=?,status=? WHERE eid=?",
      [rid, eventLog, 1, eid]
    );
  } catch (err) {
    return err;
  }
  await mysql.end();
};
