package repo

import (
	"path/filepath"

	"github.com/Conquest-Reforged/ReforgedLauncher/modpack"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/progress"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/tasks"
	"github.com/Conquest-Reforged/ReforgedLauncher/utils/version"
)

type Repository struct {
	appDir  string
	repoDir string
}

type context struct {
	baseDir string
	visited map[string]bool
}

type Dependencies map[string]string

func Open(appDir string) *Repository {
	return &Repository{
		appDir:  appDir,
		repoDir: files.MustDir(appDir, "Repository"),
	}
}

func (r *Repository) Has(mp *modpack.ModPack) bool {
	return files.Exists(filepath.Join(r.repoDir, mp.Repo.Name, mp.Version.String()))
}

func (r *Repository) PackDir(mp *modpack.ModPack) string {
	return filepath.Join(r.repoDir, mp.Repo.Name, mp.Version.String(), "pack")
}

func (r *Repository) MetaDir(mp *modpack.ModPack) string {
	return filepath.Join(r.repoDir, mp.Repo.Name, mp.Version.String(), "meta")
}

func (r *Repository) LoaderDir(mp *modpack.ModPack) string {
	return filepath.Join(r.repoDir, mp.Repo.Name, mp.Version.String(), "loader")
}

func (r *Repository) Pull(rep *modpack.Repo, v *version.Version, listener progress.Listener) error {
	release, e := rep.Get(v)
	if e != nil {
		return e
	}
	ctx := &context{
		baseDir: files.MustDir(r.repoDir, rep.Name, release.Version.String()),
		visited: map[string]bool{},
	}
	return r.pull(release, ctx, listener)
}

func (r *Repository) pull(release *modpack.Remote, ctx *context, listener progress.Listener) error {
	pack, e := r.downloadPack(release, ctx, listener)
	defer files.Del(pack)

	e = r.extractPack(release.Repo().Name, pack, ctx, listener)
	if e != nil {
		return e
	}

	e = r.pullDependencies(release.Repo().Name, ctx, listener)
	if e != nil {
		return e
	}

	return nil
}

func (r *Repository) downloadPack(release *modpack.Remote, ctx *context, listener progress.Listener) (string, error) {
	listener.TaskStatus("Downloading pack: " + release.Repo().Name)
	temp := files.TempFile(ctx.baseDir)
	return temp, tasks.Download(release.ZipUrl, temp, listener)
}

func (r *Repository) extractPack(name, path string, ctx *context, listener progress.Listener) error {
	listener.TaskStatus("Extracting pack: " + name)
	zipBall := &ZipBall{
		name: name,
		path: path,
	}
	return zipBall.Extract(ctx.baseDir, listener)
}

func (r *Repository) pullDependencies(name string, ctx *context, listener progress.Listener) error {
	listener.TaskStatus("Checking dependencies")

	var dependencies Dependencies
	deps := filepath.Join(ctx.baseDir, name+"-dep.json")
	if !files.Exists(deps) {
		return nil
	}

	e := files.ParseJsonAt(deps, &dependencies)
	if e != nil {
		return nil
	}

	listener.GlobalStatus("Downloading dependencies")
	for k, v := range dependencies {
		if _, ok := ctx.visited[k]; ok {
			continue
		}

		ctx.visited[k] = true
		rep, e := modpack.ParseRepo(k)
		if e != nil {
			continue
		}

		var remote *modpack.Remote
		if v == "latest" {
			remote, e = rep.Latest()
		} else {
			remote, e = rep.Get(version.Parse(v))
		}

		if e != nil {
			return e
		}

		e = r.pull(remote, ctx, listener)
		if e != nil {
			return e
		}
	}

	return nil
}
