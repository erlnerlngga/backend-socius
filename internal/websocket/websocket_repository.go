package websocket

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Repository struct {
	db DBTX
}

func NewRepositoryWS(db DBTX) *Repository {
	return &Repository{db: db}
}

// create new room
func (r *Repository) CreateRoom(ctx context.Context, room *RoomType) (*RoomType, error) {
	result := new(RoomType)
	room.Room_ID = uuid.New().String()
	room.Created_At = time.Now().UTC()
	room.Updated_At = time.Now().UTC()

	query := `insert into room(room_id, name_room, created_at, updated_at) values (?, ?, ?, ?);`
	err := r.db.QueryRowContext(ctx, query, room.Room_ID, room.Name_Room, room.Created_At, room.Updated_At).Scan(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// check is romm available or not
func (r *Repository) CheckRoom(ctx context.Context, roomd_id string) (*RoomType, error) {
	result := new(RoomType)

	query := `select room_id, name_room, create_at, updated_at from room where room_id = ?;`
	err := r.db.QueryRowContext(ctx, query, roomd_id).Scan(result)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("room_id isn't found")
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Get Single Room
func (r *Repository) GetRoom(ctx context.Context, room *RommTypeRes) (*RommTypeRes, error) {
	result := new(RommTypeRes)

	query := `select * from room where room_id = ?;`

	err := r.db.QueryRowContext(ctx, query, room.Room_ID).Scan(&result.Room_ID, &result.Name_Room, &result.Created_At, &result.Updated_At)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("room isn't found")
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Get All Rooms
func (r *Repository) GetRoomsByUserID(ctx context.Context, user_id string) ([]*RommTypeRes, error) {
	query := `select room_id from client where user_id = ?;`

	rows, err := r.db.QueryContext(ctx, query, user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	resultRooms := []*RommTypeRes{}
	for rows.Next() {
		ro := new(RommTypeRes)

		if err := rows.Scan(&ro.Room_ID); err != nil {
			return nil, err
		}

		room, err := r.GetRoom(ctx, ro)
		if err != nil {
			return nil, err
		}

		room, err = r.CountUnreadMessage(ctx, user_id, room)
		if err != nil {
			return nil, err
		}

		resultRooms = append(resultRooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resultRooms, nil
}

// check is that specific room there is alreay client or not
func (r *Repository) CheckClient(ctx context.Context, client_id, room_id string) (*ClientType, error) {
	result := new(ClientType)

	query := `select * from client where room_id = ? and client_id = ?;`
	err := r.db.QueryRowContext(ctx, query, room_id, client_id).Scan(result)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("client_id isn't found")
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// check is that specific room there is alreay client or not
func (r *Repository) CheckClientByUserID(ctx context.Context, user_id, room_id string) (*ClientType, error) {
	result := new(ClientType)

	query := `select * from client where room_id = ? and user_id = ?;`
	err := r.db.QueryRowContext(ctx, query, room_id, user_id).Scan(result)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user_id isn't found")
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// add client
func (r *Repository) InsertNewClient(ctx context.Context, client *ClientType) (*ClientType, error) {
	result := new(ClientType)
	client.Client_ID = uuid.New().String()
	client.Created_At = time.Now().UTC()
	client.Updated_At = time.Now().UTC()

	query := `insert into client(client_id, room_id, user_id, user_name, role, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?);`
	err := r.db.QueryRowContext(ctx, query, client.Client_ID, client.Room_ID, client.User_Name, client.Role, client.Created_At, client.Updated_At).Scan(&result.Client_ID, &result.Room_ID, &result.User_Name, &result.Role, &result.Created_At, &result.Updated_At)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// get client base on room ID
func (r *Repository) GetClients(ctx context.Context, room_id string) ([]*ClientType, error) {
	query := `select * from client where client_id = ?;`

	rows, err := r.db.QueryContext(ctx, query, room_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	clients := []*ClientType{}
	for rows.Next() {
		c := new(ClientType)

		if err := rows.Scan(&c.Client_ID, &c.Room_ID, &c.User_ID, &c.User_Name, &c.Created_At, &c.Updated_At); err != nil {
			return nil, err
		}

		clients = append(clients, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}

// remove room
func (r *Repository) RemoveRoom(ctx context.Context, room_id string) error {
	_, err := r.db.ExecContext(ctx, `delete from room where room_id = ?;`, room_id)

	if err != nil {
		return err
	}

	return nil
}

// remove client from that room
func (r *Repository) RemoveClient(ctx context.Context, user_id, room_id string) error {
	_, err := r.db.ExecContext(ctx, `delete from client where user_id = ? and room_id = ?;`, user_id, room_id)

	if err != nil {
		return err
	}

	return nil
}

// create log
func (r *Repository) CreateLog(ctx context.Context, log *LogType) (*LogType, error) {
	newLog := new(LogType)

	log.Log_ID = uuid.New().String()
	log.Created_At = time.Now().UTC()

	query := `insert into log(log_id, client_id, user_id, status_log, created_at) values (?, ?, ?, ?, ?);`
	err := r.db.QueryRowContext(ctx, query, log.Log_ID, log.Client_ID, log.User_ID, log.Status_Log, log.Created_At).Scan(&newLog.Log_ID, &newLog.Client_ID, &newLog.User_ID, &newLog.Status_Log, &newLog.Created_At)
	if err != nil {
		return nil, err
	}

	return newLog, nil
}

// get Count unread message
func (r *Repository) CountUnreadMessage(ctx context.Context, user_id string, room *RommTypeRes) (*RommTypeRes, error) {
	log := new(LogType)

	query := `select * from log where user_id = ? ans status_log = "leave" order by created_at desc;`

	err := r.db.QueryRowContext(ctx, query, user_id).Scan(&log.Log_ID, &log.Client_ID, &log.User_ID, &log.Status_Log, &log.Created_At)

	if err != nil {
		return nil, err
	}

	countQuery := "select count(*) as `unread_message` from message where created_at >= ? and room_id = ?;"
	err = r.db.QueryRowContext(ctx, countQuery, log.Created_At, room.Room_ID).Scan(&room.Unread_Message)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *Repository) CountAllUnreadMessage(ctx context.Context, user_id, room_id string) (int, error) {
	log := new(LogType)

	var number int

	query := `select * from log where user_id = ? ans status_log = "leave" order by created_at desc;`

	err := r.db.QueryRowContext(ctx, query, user_id).Scan(&log.Log_ID, &log.Client_ID, &log.User_ID, &log.Status_Log, &log.Created_At)

	if err != nil {
		return -1, err
	}

	countQuery := "select count(*) as `unread_message` from message where created_at >= ? and room_id = ?;"
	err = r.db.QueryRowContext(ctx, countQuery, log.Created_At, room_id).Scan(&number)
	if err != nil {
		return -1, err
	}

	return number, nil
}

// get all message
func (r *Repository) GetAllMessage(ctx context.Context, room_id string) ([]*MessageType, error) {
	query := `select * from message where room_id = ?;`

	rows, err := r.db.QueryContext(ctx, query, room_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	allMessage := []*MessageType{}
	for rows.Next() {
		m := new(MessageType)

		if err := rows.Scan(&m.Message_ID, &m.Room_ID, &m.User_ID, &m.Client_ID, &m.Content, &m.Created_At, &m.Updated_At); err != nil {
			return nil, err
		}

		allMessage = append(allMessage, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return allMessage, nil
}

// create message
func (r *Repository) CreateMessage(ctx context.Context, msg *MessageType) (*MessageType, error) {
	newMsg := new(MessageType)
	query := `insert into message(message_id, room_id, user_id, client_id, content, createdd_at, updated_at) values(?, ?, ?, ?, ?, ?, ?);`

	err := r.db.QueryRowContext(ctx, query, msg.Message_ID, msg.Room_ID, msg.User_ID, msg.Client_ID, msg.Content, msg.Created_At, msg.Updated_At).Scan(&newMsg.Message_ID, &newMsg.Room_ID, &newMsg.User_ID, &newMsg.Client_ID, &newMsg.Content, &newMsg.Created_At, &newMsg.Updated_At)
	if err != nil {
		return nil, err
	}

	return newMsg, nil
}
