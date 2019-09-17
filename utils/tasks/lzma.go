package tasks

import (
	"io"
	"os"

	"github.com/andrew-d/lzma"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

func Decompress(path, dest string, listener progress.Listener) error {
	out, e := os.Create(dest)
	if e != nil {
		return e
	}
	defer files.Close(out)

	in, e := os.Open(path)
	if e != nil {
		return e
	}
	defer files.Close(in)

	re := lzma.NewReader(in)
	defer files.Close(re)

	listener.TaskProgress(-1)
	_, e = io.Copy(out, re)
	return e
}
