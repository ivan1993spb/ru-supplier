package main

import (
	"errors"
	"net/http"
	"net/url"
)

func load(rawurl string) (*http.Response, error) {
	if len(rawurl) == 0 {
		return nil, errors.New("load(): passed empty string")
	}
	URL, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if !URL.IsAbs() {
		return nil, errors.New("passed url isn't absolute")
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
