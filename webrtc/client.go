package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	// egress is used to aviod concurrent writes on websocket connections

	egress chan Event
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (client *Client) readMessages() {
	defer func() {
		// cleanup connection
		client.manager.removeClient(client)
	}()

	if err := client.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("error reading message: %v", err)
		return
	}
	client.connection.SetReadLimit(512)
	client.connection.SetPongHandler(client.pongHandler)
	for {
		_, data, err := client.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		var request Event
		if err := json.Unmarshal(data, &request); err != nil {
			log.Printf("error reading message: %v", err)
			break
		}
		if err := client.manager.routeEvent(request, client); err != nil {
			log.Printf("error reading message: %v", err)
		}
	}
}

func (client *Client) WriteMessages() {
	defer func() {
		client.manager.removeClient(client)
	}()

	ticker := time.NewTicker(pingInterval)
	for {
		select {
		case message, ok := <-client.egress:
			if !ok {
				if err := client.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("connection closed  %v: ", err)
				}
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("failed to send message %v: ", err)
				return
			}
			if err := client.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("failed to send message %v: ", err)
				return
			}
			log.Println("message sent")
		case <-ticker.C:
			log.Println("ping")
			// send a ping to the client
			if err := client.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("writemsg error %v: ", err)
				return
			}
		}
	}
}

func (client *Client) pongHandler(pongMsg string) error {
	log.Println("pong")
	return client.connection.SetReadDeadline(time.Now().Add(pongWait))
}
