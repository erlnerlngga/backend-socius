package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type MysqlStore struct {
	db *sql.DB
}

func NewMysqlStore() (*MysqlStore, error) {
	// open the connection of db

	db, err := sql.Open("mysql", os.Getenv("DSN"))

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("database is running ..... ")

	return &MysqlStore{
		db: db,
	}, nil
}

// create user table
func (s *MysqlStore) CreateTableUser() error {
	createTable := `
		create table if not exists user (
			user_id varchar(100),
			user_name varchar(100) not null,
			email varchar(50) not null unique,
			primary key(user_id)
		);
	`

	_, err := s.db.Exec(createTable)

	return err
}

// create table user_friend
func (s *MysqlStore) CreateTableUser_Friend() error {
	createTabele := `
		create table if not exists user_friend (
			user_friend_id varchar(100),
			user_id varchar(100) references user(user_id),
			friend_id varchar(100) references user(user_id),
			primary key(user_friend_id)
		);
	`

	_, err := s.db.Exec(createTabele)

	return err
}

// create post table
func (s *MysqlStore) CreatePostTable() error {
	createTable := `
		create table if not exists post (
			post_id varchar(100),
			user_id varchar(100) references user(user_id),
			content varchar(500),
			created_at timestamp,
			updated_at timestamp,
			primary key(post_id)
		);
	`

	_, err := s.db.Exec(createTable)

	return err
}

// create image post table
func (s *MysqlStore) CreateImage_PostTable() error {
	createTable := `
		create table if not exists image_post (
			image_post_id varchar(100),
			post_id varchar(100) references post(post_id),
			user_id varchar(100) references user(user_id),
			image varchar(300),
			created_at timestamp,
			updated_at timestamp,
			primary key(image_post_id)
		);
	`

	_, err := s.db.Exec(createTable)

	return err
}

// create comment table
func (s *MysqlStore) CreateTableComment() error {
	createTable := `
		create table if not exists comment (
			comment_id varchar(100),
			post_id varchar(100) references post(post_id),
			comment_post_id varchar(100) references post(post_id),
			created_at timestamp,
			updated_at timestamp,
			primary key(comment_id)
		);
	`

	_, err := s.db.Exec(createTable)

	return err
}

// create table notfication
func (s *MysqlStore) CreateNotificationTable() error {
	createTable := `
		create table if not exists notification (
			notification_id varchar(100),
			issuer varchar(100) references user(user_id),
			notifier varchar(100) references user(user_id),
			notifier_name varchar(100),
			status varchar(20) not null,
			accept varchar(20),
			post_id varchar(100) references post(post_id),
			type varchar(20) not null,
			created_at timestamp, 
			updated_at timestamp,
			primary key(notification_id)
		);
	`

	_, err := s.db.Exec(createTable)

	return err
}

// create room message
func (s *MysqlStore) CreateTableRoom() error {
	createTable := `
		create table if not exists room (
			room_id varchar(100),
			romm_name varchar(50) not null,
			created_at timestamp,
			updated_at timestamp,
			primary key(room_id)
		);
	`

	_, err := s.db.Exec(createTable)

	return err
}

// create client table
func (s *MysqlStore) CreateTableClient() error {
	createTable := `
		create table if not exists client (
			client_id varchar(100), 
			room_id varchar(100) references room(room_id) on delete cascade,
			user_id varchar(100) references user(user_id),
			user_name varchar(100) references user(user_name),
			role varchar(20),
			created_at timestamp,
			updated_at timestamp,
			primary key(client_id)
		);
	`

	_, err := s.db.Exec(createTable)

	return err
}

// create message table
func (s *MysqlStore) CreateTableMessage() error {
	createTable := `
		create table if not exists message (
			message_id varchar(100),
			room_id varchar(100) references room(room_id) on delete cascade,
			user_id varchar(100) references user(user_id),
			client_id varchar(100) references client(client_id) on delete set null,
			content varchar(500) not null,
			created_at timestamp,
			upated_at timestamp,
			primary key(message_id)
		);
	`

	_, err := s.db.Exec(createTable)

	return err
}

// create table log
func (s *MysqlStore) CreateTableLog() error {
	createTable := `
		create table if not exists log (
			log_id varchar(100),
			client_id varchar(100) references client(client_id),
			user_id varchar(100) references user(user_id),
			status_log varchar(20) not null,
			created_at timestamp,
			primary key(log_id)
		);
	`

	_, err := s.db.Exec(createTable)

	return err
}

func (s *MysqlStore) InitDB() error {
	if err := s.CreateTableUser(); err != nil {
		return err
	}

	if err := s.CreateTableUser_Friend(); err != nil {
		return err
	}

	if err := s.CreatePostTable(); err != nil {
		return err
	}

	if err := s.CreateImage_PostTable(); err != nil {
		return err
	}

	if err := s.CreateTableComment(); err != nil {
		return err
	}

	if err := s.CreateNotificationTable(); err != nil {
		return err
	}

	if err := s.CreateTableRoom(); err != nil {
		return err
	}

	if err := s.CreateTableClient(); err != nil {
		return err
	}

	if err := s.CreateTableMessage(); err != nil {
		return err
	}

	if err := s.CreateTableLog(); err != nil {
		return err
	}

	return nil
}

func (s *MysqlStore) Close() {
	s.db.Close()
}

func (s *MysqlStore) GetDB() *sql.DB {
	return s.db
}
