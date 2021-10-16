const bot = new Bot();
bot.open({
  baseUrl: MiraiConfig.url,
  verifyKey: MiraiConfig.key,
  qq: BotAccount.id,
});