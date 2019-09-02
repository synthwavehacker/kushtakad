package models

import "time"

type Endpoint struct {
	ID       int64     `storm:"id,increment"`
	URI      string    `storm:"index,unique" json:"uri"`
	ApiKey   string    `storm:"index,unique" json:"api_key"`
	Created  time.Time `storm:"index"`
	Released time.Time `storm:"index"`
	Updated  time.Time
}
