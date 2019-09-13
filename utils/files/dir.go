package files

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func RelDir(path, name string) string {
	dir := filepath.Dir(path)
	return MustDir(dir, name)
}

func MustDir(path ...string) string {
	dir := filepath.Join(path...)
	e := os.MkdirAll(dir, os.ModePerm)
	if e != nil {
		panic(e)
	}
	return dir
}

func Match(dir, pattern string) []string {
	var paths []string
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		return paths
	}

	r, e := regexp.Compile(pattern)
	if e != nil {
		return paths
	}

	for _, f := range files {
		if r.MatchString(f.Name()) {
			paths = append(paths, filepath.Join(dir, f.Name()))
		}
	}
	return paths
}
