package tasks

import "io"

type Reader interface {
	io.Reader
	Len() int64
}

type sizedReader struct {
	io.Reader
	size int64
}

func (r *sizedReader) Len() int64 {
	return r.size
}

func WrapReader(reader io.Reader, size int64) Reader {
	return &sizedReader{
		Reader: reader,
		size:   size,
	}
}
