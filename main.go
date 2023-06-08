package main

import (
	"context"
	"log"
	"os"

	"github.com/erlnerlngga/backend-socius/db"
	"github.com/erlnerlngga/backend-socius/internal/user"
	"github.com/erlnerlngga/backend-socius/internal/websocket"
	"github.com/erlnerlngga/backend-socius/router"
)

func main() {

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := router.NewApiServer("0.0.0.0:"+port, userHandler, wsHandler)
	server.Run()
}
