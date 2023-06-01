package router

import (
	"log"
	"net/http"

	"github.com/erlnerlngga/backend-socius/internal/user"
	"github.com/erlnerlngga/backend-socius/internal/websocket"
	"github.com/erlnerlngga/backend-socius/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type APIServer struct {
	listenAddr  string
	userHandler *user.Handler
	wsHandler   *websocket.Handler
}

func NewApiServer(listenAddr string, userHandler *user.Handler, wsHandler *websocket.Handler) *APIServer {
	return &APIServer{
		listenAddr:  listenAddr,
		userHandler: userHandler,
		wsHandler:   wsHandler,
	}
}

func (s *APIServer) Run() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	router.Post("/checkEmail", util.MakeHTTPHandleFunc(s.userHandler.CheckEmail))
	router.Post("/signup", util.MakeHTTPHandleFunc(s.userHandler.SignUp))
	router.Post("/addNewFriend", util.MakeHTTPHandleFunc(s.userHandler.AddNewFriend))
	router.Delete("/removeFriend", util.MakeHTTPHandleFunc(s.userHandler.RemoveFriend))
	router.Get("/getAllFriend/{userID}", util.MakeHTTPHandleFunc(s.userHandler.GetAllFriend))
	router.Post("/createPost", util.MakeHTTPHandleFunc(s.userHandler.CreatePost))
	router.Get("/getAllPost/{userID}", util.MakeHTTPHandleFunc(s.userHandler.GetAllPost))
	router.Get("/getPost/{postID}", util.MakeHTTPHandleFunc(s.userHandler.GetPost))
	router.Get("/getAllImage/{userID}", util.MakeHTTPHandleFunc(s.userHandler.GetAllImage))
	router.Post("/createComment", util.MakeHTTPHandleFunc(s.userHandler.CreateComment))
	router.Get("/getAllComment/{postID}", util.MakeHTTPHandleFunc(s.userHandler.GetAllComment))
	router.Post("/createNotification", util.MakeHTTPHandleFunc(s.userHandler.CreateNotification))
	router.Put("/updateAddFriendNotification", util.MakeHTTPHandleFunc(s.userHandler.UpdateAddFriendNotification))
	router.Put("/updateNotificationRead/{userID}", util.MakeHTTPHandleFunc(s.userHandler.UpdateNotificationRead))
	router.Get("/getCountNotification/{userID}", util.MakeHTTPHandleFunc(s.userHandler.GetCountNotification))
	router.Get("/getAllNotification/{userID}", util.MakeHTTPHandleFunc(s.userHandler.GetAllNotification))

	router.Post("/ws/createRoom", util.MakeHTTPHandleFunc(s.wsHandler.CreateRoom))
	router.Post("/ws/addFriend", util.MakeHTTPHandleFunc(s.wsHandler.AddFriend))
	router.Get("/ws/joinRoom/{roomID}-{userID}", util.MakeHTTPHandleFunc(s.wsHandler.JoinRoom))
	router.Get("/ws/getRoomsByUser/{userID}", util.MakeHTTPHandleFunc(s.wsHandler.GetRoomByUser))
	router.Get("/ws/getAllMessage/{roomID}", util.MakeHTTPHandleFunc(s.wsHandler.GetAllMessage))
	router.Get("/ws/getAllUnreadMessage/{userID}", util.MakeHTTPHandleFunc(s.wsHandler.CountAllUnreadMessage))
	router.Delete("/ws/remove/{roomID}-{userID}", util.MakeHTTPHandleFunc(s.wsHandler.Remove))

	log.Println("server runnng in port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}
