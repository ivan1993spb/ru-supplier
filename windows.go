package main

import ()
import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	_PROGRAM_TITLE = "Внимательный Поставшик"
)

func StartInterface() {
	var mw *walk.MainWindow
	if _, err := (MainWindow{
		AssignTo: &mw,
		Title:    "Walk Data Binding Example",
		MinSize:  Size{300, 200},
		Layout:   VBox{},
		Children: []Widget{
			PushButton{
				Text: "Edit Animal",
				OnClicked: func() {

				},
			},
		},
	}.Run()); err != nil {
		log.Error.Fatal(err)
	}
}
