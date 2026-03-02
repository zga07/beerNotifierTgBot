package main

import (
	"log"
	"time"

	"gopkg.in/telebot.v3"
)

func main() {
	pref := telebot.Settings{
		Token:  "token",
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal("Ошибка при создании бота:", err)
		return
	}

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		log.Printf("[%s] написал: %s", c.Sender().Username, c.Text())
		return c.Send(c.Text())
	})

	bot.Start()
}
