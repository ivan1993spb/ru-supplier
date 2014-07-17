package main

import (
	"fmt"
	"os"
)

const (
	_HASH_STORE_FILE_NAME  = "cache.json"
	_FILTERS_FILE_NAME     = "filters.json"
	_CONFIG_FILE_NAME      = "config.json"
	_LOG_ERROR_FILE_NAME   = "error.log"
	_LOG_WARNING_FILE_NAME = "warning.log"
)

var (
	log       *Log
	hashstore *HashStore
	config    *Config
	filter    *Filter
	server    *Server
)

func main() {
	var err error
	log, err = NewLog(_LOG_ERROR_FILE_NAME, _LOG_WARNING_FILE_NAME)
	if err != nil {
		fmt.Println("can't create log file:", err)
		os.Exit(1)
	}
	hashstore, err = LoadHashStore(_HASH_STORE_FILE_NAME)
	if err != nil {
		fmt.Println("can't load hash store:", err)
		log.Error.Println("can't load hash store:", err)
		err = nil
	}
	config, err = LoadConfig(_CONFIG_FILE_NAME)
	if err != nil {
		fmt.Println("can't load configs:", err, "; will used default")
		log.Error.Println("can't load configs:", err,
			"; will used default")
		err = nil
	}
	filter, err = LoadFilter(_FILTERS_FILE_NAME)
	if err != nil {
		fmt.Println("can't load filters:", err)
		log.Error.Println("can't load filters:", err)
		config.FilterEnabled = false
		err = nil
	}
	server = NewServer()
	StartInterface()
}
