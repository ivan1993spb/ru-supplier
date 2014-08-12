package main

import (
	"net/url"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const _LOCAL_PROXY_DEF_HOST = "proxy-zakupki-gov-ru.local"

const (
	_WIN_TITLE       = "Генератор ссылок"
	_WIN_GEN_BUTTON  = "Генерировать"
	_WIN_COPY_BUTTON = "Копировать"
	_WIN_LABEL_HOST  = "Локальный хост"
	_WIN_LABEL_LINK  = "Ссылка на страницу с закупками"
)

const (
	_URL_REQUIRED_SCHEME               = "http"
	_URL_REQUIRED_HOST                 = "zakupki.gov.ru"
	_URL_REQUIRED_SORTING_TYPE         = "PUBLISH_DATE"
	_URL_REQUIRED_SORTING_DIRECTION    = "false"
	_URL_REQUIRED_QUICK_SEARCH_PATH    = "/epz/order/orderCsvSettings/quickSearch/download.html"
	_URL_REQUIRED_EXTENDED_SEARCH_PATH = "/epz/order/orderCsvSettings/extendedSearch/download.html"
)

var Paths = map[string]string{
	"/epz/order/extendedsearch/search.html": _URL_REQUIRED_EXTENDED_SEARCH_PATH,
	"/epz/order/quicksearch/search.html":    _URL_REQUIRED_QUICK_SEARCH_PATH,
	"/epz/order/quicksearch/update.html":    _URL_REQUIRED_QUICK_SEARCH_PATH,
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
				Text: _WIN_LABEL_HOST,
			},
			LineEdit{
				AssignTo: &le,
				Text:     _LOCAL_PROXY_DEF_HOST,
			},
			Label{
				Text: _WIN_LABEL_LINK,
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
								return
							}
							genURL := gen(URL, le.Text())
							if genURL != nil {
								te.SetText(genURL.String())
							}
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

func gen(URL *url.URL, host string) *url.URL {
	if !URL.IsAbs() {
		return nil
	}
	if URL.Scheme != _URL_REQUIRED_SCHEME {
		return nil
	}
	if URL.Host != _URL_REQUIRED_HOST {
		return nil
	}
	if path, ok := Paths[URL.Path]; ok {
		URL.Path = path
	} else {
		return nil
	}

	if URL.Path == _URL_REQUIRED_QUICK_SEARCH_PATH {
		URL.Query().Set("quickSearch", "true")
	} else {
		URL.Query().Set("quickSearch", "false")
	}
	if URL.Query().Get("sortBy") != _URL_REQUIRED_SORTING_TYPE {
		URL.Query().Set(
			"sortBy",
			_URL_REQUIRED_SORTING_TYPE,
		)
	}
	if URL.Query().Get("sortDirection") !=
		_URL_REQUIRED_SORTING_DIRECTION {

		URL.Query().Set(
			"sortDirection",
			_URL_REQUIRED_SORTING_DIRECTION,
		)
	}
	URL.Query().Set("userId", "null")
	URL.Query().Set("conf", "true;true;true;true;true;true;true;"+
		"true;true;true;true;true;true;true;true;true;true;")

	return &url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "/rss",
		RawQuery: url.Values{
			"url": {URL.String()},
		}.Encode(),
	}
}
