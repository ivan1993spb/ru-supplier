package main

import (
	"log"
	"os"
)

const (
	_LOG_FILE_NAME     = "prog.log"
	_CONFIG_FILE_NAME  = "config.json"
	_FILTERS_FILE_NAME = "filters.json"
)

func main() {
	if logfile, err := os.OpenFile(
		_LOG_FILE_NAME,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		os.ModePerm,
	); err != nil {
		log.Fatal(err)
	} else {
		log.SetFlags(log.LstdFlags)
		log.SetOutput(logfile)
	}
	config, err := LoadConfig(_CONFIG_FILE_NAME)
	if config == nil {
		if err != nil {
			log.Fatal("cannot load configs:", err)
		} else {
			panic("config object is nil")
		}
	}
	if err != nil {
		log.Println("config:", err)
	}
	filter, err := LoadFilter(_FILTERS_FILE_NAME)
	if filter == nil {
		if err != nil {
			log.Fatal("cannot load filters:", err)
		} else {
			panic("filter object is nil")
		}
	}
	if err != nil {
		log.Println("filter:", err)
	}
	if err = InterfaceStart(
		NewServer(config, filter),
		config,
	); err != nil {
		log.Fatal("interface fatal error:", err)
	}
}
