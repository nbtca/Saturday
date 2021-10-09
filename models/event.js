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

exports.create = async (event) => {
  var sql =
    "INSERT INTO EVENT ( uid, model, eqq, ephone, econtact_preference,user_description,event_log,status,gmt_create, gmt_modified)VALUES (?,?,?,?,?,?,?,0,SYSDATE(),SYSDATE())";
  var sqlParams = [
    event.uid,
    event.model,
    event.qq,
    event.phone,
    event.preference,
    event.description,
    event.eventLog,
  ];
  try {
    await mysql.query(sql, sqlParams);
  } catch (error) {
    return error;
  }
};
exports.update = async (event) => {
  let sql =
    "UPDATE EVENT SET rid=?, model =?,ephone=?,eqq=?,user_description =?,econtact_preference=?,status=?,repair_description=?,event_log=?,gmt_modified = SYSDATE() WHERE eid =?";
  let sqlParams = [
    event.rid,
    event.model,
    event.ephone,
    event.eqq,
    event.user_description,
    event.econtact_preference,
    event.status,
    event.repair_description,
    event.event_log,
    event.eid,
  ];
  try {
    await mysql.query(sql, sqlParams);
  } catch (error) {
    return error;
  }
};

// exports.accept = async (rid, eventLog, eid) => {
//   try {
//     await mysql.query(
//       "UPDATE `event` SET rid=?,event_log=?,status=1 WHERE eid=?",
//       [rid, eventLog, eid]
//     );
//   } catch (err) {
//     return err;
//   }
//   await mysql.end();
// };

// exports.submit = async (eventLog, description, req) => {
//   try {
//     await mysql.query(
//       "UPDATE `event` SET event_log=?,repair_description=?,status=2 WHERE eid=?",
//       [eventLog, description, req.body.eid]
//     );
//   } catch (err) {
//     return err;
//   }
//   await mysql.end();
// };

// exports.cancel = async (eventLog, eid) => {
//   try {
//     await mysql.query(
//       "UPDATE `event` SET rid=?,event_log=?,status=? WHERE eid=?",
//       [null, eventLog, 0, eid]
//     );
//   } catch (err) {
//     return err;
//   }
//   await mysql.end();
// };

// exports.close = async (aid, eventLog, status, eid) => {
//   try {
//     await mysql.query(
//       "UPDATE `event` SET aid=?,event_log=?,status=? WHERE eid=?",
//       [aid, eventLog, status, eid]
//     );
//   } catch (err) {
//     return err;
//   }
//   await mysql.end();
// };
// exports.assign = async (rid, eventLog, eid) => {
//   try {
//     await mysql.query(
//       "UPDATE `event` SET rid=?,event_log=?,status=? WHERE eid=?",
//       [rid, eventLog, 1, eid]
//     );
//   } catch (err) {
//     return err;
//   }
//   await mysql.end();
// };

// exports.delete = async (eventLog, eid) => {
//   try {
//     await mysql.query("update event set status=?,event_log=? where eid=?", [
//       -1,
//       eventLog,
//       eid,
//     ]);
//     await mysql.end();
//   } catch (error) {
//     return error;
//   }
// };
