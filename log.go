package main

import (
	logpkg "log"
	"os"
)

const (
	_ERROR_LOG_FILE_NAME   = "error.zakup.log"
	_WARNING_LOG_FILE_NAME = "warning.zakup.log"
)

type Log struct {
	Error, Warning *logpkg.Logger
}

func NewLog() (*Log, error) {
	var err error
	var efile *os.File
	efile, err = os.OpenFile(
		_ERROR_LOG_FILE_NAME,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		os.ModePerm,
	)
	if err == nil {
		var wfile *os.File
		wfile, err = os.OpenFile(
			_WARNING_LOG_FILE_NAME,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			os.ModePerm,
		)
		if err == nil {
			return &Log{
				logpkg.New(efile, "", logpkg.LstdFlags),
				logpkg.New(wfile, "", logpkg.LstdFlags),
			}, nil
		}
	}
	return nil, err
}
