package launcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"time"

	"github.com/Conquest-Reforged/ReforgedLauncher/instance"
	"github.com/Conquest-Reforged/ReforgedLauncher/instance/repo"
	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
)

func (l *Launcher) instances(w http.ResponseWriter, r *http.Request) {
	instances, e := l.Instances()
	if e != nil {
		fail(w, e)
		return
	}

	rep := repo.Open(l.AppDir)
	var metas []*instance.Meta
	for _, i := range instances {
		if rep.Has(i.ModPack) {
			metas = append(metas, &instance.Meta{
				ModPack: i.ModPack,
				Name:    i.Name,
				Pack:    i.ModPack.Repo.String(),
				Cover1:  i.Image,
				Cover2:  filepath.Join(rep.MetaDir(i.ModPack), "cover.jpg"),
			})
		}
	}

	success(w, metas)
}

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
	inst, e := l.Instance(r.URL.Path)
	if e == nil {
		success(w, inst)
	} else {
		fail(w, e)
	}
}

func (l *Launcher) deleteInstance(w http.ResponseWriter, r *http.Request) {
	instances, e := l.LoadInstances()
	if e == nil {
		delete(instances, r.URL.Path)
		e = l.SaveInstances(instances)
		if e == nil {
			success(w, nil)
			return
		}
	}
	fail(w, e)
}

func (l *Launcher) postInstance(w http.ResponseWriter, r *http.Request) {
	inst, e := l.Instance(r.URL.Path)
	if e == nil {
		var opts instance.Instance
		e = json.NewDecoder(r.Body).Decode(&opts)
		if e == nil {
			inst.GameDir = opts.GameDir
			inst.Options = opts.Options
			inst.Image = opts.Image
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
		Name:     id,
		Options:  nil,
		AppDir:   l.AppDir,
		ModPack:  remote.AsModPack(),
		LastUsed: time.Now(),
		GameDir:  instance.DefaultGameDir(id),
	}
}

func (l *Launcher) Instances() ([]*instance.Instance, error) {
	instances, e := l.LoadInstances()
	if e != nil {
		return nil, e
	}

	result := make([]*instance.Instance, len(instances))
	i := 0
	for _, inst := range instances {
		result[i] = inst
		i++
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].LastUsed.After(result[j].LastUsed)
	})

	return result, nil
}

func (l *Launcher) LoadInstances() (map[string]*instance.Instance, error) {
	instances := map[string]*instance.Instance{}
	defer initializeInstances(l, instances)

	path := files.MustFile(l.AppDir, "instances.json")
	e := files.ParseJsonAt(path, &instances)
	if e != nil {
		e := l.SaveInstances(instances)
		if e != nil {
			return nil, e
		}
	}

	return instances, nil
}

func (l *Launcher) SaveInstances(instances map[string]*instance.Instance) error {
	path := files.MustFile(l.AppDir, "instances.json")
	e := files.WriteJsonAt(path, &instances)
	if e != nil {
		return e
	}
	return nil
}

func (l *Launcher) SaveInstance(i *instance.Instance) error {
	instances, e := l.LoadInstances()
	if e != nil {
		return e
	}
	instances[i.Name] = i
	return l.SaveInstances(instances)
}

func initializeInstances(l *Launcher, i map[string]*instance.Instance) {
	for _, v := range i {
		v.AppDir = l.AppDir
	}
}
