package main

import (
	"fmt"
	"log"
	"os"
	"tgBot/internal/postgresDB"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gopkg.in/telebot.v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка подгрузки окружения: ", err)
	}

	db := postgresDB.InitDB()
	defer db.Close()
	postgresDB.CreateTable(db)

	pref := telebot.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal("Ошибка при создании бота:", err)
		return
	}

	register := func(c telebot.Context) {
		postgresDB.SaveUser(db, c.Sender().ID, c.Sender().Username)
	}

	bot.Handle("/beer", func(c telebot.Context) error {
		register(c)

		log.Printf("[%s] написал: %s", c.Sender().Username, c.Text())
		allUsers := postgresDB.GetAllUsers(db)
		notification := fmt.Sprintf("‼️Минуточку внимания, @%s открыл бутылочку хмельного‼️", c.Sender().Username)
		if c.Sender().Username == "" {
			notification = fmt.Sprintf("‼️Минуточку внимания, %s открыл бутылочку хмельного‼️", c.Sender().FirstName)
		}

		for _, id := range allUsers {
			if id == c.Sender().ID {
				continue
			}
			bot.Send(&telebot.User{ID: id}, notification)
		}
		return c.Send("Всем пришло уведомление о твоём намерении выпить пива")
	})

	bot.Handle(telebot.OnSticker, func(c telebot.Context) error {
		register(c)
		log.Printf("[%s] написал: %s", c.Sender().Username, "стикер")
		return c.Send(c.Message().Sticker)
	})

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		register(c)
		log.Printf("[%s] написал: %s", c.Sender().Username, c.Text())
		return c.Send(c.Message().Text)
	})

	bot.Start()
}
