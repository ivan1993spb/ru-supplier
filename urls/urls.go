package main

import (
	"fmt"
	"net/url"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	_WIN_TITLE       = "Генератор ссылок"
	_WIN_GEN_BUTTON  = "Генерировать"
	_WIN_COPY_BUTTON = "Копировать"
)

const (
	_LABEL_HOST = "Локальный хост"
	_LABEL_LINK = "Ссылка на страницу с закупками"
)

const _LOCAL_PROXY_DEF_HOST = "proxy-zakupki-gov-ru.local"

const (
	_URL_REQUIRED_SCHEME = "http"
	_URL_REQUIRED_HOST   = "zakupki.gov.ru"

	_URL_REQUIRED_SORTING_TYPE      = "PUBLISH_DATE"
	_URL_REQUIRED_SORTING_DIRECTION = "false"
)

var Paths = map[string]string{
	"/epz/order/extendedsearch/search.html": "/epz/order/orderCsvSettings/extendedSearch/download.html",
	"/epz/order/quicksearch/search.html":    "/epz/order/orderCsvSettings/quickSearch/download.html",
	"/epz/order/quicksearch/update.html":    "/epz/order/orderCsvSettings/extendedSearch/download.html",
}

func main() {
	var (
		te *walk.TextEdit
		le *walk.LineEdit
	)

	MainWindow{
		Title:   _WIN_TITLE,
		MinSize: Size{300, 400},
		Layout:  VBox{},
		Children: []Widget{
			Label{
				Text: _LABEL_HOST,
			},
			LineEdit{
				AssignTo: &le,
				Text:     _LOCAL_PROXY_DEF_HOST,
			},
			Label{
				Text: _LABEL_LINK,
			},
			TextEdit{
				AssignTo: &te,
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: _WIN_GEN_BUTTON,
						OnClicked: func() {
							URL, err := url.Parse(te.Text())
							if err != nil {
								fmt.Println(1, err)
								return
							}
							if !URL.IsAbs() {
								fmt.Println(2, err)
								return
							}
							if URL.Scheme != _URL_REQUIRED_SCHEME {
								fmt.Println(3, err)
								return
							}
							if URL.Host != _URL_REQUIRED_HOST {
								fmt.Println(4, err)
								return
							}
							if redirect, ok := Paths[URL.Path]; ok {
								URL.Path = redirect
							} else {
								fmt.Println(5, err)
								return
							}
							if URL.Query().Get("sortBy") != _URL_REQUIRED_SORTING_TYPE {
								URL.Query().Set("sortBy", _URL_REQUIRED_SORTING_TYPE)
							}
							if URL.Query().Get("sortDirection") != _URL_REQUIRED_SORTING_DIRECTION {
								URL.Query().Set("sortDirection", _URL_REQUIRED_SORTING_DIRECTION)
							}
							te.SetText((&url.URL{
								Scheme: "http",
								Host:   le.Text(),
								Path:   "/rss",
								RawQuery: url.Values{
									"url": {URL.String()},
								}.Encode(),
							}).String())
						},
					},
					PushButton{
						Text: _WIN_COPY_BUTTON,
						OnClicked: func() {
							walk.Clipboard().SetText(te.Text())
						},
					},
				},
			},
		},
	}.Run()
}
