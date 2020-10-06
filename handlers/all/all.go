package all

import (
	"GithubRunnerGo/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"go.uber.org/zap"
)

func AllHandler(b ext.Bot, u *gotgbot.Update) error {
	if !utils.IsUserOwner(u.EffectiveUser.Id){
		_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, "You are not allowed to use this bot", u.EffectiveMessage.MessageId)
		return nil
	}
	message := fmt.Sprintf("<b>Here's the list of running github actions runner right now :-</b>\n\n")
	for i, v := range utils.GetAllInfo(){
		if len(v.Repo) != 0{
			message += fmt.Sprintf("<b>%v</b>: <code>%v</code>\n)", i, v.Repo)
		}
	}
	_, _ = b.ReplyHTML(u.EffectiveChat.Id, message, u.EffectiveMessage.MessageId)
	return nil
}

func LoadALLHandler(updater *gotgbot.Updater, l *zap.SugaredLogger) {
	defer l.Info("All runners Module Loaded.")
	updater.Dispatcher.AddHandler(handlers.NewCommand(utils.GetAllCommand(), AllHandler))
}

