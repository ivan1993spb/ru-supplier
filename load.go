package main

import (
	"errors"
	"net/http"
	"net/url"
)

func load(URL *url.URL) (*http.Response, error) {
	if URL == nil {
		return nil, errors.New("load(): passed nil url")
	}
	return http.DefaultClient.Do(&http.Request{
		URL:   URL,
		Proto: "HTTP/1.1",
		Header: http.Header{
			"User-Agent": {random_user_agent()},
		},
		Host: URL.Host,
	})
}
