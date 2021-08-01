package r53_test

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"

	r53 "dns-exporter/internal/pkg/route53"
	"dns-exporter/mocks"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestFetch(t *testing.T) {
	c := mocks.Route53{}

	z := r53.Zones{
		Public:  make(map[string]string),
		Private: make(map[string]string),
	}

	errs := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(1)

	reply1 := &route53.ListHostedZonesOutput{
		IsTruncated: aws.Bool(true),
		HostedZones: []*route53.HostedZone{
			&route53.HostedZone{
				Config: &route53.HostedZoneConfig{
					PrivateZone: aws.Bool(false),
				},
				Id:   aws.String("/hostedzone/A1M9OJ3HY2SUQY"),
				Name: aws.String("domain1.com"),
			},
		},
	}

	reply2 := &route53.ListHostedZonesOutput{
		IsTruncated: aws.Bool(false),
		HostedZones: []*route53.HostedZone{
			&route53.HostedZone{
				Config: &route53.HostedZoneConfig{
					PrivateZone: aws.Bool(false),
				},
				Id:   aws.String("/hostedzone/B2M9OJ3HY2SUQY"),
				Name: aws.String("domain2.com"),
			},
			&route53.HostedZone{
				Config: &route53.HostedZoneConfig{
					PrivateZone: aws.Bool(true),
				},
				Id:   aws.String("/hostedzone/C3M9OJ3HY2SUQY"),
				Name: aws.String("domain1.local"),
			},
		},
	}

	expected := r53.Zones{
		Public: map[string]string{
			"domain1.com": "/hostedzone/A1M9OJ3HY2SUQY",
			"domain2.com": "/hostedzone/B2M9OJ3HY2SUQY",
		},
		Private: map[string]string{
			"domain1.local": "/hostedzone/C3M9OJ3HY2SUQY",
		},
	}

	c.On("ListHostedZones").Return(reply1, nil).Once()
	c.On("ListHostedZones").Return(reply2, nil).Once()

	z.Fetch(&c, errs, &wg)

	if !reflect.DeepEqual(z, expected) {
		t.Errorf("\nEXPECTED provider zones: \n%+v\n\nGOT provider zones: \n%+v\n\n", expected, z)
	}

	err := <-errs
	if err != nil {
		t.Fatal("\nEXPECTED error: \n<nil>\n\nGOT error:", err)
	}
}

func TestExport(t *testing.T) {
	log.SetLevel(log.ErrorLevel)

	c := mocks.Route53{}

	fs := afero.NewMemMapFs()

	z := r53.Zones{
		Public: map[string]string{
			"domain.com": "1",
		},
		Private: map[string]string{
			"local": "2",
		},
	}

	zonefiles := map[string]string{
		"1": `;; SOA Record
domain.com.	900	IN	SOA	ns-265.awsdns-33.com. awsdns-hostmaster.amazon.com. 1 7200 900 1209600 86400

;; A Records
domain.com.	300	IN	A	192.168.1.51

;; CNAME Records
www.domain.com.	300	IN	CNAME	domain.com

`,
		"2": `;; SOA Record
local.	900	IN	SOA	ns-265.awsdns-33.com. awsdns-hostmaster.amazon.com. 1 7200 900 1209600 86400

;; A Records
local.	1	IN	A	192.168.1.52

`}

	replies := map[string]*route53.ListResourceRecordSetsOutput{
		"1": &route53.ListResourceRecordSetsOutput{
			IsTruncated: aws.Bool(false),
			ResourceRecordSets: []*route53.ResourceRecordSet{
				&route53.ResourceRecordSet{
					Name: aws.String("domain.com."),
					Type: aws.String("SOA"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("ns-265.awsdns-33.com. awsdns-hostmaster.amazon.com. 1 7200 900 1209600 86400"),
						},
					},
					TTL: aws.Int64(900),
				},
				&route53.ResourceRecordSet{
					Name: aws.String("domain.com."),
					Type: aws.String("A"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("192.168.1.51"),
						},
					},
					TTL: aws.Int64(300),
				},
				&route53.ResourceRecordSet{
					Name: aws.String("www.domain.com."),
					Type: aws.String("CNAME"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("domain.com"),
						},
					},
					TTL: aws.Int64(300),
				},
			},
		},
		"2": &route53.ListResourceRecordSetsOutput{
			IsTruncated: aws.Bool(false),
			ResourceRecordSets: []*route53.ResourceRecordSet{
				&route53.ResourceRecordSet{
					Name: aws.String("local."),
					Type: aws.String("SOA"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("ns-265.awsdns-33.com. awsdns-hostmaster.amazon.com. 1 7200 900 1209600 86400"),
						},
					},
					TTL: aws.Int64(900),
				},
				&route53.ResourceRecordSet{
					Name: aws.String("local."),
					Type: aws.String("A"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("192.168.1.52"),
						},
					},
					TTL: aws.Int64(1),
				},
			},
		},
	}

	errs := make(chan error, len(zonefiles))

	var wg sync.WaitGroup
	wg.Add(len(zonefiles))

	c.On("ListResourceRecordSets").Return(replies["1"], nil).Once()
	c.On("ListResourceRecordSets").Return(replies["2"], nil).Once()

	z.Export(&c, 0, errs, &wg, ".", fs)

	err := <-errs
	if err != nil {
		t.Fatal("\nEXPECTED error: \n<nil>\n\nGOT error:", err)
	}

	for d, i := range z.Public {
		content, err := afero.ReadFile(fs, fmt.Sprintf("./Route53/Public/%s.txt", strings.TrimSuffix(strings.Replace(d, ".", "-", 1), "-")))
		if err != nil {
			t.Fatal("error reading exported zonefile:", err)
		}

		if !reflect.DeepEqual([]byte(zonefiles[i]), content) {
			t.Errorf("\nEXPECTED content: \n'%+v'\n\nGOT content: \n'%+v'\n\n", zonefiles[i], string(content))
		}
	}

	for d, i := range z.Private {
		content, err := afero.ReadFile(fs, fmt.Sprintf("./Route53/Private/%s.txt", strings.TrimSuffix(strings.Replace(d, ".", "-", 1), "-")))
		if err != nil {
			t.Fatal("error reading exported zonefile:", err)
		}

		if !reflect.DeepEqual([]byte(zonefiles[i]), content) {
			t.Errorf("\nEXPECTED content: \n'%+v'\n\nGOT content: \n'%+v'\n\n", zonefiles[i], string(content))
		}
	}
}

func TestGetRecords(t *testing.T) {
	c := mocks.Route53{}
	id := "/hostedzone/A1M9OJ3HY2SUQY"

	reply1 := &route53.ListResourceRecordSetsOutput{
		IsTruncated:    aws.Bool(true),
		NextRecordType: aws.String("AAAA"),
		NextRecordName: aws.String("aaaa.domain.com."),
		ResourceRecordSets: []*route53.ResourceRecordSet{
			&route53.ResourceRecordSet{
				Name: aws.String("a.domain.com."),
				Type: aws.String("A"),
				ResourceRecords: []*route53.ResourceRecord{
					&route53.ResourceRecord{
						Value: aws.String("192.168.1.51"),
					},
				},
				TTL: aws.Int64(300),
			},
		},
	}

	reply2 := &route53.ListResourceRecordSetsOutput{
		IsTruncated: aws.Bool(false),
		ResourceRecordSets: []*route53.ResourceRecordSet{
			&route53.ResourceRecordSet{
				Name: aws.String("aaaa.domain.com."),
				Type: aws.String("AAAA"),
				ResourceRecords: []*route53.ResourceRecord{
					&route53.ResourceRecord{
						Value: aws.String("fe80:0:0:0:202:b3ff:fe1e:8329"),
					},
				},
				TTL: aws.Int64(300),
			},
		},
	}

	expected := r53.Records{
		A: []r53.Record{
			r53.Record{
				Name:  "a.domain.com.",
				Value: []string{"192.168.1.51"},
				Alias: false,
				TTL:   300,
			},
		},
		AAAA: []r53.Record{
			r53.Record{
				Name:  "aaaa.domain.com.",
				Value: []string{"fe80:0:0:0:202:b3ff:fe1e:8329"},
				Alias: false,
				TTL:   300,
			},
		},
	}

	c.On("ListResourceRecordSets").Return(reply1, nil).Once()
	c.On("ListResourceRecordSets").Return(reply2, nil).Once()

	r, err := r53.GetRecords(id, &c)

	if err != nil {
		t.Fatal("\nEXPECTED error: \n<nil>\n\nGOT error:", err)
	}

	if !reflect.DeepEqual(expected, r) {
		t.Errorf("\nEXPECTED records: \n%+v\n\nGOT records: \n%+v\n\n", expected, r)
	}
}
