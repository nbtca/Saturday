const { Bot, Message, Middleware } = require("mirai-js");
const bot = require("./config/mirai");
// const bot = getBot();

const GroupID = "960601785"; // test
// const GroupID = "***REMOVED***";

const NewEventTemplate = event => {
  return new Message().addText("New Event").addText("New Event");
};

const GroupFilter = new Middleware().textProcessor().groupFilter([GroupID]);

bot.on(
  "GroupMessage",
  GroupFilter.done(async data => {
    bot.sendMessage({
      group: data.sender.group.id,
      message: new Message().addText(data.text),
    });
  })
);

// have to send manually once before use
async function sendMsg() {
  try {
    await bot.sendMessage({
      group: GroupID,
      message: NewEventTemplate(),
    });
  } catch (error) {
    console.log(error);
  }
}

setTimeout(() => {
  sendMsg();
}, 2000);
