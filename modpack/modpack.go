package modpack

import (
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/version"
)

type ModPack struct {
	Repo    *Repo            `json:"repo"`
	Version *version.Version `json:"version"`
}

type Remote struct {
	repo    *Repo
	Version *version.Version `json:"tag_name"`
	ZipUrl  string           `json:"zipball_url"`
}

func (r *Remote) AsModPack() *ModPack {
	return &ModPack{
		Repo:    r.repo,
		Version: r.Version,
	}
}

func (r *Remote) Repo() *Repo {
	return r.repo
}
