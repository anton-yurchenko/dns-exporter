package app

import (
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// fetch hosted zones from configured providers
func (c *Configuration) fetch(p *Providers) error {
	errs := make(chan error, len(c.Providers))

	var wg sync.WaitGroup
	wg.Add(len(c.Providers))

	// fetch each provide in a separate routine
	for _, provider := range c.Providers {
		log.WithFields(log.Fields{
			"provider": provider,
		}).Info("fetching zones")

		if provider == "CloudFlare" {
			go p.CloudFlare.Fetch(c.Clients.CloudFlare, errs, &wg)
		}

		if provider == "Route53" {
			go p.Route53.Fetch(c.Clients.Route53, errs, &wg)
		}
	}

	wg.Wait()

	// print all errors and exit
	f := false
	for i := 0; i <= len(c.Providers)-1; i++ {
		r := <-errs
		if r != nil {
			log.Error(r)
			f = true
		}
	}
	if f {
		return errors.New("errors encountered during zone fetching")
	}

	return nil
}

// export zonefiles from configured providers
func (c *Configuration) export(p *Providers) error {
	errs := make(chan error, len(c.Providers))

	var wg sync.WaitGroup
	wg.Add(len(c.Providers))

	// fetch each provide in a sepparate routine
	for _, provider := range c.Providers {
		if provider == "CloudFlare" {
			go p.CloudFlare.Export(c.Clients.HTTP, c.Delay, errs, &wg, "./data", c.FileSystem.Global)
		}

		if provider == "Route53" {
			go p.Route53.Export(c.Clients.Route53, c.Delay, errs, &wg, "./data", c.FileSystem.Global)
		}
	}

	wg.Wait()

	// print all errors and exit
	f := false
	for i := 0; i <= len(c.Providers)-1; i++ {
		r := <-errs
		if r != nil {
			log.Error(r)
			f = true
		}
	}
	if f {
		return errors.New("errors encountered during zone fetching")
	}

	return nil
}
