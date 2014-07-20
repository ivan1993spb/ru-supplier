package main

import (
	"log"
	"os"
)

const (
	_LOG_FILE_NAME        = "errors.log"
	_HASH_STORE_FILE_NAME = "cache.json"
	_CONFIG_FILE_NAME     = "config.json"
	_FILTERS_FILE_NAME    = "filters.json"
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
	hashstore, err := LoadHashStore(_HASH_STORE_FILE_NAME)
	if err != nil && hashstore == nil {
		log.Fatal("cannot load hashstore:", err)
	} else if err != nil {
		log.Println("hashstore:", err)
	}
	config, err := LoadConfig(_CONFIG_FILE_NAME)
	if err != nil && config == nil {
		log.Fatal("cannot load configs:", err)
	} else if err != nil {
		log.Println("configs:", err)
	}
	filter, err := LoadFilter(_FILTERS_FILE_NAME)
	if err != nil && filter == nil {
		log.Fatal("cannot load filters:", err)
	} else if err != nil {
		log.Println("filters:", err)
	}
	NewServer(config, filter, hashstore)
	// StartInterface(server, config, filter, hashstore)
}
