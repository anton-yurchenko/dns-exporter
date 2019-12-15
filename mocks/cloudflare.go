package mocks

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/stretchr/testify/mock"
)

type Cloudflare struct {
	mock.Mock
}

func (c *Cloudflare) ListZones(...string) ([]cloudflare.Zone, error) {
	args := c.Called()
	return args.Get(0).([]cloudflare.Zone), args.Error(1)
}
