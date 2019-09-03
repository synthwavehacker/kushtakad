package handlers

import (
	"fmt"
	"log"
	"net/http"
	"github.com/kushtaka/kushtakad/models"
	"github.com/kushtaka/kushtakad/state"
)

func GetSensor(w http.ResponseWriter, r *http.Request) {
	log.Println("GetSensor()")
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
	redirUrl := "/kushtaka/dashboard"
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
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	app.View.Sensors = sensors 
	app.Render.HTML(w, http.StatusOK, "admin/pages/sensors", app.View)
	return
}


func PostSensors(w http.ResponseWriter, r *http.Request) {
	redirUrl := "/kushtaka/sensors/page/1/limit/100"
	app, err := state.Restore(r)
	if err != nil {
		log.Fatal(err)
	}

	name := r.FormValue("name")
	sensor := &models.Sensor{Name: name}
	app.View.Forms.Sensor = sensor
	err = sensor.ValidateCreate()
	if err != nil {
		app.Fail(err.Error())
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	tx, err := app.DB.Begin(true)
	if err != nil {
		tx.Rollback()
		app.Fail(err.Error())
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	tx.One("Name", name, sensor)
	if sensor.ID > 0 {
		tx.Rollback()
		app.Fail("Sensor using that name already exists.")
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	err = tx.Save(sensor)
	if err != nil {
		tx.Rollback()
		app.Fail(err.Error())
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		app.Fail(err.Error())
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	app.View.Forms = state.NewForms()
	app.Success(fmt.Sprintf("The sensor [%s] was created successfully.", sensor.Name))
	http.Redirect(w, r, redirUrl, 301)
	return
}