package minecraft

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
)

const url = "http://launchermeta.mojang.com/mc-staging/launcher.json"

type Meta struct {
	Linux   *AppMeta `json:"linux"`
	OSX     *AppMeta `json:"osx"`
	Windows *AppMeta `json:"windows"`
}

type AppMeta struct {
	X64          *JavaMeta `json:"64"`
	AppLink      string    `json:"applink"`
	AppHash      string    `json:"apphash"`
	DownloadHash string    `json:"downloadhash"`
}

type JavaMeta struct {
	JRE *RuntimeMeta `json:"jre"`
}

type RuntimeMeta struct {
	Sha1    string `json:"sha1"`
	Url     string `json:"url"`
	Version string `json:"version"`
}

func getMeta() (*Meta, error) {
	var meta Meta

	rs, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer files.Close(rs.Body)

	e = json.NewDecoder(rs.Body).Decode(&meta)
	if e != nil {
		return nil, e
	}

	return &meta, e
}

func (m *Meta) getAppLink() *AppMeta {
	switch runtime.GOOS {
	case "windows":
		return m.Windows
	case "linux":
		return m.Linux
	case "darwin":
		return m.OSX
	default:
		return nil
	}
}
