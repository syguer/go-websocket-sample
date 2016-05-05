package main

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
)

type Client struct {
	ws        *websocket.Conn
	server    *Server
	send      chan *Message
	closeConn chan bool
}

func NewClient(ws *websocket.Conn, s *Server) *Client {
	return &Client{
		ws:        ws,
		server:    s,
		send:      make(chan *Message),
		closeConn: make(chan bool),
	}
}

func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

func (c *Client) listenWrite() {
	for {
		select {
		case msg := <-c.send:
			websocket.JSON.Send(c.ws, msg)
		case <-c.closeConn:
			c.server.Del(c)
			c.ws.Close()
			return
		}
	}
}

func (c *Client) Write(msg *Message) {
	log.Print("write: ")
	log.Println(msg)
	c.send <- msg
}

func (c *Client) listenRead() {
	for {
		select {
		case <-c.closeConn:
			c.server.Del(c)
			c.ws.Close()
			return
		default:
			var msg Message
			err := websocket.JSON.Receive(c.ws, &msg)
			log.Print("read: ")
			log.Println(msg)

			if err == io.EOF {
				log.Print("eof")
				c.closeConn <- true
			} else if err != nil {
				log.Print("err")
				c.server.Err(err)
			} else {
				log.Print("read: sendAll")
				c.server.sendAll <- &msg
			}
		}
	}
}
