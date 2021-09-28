package newznab

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"strings"
)

type Newznab struct {
	ApiKey string
	Host   string
}

type SearchResponseItem struct {
	GUID  string `json:"guid"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (n Newznab) SearchImdb(imdbId string) ([]SearchResponseItem, error) {
	url := fmt.Sprintf("http://%s/api?imdbid=%s&apikey=%s&t=movie&extended=1", n.Host, imdbId, n.ApiKey)

	fp := gofeed.NewParser()
	fmt.Printf("HTTP request: %s\n", strings.ReplaceAll(url, n.ApiKey, "xxx"))
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("fp.ParseURL: %w", err)
	}

	var itemList []SearchResponseItem

	for _, item := range feed.Items {
		searchItem := SearchResponseItem{Title: item.Title, GUID: item.GUID}
		if len(item.Enclosures) > 0 {
			searchItem.URL = item.Enclosures[0].URL
		}

		itemList = append(itemList, searchItem)
	}

	return itemList, nil
}
