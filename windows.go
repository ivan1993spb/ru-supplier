package main

import (
	. "github.com/lxn/walk/declarative"
)

func InterfaceStart(server *Server, config *Config, filter *Filter) error {
	if server == nil {
		panic("interface error: passed nil server")
	}
	if config == nil {
		panic("interface error: passed nil config")
	}
	if filter == nil {
		panic("interface error: passed nil filter")
	}

	// ...

}
