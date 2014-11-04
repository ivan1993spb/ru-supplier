// +build linux darwin

package main

import (
// "fmt"
// "time"

// "github.com/salviati/go-qt5/qt5"
)

func InterfaceStart(server ZakupkiProxyServer,
	config ServerConfig) (err error) {

	if server == nil {
		panic("interface error: passed nil server")
	}
	if config == nil {
		panic("interface error: passed nil config")
	}

	// qt5.Main(func() {
	// 	// w := qt5.NewWidget()
	// 	// w.SetWindowTitle(qt5.Version())
	// 	// w.SetSizev(300, 200)
	// 	// defer w.Close()

	// 	icon := qt5.NewIcon()
	// 	icon.Init()
	// 	// icon.InitWithFile("/home/ivan/gocode/src/bitbucket.org/pushkin_ivan/ru-supplier/src/eagle.ico")
	// 	defer icon.Close()

	// 	tray := qt5.NewSystemTray()
	// 	tray.Init()
	// 	defer tray.Close()

	// 	tray.SetIcon(icon)
	// 	time.Sleep(time.Second * 10)
	// 	tray.SetVisible(true)
	// 	time.Sleep(time.Second * 10)
	// 	tray.SetVisible(false)

	// 	fmt.Println(tray.IsVisible())
	// 	// tray.

	// 	// w.Show()
	// 	qt5.Run()
	// })

	// icon := qt5.NewIconWithFile("src/eagle.ico")
	// // icon.Init()
	// tray := qt5.NewSystemTray()
	// tray.SetIcon(icon)
	// // tray.Init()

	// time.Sleep(time.Minute)

	// tray.Close()
	// icon.Close()
	return
}
