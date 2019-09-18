package server

import (
	"context"

	"github.com/kushtaka/kushtakad/angel"
)

type ServerAngel struct {
	AngelCtx     context.Context
	AngelCancel  context.CancelFunc
	ServerCtx    context.Context
	ServerCancel context.CancelFunc
}

func NewServerAngel() *ServerAngel {
	angelCtx, angelCancel := context.WithCancel(context.Background())
	sensorCtx, sensorCancel := context.WithCancel(context.Background())
	serviceAngel := &ServerAngel{
		AngelCtx:     angelCtx,
		AngelCancel:  angelCancel,
		ServerCtx:    sensorCtx,
		ServerCancel: sensorCancel,
	}
	angel.Interuptor(serviceAngel.AngelCancel)
	return serviceAngel
}

func Run() {
	angel := NewServerAngel()
	RunServer()
	for {
		select {
		case <-angel.AngelCtx.Done(): // if the angel's context is closed
			angel.ServerCtx.Done() // close the sensor's
			log.Info("shutting down ServerAngel...done.")
			return
		case <-angel.ServerCtx.Done(): // if the angel's context is closed
			log.Info("shutting down ServerAngel...done.")
			return
			//default:
			//https://medium.com/@ashishstiwari/dont-simply-run-forever-loop-for-1594464040b1
			// is this needed?
			// my CPU really tops out
			//time.Sleep(100 * time.Millisecond)
		}
	}

}
