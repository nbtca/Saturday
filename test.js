const { Bot, Message, Middleware } = require("mirai-js");
const bot = require("./config/mirai");
// const bot = getBot();

const GroupID = "960601785"; // test
// const GroupID = "***REMOVED***";

const eventTemplate = event => {
  return new Message().addText("New Event").addText("New Event");
};

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
async function sendMsg() {
  try {
    await bot.sendMessage({
      group: GroupID,
      message: eventTemplate(),
    });
  } catch (error) {
    console.log(error);
  }
}

setTimeout(() => {
  sendMsg();
}, 2000);

// exports.msgTest = () => {
//   bot.sendMessage({
//     group: GroupID,
//     message: new Message().addText("test"),
//   });
// };
