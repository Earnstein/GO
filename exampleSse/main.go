package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Broker struct {
	Notifier       chan []byte
	newClients     chan chan []byte
	closingClients chan chan []byte
	clients        map[chan []byte]bool
	mu             sync.Mutex
}

func (b *Broker) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" && req.URL.Path == "/events" {
		b.handleSSE(res, req)
	} else if req.Method == "GET" && req.URL.Path == "/poll" {
		b.handleLongPoll(res, req)
	} else {
		http.NotFound(res, req)
	}
}

func (b *Broker) handleSSE(res http.ResponseWriter, req *http.Request) {
	flusher, ok := res.(http.Flusher)
	if !ok {
		http.Error(res, "Streaming not supported!", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	messageChan := make(chan []byte)

	b.mu.Lock()
	b.newClients <- messageChan
	b.mu.Unlock()

	defer func() {
		b.mu.Lock()
		b.closingClients <- messageChan
		b.mu.Unlock()
	}()

	notify := req.Context().Done()
	go func() {
		<-notify
		b.mu.Lock()
		b.closingClients <- messageChan
		b.mu.Unlock()
	}()

	for {
		select {
		case msg := <-messageChan:
			fmt.Fprintf(res, "data: %s\n\n", msg)
			flusher.Flush()
		case <-notify:
			return
		}
	}
}

func (b *Broker) handleLongPoll(res http.ResponseWriter, req *http.Request) {
	messageChan := make(chan []byte)

	b.mu.Lock()
	b.newClients <- messageChan
	b.mu.Unlock()

	defer func() {
		b.mu.Lock()
		b.closingClients <- messageChan
		b.mu.Unlock()
	}()

	select {
	case msg := <-messageChan:
		res.Header().Set("Content-Type", "application/json")
		res.Write(msg)
	case <-time.After(30 * time.Second):
		res.Header().Set("Content-Type", "application/json")
		res.Write([]byte(`{"error": "timeout"}`))
	}
}

func NewServer() *Broker {
	broker := &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	go broker.listen()
	return broker
}

func (b *Broker) listen() {
	for {
		select {
		case s := <-b.newClients:
			b.mu.Lock()
			b.clients[s] = true
			b.mu.Unlock()
			log.Printf("Client added to broker, total: %d\n", len(b.clients))
		case s := <-b.closingClients:
			b.mu.Lock()
			delete(b.clients, s)
			b.mu.Unlock()
			log.Printf("Removed client from broker, total: %d\n", len(b.clients))
		case event := <-b.Notifier:
			b.mu.Lock()
			for clientMessageChan := range b.clients {
				clientMessageChan <- event
			}
			b.mu.Unlock()
		}
	}
}

func main() {
	broker := NewServer()

	go func() {
		for {
			time.Sleep(3 * time.Second)
			name := randomName()
			eventString := fmt.Sprintf("The time is %v", name)
			log.Println("Receiving event")
			broker.Notifier <- []byte(eventString)
		}
	}()

	http.Handle("/", broker)
	log.Fatal("HTTP server Error", http.ListenAndServe(":8080", nil))
}

func randomName() string {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	name := ""
	for i := 0; i < 10; i++ {
		name += string(alphabet[rand.Intn(len(alphabet))])
	}
	return name
}                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              