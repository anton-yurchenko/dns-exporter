package vcs

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

// GitClone is required for stubbing
var GitClone = git.Clone

// GitNewRemote is required for stubbing
var GitNewRemote = git.NewRemote

// Init a new Git project
func Init(meta, data billy.Filesystem) error {
	_, err := git.Init(
		filesystem.NewStorage(meta, cache.NewObjectLRU(cache.DefaultMaxSize)),
		data,
	)
	if err != nil {
		return errors.Wrap(err, "error initializing git repository")
	}

	return nil
}

// Clone repository from origin
func (p Project) Clone(meta, data billy.Filesystem) error {
	_, err := GitClone(
		filesystem.NewStorage(meta, cache.NewObjectLRU(cache.DefaultMaxSize)),
		data,
		&git.CloneOptions{
			URL:           p.Remote.URL,
			ReferenceName: plumbing.NewBranchReferenceName(p.Remote.Branch),
			SingleBranch:  true,
			Depth:         1,
			Auth:          &http.BasicAuth{Username: os.Getenv("GIT_USER"), Password: os.Getenv("GIT_TOKEN")},
		})
	if err != nil {
		return errors.Wrap(err, "error cloning remote repository")
	}

	return nil
}

// Pull from origin
func (p Project) Pull(meta, data billy.Filesystem) error {
	repo, err := git.Open(
		filesystem.NewStorage(meta, cache.NewObjectLRU(cache.DefaultMaxSize)),
		data,
	)
	if err != nil {
		return errors.Wrap(err, "error opening repository")
	}

	w, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "error retreiving git worktree")
	}

	// Fetch
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Auth:       &http.BasicAuth{Username: os.Getenv("GIT_USER"), Password: os.Getenv("GIT_TOKEN")},
	})
	if err != nil && err.Error() != "already up-to-date" {
		return errors.Wrap(err, "error fetching from 'origin'")
	}

	if err.Error() == "already up-to-date" {
		return err
	}

	// Pull
	err = w.Pull(&git.PullOptions{
		RemoteName:   "origin",
		SingleBranch: true,
		Auth:         &http.BasicAuth{Username: os.Getenv("GIT_USER"), Password: os.Getenv("GIT_TOKEN")},
	})
	if err != nil {
		return errors.Wrap(err, "error pulling from 'origin'")
	}

	return nil
}

// Commit to local repository
func (p Project) Commit(timestamp time.Time, meta, data billy.Filesystem) error {
	repo, err := git.Open(
		filesystem.NewStorage(meta, cache.NewObjectLRU(cache.DefaultMaxSize)),
		data,
	)
	if err != nil {
		return errors.Wrap(err, "error opening repository")
	}

	w, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "error retreiving git worktree")
	}

	_, err = w.Add("")
	if err != nil {
		return errors.Wrap(err, "error adding modified files")
	}

	s, err := w.Status()
	if err != nil {
		return errors.Wrap(err, "error retreiving git status")
	}

	if s.IsClean() {
		return errors.New("nothing to commit, working tree clean")
	}

	_, err = w.Commit(timestamp.Format("2006-01-02T15:04:05"), &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  p.AuthorName,
			Email: p.AuthorEmail,
			When:  timestamp,
		},
		Committer: &object.Signature{
			Name:  p.AuthorName,
			Email: p.AuthorEmail,
			When:  timestamp,
		},
	})
	if err != nil {
		return errors.Wrap(err, "error commiting changes")
	}

	return nil
}

// Push to origin
func (p Project) Push(meta billy.Filesystem) error {
	r := GitNewRemote(
		filesystem.NewStorage(meta, cache.NewObjectLRU(cache.DefaultMaxSize)),
		&config.RemoteConfig{
			Name: "origin",
			URLs: []string{p.Remote.URL},
		},
	)

	err := r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{Username: os.Getenv("GIT_USER"), Password: os.Getenv("GIT_TOKEN")},
	})
	if err != nil {
		return errors.Wrap(err, "error pushing to 'origin'")
	}

	return nil
}
