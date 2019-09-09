package models

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Sensor struct {
	ID     int64        `storm:"id,increment,index"`
	TeamID int64        `storm:"id,index"`
	Name   string       `storm:"index,unique" json:"name"`
	ApiKey string       `storm:"index,unique" json:"api_key"`
	Cfgs   []ServiceCfg `storm:"index" json:"service_configs`
	mu     sync.Mutex
}

type ServiceCfg struct {
	ID        int64  `storm:"id,increment,index" json:"ID"` // we name the ID something different so that a json marshal/unmarshal doesn't accidentally inflate it
	SensorID  int64  `storm:"index" json:"sensorId"`
	ServiceID int64  `storm:"index" json:"serviceId"`
	Port      int    `storm:"index" json:"port"`
	Type      string `storm:"index" json:"type"`
}

type Service interface {
	Handle(ctx context.Context, conn net.Conn) error
}

func NewSensor() *Sensor {
	return &Sensor{ApiKey: GenerateSecureKey()}
}

func GenerateSecureKey() string {
	k := make([]byte, 32)
	io.ReadFull(rand.Reader, k)
	return fmt.Sprintf("%x", k)
}

func (s *Sensor) Wash() {
	s.Name = strings.TrimSpace(s.Name)
	s.Name = Strip(s.Name)
}

func (s *Sensor) ValidateCreate() error {
	s.Wash()
	return validation.Errors{
		"Name": validation.Validate(
			&s.Name,
			validation.Required,
			validation.Length(4, 64).Error("must be between 4-64 characters")),
		"TeamID": validation.Validate(
			&s.TeamID,
			validation.Required.Error("is required")),
	}.Filter()
}
