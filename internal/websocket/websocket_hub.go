package websocket

import (
	"context"
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
	ctx, cancel := context.WithTimeout(c, h.timeout)
	defer cancel()

	for {
		select {
		case cl := <-h.Register:
			// check room
			room, err := h.Repository.CheckRoom(ctx, cl.Room_ID)
			_, ok := h.Rooms[cl.Room_ID]
			if err == nil && ok {
				// check is client is not there
				_, err := h.Repository.CheckClient(ctx, cl.Client_ID, room.Room_ID)
				_, ok := h.Rooms[cl.Room_ID].Clients[cl.Client_ID]
				if err.Error() == "client_id isn't found" && !ok {

					log := &LogType{
						Client_ID:  cl.Client_ID,
						User_ID:    cl.User_ID,
						Status_Log: "online",
					}
					// add client to that room
					h.Repository.CreateLog(ctx, log)
					h.Rooms[cl.Room_ID].Clients[cl.Client_ID] = cl
				}
			}

		case cl := <-h.Unregister:
			room, err := h.Repository.CheckRoom(ctx, cl.Room_ID)
			_, ok := h.Rooms[cl.Room_ID]

			if err == nil && ok {
				_, err := h.Repository.CheckClient(ctx, cl.Client_ID, room.Room_ID)
				_, ok := h.Rooms[cl.Room_ID].Clients[cl.Client_ID]

				if err == nil && ok {
					log := &LogType{
						Client_ID:  cl.Client_ID,
						User_ID:    cl.User_ID,
						Status_Log: "leave",
					}

					h.Repository.CreateLog(ctx, log)
					delete(h.Rooms[cl.Room_ID].Clients, cl.Client_ID)
					close(cl.Message)
				}
			}

		case m := <-h.Broadcast:
			_, err := h.Repository.CheckRoom(ctx, m.Room_ID)
			_, ok := h.Rooms[m.Room_ID]

			if err == nil && ok {
				h.Repository.CreateMessage(ctx, m)
				for _, cl := range h.Rooms[m.Room_ID].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
