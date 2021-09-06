const {
  DataTypes
} = require('sequelize');

module.exports = sequelize => {
  const attributes = {
    uid: {
      type: DataTypes.BIGINT,
      allowNull: false,
      defaultValue: null,
      primaryKey: true,
      autoIncrement: true,
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
    indexes: []
  };
  const UserModel = sequelize.define("user_model", attributes, options);
  return UserModel;
};