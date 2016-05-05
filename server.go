package main

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

type Server struct {
	clients      map[*Client]bool
	addClient    chan *Client
	deleteClient chan *Client
	sendAll      chan *Message
	errCh        chan error
}

func NewServer() *Server {
	return &Server{
		clients:   make(map[*Client]bool),
		addClient: make(chan *Client),
		sendAll:   make(chan *Message),
		errCh:     make(chan error),
	}
}

func (s *Server) Del(c *Client) {
	s.clients[c] = false
}

func (s *Server) Err(err error) {
	log.Println("Error:", err.Error())
	s.errCh <- err
}

func (s *Server) HandleConnection(ws *websocket.Conn) {
	client := NewClient(ws, s)
	s.addClient <- client
}

func (s *Server) Listen() {
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.addClient <- client
		client.Listen()
	}

	http.Handle("/connection", websocket.Handler(onConnected))

	for {
		select {
		case c := <-s.addClient:
			s.clients[c] = true

		case c := <-s.deleteClient:
			s.clients[c] = false

		case msg := <-s.sendAll:
			log.Print("sendAll: ")
			log.Println(msg)
			for client := range s.clients {
				if s.clients[client] {
					client.Write(msg)
				}
			}

		case err := <-s.errCh:
			log.Println("Error:", err.Error())
		}
	}
}
