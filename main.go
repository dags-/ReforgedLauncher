package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/GeertJohan/go.rice"
	"github.com/marcsauter/single"

	"github.com/Conquest-Reforged/ReforgedLauncher/launcher"
	"github.com/Conquest-Reforged/ReforgedLauncher/ui"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
)

const (
	BRANDING = "ReforgedLauncher"
	MODPACKS = "https://io.conquestreforged.com/modpacks/modpacks.json"
)

var (
	w = flag.String("w", "", "")
)

func main() {
	flag.Parse()

	if *w == "" {
		startMaster()
	} else {
		startSlave()
	}
}

func startMaster() {
	// get user
	u, e := user.Current()
	errs.Panic("User", e)

	// init properties
	var properties launcher.Properties
	properties.Branding = BRANDING
	properties.ModPacksURL = MODPACKS
	properties.AppDir = files.MustDir(u.HomeDir, "AppData", "Local", BRANDING)

	// ensure only one instance of the app - trigger the /api/open endpoint if already running
	s := holdLock(properties)
	defer releaseLock(s)

	// set up logging
	logFile, e := os.Create(files.MustFile(properties.AppDir, "launcher.log"))
	errs.Panic("Open log file", e)
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	// setup & extract resources
	public := rice.MustFindBox("assets/public")
	private := rice.MustFindBox("assets/private")
	e = files.Extract("", properties.AppDir, private)
	errs.Panic("Extract assets", e)

	// launch
	l := launcher.NewLauncher(&properties, public)
	l.Run()
}

func startSlave() {
	var settings ui.Settings
	e := json.Unmarshal([]byte(*w), &settings)
	errs.Panic("Parse window settings", e)
	ui.Open(&settings)
}

func holdLock(props launcher.Properties) *single.Single {
	s := single.New(props.Branding)
	e := s.CheckLock()
	if e != nil {
		c := ui.Load(props.AppDir)
		r, e := http.Get(c.LastURL + "/api/open/window")
		if e == nil {
			defer files.Close(r.Body)
			os.Exit(0)
		} else {
			panic(e)
		}
	}
	return s
}

func releaseLock(s *single.Single) {
	errs.Panic("Release lock", s.TryUnlock())
}
