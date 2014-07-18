package main

import (
	. "github.com/lxn/walk/declarative"
)

func StartInterface() {
	if _, err := (MainWindow{
		Title:  _PROGRAM_TITLE,
		Size:   Size{200, 200},
		Layout: VBox{},
		Children: []Widget{
			TabWidget{
				Pages: []TabPage{
					{
						Title:       _TAB_TITLE_SERVER_SETTINGS,
						ToolTipText: _TAB_TOOL_TIP_TEXT_SERVER_SETTINGS,
						Layout:      VBox{},
						Children: []Widget{
							HSplitter{
								Children: []Widget{},
							},
						},
					},
					{
						Title:       _TAB_TITLE_LINKS,
						ToolTipText: _TAB_TOOL_TIP_TEXT_LINKS,
						Layout:      VBox{},
						Children:    []Widget{},
					},
				},
			},
		},
	}.Run()); err != nil {
		log.Error.Fatal(err)
	}
}
