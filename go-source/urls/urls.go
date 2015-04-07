package main

import "net/url"

const (
	_URL_REQUIRED_SCHEME               = "http"
	_URL_REQUIRED_HOST                 = "zakupki.gov.ru"
	_URL_REQUIRED_ALIAS_HOST           = "www.zakupki.gov.ru"
	_URL_REQUIRED_SORTING_TYPE         = "PUBLISH_DATE"
	_URL_REQUIRED_SORTING_DIRECTION    = "false"
	_URL_REQUIRED_QUICK_SEARCH_PATH    = "/epz/order/quicksearch/orderCsvSettings/quickSearch/download.html"
	_URL_REQUIRED_EXTENDED_SEARCH_PATH = "/epz/order/extendedsearch/orderCsvSettings/extendedSearch/download.html"
)

var Paths = map[string]string{
	"/epz/order/extendedsearch/search.html": _URL_REQUIRED_EXTENDED_SEARCH_PATH,
	"/epz/order/quicksearch/search.html":    _URL_REQUIRED_QUICK_SEARCH_PATH,
	"/epz/order/quicksearch/update.html":    _URL_REQUIRED_QUICK_SEARCH_PATH,
}

func generateURL(URL *url.URL, host string) *url.URL {
	if !URL.IsAbs() {
		return nil
	}
	if URL.Scheme != _URL_REQUIRED_SCHEME {
		return nil
	}
	if URL.Host != _URL_REQUIRED_HOST {
		if URL.Host != _URL_REQUIRED_ALIAS_HOST {
			return nil
		}
	}
	if path, ok := Paths[URL.Path]; ok {
		URL.Path = path
	} else {
		return nil
	}

	vals := URL.Query()
	if URL.Path == _URL_REQUIRED_QUICK_SEARCH_PATH {
		vals.Set("quickSearch", "true")
	} else {
		vals.Set("quickSearch", "false")
	}
	vals.Set("sortBy", _URL_REQUIRED_SORTING_TYPE)
	vals.Set("sortDirection", _URL_REQUIRED_SORTING_DIRECTION)
	vals.Set("userId", "null")
	vals.Set("conf", "true;true;true;true;true;true;true;"+
		"true;true;true;true;true;true;true;true;true;true;")
	URL.RawQuery = vals.Encode()

	return &url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "/rss",
		RawQuery: url.Values{
			"url": {URL.String()},
		}.Encode(),
	}
}
