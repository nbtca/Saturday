const { DataTypes } = require("sequelize");
const sequelize = require("../db");

module.exports = () => {
  const attributes = {
    setting: {
      type: DataTypes.STRING(10000),
      allowNull: false,
      defaultValue: null,
      primaryKey: true,
      autoIncrement: false,
      comment: null,
      field: "setting",
    },
  };
  const options = {
    tableName: "setting",
    comment: "",
    timestamps: false,
    indexes: [
    ],
  };
  const SettingModel = sequelize.define("setting_model", attributes, options);
  return SettingModel;
};
