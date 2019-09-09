package handlers

import (
	"fmt"
	"log"
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
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	sensor := &models.Sensor{}
	err = app.DB.One("ID", id, sensor)
	if err != nil {
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
		app.Fail(err.Error())
		http.Redirect(w, r, redir, 302)
		return
	}
	app.View.Teams = teams

	log.Println(teams)

	app.View.Links.Sensors = "active"
	app.View.AddCrumb("Sensors", "/kushtaka/sensors/page/1/limit/100")
	app.View.AddCrumb(sensor.Name, "#")
	app.Render.HTML(w, http.StatusOK, "admin/pages/sensor", app.View)
	return
}

func PostSensor(w http.ResponseWriter, r *http.Request) {
	log.Println("PostSensor()")
	return
}

func UpdateSensor(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateSensor()")
	return
}

func DeleteSensor(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteSensor()")
	return
}

func GetSensors(w http.ResponseWriter, r *http.Request) {
	redir := "/kushtaka/dashboard"
	app, err := state.Restore(r)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, "/404", 404)
		return
	}
	app.View.Links.Sensors = "active"

	var sensors []models.Sensor
	err = app.DB.All(&sensors)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, redir, 302)
		return
	}
	app.View.Sensors = sensors

	var teams []models.Team
	err = app.DB.All(&teams)
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, redir, 302)
		return
	}
	app.View.Teams = teams

	log.Println(teams)
	app.View.AddCrumb("Sensors", "#")
	app.Render.HTML(w, http.StatusOK, "admin/pages/sensors", app.View)
	return
}

func PostSensors(w http.ResponseWriter, r *http.Request) {
	redir := "/kushtaka/sensors/page/1/limit/100"
	app, err := state.Restore(r)
	if err != nil {
		log.Fatal(err)
	}

	name := r.FormValue("name")
	tid, err := strconv.ParseInt(r.FormValue("teamId"), 10, 64)
	if err != nil {
		app.Fail("Please select a team to notify")
		http.Redirect(w, r, redir, 302)
		return
	}

	sensor := &models.Sensor{Name: name, TeamID: tid}
	app.View.Forms.Sensor = sensor
	err = sensor.ValidateCreate()
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, redir, 302)
		return
	}

	tx, err := app.DB.Begin(true)
	if err != nil {
		tx.Rollback()
		app.Fail(err.Error())
		http.Redirect(w, r, redir, 302)
		return
	}

	tx.One("Name", name, sensor)
	if sensor.ID > 0 {
		tx.Rollback()
		app.Fail("Sensor using that name already exists.")
		http.Redirect(w, r, redir, 302)
		return
	}

	err = tx.Save(sensor)
	if err != nil {
		tx.Rollback()
		app.Fail(err.Error())
		http.Redirect(w, r, redir, 302)
		return
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		app.Fail(err.Error())
		http.Redirect(w, r, redir, 302)
		return
	}

	app.View.Forms = state.NewForms()
	app.Success(fmt.Sprintf("The sensor [%s] was created successfully.", sensor.Name))
	http.Redirect(w, r, redir, 302)
	return
}
