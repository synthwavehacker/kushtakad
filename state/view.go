package state

import (
	"github.com/kushtaka/kushtakad/models"
)

// View is built, rendered, and cleared every HTTP request
type View struct {
	FlashSuccess   []string
	FlashFail      []string
	URI            string
	Crumbs         []*Crumb
	Endpoints      []*models.Endpoint
	Endpoint       *models.Endpoint
	User           *models.User
	State          *models.State
	Links          *Links
	Forms          *Forms
	Team           *models.Team
	Teams          []models.Team
	Token          *models.Token
	Tokens         []models.Token
	Sensor         *models.Sensor
	Sensors        []models.Sensor
	Users          []models.User
	SensorServices []models.ServiceCfg
}

type Crumb struct {
	Name string
	Link string
}

type Links struct {
	Setup bool
	Login bool

	Dashboard string
	Tokens    string
	Sensors   string
	Users     string
	Teams     string
}

func NewView() *View {
	var endpoints []*models.Endpoint
	var ff []string
	var fs []string
	var tm []models.Team
	var users []models.User
	var crumbs []*Crumb
	return &View{
		FlashFail:    ff,
		FlashSuccess: fs,
		Teams:        tm,
		Users:        users,
		Endpoints:    endpoints,
		Crumbs:       crumbs,
		Endpoint:     &models.Endpoint{},
		Links:        &Links{},
		Forms:        NewForms(),
		Team:         models.NewTeam(),
		Token:        models.NewToken(),
	}
}

func (v *View) AddCrumb(name, link string) {
	c := &Crumb{Name: name, Link: link}
	v.Crumbs = append(v.Crumbs, c)
}
