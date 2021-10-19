const jwt = require("jsonwebtoken");
const { cert } = require("../config/config");
exports.jsonPush = (str, data) => {
  if (!str) {
    str = "[]";
  }
  let temp = JSON.parse(str);
  temp.push(data);
  return JSON.stringify(temp);
};
exports.respond = (res, code, msg, data) => {
  let returnObj = {
    resultCode: code,
    resultMsg: msg,
    data: data,
  };
  if (code == 0) {
    returnObj = {
      resultCode: 0,
      resultMsg: "Success",
      data: data,
    };
  }
  res.send(returnObj);
};
// TODO ???????????
exports.dateToStr = (date, time) => {
  let ans;
  var year = date.getFullYear(); //年
  var month = date.getMonth(); //月
  var day = date.getDate(); //日
  ans =
    year +
    "-" +
    (month + 1 > 9 ? month + 1 : "0" + (month + 1)) +
    "-" +
    (day > 9 ? day : "0" + day);
  if (time) {
    var hours = date.getHours(); //时
    var min = date.getMinutes(); //分
    var second = date.getSeconds(); //秒
    ans +=
      " " +
      (hours > 9 ? hours : "0" + hours) +
      ":" +
      (min > 9 ? min : "0" + min) +
      ":" +
      (second > 9 ? second : "0" + second);
  }

  return ans;
};
exports.uuid = () => {
  var d = new Date().getTime();
  var uuid = "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(
    /[xy]/g,
    function (c) {
      var r = (d + Math.random() * 16) % 16 | 0;
      d = Math.floor(d / 16);
      return (c == "x" ? r : (r & 0x3) | 0x8).toString(16);
    }
  );
  return uuid;
};

const OSS = require("ali-oss");
const { ossConfig } = require("../config/config");

let client = new OSS(ossConfig);

exports.put = async (fileName, path) => {
  try {
    // 填写OSS文件完整路径和本地文件的完整路径。OSS文件完整路径中不能包含Bucket名称。
    // 如果本地文件的完整路径中未指定本地路径，则默认从示例程序所属项目对应本地路径中上传文件。
    return await client.put(fileName, path);
  } catch (e) {
    console.log(e);
  }
};

exports.createToken = (day, data) => {
  try {
    let token = jwt.sign(
      {
        exp: Math.floor(Date.now() / 1000) + day*24 * 60 * 60,
        data: data,
      },
      cert
    );
    return token;
  } catch (e) {
    console.log(e);
  }
}
