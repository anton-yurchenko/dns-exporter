package app

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	vcs "github.com/anton-yurchenko/dns-exporter/internal/pkg/git"
	"github.com/anton-yurchenko/dns-exporter/internal/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-billy.v4/osfs"
)

var conf Configuration
var providers Providers

func init() {
	v := viper.New()

	vars := []string{
		"DELAY",
		"GIT_REMOTE_ENABLED",
		"GIT_URL",
		"GIT_BRANCH",
		"GIT_USER",
		"GIT_EMAIL",
		"GIT_TOKEN",
		"CLOUDFLARE_ENABLED",
		"CLOUDFLARE_EMAIL",
		"CLOUDFLARE_TOKEN",
		"ROUTE53_ENABLED",
		"AWS_REGION",
	}

	for _, variable := range vars {
		err := v.BindEnv(variable)
		if err != nil {
			log.Fatal(err)
		}
	}

	conf = Configuration{
		Providers: []string{},
		Clients: &Clients{
			HTTP: &http.Client{},
		},
		FileSystem: &Filesystems{
			Global: afero.NewOsFs(),
			Meta:   osfs.New("./data/.git"),
			Data:   osfs.New("./data"),
		},
		Project: &vcs.Project{
			Remote: &vcs.Origin{},
		},
	}

	if v.GetInt("DELAY") != 0 {
		conf.Delay = v.GetInt("DELAY")
	} else {
		conf.Delay = 1
	}

	initCloudflare(v)

	initRoute53(v)

	initGit(v)

	if len(conf.Providers) == 0 {
		log.Fatal("no enabled DNS providers")
	}
}

func initCloudflare(v *viper.Viper) {
	if v.GetBool("CLOUDFLARE_ENABLED") {
		var err error

		if !v.IsSet("CLOUDFLARE_EMAIL") {
			log.Fatal("missing env.var 'CLOUDFLARE_EMAIL'")
		}

		if !v.IsSet("CLOUDFLARE_TOKEN") {
			log.Fatal("missing env.var 'CLOUDFLARE_TOKEN'")
		}

		conf.Clients.CloudFlare, err = newCloudFlareClient(v.Get("CLOUDFLARE_EMAIL"), v.Get("CLOUDFLARE_TOKEN"))
		if err != nil {
			log.Fatal(err)
		}

		providers.CloudFlare.Public = make(map[string]string)
		conf.Providers = append(conf.Providers, "CloudFlare")
	}
}

func initRoute53(v *viper.Viper) {
	if v.GetBool("ROUTE53_ENABLED") {
		var err error

		if !v.IsSet("AWS_REGION") {
			log.Fatal("missing env.var 'AWS_REGION'")
		}

		conf.Clients.Route53, err = newRoute53Client()
		if err != nil {
			log.Fatal(err)
		}

		providers.Route53.Public = make(map[string]string)
		providers.Route53.Private = make(map[string]string)
		conf.Providers = append(conf.Providers, "Route53")
	}
}

func initGit(v *viper.Viper) {
	if v.GetBool("GIT_REMOTE_ENABLED") {
		if !v.IsSet("GIT_URL") {
			log.Fatal("missing env.var 'GIT_URL'")
		}
		conf.Project.Remote.URL = fmt.Sprintf("%v", v.Get("GIT_URL"))

		re := regexp.MustCompile(`https://[a-zA-Z0-9]*.[a-zA-Z]*/[a-zA-Z-_]*/[a-zA-Z-_]*.git`)
		if !re.MatchString(conf.Project.Remote.URL) {
			log.Fatal("provided 'GIT_URL' does not match the expected regex: 'https://[a-zA-Z0-9]*.[a-zA-Z]*/[a-zA-Z-_]*/[a-zA-Z-_]*.git'")
		}

		gitRepo := strings.Split(conf.Project.Remote.URL, "/")
		conf.Project.Name = strings.TrimSuffix(gitRepo[len(gitRepo)-1], ".git")

		if !v.IsSet("GIT_USER") {
			log.Fatal("missing env.var 'GIT_USER'")
		}
		conf.Project.AuthorName = fmt.Sprintf("%v", v.Get("GIT_USER"))

		if !v.IsSet("GIT_EMAIL") {
			log.Fatal("missing env.var 'GIT_EMAIL'")
		}
		conf.Project.AuthorEmail = fmt.Sprintf("%v", v.Get("GIT_EMAIL"))

		if !v.IsSet("GIT_TOKEN") {
			log.Fatal("missing env.var 'GIT_TOKEN'")
		}

		if v.Get("GIT_BRANCH") != "" {
			conf.Project.Remote.Branch = fmt.Sprintf("%v", v.Get("GIT_BRANCH"))
		} else {
			conf.Project.Remote.Branch = "master"
		}

	} else {
		conf.Project.AuthorName = "DNS-EXPORTER"
		conf.Project.AuthorEmail = "no-email@dns-exporter.com"
	}
}

// Entrypoint of an application
func Entrypoint(version string) {
	log.Info(fmt.Sprintf("dns-exporter v%s", version))

	dir, err := utils.ValidateDir(fmt.Sprintf("./%v/.git", "data"), false, conf.FileSystem.Global)
	if err != nil {
		log.Fatal(err)
	}

	if conf.Project.Remote.URL != "" {
		if dir {
			// Pull repository
			log.Info("pulling remote git repository")
			err := conf.Project.Pull(conf.FileSystem.Meta, conf.FileSystem.Data)
			if err != nil && err.Error() != "already up-to-date" {
				log.Fatal(err)
			} else if err.Error() == "already up-to-date" {
				log.Info("local repository is up-to-date with 'origin'")
			}

		} else {
			// Clone repository
			log.Info("cloning remote git repository")
			err := conf.Project.Clone(conf.FileSystem.Meta, conf.FileSystem.Data)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		if dir {
			log.Info("using local git repository")
		} else {
			// Init new repository
			log.Info("creating local git repository")
			err := vcs.Init(conf.FileSystem.Meta, conf.FileSystem.Data)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := conf.fetch(&providers); err != nil {
		log.Fatal(err)
	}

	if err := conf.export(&providers); err != nil {
		log.Fatal(err)
	}

	log.Info("commiting changes to local git repository")
	err = conf.Project.Commit(time.Now(), conf.FileSystem.Meta, conf.FileSystem.Data)
	if err == nil {
		if conf.Project.Remote.URL != "" {
			log.Info("pushing to remote git repository")
			if err := conf.Project.Push(conf.FileSystem.Meta); err != nil {
				log.Fatal(err)
			}
		}
	} else if err.Error() == "nothing to commit, working tree clean" {
		log.Info(err)
	} else {
		log.Fatal(err)
	}
}
