package mocks

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"github.com/stretchr/testify/mock"
)

type Cloudflare struct {
	mock.Mock
}

func (c *Cloudflare) ListZones(context.Context, ...string) ([]cloudflare.Zone, error) {
	args := c.Called()
	return args.Get(0).([]cloudflare.Zone), args.Error(1)
}
