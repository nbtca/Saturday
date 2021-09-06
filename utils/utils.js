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
// exports.wiper = (target) => {
//   t;
// };
exports.dateToStr = (date, time) => {
  let ans;
  var year = date.getFullYear(); //年
  var month = date.getMonth(); //月
  var day = date.getDate(); //日
  ans = year + "-" + (month + 1 > 9 ? month + 1 : "0" + (month + 1)) + "-" + (day > 9 ? day : "0" + day);
  if (time) {
    var hours = date.getHours(); //时
    var min = date.getMinutes(); //分
    var second = date.getSeconds(); //秒
    ans += " " + (hours > 9 ? hours : "0" + hours) + ":" + (min > 9 ? min : "0" + min) + ":" + (second > 9 ? second : "0" + second);
  }

  return ans;
};
