package websocket

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RoomType struct {
	Room_ID    string    `json:"room_id"`
	Name_Room  string    `json:"name_room"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

type RoomTypeRes struct {
	Room_ID        string    `json:"room_id"`
	Name_Room      string    `json:"name_room"`
	Created_At     time.Time `json:"created_at"`
	Updated_At     time.Time `json:"updated_at"`
	Unread_Message int       `json:"unread_message"`
}

type CreateNewRoomType struct {
	Name_Room string `json:"name_room"`
	User_ID   string `json:"user_id"`
	User_Name string `json:"user_name"`
}

type Client struct {
	Conn      *websocket.Conn
	Message   chan *MessageType
	Client_ID string `json:"client_id"`
	User_ID   string `json:"user_id"`
	Room_ID   string `json:"room_id"`
	User_Name string `json:"user_name"`
}

type ClientType struct {
	Client_ID  string    `json:"client_id"`
	Room_ID    string    `json:"room_id"`
	User_ID    string    `json:"user_id"`
	User_Name  string    `json:"user_name"`
	Role       string    `json:"role"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

type MessageType struct {
	Message_ID    string    `json:"message_id"`
	Room_ID       string    `json:"room_id"`
	User_ID       string    `json:"user_id"`
	User_Name     string    `json:"user_name"`
	Photo_Profile string    `json:"photo_profile"`
	Client_ID     string    `json:"client_id"`
	Content       string    `json:"content"`
	Created_At    time.Time `json:"created_at"`
	Updated_At    time.Time `json:"updated_at"`
}

type LogType struct {
	Log_ID     string    `json:"log_id"`
	Client_ID  string    `json:"client_id"`
	User_ID    string    `json:"user_id"`
	Status_Log string    `json:"status_log"`
	Created_At time.Time `json:"created_at"`
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("1. readMessage", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			break
		}

		msg := &MessageType{
			Message_ID: uuid.New().String(),
			Room_ID:    c.Room_ID,
			User_ID:    c.User_ID,
			Client_ID:  c.Client_ID,
			Content:    string(m),
			Created_At: time.Now().UTC(),
			Updated_At: time.Now().UTC(),
		}

		hub.Broadcast <- msg
	}
}
