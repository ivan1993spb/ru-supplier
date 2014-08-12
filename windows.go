package main

import (
	"log"
	"os/exec"
	"time"

	"github.com/lxn/walk"
)

const (
	// true to allow run server on application start
	_RUN_SERVER_ON_STARTING = true
	// time for which proxy must start
	_START_SERVER_TIMEOUT = time.Second * 2
)

const (
	_PROG_ICON_FILE_NAME        = "eagle.ico"
	_PROG_DESCRIPTION_FILE_NAME = "docs/index.html"
)

const (
	_PROG_NAME    = "Внимательный Поставщик"
	_PROG_VERSION = "2.0"
	_PROG_TITLE   = _PROG_NAME + " " + _PROG_VERSION
)

const (
	_ACTION_TITLE_START_SERVER    = "Запустить прокси"
	_ACTION_TITLE_STOP_SERVER     = "Остановить прокси"
	_ACTION_TITLE_FILTER_ENABLED  = "Включить фильтр"
	_ACTION_TITLE_FILTER_DISABLED = "Выключить фильтр"
	_ACTION_TITLE_REMOVE_CACHE    = "Сбросить кэш"
	_ACTION_TITLE_OPEN_URL_GEN    = "Генератор ссылок"
	_ACTION_TITLE_OPEN_DIR        = "Папка настроек"
	_ACTION_TITLE_OPEN_README     = "Инструкция"
	_ACTION_TITLE_EXIT            = "Выход"
)

const (
	_NOTICE_APP_START     = "Здравствуйте, программа запущена"
	_NOTICE_CACHE_REMOVED = "Кэш удален, программа некоторое время " +
		"будет загружать и обрабатывать все доступные закупки"
	_NOTICE_CONFIGS = "Все настройки программы и фильтры храняться " +
		"в файлах с расширением .json"
	_NOTICE_ENABLED_FILTERS  = "Фильтр включен"
	_NOTICE_DISABLED_FILTERS = "Фильтр выключен"
	_NOTICE_PROXY_ENABLED    = "Локальный прокси запущен"
	_NOTICE_PROXY_DISABLED   = "Локальный прокси остановлен"
)

func InterfaceStart(server *Server, config *Config) (err error) {
	if server == nil {
		panic("interface error: passed nil server")
	}
	if config == nil {
		panic("interface error: passed nil config")
	}

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                        INITIALIZATION                       *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	startServer := func() {
		if !server.IsRunning() {
			go func() {
				// start server
				if err := server.Start(); err != nil {
					log.Println("cannot start server:", err)
				}
			}()
			time.Sleep(_START_SERVER_TIMEOUT)
		}
	}
	stopServer := func() {
		if server.IsRunning() {
			// shut down server
			if err = server.ShutDown(); err != nil {
				log.Println("cannot shut down server", err)
			}
		}
	}
	if _RUN_SERVER_ON_STARTING && !server.IsRunning() {
		startServer()
	}
	defer func() {
		if server.IsRunning() {
			stopServer()
		}
	}()
	defer func() {
		if err = config.Save(); err != nil {
			log.Println("cannot save configures:", err)
		}
	}()

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                      END INITIALIZATION                     *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	mw, err := walk.NewMainWindow()
	if err != nil {
		return
	}
	defer mw.Dispose()
	ni, err := walk.NewNotifyIcon()
	if err != nil {
		return
	}
	defer ni.Dispose()
	if err = ni.SetVisible(true); err != nil {
		return
	}
	if err = ni.SetToolTip(_PROG_TITLE); err != nil {
		return
	}
	ni.ShowMessage(_PROG_TITLE, _NOTICE_APP_START)
	// create image icon
	if icon, err :=
		walk.NewIconFromFile(_PROG_ICON_FILE_NAME); err != nil {
		log.Println("cannot load icon from file:", err)
	} else {
		defer icon.Dispose()
		if err = ni.SetIcon(icon); err != nil {
			log.Println("cannot bind img with notify icon:", err)
		}
	}

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                            ACTIONS                          *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	startServerAction := walk.NewAction()
	err = startServerAction.SetText(_ACTION_TITLE_START_SERVER)
	if err != nil {
		return
	}
	err = startServerAction.SetVisible(!server.IsRunning())
	if err != nil {
		return
	}
	err = ni.ContextMenu().Actions().Add(startServerAction)
	if err != nil {
		return
	}

	stopServerAction := walk.NewAction()
	err = stopServerAction.SetText(_ACTION_TITLE_STOP_SERVER)
	if err != nil {
		return
	}
	err = stopServerAction.SetVisible(server.IsRunning())
	if err != nil {
		return
	}
	err = ni.ContextMenu().Actions().Add(stopServerAction)
	if err != nil {
		return
	}

	filterEnableAction := walk.NewAction()
	err = filterEnableAction.SetText(_ACTION_TITLE_FILTER_ENABLED)
	if err != nil {
		return
	}
	err = filterEnableAction.SetVisible(!config.FilterEnabled)
	if err != nil {
		return
	}
	err = ni.ContextMenu().Actions().Add(filterEnableAction)
	if err != nil {
		return
	}

	filterDisabledAction := walk.NewAction()
	err = filterDisabledAction.SetText(_ACTION_TITLE_FILTER_DISABLED)
	if err != nil {
		return
	}
	err = filterDisabledAction.SetVisible(config.FilterEnabled)
	if err != nil {
		return
	}
	err = ni.ContextMenu().Actions().Add(filterDisabledAction)
	if err != nil {
		return
	}

	removeCacheAction := walk.NewAction()
	err = removeCacheAction.SetText(_ACTION_TITLE_REMOVE_CACHE)
	if err != nil {
		return
	}
	err = ni.ContextMenu().Actions().Add(removeCacheAction)
	if err != nil {
		return
	}

	err = ni.ContextMenu().Actions().Add(walk.NewSeparatorAction())
	if err != nil {
		return
	}

	openURLGenAction := walk.NewAction()
	err = openURLGenAction.SetText(_ACTION_TITLE_OPEN_URL_GEN)
	if err != nil {
		return
	}
	err = ni.ContextMenu().Actions().Add(openURLGenAction)
	if err != nil {
		return
	}

	openDirAction := walk.NewAction()
	err = openDirAction.SetText(_ACTION_TITLE_OPEN_DIR)
	if err != nil {
		return
	}
	err = ni.ContextMenu().Actions().Add(openDirAction)
	if err != nil {
		return
	}

	openReadMeAction := walk.NewAction()
	err = openReadMeAction.SetText(_ACTION_TITLE_OPEN_README)
	if err != nil {
		return
	}
	err = ni.ContextMenu().Actions().Add(openReadMeAction)
	if err != nil {
		return
	}

	err = ni.ContextMenu().Actions().Add(walk.NewSeparatorAction())
	if err != nil {
		return
	}

	exitAction := walk.NewAction()
	err = exitAction.SetText(_ACTION_TITLE_EXIT)
	if err != nil {
		return
	}
	err = ni.ContextMenu().Actions().Add(exitAction)
	if err != nil {
		return
	}

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                        EVENT HANDLERS                       *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	updateServerButtons := func() {
		if server.IsRunning() {
			if err = startServerAction.SetVisible(false); err != nil {
				log.Println(err)
			}
			if err = stopServerAction.SetVisible(true); err != nil {
				log.Println(err)
			}
		} else {
			if err = stopServerAction.SetVisible(false); err != nil {
				log.Println(err)
			}
			if err = startServerAction.SetVisible(true); err != nil {
				log.Println(err)
			}
		}
	}

	startServerAction.Triggered().Attach(func() {
		if !server.IsRunning() {
			if err = startServerAction.SetEnabled(false); err != nil {
				log.Println(err)
			}
			startServer()
			if server.IsRunning() {
				ni.ShowMessage(_PROG_TITLE, _NOTICE_PROXY_ENABLED)
			}
			if err = startServerAction.SetEnabled(true); err != nil {
				log.Println(err)
			}
		}
		updateServerButtons()
	})

	stopServerAction.Triggered().Attach(func() {
		if server.IsRunning() {
			stopServer()
			ni.ShowMessage(_PROG_TITLE, _NOTICE_PROXY_DISABLED)
		}
		updateServerButtons()
	})

	updateFilterButtons := func() {
		if config.FilterEnabled {
			err = filterEnableAction.SetVisible(false)
			if err != nil {
				log.Println(err)
			}
			err = filterDisabledAction.SetVisible(true)
			if err != nil {
				log.Println(err)
			}
		} else {
			err = filterDisabledAction.SetVisible(false)
			if err != nil {
				log.Println(err)
			}
			err = filterEnableAction.SetVisible(true)
			if err != nil {
				log.Println(err)
			}
		}
	}

	filterEnableAction.Triggered().Attach(func() {
		if !config.FilterEnabled {
			config.SetFilterEnabled(true)
			ni.ShowInfo(_PROG_TITLE, _NOTICE_ENABLED_FILTERS)
		}
		updateFilterButtons()
	})

	filterDisabledAction.Triggered().Attach(func() {
		if config.FilterEnabled {
			config.SetFilterEnabled(false)
			ni.ShowInfo(_PROG_TITLE, _NOTICE_DISABLED_FILTERS)
		}
		updateFilterButtons()
	})

	removeCacheAction.Triggered().Attach(func() {
		if err = server.RemoveCache(); err != nil {
			log.Println("cannot remove cache:", err)
		} else {
			ni.ShowInfo(_PROG_TITLE, _NOTICE_CACHE_REMOVED)
		}
	})

	openURLGenAction.Triggered().Attach(func() {
		// err = exec.Command(
		// 	"cmd", "/C", "start", _PROG_DESCRIPTION_FILE_NAME,
		// ).Start()
		// if err != nil {
		// 	log.Println("cannot open url generator:", err)
		// }
	})

	openDirAction.Triggered().Attach(func() {
		err = exec.Command("cmd", "/C", "start", ".").Start()
		if err != nil {
			log.Println("cannot open program directory:", err)
		} else {
			ni.ShowInfo(_PROG_TITLE, _NOTICE_CONFIGS)
		}
	})

	openReadMeAction.Triggered().Attach(func() {
		err = exec.Command(
			"cmd", "/C", "start", _PROG_DESCRIPTION_FILE_NAME,
		).Start()
		if err != nil {
			log.Println("cannot open README:", err)
		}
	})

	exitAction.Triggered().Attach(func() {
		walk.App().Exit(0)
	})

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                       END EVENT HANDLERS                    *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	mw.Run()
	return
}
