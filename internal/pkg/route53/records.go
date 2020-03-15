package r53

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
)

// Append Record Sets to the Zone
func (r *Records) Append(set *route53.ListResourceRecordSetsOutput) error {
	s := set
	for _, record := range s.ResourceRecordSets {
		var val []string
		var rec Record
		if record.AliasTarget == nil {
			for _, v := range record.ResourceRecords {
				val = append(val, *v.Value)
			}

			rec.Name = *record.Name
			rec.Value = val
			rec.TTL = *record.TTL
		} else {
			rec.Name = *record.Name
			rec.Value = []string{*record.AliasTarget.DNSName}
			rec.Alias = true
		}

		// add record
		switch *record.Type {
		case "SOA":
			r.SOA = rec
		case "NS":
			r.NS = append(r.NS, rec)
		case "MX":
			r.MX = append(r.MX, rec)
		case "A":
			r.A = append(r.A, rec)
		case "AAAA":
			r.AAAA = append(r.AAAA, rec)
		case "CNAME":
			r.CNAME = append(r.CNAME, rec)
		case "TXT":
			r.TXT = append(r.TXT, rec)
		case "SRV":
			r.SRV = append(r.SRV, rec)
		case "PTR":
			r.PTR = append(r.PTR, rec)
		case "SPF":
			r.SPF = append(r.SPF, rec)
		case "NAPTR":
			r.NAPTR = append(r.NAPTR, rec)
		case "CAA":
			r.CAA = append(r.CAA, rec)
		default:
			return fmt.Errorf("not supported record type: %s", *record.Type)
		}
	}

	return nil
}

// ConvertToZonefile returns formatted zonefile
func (r *Records) ConvertToZonefile() (bytes.Buffer, error) {
	b := bytes.Buffer{}
	p := bytes.Buffer{}

	var aliasHeaderMarker bool

	// SOA
	_, err := b.WriteString(fmt.Sprintf(";; SOA Record\n%s\t%v\tIN\tSOA\t%s\n\n", r.SOA.Name, r.SOA.TTL, r.SOA.Value[0]))
	if err != nil {
		return bytes.Buffer{}, errors.Wrap(err, "error formatting records of type 'SOA'")
	}

	// NS
	if len(r.NS) > 0 {
		t := "NS"
		err := format(&aliasHeaderMarker, &r.NS, &b, &p, t)
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, fmt.Sprintf("error formatting records of type '%s'", t))
		}
	}

	// MX
	if len(r.MX) > 0 {
		t := "MX"
		err := format(&aliasHeaderMarker, &r.MX, &b, &p, t)
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, fmt.Sprintf("error formatting records of type '%s'", t))
		}
	}

	// A
	if len(r.A) > 0 {
		t := "A"
		err := format(&aliasHeaderMarker, &r.A, &b, &p, t)
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, fmt.Sprintf("error formatting records of type '%s'", t))
		}
	}

	// AAAA
	if len(r.AAAA) > 0 {
		t := "AAAA"
		err := format(&aliasHeaderMarker, &r.AAAA, &b, &p, t)
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, fmt.Sprintf("error formatting records of type '%s'", t))
		}
	}

	// CNAME
	if len(r.CNAME) > 0 {
		t := "CNAME"
		err := format(&aliasHeaderMarker, &r.CNAME, &b, &p, "CNAME")
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, fmt.Sprintf("error formatting records of type '%s'", t))
		}
	}

	// TXT
	if len(r.TXT) > 0 {
		t := "TXT"
		err := format(&aliasHeaderMarker, &r.TXT, &b, &p, "TXT")
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, fmt.Sprintf("error formatting records of type '%s'", t))
		}
	}

	// SRV
	if len(r.SRV) > 0 {
		t := "SRV"
		err := format(&aliasHeaderMarker, &r.SRV, &b, &p, "SRV")
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, fmt.Sprintf("error formatting records of type '%s'", t))
		}
	}

	// PTR
	if len(r.PTR) > 0 {
		t := "PTR"
		err := format(&aliasHeaderMarker, &r.PTR, &b, &p, "PTR")
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, fmt.Sprintf("error formatting records of type '%s'", t))
		}
	}

	// SPF
	if len(r.SPF) > 0 {
		t := "SPF"
		err := format(&aliasHeaderMarker, &r.SPF, &b, &p, "SPF")
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, fmt.Sprintf("error formatting records of type '%s'", t))
		}
	}

	// NAPTR
	if len(r.NAPTR) > 0 {
		t := "NAPTR"
		err := format(&aliasHeaderMarker, &r.NAPTR, &b, &p, "NAPTR")
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, fmt.Sprintf("error formatting records of type '%s'", t))
		}
	}

	// CAA
	if len(r.CAA) > 0 {
		t := "CAA"
		err := format(&aliasHeaderMarker, &r.CAA, &b, &p, "CAA")
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, fmt.Sprintf("error formatting records of type '%s'", t))
		}
	}

	_, err = b.WriteString(p.String())
	if err != nil {
		return bytes.Buffer{}, errors.Wrap(err, "error appending Route53 Aliases")
	}

	return b, nil
}

// format zonefile content
func format(marker *bool, r *[]Record, b, p *bytes.Buffer, t string) error {
	_, err := b.WriteString(fmt.Sprintf(";; %s Records\n", t))
	if err != nil {
		return errors.New("error formatting zonefile")
	}

	for _, record := range *r {
		if !record.Alias {
			for _, rec := range record.Value {
				_, err := b.WriteString(fmt.Sprintf("%s\t%v\tIN\t%s\t%s\n", record.Name, record.TTL, t, rec))
				if err != nil {
					return fmt.Errorf("error formatting records of type: %s", t)
				}
			}
		} else {
			if !*marker {
				_, err := p.WriteString(";; Route53 Alias Records\n")
				if err != nil {
					return errors.New("error adding alias title")
				}
				*marker = true
			}

			for _, rec := range record.Value {
				_, err := p.WriteString(fmt.Sprintf("%s\t%v\tIN\t%s\t%s\n", record.Name, record.TTL, t, rec))
				if err != nil {
					return fmt.Errorf("error formatting records of type: %s", t)
				}
			}
		}
	}

	_, err = b.WriteString("\n")
	if err != nil {
		return fmt.Errorf("error formatting records of type: %s", t)
	}
	return nil
}
