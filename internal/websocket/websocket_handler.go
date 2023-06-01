package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/erlnerlngga/backend-socius/util"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewWSHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) error {
	newRoom := new(CreateNewRoomType)
	if err := json.NewDecoder(r.Body).Decode(newRoom); err != nil {
		return err
	}

	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()

	rm := &RoomType{
		Name_Room: newRoom.Name_Room,
	}
	room, err := h.hub.Repository.CreateRoom(ctx, rm)
	if err != nil {
		return err
	}

	cl := &ClientType{
		Room_ID:   room.Room_ID,
		User_ID:   newRoom.User_ID,
		User_Name: newRoom.User_Name,
		Role:      "admin",
	}

	client, err := h.hub.InsertNewClient(ctx, cl)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, client)
}

type GetRoomByUserTypeReq struct {
	User_ID string `json:"user_id"`
	Room_ID string `json:"room_id"`
}

func (h *Handler) GetRoomByUser(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	defer r.Body.Close()
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()
	rooms, err := h.hub.Repository.GetRoomsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, rooms)
}

func (h *Handler) AddFriend(w http.ResponseWriter, r *http.Request) error {
	friend := new(ClientType)

	if err := json.NewDecoder(r.Body).Decode(friend); err != nil {
		return err
	}

	defer r.Body.Close()

	friend.Role = "user"

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()
	client, err := h.hub.Repository.InsertNewClient(ctx, friend)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, client)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	roomID := chi.URLParam(r, "roomID")
	userID := chi.URLParam(r, "userID")

	var cl *Client

	// check client
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()
	res, err := h.hub.Repository.CheckClient(ctx, userID, roomID)
	if err != nil {
		return err
	}

	cl = &Client{
		Conn:      conn,
		Message:   make(chan *MessageType, 10),
		Client_ID: res.Client_ID,
		User_ID:   res.User_ID,
		Room_ID:   res.Room_ID,
		User_Name: res.User_Name,
	}

	h.hub.Register <- cl

	go cl.writeMessage()
	cl.readMessage(h.hub)

	return nil
}

func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")
	roomID := chi.URLParam(r, "roomID")

	// check client
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()
	res, err := h.hub.Repository.CheckClient(ctx, userID, roomID)
	if err != nil {
		return err
	}

	err = h.hub.RemoveClient(ctx, userID, roomID)
	if err != nil {
		return err
	}

	if res.Role == "admin" {
		err := h.hub.RemoveRoom(ctx, roomID)
		if err != nil {
			return err
		}
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (h *Handler) GetAllMessage(w http.ResponseWriter, r *http.Request) error {
	roomID := chi.URLParam(r, "roomID")

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()

	message, err := h.hub.GetAllMessage(ctx, roomID)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, message)
}

func (h *Handler) CountAllUnreadMessage(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")
	result := 0

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()

	rooms, err := h.hub.GetRoomsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, val := range rooms {
		num, err := h.hub.CountAllUnreadMessage(ctx, userID, val.Room_ID)
		if err != nil {
			return err
		}

		result = result + num
	}

	return util.WriteJSON(w, http.StatusOK, map[string]int{"unread_message": result})
}
