// +build windows

package main

import (
	"net/url"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	_WIN_TITLE       = "Генератор ссылок"
	_WIN_GEN_BUTTON  = "Генерировать"
	_WIN_COPY_BUTTON = "Копировать"
	_WIN_LABEL_HOST  = "Локальный хост"
	_WIN_LABEL_LINK  = "Ссылка на страницу с закупками"
)

func main() {
	var (
		te *walk.TextEdit
		le *walk.LineEdit
	)

	MainWindow{
		Title:   _WIN_TITLE,
		MinSize: Size{300, 400},
		Layout:  VBox{},
		Children: []Widget{
			Label{
				Text: _WIN_LABEL_HOST,
			},
			LineEdit{
				AssignTo: &le,
				Text:     GetHttpHost(),
			},
			Label{
				Text: _WIN_LABEL_LINK,
			},
			TextEdit{
				AssignTo: &te,
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: _WIN_GEN_BUTTON,
						OnClicked: func() {
							URL, err := url.Parse(te.Text())
							if err != nil {
								return
							}

							genURL := generateURL(URL, le.Text())
							if genURL != nil {
								te.SetText(genURL.String())
							}
						},
					},
					PushButton{
						Text: _WIN_COPY_BUTTON,
						OnClicked: func() {
							walk.Clipboard().SetText(te.Text())
						},
					},
				},
			},
		},
	}.Run()
}
