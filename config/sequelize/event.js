const {
  DataTypes
} = require('sequelize');
const sequelize = require("../db");

module.exports = () => {
  const attributes = {
    uid: {
      type: DataTypes.STRING(40),
      allowNull: false,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "uid",
      references: {
        key: "uid",
        model: "user_model"
      }
    },
    eid: {
      type: DataTypes.STRING(40),
      allowNull: false,
      defaultValue: null,
      primaryKey: true,
      autoIncrement: false,
      comment: null,
      field: "eid"
    },
    model: {
      type: DataTypes.STRING(20),
      allowNull: true,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "model"
    },
    ephone: {
      type: DataTypes.STRING(11),
      allowNull: false,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "ephone"
    },
    eqq: {
      type: DataTypes.STRING(20),
      allowNull: true,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "eqq"
    },
    econtact_preference: {
      type: DataTypes.INTEGER(4),
      allowNull: false,
      defaultValue: "0",
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "econtact_preference"
    },
    user_description: {
      type: DataTypes.STRING(500),
      allowNull: true,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "user_description"
    },
    repair_description: {
      type: DataTypes.STRING(1000),
      allowNull: true,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "repair_description"
    },
    event_log: {
      type: DataTypes.STRING(3000),
      allowNull: true,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "event_log"
    },
    status: {
      type: DataTypes.INTEGER(4),
      allowNull: false,
      defaultValue: "0",
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "status"
    },
    rid: {
      type: DataTypes.CHAR(10),
      allowNull: true,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "rid",
      references: {
        key: "rid",
        model: "repairelements_model"
      }
    },
    closed_by: {
      type: DataTypes.CHAR(10),
      allowNull: true,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "closed_by",
      references: {
        key: "rid",
        model: "repairelements_model"
      }
    },
    gmt_create: {
      type: DataTypes.DATE,
      allowNull: false,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "gmt_create"
    },
    gmt_modified: {
      type: DataTypes.DATE,
      allowNull: false,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "gmt_modified"
    }
  };
  const options = {
    tableName: "event",
    comment: "",
    timestamps: false,
    indexes: [{
      name: "fk_Event_Admin_2",
      unique: false,
      type: "BTREE",
      fields: ["closed_by"]
    }, {
      name: "fk_Event_User_1",
      unique: false,
      type: "BTREE",
      fields: ["uid"]
    }, {
      name: "fk_Event_repairElements_1",
      unique: false,
      type: "BTREE",
      fields: ["rid"]
    }]
  };
  const EventModel = sequelize.define("event_model", attributes, options);
  return EventModel;
};