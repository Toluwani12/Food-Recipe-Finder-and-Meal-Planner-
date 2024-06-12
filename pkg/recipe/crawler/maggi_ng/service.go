package maggi_ng

import (
	"Food/pkg/clients/maggi_ng"
	"Food/pkg/recipe/model"
)

type MaggiCrawler struct {
	Client *maggi_ng.Client
}

func NewMaggiCrawler(client *maggi_ng.Client) *MaggiCrawler {
	return &MaggiCrawler{Client: client}
}

func (mc *MaggiCrawler) CrawlRecipe() (*[]model.RequestData, error) {
	params := maggi_ng.SearchParams{
		ContentType: "",
		Range:       0,
		SearchPage:  false,
		Filters:     false,
	}
	data, err := mc.Client.FetchLinks(params)
	if err != nil {
		return nil, err
	}

	return data, nil
}
