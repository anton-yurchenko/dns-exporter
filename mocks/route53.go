package mocks

import (
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/stretchr/testify/mock"
)

type Route53 struct {
	mock.Mock
}

func (c *Route53) ListHostedZones(*route53.ListHostedZonesInput) (*route53.ListHostedZonesOutput, error) {
	args := c.Called()
	return args.Get(0).(*route53.ListHostedZonesOutput), args.Error(1)
}

func (c *Route53) ListResourceRecordSets(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error) {
	args := c.Called()
	return args.Get(0).(*route53.ListResourceRecordSetsOutput), args.Error(1)
}
