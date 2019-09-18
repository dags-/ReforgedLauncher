package launcher

import (
	"net/http"
	"sort"

	"github.com/Conquest-Reforged/ReforgedLauncher/instance"
	"github.com/Conquest-Reforged/ReforgedLauncher/instance/repo"
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
				Cover1:  i.UserImage,
				Cover2:  i.RepoImage,
			})
		}
	}

	success(w, metas)
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

func initializeInstances(l *Launcher, i map[string]*instance.Instance) {
	for _, v := range i {
		v.AppDir = l.AppDir
	}
}
