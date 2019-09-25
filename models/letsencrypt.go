package models

import (
	"github.com/boltdb/bolt"
	"github.com/mholt/certmagic"
)

type LE struct {
	Magic   *certmagic.Config
	Domains []string
	DB      *bolt.DB
}

func NewStageLE(email string, domains []string) LE {
	return LE{
		Magic:   leStageCfg(email),
		Domains: domains,
	}
}

func leStageCfg(email string) *certmagic.Config {
	cert := certmagic.NewDefault()
	cert.CA = certmagic.LetsEncryptStagingCA
	cert.Email = email
	cert.Agreed = true
	return cert
}
