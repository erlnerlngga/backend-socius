package websocket

import (
	"encoding/json"
	"log"
	"net/http"

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

	rm := &RoomType{
		Name_Room: newRoom.Name_Room,
	}
	newRoomRes, err := h.hub.Repository.CreateRoom(rm)
	if err != nil {
		return err
	}

	cl := &ClientType{
		Room_ID:   newRoomRes.Room_ID,
		User_ID:   newRoom.User_ID,
		User_Name: newRoom.User_Name,
		Role:      "admin",
	}

	err = h.hub.InsertNewClient(cl)
	if err != nil {
		return err
	}

	h.hub.Rooms[newRoomRes.Room_ID] = &Room{
		Room_ID:   newRoomRes.Room_ID,
		Room_Name: newRoomRes.Name_Room,
		Clients:   make(map[string]*Client),
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) UpdateRoomName(w http.ResponseWriter, r *http.Request) error {
	upRoom := new(RoomType)

	if err := json.NewDecoder(r.Body).Decode(upRoom); err != nil {
		return err
	}

	defer r.Body.Close()

	err := h.hub.Repository.UpdateRoomName(upRoom)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type GetRoomByUserTypeReq struct {
	User_ID string `json:"user_id"`
	Room_ID string `json:"room_id"`
}

func (h *Handler) GetRoomByUser(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	defer r.Body.Close()

	rooms, err := h.hub.Repository.GetRoomsByUserID(userID)
	if err != nil {
		return err
	}

	for _, val := range rooms {
		_, ok := h.hub.Rooms[val.Room_ID]

		if !ok {
			h.hub.Rooms[val.Room_ID] = &Room{
				Room_ID:   val.Room_ID,
				Room_Name: val.Name_Room,
				Clients:   make(map[string]*Client),
			}
		}
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

	err := h.hub.Repository.InsertNewClient(friend)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:3000"
	},
}

func (h *Handler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panicln("UPGRADER CONNECTION", err)
	}

	roomID := chi.URLParam(r, "roomID")
	userID := chi.URLParam(r, "userID")

	var cl *Client

	// check client

	res, err := h.hub.Repository.CheckClient(userID, roomID)
	if err != nil {
		log.Panicln("CheckClient Handler", err)
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
}

func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")
	roomID := chi.URLParam(r, "roomID")

	// check client

	res, err := h.hub.Repository.CheckClient(userID, roomID)
	if err != nil {
		return err
	}

	err = h.hub.RemoveClient(userID, roomID)
	if err != nil {
		return err
	}

	if res.Role == "admin" {
		err := h.hub.RemoveRoom(roomID)
		if err != nil {
			return err
		}

		delete(h.hub.Rooms, roomID)
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (h *Handler) GetAllMessage(w http.ResponseWriter, r *http.Request) error {
	roomID := chi.URLParam(r, "roomID")

	message, err := h.hub.GetAllMessage(roomID)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, message)
}

func (h *Handler) CountAllUnreadMessage(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")
	result := 0

	rooms, err := h.hub.GetRoomsByUserID(userID)
	if err != nil {
		return err
	}

	for _, val := range rooms {
		num, err := h.hub.CountAllUnreadMessage(userID, val.Room_ID)
		if err != nil {
			return err
		}

		result = result + num
	}

	return util.WriteJSON(w, http.StatusOK, map[string]int{"unread_message": result})
}
