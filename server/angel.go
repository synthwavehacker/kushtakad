package server

import (
	"context"
	"net/http"

	"github.com/kushtaka/kushtakad/angel"
)

type ServerAngel struct {
	AngelCtx     context.Context
	AngelCancel  context.CancelFunc
	ServerCtx    context.Context
	ServerCancel context.CancelFunc
	Server       *http.Server
	Reboot       chan bool
}

func NewServer(sa *ServerAngel) *http.Server {
	sctx, scancel := context.WithCancel(context.Background())
	srv := RunServer(sa.Reboot)
	sa.ServerCtx = sctx
	sa.ServerCancel = scancel
	return srv
}

func NewServerAngel() *ServerAngel {
	reboot := make(chan bool)
	actx, acancel := context.WithCancel(context.Background())
	sctx, scancel := context.WithCancel(context.Background())
	srv := RunServer(reboot)
	sa := &ServerAngel{
		AngelCtx:     actx,
		AngelCancel:  acancel,
		ServerCtx:    sctx,
		ServerCancel: scancel,
		Server:       srv,
		Reboot:       reboot,
	}
	angel.Interuptor(sa.AngelCancel)
	return sa
}

func Run() {
	sa := NewServerAngel()
	for {
		select {
		case <-sa.Reboot:
			log.Info("Reboot channel signaled...")
			log.Info("Shutting down server...")
			sa.Server.Shutdown(sa.ServerCtx)
			log.Info("Done.")
			log.Info("Rebooting server...")
			sa.Server = NewServer(sa)
			log.Info("Done.")
		case <-sa.AngelCtx.Done(): // if the angel's context is closed
			log.Info("shutting down Angel...done.")
			sa.Server.Shutdown(sa.ServerCtx)
			log.Info("shutting down Server...done.")
			return
		case <-sa.ServerCtx.Done(): // if the angel's context is closed
			log.Info("shutting down ServerAngel...done.")
			return
		}
		//default:
		//https://medium.com/@ashishstiwari/dont-simply-run-forever-loop-for-1594464040b1
		// is this needed?
		// my CPU really tops out
		//time.Sleep(100 * time.Millisecond)
	}

}
