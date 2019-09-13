package ui

import (
	"runtime"

	rice "github.com/GeertJohan/go.rice"
	"github.com/dags-/systray"
)

func getIcon(box *rice.Box) []byte {
	if runtime.GOOS == "windows" {
		return box.MustBytes("assets/image/tray.ico")
	} else {
		return box.MustBytes("assets/image/tray.png")
	}
}

func (m *Manager) handleTray() {
	systray.Run(m.trayReady, m.trayExit)
}

func (m *Manager) trayReady() {
	systray.SetTitle("Launcher")
	systray.SetIcon(m.icon)
	c := Load(m.appDir)
	auto := systray.AddMenuItem("Auto Launch", "")
	exit := systray.AddMenuItem("Quit", "")

	if c.AutoLaunch {
		auto.Check()
	} else {
		auto.Uncheck()
	}

	for {
		select {
		case <-systray.ClickedCh:
			m.NewSizedWindow(Settings{
				Path:      "/home",
				Resizable: true,
			})
			systray.Quit()
			return
		case <-auto.ClickedCh:
			if auto.Checked() {
				auto.Uncheck()
			} else {
				auto.Check()
			}
			c.AutoLaunch = auto.Checked()
			Save(m.appDir, c)
			break
		case <-exit.ClickedCh:
			systray.Quit()
			go m.Exit()
			return
		case <-m.open:
			m.NewSizedWindow(Settings{
				Path:      "/home",
				Resizable: true,
			})
			systray.Quit()
			return
		}
	}
}

func (m *Manager) trayExit() {}
