package files

import (
	"encoding/json"
	"io"
	"os"
)

func ParseJson(r io.Reader, i interface{}) error {
	return json.NewDecoder(r).Decode(i)
}

func ParseJsonAt(path string, i interface{}) error {
	in, e := os.Open(path)
	if e != nil {
		return e
	}
	defer Close(in)
	return ParseJson(in, i)
}

func WriteJsonAt(path string, i interface{}) error {
	path = MustFile(path)
	out, e := os.Create(path)
	if e != nil {
		return e
	}
	defer Close(out)
	return WriteJson(out, i)
}

func WriteJson(w io.Writer, i interface{}) error {
	en := json.NewEncoder(w)
	en.SetIndent("", "  ")
	return en.Encode(i)
}
