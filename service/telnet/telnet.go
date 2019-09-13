package telnet

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/kushtaka/kushtakad/events"
	"github.com/op/go-logging"
	"github.com/rs/xid"
)

var (
	motd = `********************************************************************************
*             Copyright(C) 2008-2015 Huawei Technologies Co., Ltd.             *
*                             All rights reserved                              *
*                  Without the owner's prior written consent,                  *
*           no decompiling or reverse-engineering shall be allowed.            *
* Notice:                                                                      *
*                   This is a private communication system.                    *
*             Unauthorized access or use may lead to prosecution.              *
********************************************************************************

Warning: Telnet is not a secure protocol, and it is recommended to use STelnet. 

Login authentication


`
	prompt = `$ `
)

var log = logging.MustGetLogger("telnet")

// Telnet is a placeholder
func Telnet() *TelnetService {
	s := &TelnetService{
		Emulate: motd,
		Prompt:  prompt,
	}

	return s
}

type TelnetService struct {
	ID       int64  `storm:"id,increment,index" json:"id"`
	SensorID int64  `storm:"index" json:"sensorId"`
	Port     int    `json:"port"`
	Prompt   string `json:"prompt"`
	Emulate  string `json:"emulate"`
	Type     string `json:"type"`

	Host   string
	ApiKey string
}

func (s TelnetService) SetHost(h string) {
	s.Host = h
}

func (s TelnetService) SetApiKey(k string) {
	s.ApiKey = k
}

func (s TelnetService) Handle(ctx context.Context, conn net.Conn) error {
	em := events.NewEventManager(s.Type, s.Port, s.SensorID)

	log.Debugf("Handle %s %s", s.Host, s.ApiKey)

	err := em.SendEvent("new", s.Host, s.ApiKey, conn.RemoteAddr())
	if err != nil {
		log.Debug(err)
	}
	id := xid.New()

	defer conn.Close()

	term := NewTerminal(conn, s.Prompt)

	term.Write([]byte(s.Emulate + "\n"))

	term.SetPrompt("Username: ")
	username, err := term.ReadLine()
	if err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}

	password, err := term.ReadPassword("Password: ")
	if err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}

	log.Debug(id, username, password)

	term.SetPrompt(s.Prompt)

	for {
		line, err := term.ReadLine()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		if line == "" {
			continue
		}

		term.Write([]byte(fmt.Sprintf("sh: %s: command not found\n", line)))
	}
}
