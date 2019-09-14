package main

import (
	"flag"
	"math/rand"
	"os"
	"time"

	"github.com/kushtaka/kushtakad/server"
	"github.com/kushtaka/kushtakad/service"
	"github.com/op/go-logging"
)

const empty = ""

var log = logging.MustGetLogger("main")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	// For demo purposes, create two backend for os.Stderr.
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)
}

func main() {

	host := flag.String("host", empty, "the hostname of the kushtakad orchestrator server (string)")
	apikey := flag.String("apikey", empty, "the api key of the sensor, create from the kushtaka dashboard. (string)")
	sensor := flag.Bool("sensor", false, "would you like this instance to be a sensor? (bool)")
	flag.Parse()

	if *sensor {
		service.Run(*host, *apikey)
	} else {
		server.Run()
	}
}
