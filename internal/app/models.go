package app

import (
	cf "github.com/anton-yurchenko/dns-exporter/internal/pkg/cloudflare"
	vcs "github.com/anton-yurchenko/dns-exporter/internal/pkg/git"
	r53 "github.com/anton-yurchenko/dns-exporter/internal/pkg/route53"
	"github.com/spf13/afero"
	"gopkg.in/src-d/go-billy.v4"
)

// Configuration contains app runtime config
type Configuration struct {
	Providers  []string
	Project    *vcs.Project
	FileSystem *Filesystems
	Clients    *Clients
	Delay      int
}

// Filesystems contains different filesystems abstractions
type Filesystems struct {
	Global afero.Fs
	Meta   billy.Filesystem
	Data   billy.Filesystem
}

// Clients contains authenticated DNS provider clients
type Clients struct {
	CloudFlare cf.Client
	HTTP       cf.HTTPClient
	Route53    r53.Client
}

// Providers contains fetched providers Zones
type Providers struct {
	CloudFlare cf.Zones
	Route53    r53.Zones
}
