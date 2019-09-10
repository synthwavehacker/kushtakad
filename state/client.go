package state

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	heartWait      = 60 * time.Second
	heartBeat      = (heartWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ClientMsg struct {
	SensorID int64
	Msg      string
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	SensorID int64
	hub      *ServerHub
	conn     *websocket.Conn
	sender   chan *ClientMsg
	receiver chan *ClientMsg
}

func (c *Client) startReceiver() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(heartWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(heartWait)); return nil })
	for {
		var clientMsg *ClientMsg
		err := c.conn.ReadJSON(clientMsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.receiver <- clientMsg
	}
}

func (c *Client) startSender() {
	ticker := time.NewTicker(heartBeat)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case clientMsg, ok := <-c.sender:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteJSON(clientMsg)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			msg := &ClientMsg{}
			err := c.conn.WriteJSON(msg)
			if err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func (hub *ServerHub) ServeWs(w http.ResponseWriter, r *http.Request) {
	var apiKey string
	token, ok := r.Header["Authorization"]
	if ok && len(token) >= 1 {
		apiKey = token[0]
		apiKey = strings.TrimPrefix(apiKey, "Bearer ")
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, sender: make(chan *ClientMsg), receiver: make(chan *ClientMsg)}
	client.hub.register <- client
	go client.startReceiver()
	go client.startSender()
}
