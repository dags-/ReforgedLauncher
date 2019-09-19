package ui

import (
	"io/ioutil"
	"log"

	"github.com/dags-/systray"

	"github.com/Conquest-Reforged/ReforgedLauncher/metainfo"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
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

	info := systray.AddMenuItem(metainfo.NAME+" v"+metainfo.VERSION, "")
	info.Disable()
	systray.AddSeparator()
	auto := systray.AddMenuItem("Quick Launch", "")
	systray.AddSeparator()
	launch := systray.AddMenuItem("Launch", "")
	open := systray.AddMenuItem("Open", "")
	systray.AddSeparator()
	exit := systray.AddMenuItem("Quit", "")

	if Load(m.appDir).AutoLaunch {
		auto.Check()
	}

	for {
		select {
		case <-systray.ClickedCh:
			if !m.HasWindow() {
				c := Load(m.appDir)
				tasks.Trigger(c.LastURL + "/api/window/open")
			}
			break
		case <-open.ClickedCh:
			if !m.HasWindow() {
				m.Home()
			}
			break
		case <-launch.ClickedCh:
			if !m.HasWindow() {
				c := Load(m.appDir)
				tasks.Trigger(c.LastURL + "/api/window/open?quick=true")
			}
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
