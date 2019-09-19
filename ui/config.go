package ui

import (
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/errs"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
)

type Config struct {
	WindowWidth  int    `json:"window_width"`
	WindowHeight int    `json:"window_height"`
	LastURL      string `json:"last_url"`
	AutoLaunch   bool   `json:"auto_launch"`
	HomeScale    int    `json:"home_scale"`
}

func FirstLoad(appDir, addr string) *Config {
	config := Load(appDir)
	config.LastURL = addr
	Save(appDir, config)
	return config
}

func Load(appDir string) *Config {
	config := &Config{
		WindowWidth:  800,
		WindowHeight: 480,
		HomeScale:    7,
		AutoLaunch:   true,
	}
	path := files.MustFile(appDir, "config.json")
	e := files.ParseJsonAt(path, config)
	if e != nil {
		Save(appDir, config)
	}
	return config
}

func Save(appDir string, c *Config) {
	path := files.MustFile(appDir, "config.json")
	e := files.WriteJsonAt(path, c)
	errs.Log("Save config", e)
}
