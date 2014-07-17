package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	_PROGRAM_TITLE = "Внимательный Поставшик"
)

type ServerMainWindow struct {
	*walk.MainWindow
}

func StartInterface() {
	mw := new(ServerMainWindow)
	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    _PROGRAM_TITLE,
		MinSize:  Size{100, 100},
		Size:     Size{120, 120},
		MaxSize:  Size{150, 150},
		Layout:   VBox{},
	}.Run()); err != nil {
		log.Error.Fatal("interface error: ", err)
	}
}
