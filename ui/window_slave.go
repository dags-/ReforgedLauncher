package ui

import (
	"bufio"
	"encoding/json"
	"github.com/dags-/webview"
	"os"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
)

func Open(settings *Settings) {
	w := webview.New(webview.Settings{
		Title:      settings.Title,
		Icon:       settings.Icon,
		URL:        settings.Url,
		Width:      settings.Width,
		Height:     settings.Height,
		Resizable:  settings.Resizable,
		Borderless: settings.Borderless,
	})

	defer w.Exit()

	go listen(w)

	w.SetColor(0, 0, 0, 255)
	w.Run()
}

func listen(w webview.WebView) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		var event Event
		e := json.Unmarshal(s.Bytes(), &event)
		errs.Log("Parse event", e)
		handle(w, &event)
	}
}

func handle(w webview.WebView, e *Event) {
	switch e.Type {
	case INVOKE:
		w.Dispatch(func() {
			e := w.Eval(e.Message)
			errs.Log("Eval", e)
		})
	case INJECT:
		w.Dispatch(func() {
			w.InjectCSS(e.Message)
		})
	case EXIT:
		w.Dispatch(func() {
			w.Exit()
		})
	}
}
