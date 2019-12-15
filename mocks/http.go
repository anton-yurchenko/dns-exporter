package mocks

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type HTTP struct {
	mock.Mock
}

func (c *HTTP) Do(req *http.Request) (*http.Response, error) {
	args := c.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}
