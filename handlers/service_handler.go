package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kushtaka/kushtakad/models"
	"github.com/kushtaka/kushtakad/service/telnet"
	"github.com/kushtaka/kushtakad/state"
)

func PostService(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(r.Body)
	var cfg models.ServiceCfg
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}

	js, err := json.Marshal(cfg)
	if err != nil {
		log.Println(err)
		return
	}

	tx, err := app.DB.Begin(true)
	if err != nil {
		tx.Rollback()
		log.Println("can't begin ", err)
		return
	}

	var sensor models.Sensor
	tx.One("ID", cfg.SensorID, &sensor)
	if sensor.ID == 0 {
		tx.Rollback()
		log.Println("zero sensor id ", err)
		return
	}

	var serviceID int64
	switch cfg.Type {
	case "telnet":
		tel := telnet.TelnetService{
			SensorID: cfg.SensorID,
			Port:     21,
			Prompt:   "$",
			Emulate:  "basic",
			Type:     "telnet",
		}

		err = tx.Save(&tel)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return
		}

		serviceID = tel.ID
		log.Println("Service ID ", serviceID)

	default:
		tx.Rollback()
		log.Println("unable to find service type")
		return
	}

	cfg.ServiceID = serviceID
	cfg.SensorID = sensor.ID
	err = tx.Save(&cfg)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
