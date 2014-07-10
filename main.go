package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

// port must be 80 becouse some rss clients require this
const _LOCAL_ADDR = "zakup-robot.ru:80"

const (
	_HASH_STORE_FILE       = "cache.json"
	_LOG_ERROR_FILE_NAME   = "error.zakup.log"
	_LOG_WARNING_FILE_NAME = "warning.zakup.log"
)

var (
	log       *Log
	hashstore *HashStore
)

func main() {
	var err error
	log, err = NewLog(_LOG_ERROR_FILE_NAME, _LOG_WARNING_FILE_NAME)
	if err != nil {
		fmt.Println("can't create log file:", err)
		os.Exit(1)
	}
	hashstore, err = LoadHashStore(_HASH_STORE_FILE)
	if err != nil {
		fmt.Println(err)
		log.Error.Println(err)
		os.Exit(1)
	}
	lis, err := net.Listen("tcp", _LOCAL_ADDR)
	if err != nil {
		fmt.Println(err)
		log.Error.Println(err)
		os.Exit(1)
	}
	err = http.Serve(NewServer(lis).Bind())
	if err != nil {
		fmt.Println(err)
		log.Error.Println(err)
		os.Exit(1)
	}
}
