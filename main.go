package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/kushtaka/kushtakad/server"
	"github.com/kushtaka/kushtakad/service"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

const empty = ""

func main() {

	host := flag.String("host", empty, "the hostname of the kushtakad orchestrator server (string)")
	apikey := flag.String("apikey", empty, "the api key of the sensor, create from the kushtaka dashboard. (string)")
	sensor := flag.Bool("sensor", false, "would you like this instance to be a sensor? (bool)")
	flag.Parse()
	log.Println(*host, *apikey, *sensor)

	if *sensor {
		service.Run(*host, *apikey)
	} else {
		server.Run()
	}
}
