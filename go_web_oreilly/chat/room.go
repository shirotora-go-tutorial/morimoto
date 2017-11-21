package main

import (
	"../trace"
	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
	"log"
	"net/http"
)

type room struct {
	//forward chan []byte
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  trace.Tracer
	avatar  Avatar
}

func newRoom(avatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("new client was joined")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("client was left")
		case msg := <-r.forward:
			//r.tracer.Trace("get message", string(msg))
			r.tracer.Trace("get message", msg.Message)
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.tracer.Trace("succes to send messege to client")
				default:
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("failed to send clean up client")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("failed to get cookie", err)
		return
	}

	client := &client{
		socket: socket,
		//send: make(chan []byte, messageBufferSize),
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.join <- client
	defer func() {
		r.leave <- client
	}()
	go client.write()
	client.read()
}
