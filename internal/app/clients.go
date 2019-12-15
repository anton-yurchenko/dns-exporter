package app

import (
	"fmt"
	cf "github.com/anton-yurchenko/dns-exporter/internal/pkg/cloudflare"
	r53 "github.com/anton-yurchenko/dns-exporter/internal/pkg/route53"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
	"os"
)

// newCloudFlareClient returns new Cloudflare client
func newCloudFlareClient(email, token interface{}) (cf.Client, error) {
	c, err := cloudflare.New(fmt.Sprintf("%v", token), fmt.Sprintf("%v", email))
	if err != nil {
		return nil, errors.Wrap(err, "error creating CloudFlare client")
	}

	return c, nil
}

// newRoute53Client returns new Route53 client
func newRoute53Client() (r53.Client, error) {
	r := os.Getenv("AWS_REGION")
	if r == "" {
		return nil, errors.New("missing required environmental variable 'AWS_REGION'")
	}

	a := os.Getenv("AWS_SDK_LOAD_CONFIG")
	if a == "" {
		err := os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
		if err != nil {
			return nil, errors.Wrap(err, "error setting environmental variable 'AWS_SDK_LOAD_CONFIG=true'")
		}
	}

	return route53.New(session.Must(session.NewSession())), nil
}
