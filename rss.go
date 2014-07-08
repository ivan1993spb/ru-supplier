package main

import "github.com/gorilla/feeds"

const _DEFAULT_TITLE = "RSS feed for russian government orders"

func OrdersToRssFeed(orders []*Order, title string) *feeds.RssFeed {
	if len(title) == 0 {
		title = _DEFAULT_TITLE
	}
	feed := &feeds.RssFeed{
		Title: title,
		Link:  "http://zakupki.gov.ru",
		Description: "RSS feed for russian government orders in" +
			"really simple form with links and filters",
		ManagingEditor: "robot@localhost (Robot)",
		WebMaster:      "pushkin13@bk.ru (Pushkin Ivan)",
	}
	for _, order := range orders {
		feed.Items = append(feed.Items,
			&feeds.RssItem{
				Title:       order.Title(),
				Link:        order.ShortLink(),
				Description: order.Description(),
				Author:      order.OrganisationName,
				PubDate:     order.PubDateRFC1123(),
			})
	}
	return feed
}
