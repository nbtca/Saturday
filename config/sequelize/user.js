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
      primaryKey: true,
      autoIncrement: false,
      comment: null,
      field: "uid"
    },
    uopenid: {
      type: DataTypes.CHAR(28),
      allowNull: true,
      defaultValue: null,
      primaryKey: false,
      autoIncrement: false,
      comment: null,
      field: "uopenid"
    },
    gmt_create: {
      type: DataTypes.DATE,
      allowNull: true,
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
    tableName: "user",
    comment: "",
    indexes: [],
    timestamps: false,
  };
  const UserModel = sequelize.define("user_model", attributes, options);
  return UserModel;
};