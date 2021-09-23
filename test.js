const OSS = require("ali-oss");

let client = new OSS({
  region: "oss-cn-hangzhou",
  accessKeyId: "LTAI5tCyeZFdHskUvpTRCyPp",
  accessKeySecret: "r79n1DQaL5Y0lpremWGguBoHFA3aky",
  bucket: "sunday-res",
});

exports.put = async (fileName, path) => {
  try {
    // 填写OSS文件完整路径和本地文件的完整路径。OSS文件完整路径中不能包含Bucket名称。
    // 如果本地文件的完整路径中未指定本地路径，则默认从示例程序所属项目对应本地路径中上传文件。
    const result = await client.put(fileName, path);
    console.log(result);
  } catch (e) {
    console.log(e);
  }
};

