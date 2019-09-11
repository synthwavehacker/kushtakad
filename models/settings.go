package models

import (
	"github.com/asdine/storm"
	"github.com/gorilla/securecookie"
)

const SettingsID = 1

type Settings struct {
	ID           int64  `storm:"id,increment" json:"id"`
	SessionHash  []byte `json:"session_hash"`
	SessionBlock []byte `json:"session_block"`
	CsrfHash     []byte `json:"csrf_hash"`
	Host         string
	Scheme       string
}

func InitSettings(db *storm.DB) (Settings, error) {
	var s Settings
	db.One("ID", SettingsID, &s)
	if len(s.SessionHash) != 32 {
		s.SessionHash = securecookie.GenerateRandomKey(32)
	}

	if len(s.SessionBlock) != 16 {
		s.SessionBlock = securecookie.GenerateRandomKey(16)
	}

	if len(s.CsrfHash) != 32 {
		s.CsrfHash = securecookie.GenerateRandomKey(32)
	}

	if len(s.Host) == 0 {
		s.Host = "localhost:8080"
	}

	if s.Scheme != "http" || s.Scheme != "https" {
		s.Scheme = "http"
	}

	err := db.Save(&s)
	if err != nil {
		return s, err
	}

	return s, nil
}

func FindSettings(db *storm.DB) (*Settings, error) {
	var s Settings
	err := db.One("ID", SettingsID, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
