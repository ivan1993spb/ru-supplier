package main

import (
	logpkg "log"
	"os"
)

type Log struct {
	Error, Warning *logpkg.Logger
}

func NewLog(efname, wfname string) (*Log, error) {
	var err error
	var efile *os.File
	efile, err = os.OpenFile(
		efname,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		os.ModePerm,
	)
	if err == nil {
		var wfile *os.File
		wfile, err = os.OpenFile(
			wfname,
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
