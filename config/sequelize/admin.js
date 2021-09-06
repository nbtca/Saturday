const { DataTypes } = require("sequelize");
const sequelize = require("../db");
module.exports = () => {
  const attributes = {
    rid: {
      type: DataTypes.CHAR(10),
      allowNull: false,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "rid",
      references: {
        key: "rid",
        model: "repairelements_model",
      },
    },
    aid: {
      type: DataTypes.BIGINT,
      allowNull: false,
      defaultValue: null,
      primaryKey: true,
      autoIncrement: true,
      comment: null,
      field: "aid",
    },
    gmt_create: {
      type: DataTypes.DATE,
      allowNull: false,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "gmt_create",
    },
    gmt_modified: {
      type: DataTypes.DATE,
      allowNull: false,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "gmt_modified",
    },
  };
  const options = {
    tableName: "admin",
    comment: "",
    timestamps: false,
    indexes: [
      {
        name: "fk_Admin_repairElements_1",
        unique: false,
        type: "BTREE",
        fields: ["rid"],
      },
    ],
  };
  const AdminModel = sequelize.define("admin_model", attributes, options);
  return AdminModel;
};
