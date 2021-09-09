const OSS = require("ali-oss");

let client = new OSS({
  region: "oss-cn-hangzhou",
  accessKeyId: "LTAI5t7vx8WjAzNAnP8xcKd4",
  accessKeySecret: "aofNwWUIeX9QjyXfiomHd8ij5ewGec",
  bucket: "sunday-res",
});

async function putBuffer(file) {
  try {
    let result = await client.put("object-name", file);
    console.log(result);
  } catch (e) {
    console.log(e);
  }
}

// putBuffer();
module.exports = putBuffer;
