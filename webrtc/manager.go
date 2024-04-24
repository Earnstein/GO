package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

type Manager struct {
	clients ClientList
	sync.RWMutex
	otps     RentensionMap
	handlers map[string]EventHandler
}

func NewManager(ctx context.Context) *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
		otps:     NewRentionMap(ctx, 5*time.Second),
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessage
	m.handlers[EventChangeRoom] = ChatRoomHandler
}

func ChatRoomHandler(event Event, client *Client) error {
	var changeRoomEvent ChangeRoomEvent
	if err := json.Unmarshal(event.Data, &changeRoomEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	client.chatroom = changeRoomEvent.Name
	return nil
}

func SendMessage(event Event, client *Client) error {
	message := SendMessageEvent{}
	if err := json.Unmarshal(event.Data, &message); err != nil {
		return fmt.Errorf("error reading message: %v", err)
	}

	var broadMessage NewMessageEvent
	broadMessage.Message = message.Message
	broadMessage.From = message.From

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	outGoingEvent := Event{
		Data: data,
		Type: EventNewMessage,
	}

	for c := range client.manager.clients {
		if c.chatroom == client.chatroom {
			c.egress <- outGoingEvent
		}
	}
	return nil
}

func (manager *Manager) routeEvent(event Event, client *Client) error {
	if handler, ok := manager.handlers[event.Type]; ok {
		if err := handler(event, client); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("event not found")
	}

}
func (manager *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {

	otp := r.URL.Query().Get("otp")
	if otp == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !manager.otps.verifyOTP(otp) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
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

func (m *Manager) LoginHandler(w http.ResponseWriter, r *http.Request) {
	type userLoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req userLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Username == "admin" && req.Password == "admin" {
		type response struct {
			OTP string `json:"otp"`
		}
		otp := m.otps.NewOTP()
		resp := response{
			OTP: otp.Key,
		}

		data, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
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

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	switch origin {
	case "https://localhost:8080":
		return true
	default:
		return false
	}
}
