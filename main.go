package main

import (
	"context"
	"fmt"

	"github.com/hurtki/school-events-bot/internal/bot"
	"github.com/hurtki/school-events-bot/internal/config"
)

func main() {
	botCfg, err := config.LoadBotConfig(config.EnvFileSource)
	if err != nil {
		fmt.Println(err)
		return
	}
	bot, err := bot.NewBot(botCfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = bot.NotifyAboutUpdate(context.Background())
	fmt.Println(err)
}
