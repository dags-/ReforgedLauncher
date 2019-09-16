package minecraft

import (
	"encoding/json"
	"errors"
	"net/http"
	"runtime"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
)

const url = "http://launchermeta.mojang.com/mc-staging/launcher.json"

type Meta struct {
	Linux   *AppLink `json:"linux"`
	OSX     *AppLink `json:"osx"`
	Windows *AppLink `json:"windows"`
}

type AppLink struct {
	X64          *Java  `json:"64"`
	AppLink      string `json:"applink"`
	AppHash      string `json:"apphash"`
	DownloadHash string `json:"downloadhash"`
}

type Java struct {
	JRE *Runtime `json:"jre"`
}

type Runtime struct {
	Sha1    string `json:"sha1"`
	Url     string `json:"url"`
	Version string `json:"version"`
}

func download(appdir string, listener progress.Listener) (string, error) {
	m, e := getMeta()
	if e != nil {
		return "", e
	}

	link := m.getAppLink()
	path := files.TempFile(appdir, "Launcher")

	listener.GlobalStatus("Downloading Minecraft launcher")
	e = tasks.Download(link.AppLink, path, listener)
	if e != nil {
		return "", e
	}

	listener.GlobalStatus("Checking hash")
	if !files.CheckMD5(path, link.DownloadHash) {
		return "", errors.New("invalid md5 hash")
	}
	return path, nil
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

func (m *Meta) getAppLink() *AppLink {
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
