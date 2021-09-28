package newznab

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"strconv"
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
	Size  int64  `json:"size"`
}

func parseAttrs(item *gofeed.Item) (map[string]string, error) {
	if item == nil {
		return nil, fmt.Errorf("cannot parse nil item")
	}

	items := make(map[string]string)

	newznabExt := item.Extensions["newznab"]
	if newznabExt != nil {
		attrs := newznabExt["attr"]
		if attrs != nil {
			for _, extension := range attrs {
				innerAttrs := extension.Attrs

				var name, value string

				for attrName, attrValue := range innerAttrs {
					if attrName == "name" {
						name = attrValue
					}
					if attrName == "value" {
						value = attrValue
					}
				}

				if name != "" {
					items[name] = value
				}
			}
		}
	}

	return items, nil
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
		attrs, err := parseAttrs(item)
		if err != nil {
			return nil, fmt.Errorf("parseAttrs: %w", err)
		}
		searchItem := SearchResponseItem{Title: item.Title, GUID: item.GUID}
		if attrs["size"] != "" {
			size, err := strconv.Atoi(attrs["size"])
			if err != nil {
				fmt.Printf("invalid size format. cannot convert '%s' to integer", attrs["size"])
			} else {
				searchItem.Size = int64(size)
			}
		}
		if len(item.Enclosures) > 0 {
			searchItem.URL = item.Enclosures[0].URL
		}

		itemList = append(itemList, searchItem)
	}

	return itemList, nil
}
