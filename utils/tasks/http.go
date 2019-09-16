package tasks

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

func Trigger(url string) {
	r, e := http.Get(url)
	if e == nil {
		defer files.Close(r.Body)
	}
}

func Download(url, path string, listener progress.Listener) error {
	out, e := os.Create(files.MustFile(path))
	if e != nil {
		return e
	}
	defer files.Close(out)

	for {
		in, e := http.Get(url)
		if e != nil {
			return e
		}

		if in.ContentLength <= 0 {
			files.Close(in.Body)
			time.Sleep(time.Millisecond * 500)
			continue
		}

		e = Copy(WrapReader(in.Body, in.ContentLength), out, listener)
		files.Close(in.Body)
		return e
	}
}

func GetJson(url string, i interface{}) error {
	r, e := http.Get(url)
	if e != nil {
		return e
	}
	defer files.Close(r.Body)
	return json.NewDecoder(r.Body).Decode(i)
}
