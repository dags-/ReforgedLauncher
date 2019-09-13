package modpack

import (
	"fmt"
	"strings"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/version"
)

const (
	releasesUrl = `https://api.github.com/repos/%s/releases`
)

type Repo struct {
	Owner string
	Name  string
}

func NewRepo(owner, name string) *Repo {
	return &Repo{
		Owner: owner,
		Name:  name,
	}
}

func ParseRepo(path string) (*Repo, error) {
	repo := &Repo{}
	e := parseRepo(path, repo)
	if e != nil {
		return nil, e
	}
	return repo, nil
}

func (r *Repo) Latest() (*Remote, error) {
	releases, e := r.Releases()
	if e != nil {
		return nil, e
	}
	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases available")
	}
	return releases[0], nil
}

func (r *Repo) Releases() ([]*Remote, error) {
	var releases []*Remote
	e := tasks.GetJson(fmt.Sprintf(releasesUrl, r.String()), &releases)
	if e != nil {
		return nil, e
	}
	for _, rel := range releases {
		rel.repo = r
	}
	return releases, nil
}

func (r *Repo) Get(v *version.Version) (*Remote, error) {
	rls, e := r.Releases()
	if e != nil {
		return nil, e
	}
	for _, rl := range rls {
		if rl.Version.Matches(v) {
			return rl, nil
		}
	}
	return nil, fmt.Errorf("version not found: %s", v)
}

func (r *Repo) String() string {
	return r.Owner + "/" + r.Name
}

func (r *Repo) MarshalJSON() ([]byte, error) {
	return []byte(`"` + r.String() + `"`), nil
}

func (r *Repo) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = s[1 : len(s)-1]
	return parseRepo(s, r)
}

func parseRepo(path string, r *Repo) error {
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid path: %s", path)
	} else {
		r.Owner = parts[0]
		r.Name = parts[1]
		return nil
	}
}
