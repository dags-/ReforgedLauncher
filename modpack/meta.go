package modpack

import (
	"fmt"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
)

type Meta struct {
	Repo  *Repo  `json:"repo"`
	Name  string `json:"name"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
}

type GithubRepo struct {
	Description string     `json:"description"`
	Name        string     `json:"name"`
	Owner       GithubUser `json:"owner"`
}

type GithubUser struct {
	Icon string `json:"avatar_url"`
}

const (
	repoUrl = "http://api.github.com/repos/%s/%s"
)

func (r *Repo) GetMeta() (*Meta, error) {
	var repo GithubRepo
	e := tasks.GetJson(fmt.Sprintf(repoUrl, r.Owner, r.Name), &repo)
	if e != nil {
		return nil, e
	}
	return &Meta{
		Repo:  r,
		Name:  repo.Name,
		Title: repo.Description,
		Icon:  repo.Owner.Icon,
	}, nil
}
