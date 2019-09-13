// Copyright 2016-2019 DutchSec (https://dutchsec.com/)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package telnet

import (
	"context"
	"io"
	"net"

	"github.com/rs/xid"

	lua "github.com/yuin/gopher-lua"
)

// TelnetLua is a placeholder
func TelnetLua() *telnetLuaService {
	s := &telnetLuaService{
		MOTD:   motd,
		Prompt: prompt,
	}

	return s
}

type telnetLuaService struct {
	File string `toml:"file"`

	Prompt string `toml:"prompt"`
	MOTD   string `toml:"motd"`
}

func (s *telnetLuaService) Handle(ctx context.Context, conn net.Conn) error {
	id := xid.New()

	defer conn.Close()

	term := NewTerminal(conn, s.Prompt)
	term.Write([]byte(s.MOTD + "\n"))
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

	L := lua.NewState()

	mt := L.NewTypeMetatable("terminal")
	L.SetGlobal("terminal", mt)

	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(),
		map[string]lua.LGFunction{
			"log_event": func(L *lua.LState) int {
				ud := L.CheckUserData(1)

				_, ok := ud.Value.(*Terminal)
				if !ok {
					L.ArgError(1, "terminal expected")
					return 1
				}

				if L.GetTop() != 2 {
					L.ArgError(1, "table expected")
					return 0
				}

				_, ok = FromLUA(L.Get(2)).(map[string]interface{})
				if !ok {
					//log.Errorf("Unexpected type: %#+v", FromLUA(L.Get(2)))
					return 0
				}

				return 0
			},
			"read_line": func(L *lua.LState) int {
				ud := L.CheckUserData(1)

				term, ok := ud.Value.(*Terminal)
				if !ok {
					L.ArgError(1, "terminal expected")
					return 1
				}

				line, err := term.ReadLine()
				if err == io.EOF {
					return 0
				} else if err != nil {
					return 0
				}

				L.Push(lua.LString(line))
				return 1
			},
			"write": func(L *lua.LState) int {
				ud := L.CheckUserData(1)

				term, ok := ud.Value.(*Terminal)
				if !ok {
					L.ArgError(1, "terminal expected")
					return 0
				}

				if L.GetTop() != 2 {
					L.ArgError(1, "string expected")
					return 0
				}

				//log.Info(L.Get(2).String())
				term.Write([]byte(L.Get(2).String()))
				return 0
			},
			"write_line": func(L *lua.LState) int {
				ud := L.CheckUserData(1)

				term, ok := ud.Value.(*Terminal)
				if !ok {
					L.ArgError(1, "terminal expected")
					return 0
				}

				if L.GetTop() != 2 {
					L.ArgError(1, "string expected")
					return 0
				}

				//log.Info(L.Get(2).String())
				term.Write([]byte(L.Get(2).String()))
				term.Write([]byte("\n"))
				return 0
			},
		}))

	if err := L.DoFile(s.File); err != nil {
		return err
	}

	ud := L.NewUserData()
	ud.Value = term
	L.SetMetatable(ud, L.GetTypeMetatable("terminal"))

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("handle"),
		NRet:    1,
		Protect: true,
	}, ud); err != nil {
		//log.Error("Error calling lua method: %s", err.Error())
		return err
	}

	ret := L.Get(-1)
	L.Pop(1)

	_ = ret

	return nil
}
