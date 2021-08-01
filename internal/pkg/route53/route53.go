package r53

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"dns-exporter/internal/pkg/utils"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Fetch hosted zones
func (z Zones) Fetch(c Client, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	var t *string

	for {
		o, err := c.ListHostedZones(&route53.ListHostedZonesInput{
			Marker: t,
		})
		if err != nil {
			errs <- errors.Wrap(err, "Route53: error fetching zones")
			return
		}

		for _, zone := range o.HostedZones {
			if *zone.Config.PrivateZone {
				z.Private[*zone.Name] = *zone.Id
			} else {
				z.Public[*zone.Name] = *zone.Id
			}
		}

		if *o.IsTruncated {
			t = o.NextMarker
		} else {
			break
		}
	}

	errs <- nil
}

// Export hosted zones
func (z Zones) Export(c Client, delay int, errs chan error, wg *sync.WaitGroup, root string, fs afero.Fs) {
	defer wg.Done()

	parent := fmt.Sprintf("%v/Route53", root)
	public := fmt.Sprintf("%v/Route53/Public", root)
	private := fmt.Sprintf("%v/Route53/Private", root)

	// validate filetree
	for _, dir := range []string{parent, public, private} {
		_, err := utils.ValidateDir(dir, true, fs)
		if err != nil {
			errs <- errors.Wrap(err, "Route53: error exporting zones")
			return
		}
	}

	export := func(domain, id, t, dir string, c Client, fs afero.Fs) {
		log.WithFields(log.Fields{
			"provider": "Route53",
			"zone":     strings.TrimSuffix(domain, "."),
			"type":     t,
		}).Info("exporting zone")

		r, err := getRecords(id, c)
		if err != nil {
			errs <- errors.Wrap(err, fmt.Sprintf("Route53: error retrieving zone '%s' records", domain))
			return
		}

		content, err := r.ConvertToZonefile()
		if err != nil {
			errs <- errors.Wrap(err, fmt.Sprintf("Route53: error composing zonefile for '%s' zone", domain))
			return
		}

		// write zonefile
		_, err = utils.WriteToFile(domain, content.String(), dir, fs)
		if err != nil {
			errs <- errors.Wrap(err, fmt.Sprintf("CloudFlare: error exporting zone: '%s'", domain))
			return
		}

		time.Sleep(time.Duration(delay) * time.Second)
	}

	// export zonefiles
	for domain, id := range z.Public {
		export(domain, id, "public", public, c, fs)
	}

	for domain, id := range z.Private {
		export(domain, id, "private", private, c, fs)
	}

	errs <- nil
}

// getRecords returns parsed zone records struct
func getRecords(id string, c Client) (Records, error) {
	var r Records

	var t, n *string
	for {
		o, err := c.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
			HostedZoneId:    &id,
			StartRecordType: t,
			StartRecordName: n,
		})
		if err != nil {
			return Records{}, errors.Wrap(err, "error retrieving zone records")
		}

		err = r.Append(o)
		if err != nil {
			return Records{}, errors.Wrap(err, "error mapping zone records")
		}

		if *o.IsTruncated {
			t = o.NextRecordType
			n = o.NextRecordName
		} else {
			break
		}
	}

	return r, nil
}
