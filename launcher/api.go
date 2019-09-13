package launcher

import (
	"net/http"

	"github.com/Conquest-Reforged/ReforgedLauncher/ui"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

type Result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func listener(w *ui.Window) progress.Listener {
	return progress.Log(&progress.Progress{
		Status: func(s string) {
			w.Js(`status("%s");`, s)
		},
		Message: func(s string) {
			w.Js(`message("%s");`, s)
		},
		Task: func(f float64) {
			w.Js(`task(%v);`, f)
		},
		Global: func(f float64) {
			w.Js(`overall(%v);`, f)
		},
	})
}

func success(w http.ResponseWriter, i interface{}) {
	if i == nil {
		i = "empty"
	}
	e := files.WriteJson(w, &Result{
		Success: true,
		Data:    i,
	})
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func fail(w http.ResponseWriter, data error) {
	w.Header().Set("Content-Type", "application/json")
	e := files.WriteJson(w, &Result{
		Success: false,
		Data:    data.Error(),
	})

	if e == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}
