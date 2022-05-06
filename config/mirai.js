const { Bot } = require("mirai-js");
const { MiraiConfig, BotAccount } = require("../config");

const bot = new Bot();

try {
  bot.open({
    baseUrl: MiraiConfig.url,
    verifyKey: MiraiConfig.key,
    qq: BotAccount.id,
  });
} catch (error) {
  console.log(error);
}

module.exports = bot;
