# dns-exporter
[![License](https://img.shields.io/github/license/anton-yurchenko/dns-exporter?style=flat-square)](LICENSE.md) [![Release](https://img.shields.io/github/v/release/anton-yurchenko/dns-exporter?style=flat-square)](https://github.com/anton-yurchenko/dns-exporter/releases/latest) [![Docker Build](https://img.shields.io/docker/cloud/build/antonyurchenko/dns-exporter?style=flat-square)](https://hub.docker.com/r/antonyurchenko/dns-exporter) [![Docker Pulls](https://img.shields.io/docker/pulls/antonyurchenko/dns-exporter?style=flat-square)](https://hub.docker.com/r/antonyurchenko/dns-exporter)

**DNS-EXPORTER** is DNS Archiving/Documentation tool.  

## Features:
- Export all DNS records in Zonefile-like format.  
- Export to local/remote Git repository allowing easy tracking of changes.  
- Supported DNS providers: **CloudFlare, Route53**
- Supported public / private zones.  

## Example Export:
```
;; SOA Record
domain.com.	900	IN	SOA	ns-245.awsdns-33.com. awsdns-hostmaster.amazon.com. 1 7200 900 1209600 86400

;; NS Records
domain.com.	172800	IN	NS	ns-1800.awsdns-33.co.uk.
domain.com.	172800	IN	NS	ns-245.awsdns-33.com.
domain.com.	172800	IN	NS	ns-1352.awsdns-33.org.
domain.com.	172800	IN	NS	ns-736.awsdns-33.net.
dns.domain.com.	300	IN	NS	ns2.amazon.org
dns.domain.com.	300	IN	NS	ns1.amazon.com
dns.domain.com.	300	IN	NS	ns3.amazon.net
dns.domain.com.	300	IN	NS	ns4.amazon.co.uk

;; MX Records
mail.domain.com.	300	IN	MX	10 mailserver.domain.com.
mail.domain.com.	300	IN	MX	20 mailserver2.domain.com.

;; A Records
multiple-a.domain.com.	300	IN	A	192.0.2.236
multiple-a.domain.com.	300	IN	A	192.0.2.237
multiple-a.domain.com.	300	IN	A	192.0.2.238
short-a.domain.com.	60	IN	A	192.0.2.231
simple-a.domain.com.	300	IN	A	192.0.2.235

;; AAAA Records
ipv6.domain.com.	300	IN	AAAA	fe80:0:0:0:202:b2fe:fe1e:8329

;; CNAME Records
cname.domain.com.	300	IN	CNAME	multiple-a.domain.com

;; TXT Records
simple-txt.domain.com.	300	IN	TXT	"Sample Text Entries"

;; SRV Records
server.domain.com.	300	IN	SRV	1 10 5269 xmpp-server.domain.com.

;; PTR Records
pointer.domain.com.	300	IN	PTR	www.domain.com

;; SPF Records
policy.domain.com.	300	IN	SPF	"v=spf1 ip4:192.168.0.1/16-all"

;; NAPTR Records
name-auth.domain.com.	300	IN	NAPTR	100 100 "U" "" "!^.*$!sip:info@bar.example.com!" .

;; CAA Records
ca.domain.com.	300	IN	CAA	0 issue "caa.example.com"
ca.domain.com.	300	IN	CAA	0 issuewild ";"

;; Route53 Alias Records
alias-a.domain.com.	0	IN	A	a0123456789abcdef.awsglobalaccelerator.com.
alias-c.domain.com.	0	IN	CNAME	a0123456789abcdef.awsglobalaccelerator.com.
alias-caa.domain.com.	0	IN	CAA	a0123456789abcdef.awsglobalaccelerator.com.

```

## Manual:
<details><summary>Click to expand</summary>

- [Configuration](docs/configuration.md)
- Execution
  - [Docker](docs/docker.md)
  - [Kubernetes](docs/kubernetes.md)
  - [AWS Lambda](docs/aws-lambda.md)

</details>

## Remarks:
- **Exported files should never be imported directly!** Those exports does not follow DNS Zonefile format precisely and are intended to be read by human.  

## Known Issues:
- `Cloudflare max zones = 50`: this is due to CloudFlare package hardcoded max amount of listed domains:
    ```res, err = api.makeRequest("GET", "/zones?per_page=50", nil)```
- `Cloudflare zonefiles always changed`: this is because of a timestamp provided as a part of a SOA record.

## License
[MIT](LICENSE.md) © 2019-present Anton Yurchenko