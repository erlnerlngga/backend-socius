package main

import (
	"context"
	"log"

	"github.com/erlnerlngga/backend-socius/db"
	"github.com/erlnerlngga/backend-socius/internal/user"
	"github.com/erlnerlngga/backend-socius/internal/websocket"
	"github.com/erlnerlngga/backend-socius/router"
)

func main() {
	log.Println("test")

	db, err := db.NewMysqlStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.InitDB(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	userRepo := user.NewUserRepository(db.GetDB())
	userHandler := user.NewUserHandler(userRepo)

	wsRepo := websocket.NewRepositoryWS(db.GetDB())
	wsHub := websocket.NewHub(*wsRepo)
	wsHandler := websocket.NewWSHandler(wsHub)
	go wsHub.Run(context.Background())

	server := router.NewApiServer(":8080", userHandler, wsHandler)
	server.Run()
}
