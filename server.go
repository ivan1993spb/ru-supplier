package main

import (
	"encoding/xml"
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
	if err := r.ParseForm(); err != nil {
		log.Warning.Println("reading request error:", err)
	} else if rawurl := r.Form.Get("url"); len(rawurl) == 0 {
		log.Warning.Println("empty request was accepted")
	} else if URL, err := url.Parse(rawurl); err != nil {
		log.Warning.Println("passed invalid url:", err)
	} else if !URL.IsAbs() {
		log.Warning.Println("passed url isn't absolute")
	} else if resp, err := load(URL); err != nil {
		log.Error.Println("connection problems:", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Warning.Printf("server return %q\n",
				strings.ToLower(resp.Status))
			w.WriteHeader(resp.StatusCode)
		} else {
			orders, err := parse(resp.Body)
			if err != nil {
				log.Warning.Println("can't read or parse response:",
					err)
			}
			if len(orders) > 0 {
				orders = filter(orders)
			}
			w.Header().Set("Content-Type",
				"application/xml; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(xml.Header))
			err = xml.NewEncoder(w).Encode(
				OrdersToRssFeed(
					orders,
					// search string as title (if exists)
					URL.Query().Get("searchString"),
				).FeedXml(),
			)
			if err != nil {
				log.Error.Println("can't send response:", err)
			}
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
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
