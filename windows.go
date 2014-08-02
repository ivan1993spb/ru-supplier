package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"syscall"
)
import "github.com/lxn/walk"

const _PROGRAMM_ICON_FILE_NAME = "eagle.ico"

const (
	// require root privileges if true
	_ALLOW_REWRITE_HOSTS_FILE = false
	// false for debugging
	_ALLOW_FREE_CONSOLE = false
)

const (
	// _PROG_TITLE               = "Внимательный Поставщик"
	// _PROG_VERSION             = "2.0"
	// _NOTIFY_ICON_TOOL_TIP_MSG = _PROG_TITLE + " " + _PROG_VERSION

	_ACTION_TITLE_START_SERVER    = "Запустить"
	_ACTION_TITLE_STOP_SERVER     = "Остановить"
	_ACTION_TITLE_FILTER_ENABLED  = "Фильтровать"
	_ACTION_TITLE_FILTER_DISABLED = "Не фильтровать"
	_ACTION_TITLE_REMOVE_CACHE    = "Сбросить кэш"
	_ACTION_TITLE_OPEN_DIR        = "Открыть папку"
	// _ACTION_TITLE_OPEN_README  = "Открыть инструкцию"
	_ACTION_TITLE_EXIT = "Выход"

	_TOOL_TIP_MESSAGE_SERVER_RUNNING = "Работает"
	_TOOL_TIP_MESSAGE_SERVER_STOPPED = "Остановлен"
)

func FreeConsole() error {
	kernel32, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return err
	}
	freeConsole, err := syscall.GetProcAddress(
		kernel32,
		"FreeConsole",
	)
	if err != nil {
		return err
	}
	_, _, errnum := syscall.Syscall(uintptr(freeConsole), 0, 0, 0, 0)
	if errnum != 0 {
		return errors.New("syscall return error code")
	}
	return nil
}

func FileHostsPath() string {
	if root := os.Getenv("SystemRoot"); len(root) > 0 {
		return path.Join(root, "system32", "drivers", "etc", "hosts")
	}
	return ""
}

func CreateLocalHostIfNotExists(host string) error {
	if host = strings.TrimSpace(host); len(host) == 0 {
		panic("passed empty host")
	}
	if hostsPath := FileHostsPath(); len(hostsPath) > 0 {
		fhosts, err := os.OpenFile(
			hostsPath,
			os.O_RDWR|os.O_APPEND|os.O_CREATE,
			os.ModePerm,
		)
		if err != nil {
			return errors.New("cannot open hosts file: " +
				err.Error())
		}
		defer fhosts.Close()
		r := bufio.NewReader(fhosts)
		for {
			line, err := r.ReadString('\n')
			if err != nil && err != io.EOF {
				return errors.New("cannot read hosts file: " +
					err.Error())
			}
			if err == io.EOF && len(line) == 0 {
				break
			}
			// skip comments
			if comm := strings.IndexByte(line, '#'); comm > -1 {
				line = line[:comm]
			}
			if len(line) == 0 {
				continue
			}
			if fields := strings.Fields(line); len(fields) > 1 {
				for i := 1; i < len(fields); i++ {
					if strings.EqualFold(fields[i], host) {
						return nil
					}
				}
			}
			if err == io.EOF {
				break
			}
		}
		_, err = fhosts.WriteString("127.0.0.1\t" + host + "\n")
		return err
	}
	return ErrNotFoundHostFile
}

func RemoveLocalHostIfExists(host string) error {
	if host = strings.TrimSpace(host); len(host) == 0 {
		panic("passed empty host")
	}
	if hostsPath := FileHostsPath(); len(hostsPath) > 0 {
		fhosts, err := os.OpenFile(hostsPath, os.O_RDWR, os.ModePerm)
		if err != nil {
			return err
		}
		defer fhosts.Close()
		buff := bytes.NewBuffer(nil)
		rdr := bufio.NewReader(fhosts)
		exp := regexp.MustCompile("(?i)" + regexp.QuoteMeta(host))
		var line, comm string
		for {
			line, err = rdr.ReadString('\n')
			if err != nil && err != io.EOF {
				return err
			}
			if err == io.EOF && len(line) == 0 {
				break
			}
			if i := strings.IndexByte(line, '#'); i > -1 {
				line, comm = line[:i], line[i:]
			}
			if len(strings.TrimSpace(line)) == 0 {
				buff.WriteString(line + comm)
				continue
			}
			line = exp.ReplaceAllString(line, "")
			if len(strings.Fields(line)) > 1 {
				buff.WriteString(line + comm)
			} else {
				buff.WriteString(comm)
			}
			if err == io.EOF {
				break
			}
		}
		if _, err = fhosts.Seek(0, 0); err != nil {
			return err
		}
		if n, err := buff.WriteTo(fhosts); err != nil {
			return err
		} else if err = fhosts.Truncate(n); err != nil {
			return err
		}
		return nil
	}
	return ErrNotFoundHostFile
}

func InterfaceStart(server *Server, config *Config) error {
	if server == nil {
		panic("interface error: passed nil server")
	}
	if config == nil {
		panic("interface error: passed nil config")
	}

	var err error

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                        INITIALIZATION                       *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	if _ALLOW_FREE_CONSOLE {
		// free console
		if err = FreeConsole(); err != nil {
			log.Println("cannot free console:", err)
		}
	}
	if !server.IsRunning() {
		if _ALLOW_REWRITE_HOSTS_FILE {
			// edit hosts file
			err = CreateLocalHostIfNotExists(config.Host)
			if err != nil {
				log.Println("cannot create local addr:", err)
			}
		}
		// start server
		if err = server.Start(); err != nil {
			return err
		}
	}
	defer func() {
		if server.IsRunning() {
			// shut down server
			if err = server.ShutDown(); err != nil {
				log.Println("cannot shut down server", err)
			}
		}
		if _ALLOW_REWRITE_HOSTS_FILE {
			// remove local host from hosts file
			err = RemoveLocalHostIfExists(config.Host)
			if err != nil {
				log.Println("cannot remove local host:", err)
			}
		}
	}()

	/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	 *                      END INITIALIZATION                     *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

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
	 *                            ACTIONS                            *
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

	// // // // // // // // // // // // // // // // // // // //
	// // // // // // // // // // // // // // // // // // // //
	// // // // // // // // // // // // // // // // // // // //
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
	// // // // // // // // // // // // // // // // // // // //
	// // // // // // // // // // // // // // // // // // // //
	// // // // // // // // // // // // // // // // // // // //

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

	// // // // // // // // // // // // // // // // // // // //
	// // // // // // // // // // // // // // // // // // // //
	// // // // // // // // // // // // // // // // // // // //
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
	// // // // // // // // // // // // // // // // // // // //
	// // // // // // // // // // // // // // // // // // // //
	// // // // // // // // // // // // // // // // // // // //

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
	 *                       END EVENT HANDLERS                      *
	 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

	mw.Run()
	return nil
}
