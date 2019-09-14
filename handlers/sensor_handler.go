package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kushtaka/kushtakad/models"
	"github.com/kushtaka/kushtakad/state"
)

func GetSensor(w http.ResponseWriter, r *http.Request) {
	redir := "/kushtaka/sensors/page/1/limit/100"
	app, err := state.Restore(r)
	if err != nil {
		log.Fatal(err)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	sensor := &models.Sensor{}
	err = app.DB.One("ID", id, sensor)
	if err != nil {
		log.Error(err)
		app.Fail("Sensor does not exist")
		http.Redirect(w, r, redir, 302)
		return
	}

	for _, v := range sensor.Cfgs {
		app.View.SensorServices = append(app.View.SensorServices, v)
	}
	app.View.Sensor = sensor

	var teams []models.Team
	err = app.DB.All(&teams)
	if err != nil {
		log.Error(err)
		app.Fail(err.Error())
		http.Redirect(w, r, redir, 302)
		return
	}

	app.View.Teams = teams
	app.View.Links.Sensors = "active"
	app.View.AddCrumb("Sensors", "/kushtaka/sensors/page/1/limit/100")
	app.View.AddCrumb(sensor.Name, "#")
	app.Render.HTML(w, http.StatusOK, "admin/pages/sensor", app.View)
	return
}

func PostSensor(w http.ResponseWriter, r *http.Request) {
	log.Error("PostSensor()")
	return
}

func UpdateSensor(w http.ResponseWriter, r *http.Request) {
	log.Error("UpdateSensor()")
	return
}

func DeleteSensor(w http.ResponseWriter, r *http.Request) {
	log.Error("DeleteSensor()")
	return
}
