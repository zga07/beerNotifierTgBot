package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3"
)

var userStates = make(map[int64]string)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка подгрузки токена: ", err)
	}
	token := os.Getenv("BOT_TOKEN")
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal("Ошибка при создании бота:", err)
		return
	}

	bot.Handle("/books", func(c telebot.Context) error {
		userStates[c.Sender().ID] = "waiting_for_page"
		return c.Send("Какую страницу книг мне открыть? Пришли номер (от 1 до 50):")
	})

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		userID := c.Sender().ID

		if userStates[userID] == "waiting_for_page" {
			pageNum, err := strconv.Atoi(c.Text())
			if err != nil || pageNum < 1 || pageNum > 50 {
				return c.Send("Пожалуйста, введи число от 1 до 50.")
			}

			delete(userStates, userID)
			c.Send(fmt.Sprintf("Ищу книги на странице %d...", pageNum))

			result := GetBooks(pageNum)
			return c.Send(result)
		}

		log.Printf("[%s] написал: %s", c.Sender().Username, c.Text())
		return c.Send(c.Text())
	})

	bot.Handle("/help", func(c telebot.Context) error {
		log.Printf("[%s] написал: %s", c.Sender().Username, c.Text())
		response := fmt.Sprintf("привет, %s, я повторяю сообщения", c.Sender().FirstName)
		return c.Send(response)
	})

	bot.Handle(telebot.OnSticker, func(c telebot.Context) error {
		log.Printf("[%s] написал: %s", c.Sender().Username, "стикер")
		return c.Send(c.Message().Sticker)
	})

	bot.Start()
}

func GetBooks(page int) string {
	url := fmt.Sprintf("https://books.toscrape.com/catalogue/page-%d.html", page)

	response, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("Ошибка подключения к сайту: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Sprintf("Ошибка сайта: %s", err)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return fmt.Sprintf("Ошибка при чтении данных: %s", err)
	}
	var result string
	result = fmt.Sprintf("Страница №%d:\n\n", page)
	doc.Find("article.product_pod").Each(func(i int, s *goquery.Selection) {
		fullTitle, _ := s.Find("h3 a").Attr("title")
		price := s.Find(".price_color").Text()

		result += fmt.Sprintf("%d. %s\nСтоимость: %s\n\n", i+1, fullTitle, price)
	})

	return result
}
