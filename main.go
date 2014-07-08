package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

// port must be 80 becouse some rss clients require this
const _LOCAL_ADDR = "zakup-robot.ru:80"

var log *Log

func main() {
	var err error
	log, err = NewLog()
	if err != nil {
		fmt.Println("can't create log file:", err)
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
