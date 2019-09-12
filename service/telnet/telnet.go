package telnet

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

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
}

func (s TelnetService) Handle(ctx context.Context, conn net.Conn) error {
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

	log.Println(id, username, password)

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
