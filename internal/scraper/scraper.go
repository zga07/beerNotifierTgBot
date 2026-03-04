package scraper

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

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
