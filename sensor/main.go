package main

import (
	"math/rand"
	"time"

	"github.com/kushtaka/kushtakas/sensor"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {

	sensor.Run()
}
