package ui

import (
	"io/ioutil"

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
	exit := systray.AddMenuItem("Quit", "")
	for {
		select {
		case <-systray.ClickedCh:
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
			break
		case <-exit.ClickedCh:
			m.Exit()
			return
		}
	}
}

func (m *Manager) exit() {

}
