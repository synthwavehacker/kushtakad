package models

import (
	"github.com/xenolf/lego/challenge"
)

const (
	httpAuth = 0
	dnsAuth  = 1
)

var dnsProviders = [59]string{
	"acme-dns",
	"alidns",
	"azure",
	"auroradns",
	"bindman",
	"bluecat",
	"cloudflare",
	"cloudns",
	"cloudxns",
	"conoha",
	"designate",
	"digitalocean",
	"dnsimple",
	"dnsmadeeasy",
	"dnspod",
	"dode",
	"dreamhost",
	"duckdns",
	"dyn",
	"fastdns",
	"easydns",
	"exec",
	"exoscale",
	"gandi",
	"gandiv5",
	"glesys",
	"gcloud",
	"godaddy",
	"hostingde",
	"httpreq",
	"iij",
	"inwx",
	"joker",
	"lightsail",
	"linode",
	"linodev4",
	"manual",
	"mydnsjp",
	"namecheap",
	"namedotcom",
	"netcup",
	"nifcloud",
	"ns1",
	"oraclecloud",
	"otc",
	"ovh",
	"pdns",
	"rackspace",
	"route53",
	"rfc2136",
	"sakuracloud",
	"stackpath",
	"selectel",
	"transip",
	"vegadns",
	"vultr",
	"vscale",
	"zoneee",
}

type EnvVars struct {
	Name  string
	Value string
}

type Https struct {
	Provider challenge.Provider
	DnsType  int
	Email    string
}

func DoesExist(name string) bool {
	for _, v := range dnsProviders {
		if v == name {
			return true
		}
	}
	return false
}
