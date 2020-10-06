package add

import (
	"GithubRunnerGo/utils"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"go.uber.org/zap"
	"strings"
)

func ADDHandler(b ext.Bot, u *gotgbot.Update) error {
	if !utils.IsUserOwner(u.EffectiveUser.Id){
		_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, "You are not allowed to use this bot", u.EffectiveMessage.MessageId)
		return nil
	}
	repo := strings.Split(strings.Split(u.EffectiveMessage.Text, utils.GetADDCommand())[1], " ")[0]
	token := strings.Split(strings.Split(u.EffectiveMessage.Text, utils.GetADDCommand())[1], " ")[1]
	_, _ = b.DeleteMessage(u.EffectiveChat.Id, u.EffectiveMessage.MessageId)
	if !utils.CheckValid(repo){
		_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, "You already hab runner for this repo running, Please remove that first", u.EffectiveMessage.MessageId)
		return nil
	}
	ch := make(chan bool)
	ch <- false
	utils.StartRunner(ch, repo,token)
	_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, "Added runner successfully", u.EffectiveMessage.MessageId)
	return nil
}


func LoadAddHandler(updater *gotgbot.Updater, l *zap.SugaredLogger) {
	defer l.Info("Add Module Loaded.")
	updater.Dispatcher.AddHandler(handlers.NewCommand(utils.GetADDCommand(), ADDHandler))
}