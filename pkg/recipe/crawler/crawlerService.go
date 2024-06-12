package crawler

import (
	"Food/pkg/clients/maggi_ng"
	crawler2 "Food/pkg/recipe/crawler/maggi_ng"
	"Food/pkg/recipe/model"
)

type ICrawler interface {
	CrawlRecipe() ([]model.RequestData, error)
}

func AddCrawler(crawlerNames []string) []ICrawler {
	var crawlers []ICrawler
	for _, name := range crawlerNames {
		crawlers = append(crawlers, GetCrawler(name))
	}
	return crawlers
}

func GetCrawler(name string) ICrawler {
	switch name {
	case "maggi_ng":
		client := maggi_ng.NewClient()
		return crawler2.NewMaggiCrawler(client)
	default:
		panic("Unknown service type")
	}
}
