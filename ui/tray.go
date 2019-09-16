package ui

import (
	"io/ioutil"
	"log"

	"github.com/dags-/systray"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
)

func (m *Manager) RunTray() {
	systray.Run(m.ready, m.exit)
}

func (m *Manager) ready() {
	icon, e := ioutil.ReadFile(m.tray)

	if e == nil {
		systray.SetIcon(icon)
	} else {
		errs.Log("Load icon", e)
	}

	auto := systray.AddMenuItem("Auto Launch", "")
	open := systray.AddMenuItem("Open", "")
	exit := systray.AddMenuItem("Quit", "")

	if Load(m.appDir).AutoLaunch {
		auto.Check()
	}

	for {
		select {
		case <-systray.ClickedCh:
			if !m.HasWindow() {
				m.Home()
			}
			break
		case <-open.ClickedCh:
			if !m.HasWindow() {
				m.Home()
			}
			break
		case <-auto.ClickedCh:
			if auto.Checked() {
				auto.Uncheck()
			} else {
				auto.Check()
			}
			config := Load(m.appDir)
			config.AutoLaunch = auto.Checked()
			Save(m.appDir, config)
			break
		case <-exit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func (m *Manager) exit() {
	log.Println("Stopping")
	defer log.Println("Stopped")
	m.Exit()
}
