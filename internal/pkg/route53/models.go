package r53

import (
	"github.com/aws/aws-sdk-go/service/route53"
)

// Zones hosted by a DNS provider
type Zones struct {
	Public  map[string]string
	Private map[string]string
}

// Client interface
type Client interface {
	ListHostedZones(*route53.ListHostedZonesInput) (*route53.ListHostedZonesOutput, error)
	ListResourceRecordSets(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error)
}

// Records represents a zonefile content
type Records struct {
	SOA   Record
	NS    []Record
	MX    []Record
	A     []Record
	AAAA  []Record
	CNAME []Record
	TXT   []Record
	SRV   []Record
	PTR   []Record
	SPF   []Record
	NAPTR []Record
	CAA   []Record
}

// Record is a single DNS record
type Record struct {
	Name  string
	Value []string
	TTL   int64
	Alias bool
}
