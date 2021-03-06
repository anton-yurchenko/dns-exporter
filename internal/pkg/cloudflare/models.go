package cf

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"net/http"
)

// Zones hosted by a DNS provider
type Zones struct {
	Public map[string]string
}

// Client interface
type Client interface {
	ListZones(context.Context, ...string) ([]cloudflare.Zone, error)
}

// HTTPClient interface
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
