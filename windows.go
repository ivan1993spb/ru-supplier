package main

import (
	"github.com/lxn/walk"
	// . "github.com/lxn/walk/declarative"
)

const _ICON_FILE_NAME = "eagle.ico"

func InterfaceStart(server *Server, config *Config, filter *Filter) error {
	if server == nil {
		panic("interface error: passed nil server")
	}
	if config == nil {
		panic("interface error: passed nil config")
	}
	if filter == nil {
		panic("interface error: passed nil filter")
	}
	mw, err := walk.NewMainWindow()
	if err != nil {
		return err
	}
	icon, err := walk.NewIconFromFile(_ICON_FILE_NAME)
	if err != nil {
		return err
	}
	defer icon.Dispose()
	ni, err := walk.NewNotifyIcon()
	if err != nil {
		return err
	}
	defer ni.Dispose()
	if err := ni.SetIcon(icon); err != nil {
		return err
	}
	if err := ni.SetToolTip(_NOTIFY_ICON_TOOL_TIP_MSG); err != nil {
		return err
	}
	context_menu := map[string]walk.EventHandler{
		"Выход":   func() { walk.App().Exit(0) },
		"Прокси":  func() {},
		"Адреса":  func() {},
		"Фильтры": func() {},
	}
	for text, handler := range context_menu {
		action := walk.NewAction()
		if err := action.SetText(text); err != nil {
			return err
		}
		action.Triggered().Attach(handler)
		if err := ni.ContextMenu().Actions().Add(action); err != nil {
			return err
		}
	}

	if err := ni.SetVisible(true); err != nil {
		return err
	}
	if err := ni.ShowInfo("Walk NotifyIcon Example", "Click the icon to show again."); err != nil {
		return err
	}
	// Run the message loop.
	mw.Run()
	return nil
}
