package websocket

import (
	"context"
	"log"
	"time"
)

type Room struct {
	Room_ID   string             `json:"room_id"`
	Room_Name string             `json:"room_name"`
	Clients   map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *MessageType
	Repository
	timeout time.Duration
}

func NewHub(repository Repository) *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *MessageType, 5),
		Repository: repository,
		timeout:    time.Duration(2) * time.Second,
	}
}

func (h *Hub) Run(c context.Context) {

	for {
		select {
		case cl := <-h.Register:
			// check room
			_, err := h.Repository.CheckRoom(cl.Room_ID)
			_, ok := h.Rooms[cl.Room_ID]

			if err != nil {
				log.Println("1. Register", err)
			}

			if err == nil && ok {
				// check is client is not there

				_, ok := h.Rooms[cl.Room_ID].Clients[cl.Client_ID]

				if !ok {

					log := &LogType{
						Client_ID:  cl.Client_ID,
						User_ID:    cl.User_ID,
						Status_Log: "online",
					}
					// add client to that room
					h.Repository.CreateLog(log)
					h.Rooms[cl.Room_ID].Clients[cl.Client_ID] = cl
				}
			}

		case cl := <-h.Unregister:
			_, ok := h.Rooms[cl.Room_ID]

			if ok {
				_, ok := h.Rooms[cl.Room_ID].Clients[cl.Client_ID]

				if ok {
					log := &LogType{
						Client_ID:  cl.Client_ID,
						User_ID:    cl.User_ID,
						Status_Log: "leave",
					}

					h.Repository.CreateLog(log)
					delete(h.Rooms[cl.Room_ID].Clients, cl.Client_ID)
					close(cl.Message)
				}
			}

		case m := <-h.Broadcast:
			_, err := h.Repository.CheckRoom(m.Room_ID)
			_, ok := h.Rooms[m.Room_ID]
			if err != nil {
				log.Println("1. Broadcast", err)
			}
			if err == nil && ok {
				m, err = h.Repository.CreateMessage(m)
				if err != nil {
					log.Println("2. Broadcast", err)
				}

				for _, cl := range h.Rooms[m.Room_ID].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
