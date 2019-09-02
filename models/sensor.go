package models

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Sensor struct {
	ID       int64  `storm:"id,increment,index"`
	Name     string `storm:"index,unique" json:"name"`
	Services []Service
}

type Service struct {
	Port int
	Type string
}

func NewSensor() *Sensor {
	return &Sensor{}
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
	}.Filter()
}
