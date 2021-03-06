package main

import (
	"GithubRunnerGo/handlers/add"
	"GithubRunnerGo/handlers/all"
	"GithubRunnerGo/handlers/start"
	"GithubRunnerGo/utils"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"time"
)

func RegisterAllHandlers(updater *gotgbot.Updater, l *zap.SugaredLogger){
	start.LoadStartHandler(updater, l)
	add.LoadAddHandler(updater, l)
	all.LoadALLHandler(updater, l)
}

func main() {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder
	logger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), os.Stdout, zap.InfoLevel))
	defer logger.Sync() // flushes buffer, if any
	l := logger.Sugar()
	token := utils.GetBotToken()
	l.Info("Starting Bot.")
	l.Info("token: ", token)
	updater, err := gotgbot.NewUpdater(logger, token)
	l.Info("Got Updater")
	updater.UpdateGetter = ext.BaseRequester{
		Client: http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       time.Second * 65,
		},
		ApiUrl: ext.ApiUrl,
	}
	updater.Bot.Requester = ext.BaseRequester{Client: http.Client{Timeout: time.Second * 65}}
	if err != nil {
		l.Fatalw("failed to start updater", zap.Error(err))
	}
	l.Info("Starting updater")
	RegisterAllHandlers(updater, l)
	_ = updater.StartPolling()
	l.Info("Started Updater.")
	db := utils.DB{}
	db.Access(utils.GetDbUrl())
	db.GetAllConfigs()
	updater.Idle()
}

