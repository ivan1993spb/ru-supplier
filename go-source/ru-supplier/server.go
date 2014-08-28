package main

import (
	"errors"
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

// RSS protocol required port 80
const _RSS_REQUIRED_PORT = "80"

type Server struct {
	*http.ServeMux
	*sync.WaitGroup
	reader *OrderReader
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
		NewOrderReader(),
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
	if s.lis != nil {
		return errors.New("server is already running")
	}
	if s.config.Port != _RSS_REQUIRED_PORT {
		log.Println("RSS protocol required port 80")
	}
	s.lis, err = net.Listen("tcp", s.config.Host+":"+s.config.Port)
	if err != nil {
		log.Fatal("server:", err)
	}
	return http.Serve(s.lis, s)
}

func (s *Server) ShutDown() error {
	if s.lis == nil {
		return errors.New("server is already stopped")
	}
	s.Wait() // wait for all processed requests
	if s.lis == nil {
		return nil
	}
	defer func() { s.lis = nil }()
	return s.lis.Close()
}

func (s *Server) RemoveCache() error {
	return s.reader.RemoveCache()
}

func (s *Server) IsRunning() bool {
	return s.lis != nil
}

func (s *Server) RSSHandler(w http.ResponseWriter, r *http.Request) {
	s.Add(1) // signal that yet another request is processed
	defer r.Body.Close()
	var orders []*Order
	if resp, err := Load(r.FormValue("url")); err != nil {
		log.Println("loading error:", err)
	} else {
		defer resp.Body.Close()
		orders, err = s.reader.ReadOrders(resp)
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
	if URL, err := url.Parse(r.Form.Get("url")); err == nil {
		// call feed like search request
		s.render.SetTitle(URL.Query().Get("searchString"))
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if len(orders) > 0 {
		s.render.Compose(orders)
	}
	if err := s.render.WriteTo(w); err != nil {
		log.Println("can't send response:", err)
	}
	// clear feed
	s.render.Clear()
	// signal that request was processed
	s.Done()
}

func (s *Server) ShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	// redirect if order id was not passed also
	http.Redirect(w, r, MakeLink(r.FormValue("order")), http.StatusFound)
	r.Body.Close()
}
