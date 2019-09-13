package ui

import (
	"fmt"
	"html"
	"os"

	"github.com/dags-/webview"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
)

type Window struct {
	w webview.WebView
	b bool
}

type Settings struct {
	Title      string
	Width      int
	Height     int
	Path       string
	Resizable  bool
	Borderless bool
}

func (w *Window) Run() bool {
	return w.w.Loop(true)
}

func (w *Window) Js(js string, args ...interface{}) {
	for i, a := range args {
		args[i] = html.EscapeString(fmt.Sprint(a))
	}

	js = fmt.Sprintf(js, args...)
	w.w.Dispatch(func() {
		errs.Panic("Eval "+js, w.w.Eval(js))
	})
}

func (w *Window) Close() {
	w.w.Exit()
	w.w.Run()
}

func (w *Window) dispose() {
	w.w.Exit()
	w.w.Loop(false)
}

func (m *Manager) handleWindow() {
	select {
	case <-m.exit:
		os.Exit(0)
		return
	case s := <-m.create:
		m.window = m.createWindow(s)
	}

	for m.window.Run() {
		select {
		case <-m.open:
			continue
		case <-m.close:
			m.window.Close()
			return
		case <-m.exit:
			m.window.Close()
			os.Exit(0)
			return
		case s := <-m.create:
			m.window.Close()
			m.window = m.createWindow(s)
			continue
		default:
			continue
		}
	}

	m.window.dispose()
}
