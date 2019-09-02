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

	SmtpServer   string `json:"smtp_server"`
	SmtpUser     string `json:"smtp_user"`
	SmtpPassword string `json:"smtp_password"`
	SmtpPort     string `json:"smtp_port"`

	SmtpIsSetup  bool `json:"smtp_is_setup"`
	AdminIsSetup bool `json:"admin_is_setup"`
	HttpsIsSetup bool `json:"https_is_setup"`
	TeamIsSetup  bool `json:"team_is_setup"`
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
