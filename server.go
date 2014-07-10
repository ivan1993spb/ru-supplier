package main

import (
	"encoding/xml"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

const (
	_PATH_TO_RSS         = "/rss"
	_PATH_TO_SHORT_LINKS = "/open"
)

type Server struct {
	*http.ServeMux
	lis net.Listener
}

func NewServer(lis net.Listener) (s *Server) {
	s = &Server{http.NewServeMux(), lis}
	s.HandleFunc(_PATH_TO_RSS, s.RSSHandler)
	s.HandleFunc(_PATH_TO_SHORT_LINKS, ShortLinkHandler)
	return s
}

func (s *Server) Bind() (net.Listener, http.Handler) {
	return s.lis, s
}

func (s *Server) RSSHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var orders []*Order
	if err := r.ParseForm(); err != nil {
		log.Warning.Println("reading request error:", err)
	} else if resp, err := load(r.Form.Get("url")); err != nil {
		log.Warning.Println("loading error:", err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			orders, err = parse(resp)
			if err != nil && err != io.EOF {
				log.Warning.Println("can't read or parse response:",
					err)
			}
			if len(orders) > 0 {
				orders = filter(orders)
			}
		} else {
			log.Warning.Printf("server return %q\n",
				strings.ToLower(resp.Status))
		}
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(xml.Header))
	var title string
	if URL, err := url.Parse(r.Form.Get("url")); err == nil {
		title = URL.Query().Get("searchString")
	}
	err := xml.NewEncoder(w).Encode(
		OrdersToRssFeed(orders, title).FeedXml(),
	)
	if err != nil {
		log.Error.Println("can't send response:", err)
	}
}

func ShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if err := r.ParseForm(); err != nil {
		log.Warning.Println("bad request:", err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// also redirect if order id was not passed
		http.Redirect(w, r, MakeLink(r.Form.Get("order")),
			http.StatusFound)
	}
}
