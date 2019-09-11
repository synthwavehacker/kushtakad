package service

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type Angel struct {
	Auth     *Auth
	MyCtx    context.Context
	MyCancel context.CancelFunc

	SensorCtx    context.Context
	SensorCancel context.CancelFunc
}

func interuptor(cancel context.CancelFunc) {
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt)
		signal.Notify(s, syscall.SIGTERM)

		select {
		case <-s:
			cancel()
		}
	}()

}

func NewAngel(auth *Auth) *Angel {
	angelCtx, angelCancel := context.WithCancel(context.Background())
	sensorCtx, sensorCancel := context.WithCancel(context.Background())
	angel := &Angel{
		Auth:         auth,
		MyCtx:        angelCtx,
		MyCancel:     angelCancel,
		SensorCtx:    sensorCtx,
		SensorCancel: sensorCancel,
	}
	interuptor(angel.MyCancel)
	return angel
}

func Run(host, apikey string) {
	auth, err := Config(host, apikey)
	if err != nil {
		log.Error("you must pass the cli values -host and -apikey |or| have valid auth.json file.")
		log.Fatal(err)
	}
	log.Info(auth)

	angel := NewAngel(auth)
	startSensor(angel.SensorCtx)

	for {
		select {
		case <-angel.MyCtx.Done(): // if the angel's context is close
			angel.SensorCtx.Done() // close the sensor's
			log.Info("shutting down angel...done.")
			return

		default:
		}
	}

}
