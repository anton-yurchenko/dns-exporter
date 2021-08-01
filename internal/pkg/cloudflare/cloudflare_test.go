package cf_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"

	cf "dns-exporter/internal/pkg/cloudflare"

	"dns-exporter/mocks"

	"github.com/cloudflare/cloudflare-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestFetch(t *testing.T) {
	c := mocks.Cloudflare{}

	z := cf.Zones{
		Public: make(map[string]string),
	}

	errs := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(1)

	reply := []cloudflare.Zone{
		{
			ID:   "1",
			Name: "domain1.com",
		},
		{
			ID:   "2",
			Name: "domain2.com",
		},
	}

	c.On("ListZones").Return(reply, nil).Once()
	z.Fetch(&c, errs, &wg)

	// test: public zones
	expected := cf.Zones{
		Public: map[string]string{
			reply[0].Name: reply[0].ID,
			reply[1].Name: reply[1].ID,
		},
	}

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

	c := mocks.HTTP{}

	z := cf.Zones{
		Public: map[string]string{
			"domain1.com": "1",
			"domain2.com": "2",
		},
	}

	zonefiles := map[string]string{
		z.Public["domain1.com"]: `;;
;; Domain:     domain1.com.
;; Exported:   2019-10-19 18:27:26
;;
;; This file is intended for use for informational and archival
;; purposes ONLY and MUST be edited before use on a production
;; DNS server.  In particular, you must:
;;   -- update the SOA record with the correct authoritative name server
;;   -- update the SOA record with the contact e-mail address information
;;   -- update the NS record(s) with the authoritative name servers for this domain.
;;
;; For further information, please consult the BIND documentation
;; located on the following website:
;;
;; http://www.isc.org/
;;
;; And RFC 1035:
;;
;; http://www.ietf.org/rfc/rfc1035.txt
;;
;; Please note that we do NOT offer technical support for any use
;; of this zone data, the BIND name server, or any other third-party
;; DNS software.
;;
;; Use at your own risk.

;; SOA Record
domain1.com.	3600	IN	SOA	domain1.com. root.domain1.com. 2032317624 7200 3600 86400 3600

;; A Records
*.domain1.com.	1	IN	A	1.2.3.4
domain1.com.	1	IN	A	1.2.3.4

;; CNAME Records
www.domain1.com.	1	IN	CNAME	domain1.com.`,
		z.Public["domain2.com"]: `;;
;; Domain:     domain2.com.
;; Exported:   2019-10-19 18:27:26
;;
;; This file is intended for use for informational and archival
;; purposes ONLY and MUST be edited before use on a production
;; DNS server.  In particular, you must:
;;   -- update the SOA record with the correct authoritative name server
;;   -- update the SOA record with the contact e-mail address information
;;   -- update the NS record(s) with the authoritative name servers for this domain.
;;
;; For further information, please consult the BIND documentation
;; located on the following website:
;;
;; http://www.isc.org/
;;
;; And RFC 1035:
;;
;; http://www.ietf.org/rfc/rfc1035.txt
;;
;; Please note that we do NOT offer technical support for any use
;; of this zone data, the BIND name server, or any other third-party
;; DNS software.
;;
;; Use at your own risk.

;; SOA Record
domain2.com.	3600	IN	SOA	domain2.com. root.domain2.com. 2032317624 7200 3600 86400 3600

;; A Records
*.domain2.com.	1	IN	A	1.2.3.4
domain2.com.	1	IN	A	1.2.3.4

;; CNAME Records
www.domain2.com.	1	IN	CNAME	domain2.com.`,
	}

	expected := map[string]string{
		z.Public["domain1.com"]: `;; SOA Record
domain1.com.	3600	IN	SOA	domain1.com. root.domain1.com. 1 7200 3600 86400 3600

;; A Records
*.domain1.com.	1	IN	A	1.2.3.4
domain1.com.	1	IN	A	1.2.3.4

;; CNAME Records
www.domain1.com.	1	IN	CNAME	domain1.com.`,
		z.Public["domain2.com"]: `;; SOA Record
domain2.com.	3600	IN	SOA	domain2.com. root.domain2.com. 1 7200 3600 86400 3600

;; A Records
*.domain2.com.	1	IN	A	1.2.3.4
domain2.com.	1	IN	A	1.2.3.4

;; CNAME Records
www.domain2.com.	1	IN	CNAME	domain2.com.`,
	}

	fs := afero.NewMemMapFs()

	errs := make(chan error, len(zonefiles))

	var wg sync.WaitGroup
	wg.Add(len(zonefiles))

	for i, z := range zonefiles {
		request, _ := http.NewRequest("GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/export", i), nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("X-Auth-Email", "")
		request.Header.Add("X-Auth-Key", "")

		response := http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(z)),
		}

		c.On("Do", request).Return(&response, nil).Once()
	}

	z.Export(&c, 0, errs, &wg, "./", fs)

	err := <-errs
	if err != nil {
		t.Fatal("\nEXPECTED error: \n<nil>\n\nGOT error:", err)
	}

	for d, i := range z.Public {
		content, err := afero.ReadFile(fs, fmt.Sprintf("./CloudFlare/%s.txt", strings.TrimSuffix(strings.Replace(d, ".", "-", 1), "-")))
		if err != nil {
			t.Fatal("error reading exported zonefile:", err)
		}

		if !reflect.DeepEqual(expected[i], string(content)) {
			t.Errorf("\nEXPECTED content: \n%+v\n\nGOT content: \n%+v\n\n", expected[i], string(content))
		}
	}
}

func TestExportZone(t *testing.T) {
	c := mocks.HTTP{}

	// email := "user@domain.com"
	// token := "123abc"
	id := "123ab14fd747995cccc52a23b4ccc482"

	request, _ := http.NewRequest("GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/export", id), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Auth-Email", "")
	request.Header.Add("X-Auth-Key", "")

	text := `;;
;; Domain:     domain.com.
;; Exported:   2019-10-19 18:27:26
;;
;; This file is intended for use for informational and archival
;; purposes ONLY and MUST be edited before use on a production
;; DNS server.  In particular, you must:
;;   -- update the SOA record with the correct authoritative name server
;;   -- update the SOA record with the contact e-mail address information
;;   -- update the NS record(s) with the authoritative name servers for this domain.
;;
;; For further information, please consult the BIND documentation
;; located on the following website:
;;
;; http://www.isc.org/
;;
;; And RFC 1035:
;;
;; http://www.ietf.org/rfc/rfc1035.txt
;;
;; Please note that we do NOT offer technical support for any use
;; of this zone data, the BIND name server, or any other third-party
;; DNS software.
;;
;; Use at your own risk.

;; SOA Record
domain.com.	3600	IN	SOA	domain.com. root.domain.com. 2032317624 7200 3600 86400 3600

;; A Records
*.domain.com.	1	IN	A	1.2.3.4
domain.com.	1	IN	A	1.2.3.4

;; CNAME Records
www.domain.com.	1	IN	CNAME	domain.com.`

	body := bytes.NewBufferString(text)

	content := bytes.NewBufferString(text)

	response := http.Response{
		Body: ioutil.NopCloser(body),
	}

	c.On("Do", request).Return(&response, nil).Once()

	responce, err := cf.ExportZone(&c, id)

	if err != nil {
		t.Fatal("\nEXPECTED error: \n<nil>\n\nGOT error:", err)
	}

	expected, err := ioutil.ReadAll(content)

	if err != nil {
		t.Fatal("error reading responce body into []byte")
	}

	if !reflect.DeepEqual([]byte(expected), responce) {
		t.Errorf("\nEXPECTED content: \n%+v\n\nGOT content: \n%+v\n\n", string(expected), string(responce))
	}
}
