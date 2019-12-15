package r53_test

import (
	"bytes"
	"reflect"
	"testing"

	r53 "github.com/anton-yurchenko/dns-exporter/internal/pkg/route53"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

func TestAppend(t *testing.T) {
	suite := []map[*route53.ListResourceRecordSetsOutput]*r53.Records{}

	// SOA
	suite = append(suite, map[*route53.ListResourceRecordSetsOutput]*r53.Records{
		&route53.ListResourceRecordSetsOutput{
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
			},
		}: &r53.Records{
			SOA: r53.Record{
				Name:  "domain.com.",
				Value: []string{"ns-265.awsdns-33.com. awsdns-hostmaster.amazon.com. 1 7200 900 1209600 86400"},
				Alias: false,
				TTL:   900,
			},
		},
	})

	// NS
	suite = append(suite, map[*route53.ListResourceRecordSetsOutput]*r53.Records{
		&route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []*route53.ResourceRecordSet{
				&route53.ResourceRecordSet{
					Name: aws.String("domain.com."),
					Type: aws.String("NS"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("ns-1800.awsdns-33.co.uk."),
						},
					},
					TTL: aws.Int64(172800),
				},
				&route53.ResourceRecordSet{
					Name: aws.String("domain.com."),
					Type: aws.String("NS"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("ns3.amazon.net"),
						},
					},
					TTL: aws.Int64(300),
				},
			},
		}: &r53.Records{
			NS: []r53.Record{
				r53.Record{
					Name:  "domain.com.",
					Value: []string{"ns-1800.awsdns-33.co.uk."},
					Alias: false,
					TTL:   172800,
				},
				r53.Record{
					Name:  "domain.com.",
					Value: []string{"ns3.amazon.net"},
					Alias: false,
					TTL:   300,
				},
			},
		},
	})

	// MX
	suite = append(suite, map[*route53.ListResourceRecordSetsOutput]*r53.Records{
		&route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []*route53.ResourceRecordSet{
				&route53.ResourceRecordSet{
					Name: aws.String("mail1.domain.com."),
					Type: aws.String("MX"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("10 mailserver1.google.com"),
						},
					},
					TTL: aws.Int64(300),
				},
				&route53.ResourceRecordSet{
					Name: aws.String("mail2.domain.com."),
					Type: aws.String("MX"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("20 mailserver2.google.com"),
						},
					},
					TTL: aws.Int64(300),
				},
			},
		}: &r53.Records{
			MX: []r53.Record{
				r53.Record{
					Name:  "mail1.domain.com.",
					Value: []string{"10 mailserver1.google.com"},
					Alias: false,
					TTL:   300,
				},
				r53.Record{
					Name:  "mail2.domain.com.",
					Value: []string{"20 mailserver2.google.com"},
					Alias: false,
					TTL:   300,
				},
			},
		},
	})

	// A
	suite = append(suite, map[*route53.ListResourceRecordSetsOutput]*r53.Records{
		&route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []*route53.ResourceRecordSet{
				// Alias
				&route53.ResourceRecordSet{
					Name: aws.String("a-alias.domain.com."),
					Type: aws.String("A"),
					AliasTarget: &route53.AliasTarget{
						DNSName: aws.String("a0123456789abcdef.awsglobalaccelerator.com."),
					},
					TTL: aws.Int64(0),
				},
				// Round Robin
				&route53.ResourceRecordSet{
					Name: aws.String("a-multiple.domain.com."),
					Type: aws.String("A"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("192.168.1.51"),
						},
					},
					TTL: aws.Int64(300),
				},
				// Standard
				&route53.ResourceRecordSet{
					Name: aws.String("a-multiple.domain.com."),
					Type: aws.String("A"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("192.168.1.52"),
						},
					},
					TTL: aws.Int64(300),
				},
				// Low TTL
				&route53.ResourceRecordSet{
					Name: aws.String("low-ttl.domain.com."),
					Type: aws.String("A"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("192.168.1.53"),
						},
					},
					TTL: aws.Int64(60),
				},
			},
		}: &r53.Records{
			A: []r53.Record{
				r53.Record{
					Name:  "a-alias.domain.com.",
					Value: []string{"a0123456789abcdef.awsglobalaccelerator.com."},
					Alias: true,
					TTL:   0,
				},
				r53.Record{
					Name:  "a-multiple.domain.com.",
					Value: []string{"192.168.1.51"},
					Alias: false,
					TTL:   300,
				},
				r53.Record{
					Name:  "a-multiple.domain.com.",
					Value: []string{"192.168.1.52"},
					Alias: false,
					TTL:   300,
				},
				r53.Record{
					Name:  "low-ttl.domain.com.",
					Value: []string{"192.168.1.53"},
					Alias: false,
					TTL:   60,
				},
			},
		},
	})

	// AAAA
	suite = append(suite, map[*route53.ListResourceRecordSetsOutput]*r53.Records{
		&route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []*route53.ResourceRecordSet{
				&route53.ResourceRecordSet{
					Name: aws.String("ipv6.domain.com."),
					Type: aws.String("AAAA"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("fe80:0:0:0:202:b3ff:fe1e:8329"),
						},
					},
					TTL: aws.Int64(300),
				},
			},
		}: &r53.Records{
			AAAA: []r53.Record{
				r53.Record{
					Name:  "ipv6.domain.com.",
					Value: []string{"fe80:0:0:0:202:b3ff:fe1e:8329"},
					Alias: false,
					TTL:   300,
				},
			},
		},
	})

	// CNAME
	suite = append(suite, map[*route53.ListResourceRecordSetsOutput]*r53.Records{
		&route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []*route53.ResourceRecordSet{
				&route53.ResourceRecordSet{
					Name: aws.String("website.domain.com."),
					Type: aws.String("CNAME"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("www.domain.com"),
						},
					},
					TTL: aws.Int64(300),
				},
				// Alias
				&route53.ResourceRecordSet{
					Name: aws.String("www.domain.com."),
					Type: aws.String("CNAME"),
					AliasTarget: &route53.AliasTarget{
						DNSName: aws.String("c0123456789abcdef.awsglobalaccelerator.com."),
					},
					TTL: aws.Int64(0),
				},
			},
		}: &r53.Records{
			CNAME: []r53.Record{
				r53.Record{
					Name:  "website.domain.com.",
					Value: []string{"www.domain.com"},
					Alias: false,
					TTL:   300,
				},
				r53.Record{
					Name:  "www.domain.com.",
					Value: []string{"c0123456789abcdef.awsglobalaccelerator.com."},
					Alias: true,
					TTL:   0,
				},
			},
		},
	})

	// TXT
	suite = append(suite, map[*route53.ListResourceRecordSetsOutput]*r53.Records{
		&route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []*route53.ResourceRecordSet{
				&route53.ResourceRecordSet{
					Name: aws.String("txt.domain.com."),
					Type: aws.String("TXT"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("some string here"),
						},
					},
					TTL: aws.Int64(300),
				},
			},
		}: &r53.Records{
			TXT: []r53.Record{
				r53.Record{
					Name:  "txt.domain.com.",
					Value: []string{"some string here"},
					Alias: false,
					TTL:   300,
				},
			},
		},
	})

	// SRV
	suite = append(suite, map[*route53.ListResourceRecordSetsOutput]*r53.Records{
		&route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []*route53.ResourceRecordSet{
				&route53.ResourceRecordSet{
					Name: aws.String("dc.domain.com."),
					Type: aws.String("SRV"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("1 10 5269 controller.domain.com."),
						},
					},
					TTL: aws.Int64(300),
				},
			},
		}: &r53.Records{
			SRV: []r53.Record{
				r53.Record{
					Name:  "dc.domain.com.",
					Value: []string{"1 10 5269 controller.domain.com."},
					Alias: false,
					TTL:   300,
				},
			},
		},
	})

	// PTR
	suite = append(suite, map[*route53.ListResourceRecordSetsOutput]*r53.Records{
		&route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []*route53.ResourceRecordSet{
				&route53.ResourceRecordSet{
					Name: aws.String("pointer.domain.com."),
					Type: aws.String("PTR"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("www.domain.com"),
						},
					},
					TTL: aws.Int64(300),
				},
			},
		}: &r53.Records{
			PTR: []r53.Record{
				r53.Record{
					Name:  "pointer.domain.com.",
					Value: []string{"www.domain.com"},
					Alias: false,
					TTL:   300,
				},
			},
		},
	})

	// SPF
	suite = append(suite, map[*route53.ListResourceRecordSetsOutput]*r53.Records{
		&route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []*route53.ResourceRecordSet{
				&route53.ResourceRecordSet{
					Name: aws.String("policy.domain.com."),
					Type: aws.String("SPF"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("v=spf1 ip4:192.168.0.0/16-all"),
						},
					},
					TTL: aws.Int64(300),
				},
			},
		}: &r53.Records{
			SPF: []r53.Record{
				r53.Record{
					Name:  "policy.domain.com.",
					Value: []string{"v=spf1 ip4:192.168.0.0/16-all"},
					Alias: false,
					TTL:   300,
				},
			},
		},
	})

	// NAPTR
	suite = append(suite, map[*route53.ListResourceRecordSetsOutput]*r53.Records{
		&route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []*route53.ResourceRecordSet{
				&route53.ResourceRecordSet{
					Name: aws.String("name-auth.domain.com."),
					Type: aws.String("NAPTR"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("100 100 \"U\" \"\" \"!^.*$!sip:info@bar.example.com!\" ."),
						},
					},
					TTL: aws.Int64(300),
				},
			},
		}: &r53.Records{
			NAPTR: []r53.Record{
				r53.Record{
					Name:  "name-auth.domain.com.",
					Value: []string{"100 100 \"U\" \"\" \"!^.*$!sip:info@bar.example.com!\" ."},
					Alias: false,
					TTL:   300,
				},
			},
		},
	})

	// CAA
	suite = append(suite, map[*route53.ListResourceRecordSetsOutput]*r53.Records{
		&route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []*route53.ResourceRecordSet{
				&route53.ResourceRecordSet{
					Name: aws.String("ca1.domain.com."),
					Type: aws.String("CAA"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("0 issue \"caa.example.com\""),
						},
					},
					TTL: aws.Int64(300),
				},
				&route53.ResourceRecordSet{
					Name: aws.String("ca2.domain.com."),
					Type: aws.String("CAA"),
					ResourceRecords: []*route53.ResourceRecord{
						&route53.ResourceRecord{
							Value: aws.String("issuewild \";\""),
						},
					},
					TTL: aws.Int64(300),
				},
				// Alias
				&route53.ResourceRecordSet{
					Name: aws.String("ca3.domain.com."),
					Type: aws.String("CAA"),
					AliasTarget: &route53.AliasTarget{
						DNSName: aws.String("b0123456789abcdef.awsglobalaccelerator.com."),
					},
					TTL: aws.Int64(0),
				},
			},
		}: &r53.Records{
			CAA: []r53.Record{
				r53.Record{
					Name:  "ca1.domain.com.",
					Value: []string{"0 issue \"caa.example.com\""},
					Alias: false,
					TTL:   300,
				},
				r53.Record{
					Name:  "ca2.domain.com.",
					Value: []string{"issuewild \";\""},
					Alias: false,
					TTL:   300,
				},
				r53.Record{
					Name:  "ca3.domain.com.",
					Value: []string{"b0123456789abcdef.awsglobalaccelerator.com."},
					Alias: true,
					TTL:   0,
				},
			},
		},
	})

	for _, test := range suite {
		for input, expected := range test {
			got := r53.Records{}

			err := got.Append(input)

			if err != nil {
				t.Fatal("\nEXPECTED error: \n<nil>\n\nGOT error:", err)
			}

			if !reflect.DeepEqual(&got, expected) {
				t.Errorf("\nEXPECTED provider zones: \n%+v\n\nGOT provider zones: \n%+v\n\n", expected, &got)
			}
		}
	}
}

func TestConvertToZonefile(t *testing.T) {

	records := r53.Records{
		SOA: r53.Record{
			Name:  "domain.com.",
			Value: []string{"ns-265.awsdns-33.com. awsdns-hostmaster.amazon.com. 1 7200 900 1209600 86400"},
			Alias: false,
			TTL:   900,
		},
		NS: []r53.Record{
			r53.Record{
				Name:  "domain.com.",
				Value: []string{"ns-1800.awsdns-33.co.uk."},
				Alias: false,
				TTL:   172800,
			},
			r53.Record{
				Name:  "domain.com.",
				Value: []string{"ns3.amazon.net"},
				Alias: false,
				TTL:   300,
			},
		},
		MX: []r53.Record{
			r53.Record{
				Name:  "mail1.domain.com.",
				Value: []string{"10 mailserver1.google.com"},
				Alias: false,
				TTL:   300,
			},
			r53.Record{
				Name:  "mail2.domain.com.",
				Value: []string{"20 mailserver2.google.com"},
				Alias: false,
				TTL:   300,
			},
		},
		A: []r53.Record{
			r53.Record{
				Name:  "a-alias.domain.com.",
				Value: []string{"a0123456789abcdef.awsglobalaccelerator.com."},
				Alias: true,
				TTL:   0,
			},
			r53.Record{
				Name:  "a-multiple.domain.com.",
				Value: []string{"192.168.1.51"},
				Alias: false,
				TTL:   300,
			},
			r53.Record{
				Name:  "a-multiple.domain.com.",
				Value: []string{"192.168.1.52"},
				Alias: false,
				TTL:   300,
			},
			r53.Record{
				Name:  "low-ttl.domain.com.",
				Value: []string{"192.168.1.53"},
				Alias: false,
				TTL:   60,
			},
		},
		AAAA: []r53.Record{
			r53.Record{
				Name:  "ipv6.domain.com.",
				Value: []string{"fe80:0:0:0:202:b3ff:fe1e:8329"},
				Alias: false,
				TTL:   300,
			},
		},
		CNAME: []r53.Record{
			r53.Record{
				Name:  "website.domain.com.",
				Value: []string{"www.domain.com"},
				Alias: false,
				TTL:   300,
			},
			r53.Record{
				Name:  "www.domain.com.",
				Value: []string{"c0123456789abcdef.awsglobalaccelerator.com."},
				Alias: true,
				TTL:   0,
			},
		},
		TXT: []r53.Record{
			r53.Record{
				Name:  "txt.domain.com.",
				Value: []string{"some string here"},
				Alias: false,
				TTL:   300,
			},
		},
		SRV: []r53.Record{
			r53.Record{
				Name:  "dc.domain.com.",
				Value: []string{"1 10 5269 controller.domain.com."},
				Alias: false,
				TTL:   300,
			},
		},
		PTR: []r53.Record{
			r53.Record{
				Name:  "pointer.domain.com.",
				Value: []string{"www.domain.com"},
				Alias: false,
				TTL:   300,
			},
		},
		SPF: []r53.Record{
			r53.Record{
				Name:  "policy.domain.com.",
				Value: []string{"v=spf1 ip4:192.168.0.0/16-all"},
				Alias: false,
				TTL:   300,
			},
		},
		NAPTR: []r53.Record{
			r53.Record{
				Name:  "name-auth.domain.com.",
				Value: []string{"100 100 \"U\" \"\" \"!^.*$!sip:info@bar.example.com!\" ."},
				Alias: false,
				TTL:   300,
			},
		},
		CAA: []r53.Record{
			r53.Record{
				Name:  "ca1.domain.com.",
				Value: []string{"0 issue \"caa.example.com\""},
				Alias: false,
				TTL:   300,
			},
			r53.Record{
				Name:  "ca2.domain.com.",
				Value: []string{"issuewild \";\""},
				Alias: false,
				TTL:   300,
			},
			r53.Record{
				Name:  "ca3.domain.com.",
				Value: []string{"b0123456789abcdef.awsglobalaccelerator.com."},
				Alias: true,
				TTL:   0,
			},
		},
	}

	expected := bytes.NewBuffer([]byte(`;; SOA Record
domain.com.	900	IN	SOA	ns-265.awsdns-33.com. awsdns-hostmaster.amazon.com. 1 7200 900 1209600 86400

;; NS Records
domain.com.	172800	IN	NS	ns-1800.awsdns-33.co.uk.
domain.com.	300	IN	NS	ns3.amazon.net

;; MX Records
mail1.domain.com.	300	IN	MX	10 mailserver1.google.com
mail2.domain.com.	300	IN	MX	20 mailserver2.google.com

;; A Records
a-multiple.domain.com.	300	IN	A	192.168.1.51
a-multiple.domain.com.	300	IN	A	192.168.1.52
low-ttl.domain.com.	60	IN	A	192.168.1.53

;; AAAA Records
ipv6.domain.com.	300	IN	AAAA	fe80:0:0:0:202:b3ff:fe1e:8329

;; CNAME Records
website.domain.com.	300	IN	CNAME	www.domain.com

;; TXT Records
txt.domain.com.	300	IN	TXT	some string here

;; SRV Records
dc.domain.com.	300	IN	SRV	1 10 5269 controller.domain.com.

;; PTR Records
pointer.domain.com.	300	IN	PTR	www.domain.com

;; SPF Records
policy.domain.com.	300	IN	SPF	v=spf1 ip4:192.168.0.0/16-all

;; NAPTR Records
name-auth.domain.com.	300	IN	NAPTR	100 100 "U" "" "!^.*$!sip:info@bar.example.com!" .

;; CAA Records
ca1.domain.com.	300	IN	CAA	0 issue "caa.example.com"
ca2.domain.com.	300	IN	CAA	issuewild ";"

;; Route53 Alias Records
a-alias.domain.com.	0	IN	A	a0123456789abcdef.awsglobalaccelerator.com.
www.domain.com.	0	IN	CNAME	c0123456789abcdef.awsglobalaccelerator.com.
ca3.domain.com.	0	IN	CAA	b0123456789abcdef.awsglobalaccelerator.com.
`))

	zonefile, err := records.ConvertToZonefile()

	if err != nil {
		t.Fatal("\nEXPECTED error: \n<nil>\n\nGOT error:", err)
	}

	if !reflect.DeepEqual(expected.Bytes(), zonefile.Bytes()) {
		t.Errorf("\nEXPECTED zonefile: \n'%+v'\n\nGOT zonefile: \n'%+v'\n\n", expected.String(), zonefile.String())
	}
}

func TestFormat(t *testing.T) {
	suite := []map[*[]r53.Record]map[string]*bytes.Buffer{}

	// SOA
	suite = append(suite, map[*[]r53.Record]map[string]*bytes.Buffer{
		&[]r53.Record{
			r53.Record{
				Name:  "domain.com.",
				Value: []string{"ns-265.awsdns-33.com. awsdns-hostmaster.amazon.com. 1 7200 900 1209600 86400"},
				Alias: false,
				TTL:   900,
			},
		}: map[string]*bytes.Buffer{
			"type": bytes.NewBuffer([]byte("SOA")),
			"records": bytes.NewBuffer([]byte(`;; SOA Records
domain.com.	900	IN	SOA	ns-265.awsdns-33.com. awsdns-hostmaster.amazon.com. 1 7200 900 1209600 86400

`)),
			"aliases": bytes.NewBuffer([]byte(``)),
		}})

	// NS
	suite = append(suite, map[*[]r53.Record]map[string]*bytes.Buffer{
		&[]r53.Record{
			r53.Record{
				Name:  "domain.com.",
				Value: []string{"ns-1800.awsdns-33.co.uk."},
				Alias: false,
				TTL:   172800,
			},
			r53.Record{
				Name:  "domain.com.",
				Value: []string{"ns3.amazon.net"},
				Alias: false,
				TTL:   300,
			},
		}: map[string]*bytes.Buffer{
			"type": bytes.NewBuffer([]byte("NS")),
			"records": bytes.NewBuffer([]byte(`;; NS Records
domain.com.	172800	IN	NS	ns-1800.awsdns-33.co.uk.
domain.com.	300	IN	NS	ns3.amazon.net

`)),
			"aliases": bytes.NewBuffer([]byte(``)),
		}})

	// MX
	suite = append(suite, map[*[]r53.Record]map[string]*bytes.Buffer{
		&[]r53.Record{
			r53.Record{
				Name:  "mail1.domain.com.",
				Value: []string{"10 mailserver1.google.com"},
				Alias: false,
				TTL:   300,
			},
			r53.Record{
				Name:  "mail2.domain.com.",
				Value: []string{"20 mailserver2.google.com"},
				Alias: false,
				TTL:   300,
			},
		}: map[string]*bytes.Buffer{
			"type": bytes.NewBuffer([]byte("MX")),
			"records": bytes.NewBuffer([]byte(`;; MX Records
mail1.domain.com.	300	IN	MX	10 mailserver1.google.com
mail2.domain.com.	300	IN	MX	20 mailserver2.google.com

`)),
			"aliases": bytes.NewBuffer([]byte(``)),
		}})

	// A
	suite = append(suite, map[*[]r53.Record]map[string]*bytes.Buffer{
		&[]r53.Record{
			r53.Record{
				Name:  "a-alias.domain.com.",
				Value: []string{"a0123456789abcdef.awsglobalaccelerator.com."},
				Alias: true,
				TTL:   0,
			},
			r53.Record{
				Name:  "a-multiple.domain.com.",
				Value: []string{"192.168.1.51"},
				Alias: false,
				TTL:   300,
			},
			r53.Record{
				Name:  "a-multiple.domain.com.",
				Value: []string{"192.168.1.52"},
				Alias: false,
				TTL:   300,
			},
			r53.Record{
				Name:  "low-ttl.domain.com.",
				Value: []string{"192.168.1.53"},
				Alias: false,
				TTL:   60,
			},
		}: map[string]*bytes.Buffer{
			"type": bytes.NewBuffer([]byte("A")),
			"records": bytes.NewBuffer([]byte(`;; A Records
a-multiple.domain.com.	300	IN	A	192.168.1.51
a-multiple.domain.com.	300	IN	A	192.168.1.52
low-ttl.domain.com.	60	IN	A	192.168.1.53

`)),
			"aliases": bytes.NewBuffer([]byte(`;; Route53 Alias Records
a-alias.domain.com.	0	IN	A	a0123456789abcdef.awsglobalaccelerator.com.
`)),
		}})

	// AAAA
	suite = append(suite, map[*[]r53.Record]map[string]*bytes.Buffer{
		&[]r53.Record{
			r53.Record{
				Name:  "ipv6.domain.com.",
				Value: []string{"fe80:0:0:0:202:b3ff:fe1e:8329"},
				Alias: false,
				TTL:   300,
			},
		}: map[string]*bytes.Buffer{
			"type": bytes.NewBuffer([]byte("AAAA")),
			"records": bytes.NewBuffer([]byte(`;; AAAA Records
ipv6.domain.com.	300	IN	AAAA	fe80:0:0:0:202:b3ff:fe1e:8329

`)),
			"aliases": bytes.NewBuffer([]byte(``)),
		}})

	// CNAME
	suite = append(suite, map[*[]r53.Record]map[string]*bytes.Buffer{
		&[]r53.Record{
			r53.Record{
				Name:  "website.domain.com.",
				Value: []string{"www.domain.com"},
				Alias: false,
				TTL:   300,
			},
			r53.Record{
				Name:  "www.domain.com.",
				Value: []string{"c0123456789abcdef.awsglobalaccelerator.com."},
				Alias: true,
				TTL:   0,
			},
		}: map[string]*bytes.Buffer{
			"type": bytes.NewBuffer([]byte("CNAME")),
			"records": bytes.NewBuffer([]byte(`;; CNAME Records
website.domain.com.	300	IN	CNAME	www.domain.com

`)),
			"aliases": bytes.NewBuffer([]byte(`;; Route53 Alias Records
www.domain.com.	0	IN	CNAME	c0123456789abcdef.awsglobalaccelerator.com.
`)),
		}})

	// TXT
	suite = append(suite, map[*[]r53.Record]map[string]*bytes.Buffer{
		&[]r53.Record{
			r53.Record{
				Name:  "txt.domain.com.",
				Value: []string{"some string here"},
				Alias: false,
				TTL:   300,
			},
		}: map[string]*bytes.Buffer{
			"type": bytes.NewBuffer([]byte("TXT")),
			"records": bytes.NewBuffer([]byte(`;; TXT Records
txt.domain.com.	300	IN	TXT	some string here

`)),
			"aliases": bytes.NewBuffer([]byte(``)),
		}})

	// SRV
	suite = append(suite, map[*[]r53.Record]map[string]*bytes.Buffer{
		&[]r53.Record{
			r53.Record{
				Name:  "dc.domain.com.",
				Value: []string{"1 10 5269 controller.domain.com."},
				Alias: false,
				TTL:   300,
			},
		}: map[string]*bytes.Buffer{
			"type": bytes.NewBuffer([]byte("SRV")),
			"records": bytes.NewBuffer([]byte(`;; SRV Records
dc.domain.com.	300	IN	SRV	1 10 5269 controller.domain.com.

`)),
			"aliases": bytes.NewBuffer([]byte(``)),
		}})

	// PTRC
	suite = append(suite, map[*[]r53.Record]map[string]*bytes.Buffer{
		&[]r53.Record{
			r53.Record{
				Name:  "pointer.domain.com.",
				Value: []string{"www.domain.com"},
				Alias: false,
				TTL:   300,
			},
		}: map[string]*bytes.Buffer{
			"type": bytes.NewBuffer([]byte("PTR")),
			"records": bytes.NewBuffer([]byte(`;; PTR Records
pointer.domain.com.	300	IN	PTR	www.domain.com

`)),
			"aliases": bytes.NewBuffer([]byte(``)),
		}})

	// SPF
	suite = append(suite, map[*[]r53.Record]map[string]*bytes.Buffer{
		&[]r53.Record{
			r53.Record{
				Name:  "policy.domain.com.",
				Value: []string{"v=spf1 ip4:192.168.0.0/16-all"},
				Alias: false,
				TTL:   300,
			},
		}: map[string]*bytes.Buffer{
			"type": bytes.NewBuffer([]byte("SPF")),
			"records": bytes.NewBuffer([]byte(`;; SPF Records
policy.domain.com.	300	IN	SPF	v=spf1 ip4:192.168.0.0/16-all

`)),
			"aliases": bytes.NewBuffer([]byte(``)),
		}})

	// NAPTR
	suite = append(suite, map[*[]r53.Record]map[string]*bytes.Buffer{
		&[]r53.Record{
			r53.Record{
				Name:  "name-auth.domain.com.",
				Value: []string{"100 100 \"U\" \"\" \"!^.*$!sip:info@bar.example.com!\" ."},
				Alias: false,
				TTL:   300,
			},
		}: map[string]*bytes.Buffer{
			"type": bytes.NewBuffer([]byte("NAPTR")),
			"records": bytes.NewBuffer([]byte(`;; NAPTR Records
name-auth.domain.com.	300	IN	NAPTR	100 100 "U" "" "!^.*$!sip:info@bar.example.com!" .

`)),
			"aliases": bytes.NewBuffer([]byte(``)),
		}})

	// CAA
	suite = append(suite, map[*[]r53.Record]map[string]*bytes.Buffer{
		&[]r53.Record{
			r53.Record{
				Name:  "ca1.domain.com.",
				Value: []string{"0 issue \"caa.example.com\""},
				Alias: false,
				TTL:   300,
			},
			r53.Record{
				Name:  "ca2.domain.com.",
				Value: []string{"issuewild \";\""},
				Alias: false,
				TTL:   300,
			},
			r53.Record{
				Name:  "ca3.domain.com.",
				Value: []string{"b0123456789abcdef.awsglobalaccelerator.com."},
				Alias: true,
				TTL:   0,
			},
		}: map[string]*bytes.Buffer{
			"type": bytes.NewBuffer([]byte("CAA")),
			"records": bytes.NewBuffer([]byte(`;; CAA Records
ca1.domain.com.	300	IN	CAA	0 issue "caa.example.com"
ca2.domain.com.	300	IN	CAA	issuewild ";"

`)),
			"aliases": bytes.NewBuffer([]byte(`;; Route53 Alias Records
ca3.domain.com.	0	IN	CAA	b0123456789abcdef.awsglobalaccelerator.com.
`)),
		}})

	for _, test := range suite {
		for input, expected := range test {
			records := bytes.Buffer{}
			aliases := bytes.Buffer{}
			var marker bool

			err := r53.Format(&marker, input, &records, &aliases, expected["type"].String())

			if err != nil {
				t.Fatal("\nEXPECTED error: \n<nil>\n\nGOT error:", err)
			}

			if !reflect.DeepEqual(expected["records"].String(), records.String()) {
				t.Errorf("\nEXPECTED records: \n'%+v'\n\nGOT records: \n'%+v'\n\n", expected["records"].String(), records.String())
			}

			if !reflect.DeepEqual(expected["aliases"].String(), aliases.String()) {
				t.Errorf("\nEXPECTED aliases: \n'%+v'\n\nGOT aliases: \n'%+v'\n\n", expected["aliases"].String(), aliases.String())
			}
		}
	}
}
