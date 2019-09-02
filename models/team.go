package models

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	is "github.com/go-ozzo/ozzo-validation/is"
)

const DefaultTeam = "Default"

type Team struct {
	ID        int64  `storm:"id,increment,index"`
	Name      string `storm:"index,unique" json:"name"`
	IsDeleted bool   `storm:"index" json:"is_deleted"`
	Members   []string

	MemberToAdd string
}

func NewTeam() *Team {
	return &Team{}
}

func (t *Team) Wash() {
	t.Name = strings.TrimSpace(t.Name)
	t.Name = Strip(t.Name)

	for key, member := range t.Members {
		m := strings.TrimSpace(member)
		m = Strip(m)
		t.Members[key] = m
	}
}

func (t *Team) ValidateCreate() error {
	t.Wash()
	return validation.Errors{
		"Name": validation.Validate(
			&t.Name,
			validation.Required,
			validation.Length(4, 64).Error("must be between 4-64 characters")),
	}.Filter()
}

func (t *Team) ValidateAddMember(newmember string) error {
	t.Wash()

	t.MemberToAdd = newmember

	for _, member := range t.Members {
		if t.MemberToAdd == member {
			return errors.New("Email is already added to this team")
		}
	}

	err := validation.Errors{
		"Member Name": validation.Validate(
			t.MemberToAdd,
			validation.Required,
			validation.Length(5, 120),
			is.Email)}.Filter()

	if err != nil {
		return err
	}

	t.Members = append(t.Members, t.MemberToAdd)

	return nil
}
