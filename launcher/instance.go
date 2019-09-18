package launcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Conquest-Reforged/ReforgedLauncher/instance"
	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
)

func (l *Launcher) instance(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		l.getInstance(w, r)
	}

	if r.Method == http.MethodDelete {
		l.deleteInstance(w, r)
	}

	if r.Method == http.MethodPost {
		l.postInstance(w, r)
	}
}

func (l *Launcher) getInstance(w http.ResponseWriter, r *http.Request) {
	inst, e := l.Instance(unescape(r.URL.Path))
	if e == nil {
		success(w, inst)
	} else {
		fail(w, e)
	}
}

func (l *Launcher) deleteInstance(w http.ResponseWriter, r *http.Request) {
	instances, e := l.LoadInstances()
	if e == nil {
		delete(instances, unescape(r.URL.Path))
		e = l.SaveInstances(instances)
		if e == nil {
			success(w, nil)
			return
		}
	}
	fail(w, e)
}

func (l *Launcher) postInstance(w http.ResponseWriter, r *http.Request) {
	inst, e := l.Instance(unescape(r.URL.Path))
	if e == nil {
		var opts instance.Instance
		e = json.NewDecoder(r.Body).Decode(&opts)
		if e == nil {
			inst.GameDir = opts.GameDir
			inst.Options = opts.Options
			inst.UserImage = opts.UserImage
			inst.AutoUpdate = opts.AutoUpdate
			e = l.SaveInstance(inst)
			if e == nil {
				success(w, nil)
				return
			}
		}
	}
	fail(w, e)
}

func (l *Launcher) LastInstance() (*instance.Instance, error) {
	instances, e := l.LoadInstances()
	if e != nil {
		return nil, e
	}

	if len(instances) == 0 {
		return nil, fmt.Errorf("no instances installed")
	}

	var last *instance.Instance
	for _, i := range instances {
		if last == nil || i.LastUsed.Unix() > last.LastUsed.Unix() {
			last = i
		}
	}

	return last, nil
}

func (l *Launcher) Instance(id string) (*instance.Instance, error) {
	instances, e := l.LoadInstances()
	if e != nil {
		return nil, e
	}

	i, ok := instances[id]
	if !ok {
		return nil, fmt.Errorf("instance not found: %s", id)
	}

	return i, nil
}

func (l *Launcher) NewInstance(id string, remote *modpack.Remote) *instance.Instance {
	return &instance.Instance{
		Name:       id,
		Options:    nil,
		AppDir:     l.AppDir,
		AutoUpdate: true,
		ModPack:    remote.AsModPack(),
		LastUsed:   time.Now(),
		GameDir:    instance.DefaultGameDir(id),
		RepoImage:  remote.Repo().CoverImage(),
	}
}

func (l *Launcher) SaveInstance(i *instance.Instance) error {
	instances, e := l.LoadInstances()
	if e != nil {
		return e
	}
	instances[i.Name] = i
	return l.SaveInstances(instances)
}
