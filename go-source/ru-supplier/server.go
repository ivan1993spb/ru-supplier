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

type ZakupkiProxyServer interface {
	Start() error
	ShutDown() error
	IsRunning() bool
	RemoveCache() error
}

type Server struct {
	*http.ServeMux
	*sync.WaitGroup
	reader OrderParserReader
	filter OrderFilter
	config ServerConfig
	lis    net.Listener
}

func NewServer(config ServerConfig, filter OrderFilter) (s *Server) {
	if config == nil {
		panic("Server: passed nil config")
	}
	if filter == nil {
		panic("Server: passed nil filter")
	}

	s = &Server{
		http.NewServeMux(),
		&sync.WaitGroup{},
		NewOrderReader(),
		filter,
		config,
		nil,
	}

	s.HandleFunc(_PATH_TO_RSS, s.RSSHandler)
	s.HandleFunc(_PATH_TO_SHORT_LINKS, s.ShortLinkHandler)

	return s
}

func (s *Server) Start() (err error) {
	if s.lis != nil {
		return errors.New("Server is already running")
	}

	if s.config.GetPort() != _RSS_REQUIRED_PORT {
		log.Println("RSS protocol required port 80")
	}

	s.lis, err = net.Listen("tcp",
		s.config.GetHost()+":"+s.config.GetPort())

	if err != nil {
		log.Fatal("Server:", err)
	}

	return http.Serve(s.lis, s)
}

func (s *Server) ShutDown() error {
	if s.lis == nil {
		return errors.New("Server is already stopped")
	}

	s.Wait() // wait for all processed requests

	if s.lis == nil {
		return nil
	}

	defer func() {
		s.lis = nil
	}()
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
		log.Println("Loading error:", err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Println("Server return status " + resp.Status)
		} else {
			orders, err = s.reader.ReadOrders(resp)
			if err != nil && err != io.EOF {
				log.Println("Can't read or parse response: ", err)
			}
		}

		log.Printf("Loaded %d orders\n", len(orders))

		if len(orders) > 0 && s.config.IsFilterEnabled() {
			var filtered float32
			orders, filtered = s.filter.Execute(orders)
			log.Printf("%.1f%% of orders were removed by filter\n",
				filtered*100)
		}
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	render := NewRender(s.config)

	if URL, err := url.Parse(r.FormValue("url")); err == nil {
		// call feed like search request
		render.SetTitle(URL.Query().Get("searchString"))
	} else {
		log.Println("getting feed title error:", err)
	}

	if len(orders) > 0 {
		render.Compose(orders)
	}

	if err := render.WriteTo(w); err != nil {
		log.Println("can't send response:", err)
	}

	// signal that request was processed
	s.Done()
}

func (s *Server) ShortLinkHandler(w http.ResponseWriter,
	r *http.Request) {
	// redirect if order id was not passed also
	http.Redirect(w, r, MakeLink(r.FormValue("order")),
		http.StatusFound)
	r.Body.Close()
}
