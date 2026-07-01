package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"strings"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type WebSocketServer struct {
	rooms     map[string]map[*websocket.Conn]bool
	broadcast chan *Message
	mu        sync.RWMutex
}

func NewWebSocket() *WebSocketServer {
	return &WebSocketServer{
		rooms:     make(map[string]map[*websocket.Conn]bool),
		broadcast: make(chan *Message),
	}
}
func (s *WebSocketServer) HandleWebSocket(ctx *websocket.Conn, roomCode string) {
	roomCode = strings.TrimSpace(roomCode)
	if roomCode == "" {
		ctx.Close()
		return
	}

	s.addClient(roomCode, ctx)
	defer func() {
		s.removeClient(roomCode, ctx)
		ctx.Close()
	}()
	for {
		_, msg, err := ctx.ReadMessage()
		if err != nil {
			log.Println("Read Error:", err)
			break
		}
		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Unmarshal Error:", err)
			continue
		}
		message.Text = strings.TrimSpace(message.Text)
		if message.Text == "" {
			continue
		}
		if message.Room == "" {
			message.Room = roomCode
		}
		s.broadcast <- &message
	}
}
func (s *WebSocketServer) HandleMessages() {
	for {
		msg := <-s.broadcast
		clients := s.clientsInRoom(msg.Room)
		for _, client := range clients {
			err := client.WriteMessage(websocket.TextMessage, getMessageTemplate(msg))
			if err != nil {
				log.Println("Write Error:", err)
				client.Close()
				s.removeClient(msg.Room, client)
			}
		}
	}
}

func (s *WebSocketServer) addClient(roomCode string, client *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.rooms[roomCode] == nil {
		s.rooms[roomCode] = make(map[*websocket.Conn]bool)
	}
	s.rooms[roomCode][client] = true
}

func (s *WebSocketServer) removeClient(roomCode string, client *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.rooms[roomCode] == nil {
		return
	}
	delete(s.rooms[roomCode], client)
	if len(s.rooms[roomCode]) == 0 {
		delete(s.rooms, roomCode)
	}
}

func (s *WebSocketServer) clientsInRoom(roomCode string) []*websocket.Conn {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clients := make([]*websocket.Conn, 0, len(s.rooms[roomCode]))
	for client := range s.rooms[roomCode] {
		clients = append(clients, client)
	}
	return clients
}

func getMessageTemplate(msg *Message) []byte {
	tmpl, err := template.ParseFiles("views/message.html")
	if err != nil {
		log.Println("Template Parse Error:", err)
		return nil
	}
	var rendereredMessage bytes.Buffer
	if err := tmpl.Execute(&rendereredMessage, msg); err != nil {
		log.Println("Template Execute Error:", err)
		return nil
	}
	return rendereredMessage.Bytes()

}
