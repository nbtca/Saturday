const { mysql } = require("../config/config");

exports.get = async (rid) => {
  console.log("e;e");

  let dbResults;
  try {
    if (rid == -1) {
      dbResults = await mysql.query(
        "SELECT rid,ralias,name,class,gmt_create,gmt_modified,rprofile,event_count FROM repairelements WHERE rid !=0000000000 ORDER BY RAND() LIMIT 1;"
      );
    } else if (rid == null) {
      dbResults = await mysql.query(
        "SELECT rid,ralias,name,class,gmt_create,gmt_modified,rprofile,event_count FROM repairelements"
      );
    } else {
      dbResults = await mysql.query(
        "SELECT rid,ralias,name,class,gmt_create,gmt_modified,rprofile,event_count FROM repairelements WHERE rid=?",
        [rid]
      );
    }
  } catch (err) {
    console.log(err);
    return err;
  }
  await mysql.end();
  return dbResults;
};

exports.getAlias = async (rid) => {
  let dbResults;
  try {
    dbResults = await mysql.query(
      "SELECT ralias FROM repairelements WHERE rid=?",
      [rid]
    );
  } catch (err) {
    return err;
  }
  await mysql.end();
  return dbResults[0];
};
