package websocket

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type Repository struct {
	db DBTX
}

func NewRepositoryWS(db DBTX) *Repository {
	return &Repository{db: db}
}

// create new room
func (r *Repository) CreateRoom(room *RoomType) (*RoomType, error) {

	room.Room_ID = uuid.New().String()
	room.Created_At = time.Now().UTC()
	room.Updated_At = time.Now().UTC()

	query := `insert into room(room_id, name_room, created_at, updated_at) values (?, ?, ?, ?);`
	_, err := r.db.Exec(query, room.Room_ID, room.Name_Room, room.Created_At, room.Updated_At)
	if err != nil {
		log.Println("1. CreateRoom", err)
		return nil, err
	}

	newRes := new(RoomType)

	querySelect := "select * from room where room_id = ?;"
	err = r.db.QueryRow(querySelect, room.Room_ID).Scan(&newRes.Room_ID, &newRes.Name_Room, &newRes.Created_At, &newRes.Updated_At)
	if err != nil {
		log.Println("2. CreateRoom", err)
		return nil, err
	}

	return newRes, nil
}

// change room name
func (r *Repository) UpdateRoomName(room *RoomType) error {
	room.Updated_At = time.Now().UTC()

	query := `update room set name_room = ?, updated_at = ? where room_id = ?;`

	_, err := r.db.Exec(query, room.Name_Room, room.Updated_At, room.Room_ID)
	if err != nil {
		log.Println("1. UpdateRoomName", err)
		return err
	}

	return nil
}

// check is romm available or not
func (r *Repository) CheckRoom(room_id string) (*RoomType, error) {
	result := new(RoomType)

	query := `select room_id, name_room, created_at, updated_at from room where room_id = ?;`
	err := r.db.QueryRow(query, room_id).Scan(&result.Room_ID, &result.Name_Room, &result.Created_At, &result.Updated_At)

	if err == sql.ErrNoRows {
		log.Println("1. CheckRoom", err)
		return nil, fmt.Errorf("room_id isn't found")
	}

	if err != nil {
		log.Println("2. CheckRoom", err)
		return nil, err
	}

	return result, nil
}

// Get Single Room
func (r *Repository) GetRoom(room *RoomTypeRes) (*RoomTypeRes, error) {
	result := new(RoomTypeRes)

	query := `select * from room where room_id = ?;`

	err := r.db.QueryRow(query, room.Room_ID).Scan(&result.Room_ID, &result.Name_Room, &result.Created_At, &result.Updated_At)

	if err == sql.ErrNoRows {
		log.Println("1. GetRoom", err)
		return nil, fmt.Errorf("room isn't found")
	}

	if err != nil {
		log.Println("2. GetRoom", err)
		return nil, err
	}

	return result, nil
}

// Get All Rooms
func (r *Repository) GetRoomsByUserID(user_id string) ([]*RoomTypeRes, error) {
	query := `select room_id from client where user_id = ?;`

	rows, err := r.db.Query(query, user_id)

	if err != nil {
		log.Println("1. GetRoomsByUserID", err)
		return nil, err
	}

	defer rows.Close()

	resultRooms := []*RoomTypeRes{}
	for rows.Next() {
		ro := new(RoomTypeRes)

		if err := rows.Scan(&ro.Room_ID); err != nil {
			log.Println("2. GetRoomsByUserID", err)
			return nil, err
		}

		room, err := r.GetRoom(ro)
		if err != nil {
			log.Println("3. GetRoomsByUserID", err)
			return nil, err
		}

		room, err = r.CountUnreadMessage(user_id, room)
		if err != nil {
			log.Println("4. GetRoomsByUserID", err)
			return nil, err
		}

		resultRooms = append(resultRooms, room)
	}

	if err := rows.Err(); err != nil {
		log.Println("5. GetRoomsByUserID", err)
		return nil, err
	}

	return resultRooms, nil
}

// check is that specific room there is alreay client or not
func (r *Repository) CheckClient(user_id, room_id string) (*ClientType, error) {
	result := new(ClientType)

	query := `select * from client where user_id = ? and room_id = ?;`
	err := r.db.QueryRow(query, user_id, room_id).Scan(&result.Client_ID, &result.Room_ID, &result.User_ID, &result.User_Name, &result.Role, &result.Created_At, &result.Updated_At)

	if err == sql.ErrNoRows {
		log.Println("1. CheckClient", err)
		return nil, fmt.Errorf("user_id isn't found")
	}

	if err != nil {
		log.Println("2. CheckClient", err)
		return nil, err
	}

	return result, nil
}

// check is that specific room there is alreay client or not
func (r *Repository) CheckClientByClientID(client_id, room_id string) (*ClientType, error) {
	result := new(ClientType)

	query := `select * from client where client_id = ? and room_id = ?;`
	err := r.db.QueryRow(query, client_id, room_id).Scan(&result.Client_ID, &result.Room_ID, &result.User_ID, &result.User_Name, &result.Role, &result.Created_At, &result.Updated_At)

	if err == sql.ErrNoRows {
		log.Println("1. CheckClientByUserID", err)
		return nil, fmt.Errorf("client_id isn't found")
	}

	if err != nil {
		log.Println("2. CheckClientByUserID", err)
		return nil, err
	}

	return result, nil
}

// add client
func (r *Repository) InsertNewClient(client *ClientType) error {
	client.Client_ID = uuid.New().String()
	client.Created_At = time.Now().UTC()
	client.Updated_At = time.Now().UTC()

	query := `insert into client(client_id, room_id, user_id, user_name, role, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?);`
	_, err := r.db.Exec(query, client.Client_ID, client.Room_ID, client.User_ID, client.User_Name, client.Role, client.Created_At, client.Updated_At)

	if err != nil {
		log.Println("1. InsertNewClient", err)
		return err
	}

	return nil
}

// get client base on room ID
func (r *Repository) GetClients(room_id string) ([]*ClientType, error) {
	query := `select * from client where client_id = ?;`

	rows, err := r.db.Query(query, room_id)

	if err != nil {
		log.Println("1. GetClients", err)
		return nil, err
	}

	defer rows.Close()

	clients := []*ClientType{}
	for rows.Next() {
		c := new(ClientType)

		if err := rows.Scan(&c.Client_ID, &c.Room_ID, &c.User_ID, &c.User_Name, &c.Created_At, &c.Updated_At); err != nil {
			log.Println("2. GetClients", err)
			return nil, err
		}

		clients = append(clients, c)
	}

	if err := rows.Err(); err != nil {
		log.Println("3. GetClients", err)
		return nil, err
	}

	return clients, nil
}

// remove room
func (r *Repository) RemoveRoom(room_id string) error {
	_, err := r.db.Exec(`delete from room where room_id = ?;`, room_id)

	if err != nil {
		log.Println("1. RemoveRoom", err)
		return err
	}

	return nil
}

// remove client from that room
func (r *Repository) RemoveClient(user_id, room_id string) error {
	_, err := r.db.Exec(`delete from client where user_id = ? and room_id = ?;`, user_id, room_id)

	if err != nil {
		log.Println("1. RemoveClient", err)
		return err
	}

	return nil
}

// create log
func (r *Repository) CreateLog(log *LogType) error {

	log.Log_ID = uuid.New().String()
	log.Created_At = time.Now().UTC()

	query := `insert into log(log_id, client_id, user_id, status_log, created_at) values (?, ?, ?, ?, ?);`
	_, err := r.db.Exec(query, log.Log_ID, log.Client_ID, log.User_ID, log.Status_Log, log.Created_At)
	if err != nil {
		fmt.Println("1. CreateLog", err)
		return err
	}

	return nil
}

// get count is not yet join
func (r *Repository) CountMessageNotYetJoin(room_id string) (int, error) {
	var unread_message int
	query := "select count(*) as `unread_message` from message where room_id = ?;"

	err := r.db.QueryRow(query, room_id).Scan(&unread_message)

	if err == sql.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		fmt.Println("1. CountMessageNotYetJoin", err)
		return -1, err
	}

	return unread_message, nil
}

// get Count unread message
func (r *Repository) CountUnreadMessage(user_id string, room *RoomTypeRes) (*RoomTypeRes, error) {
	log := new(LogType)

	query := `select * from log where user_id = ? and status_log = "leave" order by created_at desc;`

	err := r.db.QueryRow(query, user_id).Scan(&log.Log_ID, &log.Client_ID, &log.User_ID, &log.Status_Log, &log.Created_At)

	if err == sql.ErrNoRows {
		numRes, err := r.CountMessageNotYetJoin(room.Room_ID)
		if err != nil {
			fmt.Println("0. CountUnreadMessage", err)
			return nil, err
		}
		room.Unread_Message = numRes
		return room, nil
	}

	if err != nil {
		fmt.Println("1. CountUnreadMessage", err)
		return nil, err
	}

	countQuery := "select count(*) as `unread_message` from message where created_at >= ? and room_id = ?;"
	err = r.db.QueryRow(countQuery, log.Created_At, room.Room_ID).Scan(&room.Unread_Message)

	if err == sql.ErrNoRows {
		room.Unread_Message = 0
		return room, nil
	}

	if err != nil {
		fmt.Println("2. CountUnreadMessage", err)
		return nil, err
	}

	return room, nil
}

func (r *Repository) CountAllUnreadMessage(user_id, room_id string) (int, error) {
	log := new(LogType)

	var number int

	query := `select * from log where user_id = ? and status_log = "leave" order by created_at desc;`

	err := r.db.QueryRow(query, user_id).Scan(&log.Log_ID, &log.Client_ID, &log.User_ID, &log.Status_Log, &log.Created_At)

	if err == sql.ErrNoRows {
		numRes, err := r.CountMessageNotYetJoin(room_id)
		if err != nil {
			fmt.Println("0. CountAllUnreadMessage", err)
			return -1, err
		}
		return numRes, nil
	}

	if err != nil {
		fmt.Println("1. CountAllUnreadMessage", err)
		return -1, err
	}

	countQuery := "select count(*) as `unread_message` from message where created_at >= ? and room_id = ?;"
	err = r.db.QueryRow(countQuery, log.Created_At, room_id).Scan(&number)

	if err == sql.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		fmt.Println("2. CountAllUnreadMessage", err)
		return -1, err
	}

	return number, nil
}

// get user
func (r *Repository) GetUser(u *MessageType) (*MessageType, error) {
	query := `select user_name, photo_profile from user where user_id = ?;`

	err := r.db.QueryRow(query, u.User_ID).Scan(&u.User_Name, &u.Photo_Profile)
	if err != nil {
		fmt.Println("1. GetUser", err)
		return nil, err
	}

	return u, nil
}

// get all message
func (r *Repository) GetAllMessage(room_id string) ([]*MessageType, error) {
	query := `select * from message where room_id = ?;`

	rows, err := r.db.Query(query, room_id)
	if err != nil {
		fmt.Println("1. GetAllMessage", err)
		return nil, err
	}

	defer rows.Close()

	allMessage := []*MessageType{}
	for rows.Next() {
		m := new(MessageType)

		if err := rows.Scan(&m.Message_ID, &m.Room_ID, &m.User_ID, &m.Client_ID, &m.Content, &m.Created_At, &m.Updated_At); err != nil {
			fmt.Println("2. GetAllMessage", err)
			return nil, err
		}

		m, err = r.GetUser(m)
		if err != nil {
			fmt.Println("3. GetAllMessage", err)
			return nil, err
		}

		allMessage = append(allMessage, m)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("4. GetAllMessage", err)
		return nil, err
	}

	return allMessage, nil
}

// create message
func (r *Repository) CreateMessage(msg *MessageType) (*MessageType, error) {
	query := `insert into message(message_id, room_id, user_id, client_id, content, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?);`

	_, err := r.db.Exec(query, msg.Message_ID, msg.Room_ID, msg.User_ID, msg.Client_ID, msg.Content, msg.Created_At, msg.Updated_At)
	if err != nil {
		fmt.Println("1. CreateMessage", err)
		return nil, err
	}

	msg, err = r.GetUser(msg)
	if err != nil {
		fmt.Println("2. CreateMessage", err)
		return nil, err
	}

	return msg, nil
}
