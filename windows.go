package main

import (
	"log"
	"os/exec"
)
import "github.com/lxn/walk"

const _PROGRAMM_ICON_FILE_NAME = "eagle.ico"

func InterfaceStart(server *Server, config *Config) error {
	if server == nil {
		panic("interface error: passed nil server")
	}
	if config == nil {
		panic("interface error: passed nil config")
	}
	mw, err := walk.NewMainWindow()
	if err != nil {
		return err
	}
	ni, err := walk.NewNotifyIcon()
	if err != nil {
		return err
	}
	defer ni.Dispose()
	if err := ni.SetVisible(true); err != nil {
		return err
	}
	// create image icon
	if icon, err := walk.NewIconFromFile(_PROGRAMM_ICON_FILE_NAME); err != nil {
		log.Println("cannot load icon from file:", err)
	} else {
		defer icon.Dispose()
		if err = ni.SetIcon(icon); err != nil {
			log.Println("cannot bind image icon with notify icon:", err)
		}
	}

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                            ACTIONS                      *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	startServerAction := walk.NewAction()
	if err = startServerAction.SetText(_ACTION_TITLE_START_SERVER); err != nil {
		return err
	}
	if err = ni.ContextMenu().Actions().Add(startServerAction); err != nil {
		return err
	}

	stopServerAction := walk.NewAction()
	if err = stopServerAction.SetText(_ACTION_TITLE_STOP_SERVER); err != nil {
		return err
	}
	if err = ni.ContextMenu().Actions().Add(stopServerAction); err != nil {
		return err
	}

	// hide start button if server is running else hide stop button
	updateServerButtons := func() {
		if server.IsRunning() {
			if err = ni.SetToolTip(_TOOL_TIP_MESSAGE_SERVER_RUNNING); err != nil {
				log.Println(err)
			}
			if err = startServerAction.SetVisible(false); err != nil {
				log.Println(err)
			}
			if err = stopServerAction.SetVisible(true); err != nil {
				log.Println(err)
			}
		} else {
			if err = ni.SetToolTip(_TOOL_TIP_MESSAGE_SERVER_STOPPED); err != nil {
				log.Println(err)
			}
			if err = stopServerAction.SetVisible(false); err != nil {
				log.Println(err)
			}
			if err = startServerAction.SetVisible(true); err != nil {
				log.Println(err)
			}
		}
	}
	updateServerButtons()

	filterEnableAction := walk.NewAction()
	if err = filterEnableAction.SetText(_ACTION_TITLE_FILTER_ENABLED); err != nil {
		return err
	}
	if err = ni.ContextMenu().Actions().Add(filterEnableAction); err != nil {
		return err
	}

	filterDisabledAction := walk.NewAction()
	if err = filterDisabledAction.SetText(_ACTION_TITLE_FILTER_ENABLED); err != nil {
		return err
	}
	if err = ni.ContextMenu().Actions().Add(filterDisabledAction); err != nil {
		return err
	}

	updateFilterButtons := func() {
		if config.FilterEnabled {
			if err = filterEnableAction.SetVisible(false); err != nil {
				log.Println(err)
			}
			if err = filterDisabledAction.SetVisible(true); err != nil {
				log.Println(err)
			}
		} else {
			if err = filterDisabledAction.SetVisible(false); err != nil {
				log.Println(err)
			}
			if err = filterEnableAction.SetVisible(true); err != nil {
				log.Println(err)
			}
		}
	}
	updateFilterButtons()

	removeCacheAction := walk.NewAction()
	if err = removeCacheAction.SetText(_ACTION_TITLE_REMOVE_CACHE); err != nil {
		return err
	}
	if err = ni.ContextMenu().Actions().Add(removeCacheAction); err != nil {
		return err
	}

	ni.ContextMenu().Actions().Add(walk.NewSeparatorAction())

	openDirAction := walk.NewAction()
	if err = openDirAction.SetText(_ACTION_TITLE_OPEN_DIR); err != nil {
		return err
	}
	if err = ni.ContextMenu().Actions().Add(openDirAction); err != nil {
		return err
	}

	ni.ContextMenu().Actions().Add(walk.NewSeparatorAction())

	exitAction := walk.NewAction()
	if err = exitAction.SetText(_ACTION_TITLE_EXIT); err != nil {
		return err
	}
	if err = ni.ContextMenu().Actions().Add(exitAction); err != nil {
		return err
	}

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                        EVENT HANDLERS                         *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */
	startServer := func() {
		if !server.IsRunning() {
			if err = server.Start(); err != nil {
				log.Println("cannot start server:", err)
			}
		}
	}
	stopServer := func() {
		if server.IsRunning() {
			if err = server.ShutDown(); err != nil {
				log.Println("cannot shut down server:", err)
			}
		}
	}
	startServerAction.Triggered().Attach(startServer)
	startServerAction.Triggered().Attach(updateServerButtons)

	stopServerAction.Triggered().Attach(stopServer)
	stopServerAction.Triggered().Attach(updateServerButtons)

	filterEnableAction.Triggered().Attach(func() {
		config.SetFilterEnabled(true)
		config.Save()
	})
	filterEnableAction.Triggered().Attach(updateFilterButtons)

	filterDisabledAction.Triggered().Attach(func() {
		config.SetFilterEnabled(false)
		config.Save()
	})
	filterDisabledAction.Triggered().Attach(updateFilterButtons)

	removeCacheAction.Triggered().Attach(func() {
		if err = server.RemoveCache(); err != nil {
			log.Println("cannot remove cache:", err)
		}
	})

	openDirAction.Triggered().Attach(func() {
		if err = exec.Command("cmd", "/C", "start", ".").Start(); err != nil {
			log.Println("cannot open program directory:", err)
		}
	})

	exitAction.Triggered().Attach(func() {
		stopServer()
		if err = config.Save(); err != nil {
			log.Println("cannot save configures:", err)
		}
		walk.App().Exit(0)
	})

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                             END                               *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	mw.Run()
	return nil
}
