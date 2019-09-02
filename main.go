package main

import (
	"math/rand"
	"time"

	"github.com/kushtaka/kushtakad/server"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	server.Run()
}
