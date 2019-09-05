package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kushtaka/kushtakad/state"
)

// {sensorId: 1, type: serviceType, port: 23, emulate: 'basic'}
type Service struct {
	SensorID int64  `json:"sensorId"`
	Type     string `json:"type"`
}

func PostService(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(r.Body)
	var t Service
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t)

	js, err := json.Marshal(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	tx, err := app.DB.Begin(true)
	if err != nil {
		tx.Rollback()
		app.Fail(err.Error())
		http.Redirect(w, r, redirUrl, 301)
		return
	}

	tx.One("ID", t.SensorID, sensor)
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

}
