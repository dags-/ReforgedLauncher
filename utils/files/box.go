package files

import (
	"io"
	"os"
	"strings"

	rice "github.com/GeertJohan/go.rice"
)

func Extract(path, dir string, box *rice.Box) error {
	return visit(path, box, func(p string, f *rice.File) {
		out, e := os.Create(MustFile(dir, p))
		if e != nil {
			panic(e)
		}
		defer Close(out)
		_, e = io.Copy(out, f)
		if e != nil {
			panic(e)
		}
	})
}

func join(string ...string) string {
	return strings.Join(string, "/")
}

func visit(path string, box *rice.Box, visitor func(path string, f *rice.File)) error {
	f, e := box.Open(path)
	if e != nil {
		panic(e)
	}
	defer Close(f)

	stat, e := f.Stat()
	if e != nil {
		panic(e)
	}

	if stat.IsDir() {
		fs, e := f.Readdir(-1)
		if e != nil {
			panic(e)
		}
		for _, v := range fs {
			e := visit(join(path, v.Name()), box, visitor)
			if e != nil {
				panic(e)
			}
		}
	} else {
		visitor(path, f)
	}
	return nil
}
