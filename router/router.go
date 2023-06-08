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

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://socius-jade.vercel.app", "https://socius-laannen-gmailcom.vercel.app", "https://socius-5ym9o8can-laannen-gmailcom.vercel.app", "https://socius-git-main-laannen-gmailcom.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	router.Group(func(r chi.Router) {
		r.Get("/ws/joinRoom/{roomID}/{userID}", s.wsHandler.JoinRoom)
	})

	router.Get("/", util.MakeHTTPHandleFunc(s.userHandler.Welcome))
	router.Post("/signup", util.MakeHTTPHandleFunc(s.userHandler.SignUp))
	router.Post("/signin", util.MakeHTTPHandleFunc(s.userHandler.SignIn))
	router.Get("/auth/{token}", util.MakeHTTPHandleFunc(s.userHandler.VerifySignIn))

	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(WithJWTAuth)
		r.Get("/justCheck/{token}", util.MakeHTTPHandleFunc(s.userHandler.JustCheck))
		r.Post("/checkEmail", util.MakeHTTPHandleFunc(s.userHandler.CheckEmail))
		r.Get("/getUser/{userID}", util.MakeHTTPHandleFunc(s.userHandler.GetUserByID))
		r.Put("/updateUser", util.MakeHTTPHandleFunc(s.userHandler.UpdateUser))
		r.Get("/getUserbyEmail/{email}", util.MakeHTTPHandleFunc(s.userHandler.GetUserbyEmail))
		r.Post("/addNewFriend", util.MakeHTTPHandleFunc(s.userHandler.AddNewFriend))
		r.Delete("/removeFriend/{userID}/{friendID}/{userFriendID}", util.MakeHTTPHandleFunc(s.userHandler.RemoveFriend))
		r.Get("/getAllFriend/{userID}", util.MakeHTTPHandleFunc(s.userHandler.GetAllFriend))
		r.Post("/createPost", util.MakeHTTPHandleFunc(s.userHandler.CreatePost))
		r.Get("/getAllPost/{userID}", util.MakeHTTPHandleFunc(s.userHandler.GetAllPost))
		r.Get("/getAllOwnPost/{userID}", util.MakeHTTPHandleFunc(s.userHandler.GetAllOwnPost))
		r.Get("/getPost/{postID}", util.MakeHTTPHandleFunc(s.userHandler.GetPost))
		r.Get("/getAllImage/{userID}", util.MakeHTTPHandleFunc(s.userHandler.GetAllImage))
		r.Post("/createComment", util.MakeHTTPHandleFunc(s.userHandler.CreateComment))
		r.Get("/getAllComment/{postID}", util.MakeHTTPHandleFunc(s.userHandler.GetAllComment))
		r.Post("/createNotification", util.MakeHTTPHandleFunc(s.userHandler.CreateNotification))
		r.Put("/updateAddFriendNotification", util.MakeHTTPHandleFunc(s.userHandler.UpdateAddFriendNotification))
		r.Put("/updateNotificationRead/{userID}", util.MakeHTTPHandleFunc(s.userHandler.UpdateNotificationRead))
		r.Get("/getCountNotification/{userID}", util.MakeHTTPHandleFunc(s.userHandler.GetCountNotification))
		r.Get("/getAllNotification/{userID}", util.MakeHTTPHandleFunc(s.userHandler.GetAllNotification))

		// ws
		r.Post("/ws/createRoom", util.MakeHTTPHandleFunc(s.wsHandler.CreateRoom))
		r.Put("/ws/updateRoomName", util.MakeHTTPHandleFunc(s.wsHandler.UpdateRoomName))
		r.Post("/ws/addFriend", util.MakeHTTPHandleFunc(s.wsHandler.AddFriend))
		r.Get("/ws/getRoomsByUser/{userID}", util.MakeHTTPHandleFunc(s.wsHandler.GetRoomByUser))
		r.Get("/ws/getAllMessage/{roomID}", util.MakeHTTPHandleFunc(s.wsHandler.GetAllMessage))
		r.Get("/ws/getAllUnreadMessage/{userID}", util.MakeHTTPHandleFunc(s.wsHandler.CountAllUnreadMessage))
		r.Delete("/ws/remove/{roomID}/{userID}", util.MakeHTTPHandleFunc(s.wsHandler.Remove))
	})

	log.Println("server running in port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}
