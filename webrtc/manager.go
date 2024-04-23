package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	clients ClientList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

func (manager *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")
	// upgrade regular http connection into web socket

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(ws, manager)
	manager.addClient(client)

	// start client processes

	go client.readMessages()
	go client.WriteMessages()
}


func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
}


func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	
	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}