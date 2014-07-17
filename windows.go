package main

import ()
import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	_PROGRAM_TITLE             = "Внимательный Поставшик"
	_TAB_TITLE_SERVER_SETTINGS = "Прокси"
	_TAB_TITLE_LINKS           = "Адреса"
)

func StartInterface() {
	var mw *walk.MainWindow
	if _, err := (MainWindow{
		AssignTo: &mw,
		Title:    _PROGRAM_TITLE,
		Size:     Size{200, 200},
		Layout:   VBox{},
		Children: []Widget{
			TabWidget{
				Pages: []TabPage{
					{
						Title: _TAB_TITLE_SERVER_SETTINGS,
					},
					{
						Title: _TAB_TITLE_LINKS,
						Children: []Widget{
							TextEdit{
								ToolTipText: "assasa",
								Row:         122,
								Column:      122,
								Size:        Size{110, 120},
							},
						},
					},
				},
			},
		},
	}.Run()); err != nil {
		log.Error.Fatal(err)
	}
}
