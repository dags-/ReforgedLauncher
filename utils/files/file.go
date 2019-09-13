package files

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
)

func Close(c io.Closer) {
	if c != nil {
		errs.Log("Close", c.Close())
	}
}

func TempFile(dir ...string) string {
	name := fmt.Sprint(time.Now().Unix(), ".temp")
	path := append(dir, name)
	return MustFile(path...)
}

func RenameRel(path, name string) (string, error) {
	out := RelFile(path, name)
	return out, os.Rename(path, out)
}

func RelFile(path, name string) string {
	dir := filepath.Dir(path)
	return MustFile(dir, name)
}

func MustFile(path ...string) string {
	file := filepath.Join(path...)
	dir := filepath.Dir(file)
	e := os.MkdirAll(dir, os.ModePerm)
	if e != nil {
		panic(e)
	}
	return file
}

func Del(path string) {
	if Exists(path) {
		e := os.RemoveAll(path)
		errs.Log("Del "+path, e)
	}
}

func Exists(path string) bool {
	_, e := os.Stat(path)
	return e == nil
}

func CheckMD5(path, hash string) bool {
	hasher := sha1.New()
	h, e := Hash(path, hasher)
	if e != nil {
		return false
	}
	return h == hash
}

func Hash(path string, hasher hash.Hash) (string, error) {
	f, e := os.Open(path)
	if e != nil {
		return "", e
	}
	defer Close(f)

	_, e = io.Copy(hasher, f)
	if e != nil {
		return "", e
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
