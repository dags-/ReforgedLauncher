package tasks

import (
	"archive/zip"
	"io"
	"os"
	"strings"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

func Unzip(path, dir string, listener progress.Listener) error {
	return UnzipTrimFirst(path, dir, listener)
}

func UnzipTrimFirst(path, dir string, listener progress.Listener) error {
	z, e := zip.OpenReader(path)
	if e != nil {
		return e
	}
	defer files.Close(z)

	ch := make(chan float64)
	go progress.Listen(ch, listener)

	count := float64(0)
	total := float64(len(z.File) - 1)
	for _, f := range z.File {
		e := extract(f, dir)
		if e != nil {
			return e
		}
		count++
		ch <- count / total
	}

	return nil
}

func extract(f *zip.File, dir string) error {
	if f.FileInfo().IsDir() {
		return nil
	}

	r, e := f.Open()
	if e != nil {
		return e
	}
	defer files.Close(r)

	path := f.Name[strings.Index(f.Name, "/")+1:]
	w, e := os.Create(files.MustFile(dir, path))
	if e != nil {
		return e
	}
	defer files.Close(w)

	_, e = io.Copy(w, r)
	return e
}
