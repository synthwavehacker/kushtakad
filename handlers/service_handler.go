package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kushtaka/kushtakad/models"
	"github.com/kushtaka/kushtakad/service/telnet"
	"github.com/kushtaka/kushtakad/state"
)

func PostService(w http.ResponseWriter, r *http.Request) {
	app, err := state.Restore(r)
	if err != nil {
		log.Fatal(err)
	}

	var js []byte
	vars := mux.Vars(r)
	serviceType := vars["type"]
	sensorId, err := strconv.Atoi(vars["sensor_id"])
	if err != nil {
		log.Fatal(err)

	}

	tx, err := app.DB.Begin(true)
	if err != nil {
		tx.Rollback()
		log.Println("can't begin ", err)
		return
	}

	var sensor models.Sensor
	tx.One("ID", sensorId, &sensor)
	if sensor.ID == 0 {
		tx.Rollback()
		log.Println("zero sensor id ", err)
		return
	}

	cfg := models.ServiceCfg{}
	switch serviceType {
	case "telnet":
		decoder := json.NewDecoder(r.Body)
		var tel telnet.TelnetService
		err = decoder.Decode(&tel)
		if err != nil {
			panic(err)
		}

		tel.Prompt = "$ "

		err = tx.Save(&tel)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return
		}

		js, err = json.Marshal(tel)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return
		}

		cfg.ServiceID = tel.ID
		cfg.SensorID = sensor.ID
		cfg.Type = serviceType
		cfg.Port = tel.Port

		for _, v := range sensor.Cfgs {
			if v.Port == tel.Port {
				tx.Rollback()
				log.Println("Port is already assigned to another service", err)
			}
		}

		sensor.Cfgs = append(sensor.Cfgs, cfg)

	default:
		tx.Rollback()
		log.Println("unable to find service type")
		return
	}

	err = tx.Save(&cfg)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return
	}

	err = tx.Update(&sensor)
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
