const { Message, Middleware } = require("mirai-js");
const bot = require("../config/mirai");
const { GroupID } = require("../config/config");
const TestGroupID = "960601785"; // test

class Bot {
  constructor() { }

  newEventTemplate(event) {
    // let xmlContent =
    //   "<?xml version='1.0' encoding='UTF-8' standalone='yes' ?> <msg serviceID='1' templateID='-1' action='app' actionData='com.android.camera' brief='测试' sourceMsgId='0' url='' flag='1' adverSign='0' multiMsgFlag='0'> <item><title>测试</title></item> <item layout='6' ><picture cover='https://w.wallhaven.cc/full/3z/wallhaven-3z32j3.jpg' /></item></msg>";
    return new Message()
      .addText("新维修事件:\n")
      .addText("问题: ")
      .addText(event.user_description + "\n")
      .addText("型号: ")
      .addText(event.model + "\n")
      .addAtAll();
    // .addXml(xmlContent);
  }

  async msgTest() {
    await bot.sendMessage({
      group: TestGroupID,
      message: new Message().addText("test"),
    });
  }

  async sendGroupMsg(msg) {
    try {
      await bot.sendMessage({
        group: GroupID,
        message: msg,
      });
    } catch (error) {
      console.log(error);
    }
  }
}
module.exports = new Bot();

// const GroupFilter = new Middleware().textProcessor();

// bot.on(
//   "GroupMessage",
//   GroupFilter.done(async data => {
//     bot.sendMessage({
//       group: data.sender.group.id,
//       message: new Message().addText(data.text),
//     });
//   })
// );

// have to send manually once before use
// exports.exports.sendGroupMsg = async msg => {
//   try {
//     await bot.sendMessage({
//       group: GroupID,
//       message: msg,
//     });
//   } catch (error) {
//     console.log(error);
//   }
// };

