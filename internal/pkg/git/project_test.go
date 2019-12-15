package vcs_test

import (
	"testing"

	vcs "github.com/anton-yurchenko/dns-exporter/internal/pkg/git"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage"
)

func TestInit(t *testing.T) {
	err := vcs.Init(memfs.New(), memfs.New())
	if err != nil {
		t.Fatal("\nEXPECTED error: \n<nil>\n\nGOT error:", err)
	}
}

func TestClone(t *testing.T) {
	m := vcs.Project{
		Remote: &vcs.Origin{
			URL:    "https://localhost",
			Branch: "master",
		},
	}

	vcs.GitClone = func(s storage.Storer, worktree billy.Filesystem, o *git.CloneOptions) (*git.Repository, error) {
		return &git.Repository{}, nil
	}

	err := m.Clone(memfs.New(), memfs.New())
	if err != nil {
		t.Fatal("\nEXPECTED error: \n<nil>\n\nGOT error:", err)
	}
}
