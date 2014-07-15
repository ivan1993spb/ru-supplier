package main

import (
	"encoding/xml"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
)

const (
	_PATH_TO_RSS         = "/rss"
	_PATH_TO_SHORT_LINKS = "/open"
)

type Server struct {
	*http.ServeMux
	*sync.WaitGroup
	lis    net.Listener
	log    *Log
	filter *Filter
}

func NewServer() (s *Server) {
	s = &Server{http.NewServeMux(), &sync.WaitGroup{}, nil}
	s.HandleFunc(_PATH_TO_RSS, s.RSSHandler)
	s.HandleFunc(_PATH_TO_SHORT_LINKS, s.ShortLinkHandler)
	return s
}

func (s *Server) Serve(l net.Listener) error {
	if l == nil {
		panic("server: passed nil listener")
	}
	s.lis = l
	return http.Serve(l, s)
}

func (s *Server) ShutDown() error {
	s.Wait() // wait for all processed requests
	defer func() { s.lis = nil }()
	return s.lis.Close()
}

func (s *Server) RSSHandler(w http.ResponseWriter, r *http.Request) {
	s.Add(1) // signal that yet another request is processed
	defer r.Body.Close()
	var orders []*Order
	if err := r.ParseForm(); err != nil {
		s.log.Warning.Println("reading request error:", err)
	} else if resp, err := Load(r.Form.Get("url")); err != nil {
		s.log.Warning.Println("loading error:", err)
	} else {
		defer resp.Body.Close()
		orders, err = Parse(resp)
		if err != nil && err != io.EOF {
			s.log.Warning.Println("can't read or parse response: ", err)
		}
		if len(orders) > 0 {
			orders = filter(orders)
		}
	}
	var title string
	if URL, err := url.Parse(r.Form.Get("url")); err == nil {
		// call feed like search request
		title = URL.Query().Get("searchString")
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(xml.Header))
	err := xml.NewEncoder(w).Encode(
		OrdersToRssFeed(title, orders).FeedXml(),
	)
	if err != nil {
		s.log.Error.Println("can't send response:", err)
	}
	s.Done() // signal that request was processed
}

func (s *Server) ShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if err := r.ParseForm(); err != nil {
		s.log.Warning.Println("bad request:", err)
		w.WriteHeader(http.StatusOK)
	} else {
		// redirect if order id was not passed also
		http.Redirect(w, r, MakeLink(r.Form.Get("order")),
			http.StatusFound)
	}
}
