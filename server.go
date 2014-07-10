package main

import (
	"bufio"
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
	} else if resp, err := load(r.Form.Get("url")); err != nil {
		log.Warning.Println("loading error:", err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			rdr := bufio.NewReaderSize(resp.Body, 370)
			// skip first line
			rdr.ReadString('\n')
			line, err := rdr.ReadBytes('\n')
			chunk, ok := hashstore.GetHashChunk(resp.Request.URL.String())
			if ok {
			}
		} else {
			log.Warning.Printf("server return %q\n",
				strings.ToLower(resp.Status))
		}
	}

	// 	 else {
	// 		orders, err := parse(resp.Body)
	// 		if err != nil {
	// 			log.Warning.Println("can't read or parse response:",
	// 				err)
	// 		}
	// 		if len(orders) > 0 {
	// 			orders = filter(orders)
	// 		}
	// 		w.Header().Set("Content-Type",
	// 			"application/xml; charset=utf-8")
	// 		w.WriteHeader(http.StatusOK)
	// 		w.Write([]byte(xml.Header))
	// 		err = xml.NewEncoder(w).Encode(
	// 			OrdersToRssFeed(
	// 				orders,
	// 				// search string as title (if exists)
	// 				URL.Query().Get("searchString"),
	// 			).FeedXml(),
	// 		)
	// 		if err != nil {
	// 			log.Error.Println("can't send response:", err)
	// 		}
	// 	}
	// 	return
	// }
	// w.WriteHeader(http.StatusBadRequest)
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
