package files

import (
	"path/filepath"
	"strings"
)

func IsImage(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}
