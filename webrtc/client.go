package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	// egress is used to aviod concurrent writes on websocket connections

	egress chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress: make(chan []byte),
	}
}


func (client *Client) readMessages() {
	defer func() {
		// cleanup connection
		client.manager.removeClient(client)
	}()
	for {
		messageType, p, err := client.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		for wsclient := range client.manager.clients {
			wsclient.egress <- p
		}
		log.Println(messageType)
		log.Println(string(p))
	}
}


func (client *Client) WriteMessages(){
	defer func() {
		client.manager.removeClient(client)
	}()
	for {
		select {
		case message, ok := <-client.egress:
			if !ok {
				if err := client.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("connection closed  %v: ", err)
				}   
				return
			}
			if err := client.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("failed to send message %v: ", err)
				return
			}
			log.Println("message sent")
		}
	}
}