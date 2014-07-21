package main

import (
	"io"
	"log"
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
	parser *Parser
	filter *Filter
	config *Config
	render *Render
	lis    net.Listener
}

func NewServer(config *Config, filter *Filter) (s *Server) {
	if config == nil {
		panic("server: passed nil config")
	}
	if filter == nil {
		panic("server: passed nil filter")
	}
	s = &Server{
		http.NewServeMux(),
		&sync.WaitGroup{},
		NewParser(),
		filter,
		config,
		NewRender(config),
		nil,
	}
	s.HandleFunc(_PATH_TO_RSS, s.RSSHandler)
	s.HandleFunc(_PATH_TO_SHORT_LINKS, s.ShortLinkHandler)
	return s
}

func (s *Server) Start() (err error) {
	s.lis, err = net.Listen("tcp", s.config.Host+":"+s.config.Port)
	if err != nil {
		log.Fatal("server:", err)
	}
	return http.Serve(s.lis, s)
}

func (s *Server) ShutDown() error {
	s.Wait() // wait for all processed requests
	if s.lis == nil {
		return nil
	}
	defer func() { s.lis = nil }()
	return s.lis.Close()
}

func (s *Server) RemoveCache() error {
	return s.parser.RemoveCache()
}

func (s *Server) IsRunning() bool {
	return s.lis != nil
}

func (s *Server) RSSHandler(w http.ResponseWriter, r *http.Request) {
	s.Add(1) // signal that yet another request is processed
	defer r.Body.Close()
	var orders []*Order
	if err := r.ParseForm(); err != nil {
		log.Println("reading request error:", err)
	} else if resp, err := Load(r.Form.Get("url")); err != nil {
		log.Println("loading error:", err)
	} else {
		defer resp.Body.Close()
		orders, err = s.parser.Parse(resp)
		if err != nil && err != io.EOF {
			log.Println("can't read or parse response: ", err)
		}
		log.Printf("loaded %d orders\n", len(orders))
		if s.config.FilterEnabled && len(orders) > 0 {
			var filtered float32
			orders, filtered = s.filter.Execute(orders)
			log.Printf("filtered %.1f%%\n", filtered*100)
		}
	}
	var title string
	if URL, err := url.Parse(r.Form.Get("url")); err == nil {
		// call feed like search request
		title = URL.Query().Get("searchString")
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if len(orders) > 0 {
		s.render.Compose(title, orders)
	}
	if err := s.render.WriteTo(w); err != nil {
		log.Println("can't send response:", err)
	}
	s.Done() // signal that request was processed
}

func (s *Server) ShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if err := r.ParseForm(); err != nil {
		log.Println("bad request:", err)
		w.WriteHeader(http.StatusOK)
	} else {
		// redirect if order id was not passed also
		http.Redirect(w, r, MakeLink(r.Form.Get("order")),
			http.StatusFound)
	}
}
