package helpers

import (
	//stdlib

	"bytes"
	"fmt"
	"net/smtp"
	"path/filepath"
	"text/template"
	"time"

	"github.com/asdine/storm"

	packr "github.com/gobuffalo/packr/v2"
	"github.com/kushtaka/kushtakad/models"
	email "gopkg.in/jordan-wright/email.v1"
)

type Mailer struct {
	Smtp     *models.Smtp
	Settings *models.Settings
	DB       *storm.DB
	Box      *packr.Box

	EventID      int64
	EventLink    string
	Subject      string
	Text         string
	TemplateName string
	TemplateFile string
}

func NewMailer(db *storm.DB, box *packr.Box) *Mailer {
	smtp := &models.Smtp{}
	err := db.One("ID", 1, smtp)
	if err != nil {
		log.Debugf("Smtp values are net set %v", err)
	}

	m := &Mailer{Smtp: smtp, Box: box, DB: db}
	return m
}

func (m *Mailer) SendSensorEvent(eventid, furl, hashid, state, emailtext string, tt time.Time) error {

	fname := "event_sensor.tmpl"
	fp := filepath.Join("admin", "email", fname)
	mt, err := m.Box.Find(fp)
	if err != nil {
		log.Debugf("Unable to find template in packr box %v", err)
	}

	uri := models.BuildURI(m.DB)
	m.Subject = fmt.Sprintf("%s : %s", furl, eventid)
	m.Text = fmt.Sprintf("Event: %s <br>\n\nState: %s<br>\n\n", furl, state)
	m.EventLink = uri + "/kushtaka/%s"
	m.TemplateName = "EventSensor"
	m.TemplateFile = fname
	m.EventLink = fmt.Sprintf(m.EventLink, m.EventID)
	e := email.NewEmail()
	e.From = m.Smtp.Email
	e.To = []string{emailtext}
	t, err := template.New("MT").Parse(string(mt))
	if err != nil {
		return err
	}
	var out bytes.Buffer
	err = t.ExecuteTemplate(&out, m.TemplateName, m)
	e.HTML = []byte(out.String())

	hostport := fmt.Sprintf("%s:%s", m.Smtp.Host, m.Smtp.Port)
	err = e.Send(
		hostport,
		smtp.PlainAuth(
			"",
			m.Smtp.Username,
			m.Smtp.Password,
			m.Smtp.Host,
		),
	)
	return err
}
