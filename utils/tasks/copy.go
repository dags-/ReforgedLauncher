package tasks

import (
	"io"
	"os"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
)

func Copy(reader Reader, writer io.Writer, listener progress.Listener) error {
	ch := make(chan float64)
	defer close(ch)
	go progress.Listen(ch, listener)

	count := float64(0)
	total := float64(reader.Len())
	buf := make([]byte, 8096)
	for {
		nr, e := reader.Read(buf)
		if nr == 0 && e == io.EOF {
			break
		}
		if e != nil && e != io.EOF {
			return e
		}
		_, e = writer.Write(buf[:nr])
		if e != nil {
			return e
		}
		count += float64(nr)
		ch <- count / total
	}
	return nil
}

func CopyPath(from, to string, listener progress.Listener) error {
	in, e := os.Open(from)
	defer files.Close(in)
	if e != nil {
		return e
	}

	stat, e := in.Stat()
	if e != nil {
		return e
	}

	out, e := os.Create(files.MustFile(to))
	defer files.Close(out)
	if e != nil {
		return e
	}

	if listener == nil {
		_, e = io.Copy(out, in)
		return e
	}

	return Copy(WrapReader(in, stat.Size()), out, listener)
}
