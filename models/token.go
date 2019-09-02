package models

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	packr "github.com/gobuffalo/packr/v2"
	docx "github.com/kushtaka/kushtakad/tokens/docx"
)

type Token struct {
	ID       int64  `storm:"id,increment,index"`
	Name     string `storm:"index,unique" json:"name"`
	Note     string `storm:"index" json:"note"`
	Type     string // Weblink, Pdf, Docx
	Token    string `storm:"index,unique"`
	TeamsIds []int64
	File     []byte
}

func NewToken() *Token {
	return &Token{}
}

func (t *Token) Wash() {
	t.Name = strings.TrimSpace(t.Name)
	t.Name = Strip(t.Name)
}

func (t *Token) ValidateCreate() error {
	t.Wash()
	return validation.Errors{
		"Name": validation.Validate(
			&t.Name,
			validation.Required,
			validation.Length(4, 64).Error("must be between 4-64 characters")),
	}.Filter()
}

func (t *Token) BuildDocx(box *packr.Box) (s string, err error) {
	b, err := box.Find("files/template.docx")
	if err != nil {
		return "", err
	}

	r, err := docx.ReadDocxFromMemory(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return "", err
	}
	defer r.Close()

	min := time.Date(2014, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2018, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	created := fmt.Sprintf("%s", time.Unix(sec, 0).Format("2006-01-02T15:04:05Z"))
	now := fmt.Sprintf("%s", time.Now().UTC().Format("2006-01-02T15:04:05Z"))

	d := r.Editable()
	d.ReplaceCoreRaw("aaaaaaaaaaaaaaaaaaaa", created)
	d.ReplaceCoreRaw("bbbbbbbbbbbbbbbbbbbb", now)
	d.ReplaceFooterRaw("HONEYDROP_TOKEN_URL", "http://localhost:3000/blah.png")
	filename, err := d.WriteToTmpFile()
	if err != nil {
		return "", err
	}

	return filename, nil

}
