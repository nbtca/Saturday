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
// exports.respond = (res, code, msg, data) => {
//   // let returnObj = {
//   //   resultCode: code,
//   //   resultMsg: msg,
//   //   data: data,
//   // };
//   // if (code == 0) {
//   //   returnObj = {
//   //     resultCode: 0,
//   //     resultMsg: "Success",
//   //     data: data,
//   //   };
//   // }
//   // res.send(returnObj);
//   res.status(code).send(msg || data);
// };
