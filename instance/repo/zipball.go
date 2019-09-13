package repo

import (
	"archive/zip"
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Conquest-Reforged/ReforgedLauncher/instance/optifine"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

type ExtractAction func(path string, f *zip.File) error

type ZipBall struct {
	name string
	path string
}

func (z *ZipBall) Extract(dest string, listener progress.Listener) error {
	r, e := zip.OpenReader(z.path)
	if e != nil {
		return e
	}
	defer files.Close(r)

	count := float64(0)
	total := float64(len(r.File))

	ch := make(chan float64)
	go progress.Listen(ch, listener)

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			count++
			continue
		}

		action := fixPath(dest, z.getAction(f))
		listener.TaskStatus("Extracting file: " + f.Name)
		e := action(f.Name, f)
		if e != nil {
			return e
		}
		count++
		ch <- count / total
	}
	return nil
}

func (z *ZipBall) getAction(f *zip.File) ExtractAction {
	if strings.HasSuffix(f.Name, "/dependencies.json") {
		return z.dependencyFile
	}
	if strings.HasSuffix(f.Name, "/optifine.txt") {
		return z.downloadOptifine
	}
	return z.copyFile
}

func fixPath(prefix string, action ExtractAction) ExtractAction {
	return func(path string, f *zip.File) error {
		i := strings.Index(path, "/")
		path = path[i+1:]
		path = filepath.Join(prefix, path)
		return action(path, f)
	}
}

func (z *ZipBall) copyFile(path string, f *zip.File) error {
	r, e := f.Open()
	if e != nil {
		panic(e)
	}
	defer files.Close(r)

	w, e := os.Create(files.MustFile(path))
	if e != nil {
		panic(e)
	}
	defer files.Close(w)

	_, e = io.Copy(w, r)
	return e
}

func (z *ZipBall) dependencyFile(path string, f *zip.File) error {
	dir := filepath.Dir(path)
	name := z.name + "-dep.json"
	path = filepath.Join(dir, name)
	return z.copyFile(path, f)
}

func (z *ZipBall) downloadOptifine(path string, f *zip.File) error {
	r, e := f.Open()
	if e != nil {
		panic(e)
	}
	defer files.Close(r)

	s := bufio.NewScanner(r)
	if s.Scan() {
		name := s.Text()
		dir := filepath.Dir(path)
		path = filepath.Join(dir, name)
		return optifine.Download(name, path, progress.Dummy())
	}

	return nil
}
