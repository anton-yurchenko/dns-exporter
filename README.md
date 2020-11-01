# dns-exporter

[![Release](https://img.shields.io/github/v/release/anton-yurchenko/dns-exporter)](https://github.com/anton-yurchenko/dns-exporter/releases/latest)
[![codecov](https://codecov.io/gh/anton-yurchenko/dns-exporter/branch/master/graph/badge.svg)](https://codecov.io/gh/anton-yurchenko/dns-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/anton-yurchenko/dns-exporter)](https://goreportcard.com/report/github.com/anton-yurchenko/dns-exporter)
[![Tests](https://github.com/anton-yurchenko/dns-exporter/workflows/push/badge.svg)](https://github.com/anton-yurchenko/dns-exporter/actions)
[![Docker Build](https://img.shields.io/docker/cloud/build/antonyurchenko/dns-exporter)](https://hub.docker.com/r/antonyurchenko/dns-exporter)
[![Docker Pulls](https://img.shields.io/docker/pulls/antonyurchenko/dns-exporter)](https://hub.docker.com/r/antonyurchenko/dns-exporter)
[![License](https://img.shields.io/github/license/anton-yurchenko/dns-exporter)](LICENSE.md)

You are most certainly apply proper Backup procedures for your deployed application, version control your code and so on.... **But what will you do in case one of your DNS entries was misconfigured or deleted?** It will probably mean a downtime to your application while you are trying to figure out where that missing CNAME was pointing or what was stored in that TXT record!  

**DNS-EXPORTER** is designed especially for those cases! Run it periodically against a DNS providers to document your domains configuration.

## Features

- Export all DNS records in Zonefile-like format.  
- Export to local/remote Git repository allowing easy tracking of changes.  
- Supported DNS providers: **CloudFlare, Route53**
- Supported Public / Private zones.  

## Example Export

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

## Manual

- [Configuration](docs/configuration.md)
  - [Permissions](docs/permissions.md)
- Execution
  - [Docker](docs/docker.md)
  - [Kubernetes](docs/kubernetes.md)

## Remarks

- **Exported files should never be imported directly!** Those exports does not follow DNS Zonefile format precisely and are intended to be read by human.  

## License

[MIT](LICENSE.md) © 2019-present Anton Yurchenko
