package cf

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/anton-yurchenko/dns-exporter/internal/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Fetch hosted zones
func (z Zones) Fetch(c Client, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	r, err := c.ListZones()
	if err != nil {
		errs <- errors.Wrap(err, "CloudFlare: error fetching zones")
		return
	}

	for _, zone := range r {
		z.Public[zone.Name] = zone.ID
	}

	errs <- nil
}

// Export hosted zones
func (z Zones) Export(c HTTPClient, delay int, errs chan error, wg *sync.WaitGroup, root string, fs afero.Fs) {
	defer wg.Done()

	// validate provider export dir
	dir := fmt.Sprintf("%v/CloudFlare", root)
	_, err := utils.ValidateDir(dir, true, fs)
	if err != nil {
		errs <- errors.Wrap(err, "CloudFlare: error exporting zones")
		return
	}

	// export zonefiles
	for domain, id := range z.Public {
		log.WithFields(log.Fields{
			"provider": "CloudFlare",
			"zone":     domain,
		}).Info("exporting zone")

		// export zone
		content, err := exportZone(c, id)
		if err != nil {
			errs <- errors.Wrap(err, fmt.Sprintf("CloudFlare: error exporting zone: '%s'", domain))
			return
		}

		// write zonefile
		_, err = utils.WriteToFile(
			domain,
			string(";; SOA Record\n"+strings.TrimPrefix(strings.Split(string(content), ";; SOA Record")[1], "\n")),
			dir,
			fs,
		)
		if err != nil {
			errs <- errors.Wrap(err, fmt.Sprintf("CloudFlare: error exporting zone: '%s'", domain))
			return
		}

		time.Sleep(time.Duration(delay) * time.Second)
	}

	errs <- nil
}

// exportZone returns zonefile content for specified zone id
func exportZone(c HTTPClient, zoneID string) ([]byte, error) {
	// request export
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/export", zoneID), nil)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("CloudFlare: error constructing HTTP request (id='%s')", zoneID))
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Auth-Email", os.Getenv("CLOUDFLARE_EMAIL"))
	request.Header.Add("X-Auth-Key", os.Getenv("CLOUDFLARE_TOKEN"))

	response, err := c.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("CloudFlare: error consuming '/client/v4/zones/%s/dns_records/export'", zoneID))
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "CloudFlare: error reading responce body")
	}

	return b, nil
}
