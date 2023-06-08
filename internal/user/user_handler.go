package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/erlnerlngga/backend-socius/util"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
)

type Handler struct {
	Repository *Repository
}

func NewUserHandler(r *Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) Welcome(w http.ResponseWriter, r *http.Request) error {
	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "Welcome to Socius"})
}

func (h *Handler) CheckEmail(w http.ResponseWriter, r *http.Request) error {
	acc := new(SignInType)
	if err := json.NewDecoder(r.Body).Decode(acc); err != nil {
		log.Println("1. CheckEmail", err)
		return err
	}

	defer r.Body.Close()

	email, err := h.Repository.CheckEmail(acc.Email)
	if err != nil {
		log.Println("2. CheckEmail", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, email)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) error {
	acc := new(UserType)

	if err := json.NewDecoder(r.Body).Decode(acc); err != nil {
		log.Println("1. SignUp", err)
		return err
	}

	defer r.Body.Close()

	newAcc, err := h.Repository.SignUp(acc)
	if err != nil {
		log.Println("2. SignUp", err)
		return err
	}

	tokenStr, err := util.CreateJWT(newAcc.User_ID)
	if err != nil {
		log.Println("3. SignUp", err)
		return err
	}

	if err := util.SendMAIL(newAcc.Email, newAcc.User_Name, tokenStr); err != nil {
		log.Println("4. SignUp", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) error {
	email := new(SignInType)

	if err := json.NewDecoder(r.Body).Decode(email); err != nil {
		log.Println("1. SignIn", err)
		return err
	}

	defer r.Body.Close()

	account, err := h.Repository.CheckEmail(email.Email)
	if err != nil {
		log.Println("2. SignIn", err)
		return err
	}

	// create token
	tokenStr, err := util.CreateJWT(account.User_ID)
	if err != nil {
		log.Println("3. SignIn", err)
		return err
	}

	if err := util.SendMAIL(account.Email, account.User_Name, tokenStr); err != nil {
		log.Println("4. SignIn", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (h *Handler) VerifySignIn(w http.ResponseWriter, r *http.Request) error {
	tokenStr := chi.URLParam(r, "token")

	claims := new(util.ClaimsType)

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return util.JwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("1. VerifySignIn", err)
			return util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: "signature invalid"})
		}
		log.Println("2. VerifySignIn", err)
		return util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: err.Error()})
	}

	if !token.Valid {
		return util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: "token invalid"})
	}

	user, err := h.Repository.GetUser(claims.User_ID)
	if err != nil {
		log.Println("3. VerifySignIn", err)
		return err
	}

	resultVer := &VerifyResType{
		Status: "ok",
		Token:  tokenStr,
		User:   user,
	}

	return util.WriteJSON(w, http.StatusOK, resultVer)
}

func (h *Handler) JustCheck(w http.ResponseWriter, r *http.Request) error {
	tokenStr := chi.URLParam(r, "token")

	claims := new(util.ClaimsType)

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return util.JwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("1. JustCheck", err)
			return util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: "signature invalid"})
		}
		log.Println("2. JustCheck", err)
		return util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: err.Error()})
	}

	if !token.Valid {
		return util.WriteJSON(w, http.StatusUnauthorized, util.ApiError{Error: "token invalid"})
	}

	user, err := h.Repository.GetUser(claims.User_ID)
	if err != nil {
		log.Println("3. JustCheck", err)
		return err
	}

	resultVer := &VerifyResType{
		Status: "ok",
		Token:  tokenStr,
		User:   user,
	}

	return util.WriteJSON(w, http.StatusOK, resultVer)
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	user, err := h.Repository.GetUser(userID)
	if err != nil {
		log.Println("1. GetUserByID", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	userUp := new(UserType)

	if err := json.NewDecoder(r.Body).Decode(userUp); err != nil {
		log.Println("1. UpdateUser", err)
		return err
	}

	defer r.Body.Close()

	log.Println(userUp)

	err := h.Repository.UpdateUser(userUp)
	if err != nil {
		log.Println("2. UpdateUser", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) GetUserbyEmail(w http.ResponseWriter, r *http.Request) error {
	email := chi.URLParam(r, "email")

	user, err := h.Repository.GetUserbyEmail(email)
	if err != nil {
		log.Println("1. GetUserbyEmail", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) AddNewFriend(w http.ResponseWriter, r *http.Request) error {
	newFr := new(User_FriendType)

	if err := json.NewDecoder(r.Body).Decode(newFr); err != nil {
		log.Println("1. AddNewFriend", err)
		return err
	}

	defer r.Body.Close()

	err := h.Repository.AddFriend(newFr)
	if err != nil {
		log.Println("2. AddNewFriend", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) RemoveFriend(w http.ResponseWriter, r *http.Request) error {
	userFriendId := chi.URLParam(r, "userFriendID")
	userID := chi.URLParam(r, "userID")
	friendID := chi.URLParam(r, "friendID")

	err := h.Repository.RemoveFriend(userFriendId)
	if err != nil {
		log.Println("1. RemoveFriend", err)
		return err
	}

	err = h.Repository.RemoveFriendByUserID(friendID, userID)
	if err != nil {
		log.Println("2. RemoveFriend", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (h *Handler) GetAllFriend(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	friends, err := h.Repository.GetAllFriend(userID)
	if err != nil {
		log.Println("1. GetAllFriend", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, friends)
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) error {
	newPost := new(PostReqType)

	if err := json.NewDecoder(r.Body).Decode(newPost); err != nil {
		log.Println("1. CreatePost", err)
		return err
	}

	defer r.Body.Close()

	p := &PostType{
		User_ID: newPost.User_ID,
		Content: newPost.Content,
		Type:    "main",
	}

	post_ID, err := h.Repository.CreatePost(p)
	if err != nil {
		log.Println("2. CreatePost", err)
		return err
	}

	if len(newPost.Images) > 0 {
		for _, val := range newPost.Images {
			im := &Image_PostType{
				Post_ID: post_ID,
				User_ID: p.User_ID,
				Image:   val,
			}

			err := h.Repository.CreateImagePost(im)

			if err != nil {
				log.Println("3. CreatePost", err)
				return err
			}
		}
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (h *Handler) GetAllPost(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")
	var post []*GetPostResType

	number, err := h.Repository.CheckFriend(userID)
	if err != nil {
		log.Println("1. GetAllPost", err)
		return err
	}

	if number > 0 {
		post, err = h.Repository.GetAllPost(userID)
		if err != nil {
			log.Println("2. GetAllPost", err)
			return err
		}
	} else {
		post, err = h.Repository.GetAllOwnPost(userID)
		if err != nil {
			log.Println("3. GetAllPost", err)
			return err
		}
	}

	return util.WriteJSON(w, http.StatusOK, post)
}

func (h *Handler) GetAllOwnPost(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	post, err := h.Repository.GetAllOwnPost(userID)
	if err != nil {
		log.Println("1. GetAllOwnPost", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, post)
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) error {
	postID := chi.URLParam(r, "postID")

	post, err := h.Repository.GetPost(postID)
	if err != nil {
		log.Println("1. GetPost", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, post)
}

func (h *Handler) GetAllImage(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	images, err := h.Repository.GetAllImage(userID)
	if err != nil {
		log.Println("1. GetAllImage", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, images)
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) error {
	newPost := new(CommentReqType)

	if err := json.NewDecoder(r.Body).Decode(newPost); err != nil {
		log.Println("1. CreateComment", err)
		return err
	}

	defer r.Body.Close()

	p := &PostType{
		User_ID: newPost.User_ID,
		Content: newPost.Content,
		Type:    "child",
	}

	post_ID, err := h.Repository.CreatePost(p)
	if err != nil {
		log.Println("2. CreateComment", err)
		return err
	}

	if len(newPost.Images) > 0 {
		for _, val := range newPost.Images {
			im := &Image_PostType{
				Post_ID: post_ID,
				User_ID: p.User_ID,
				Image:   val,
			}

			err := h.Repository.CreateImagePost(im)

			if err != nil {
				log.Println("3. CreateComment", err)
				return err
			}
		}
	}

	commen := &CommentType{
		Post_ID:         newPost.Post_ID,
		Comment_Post_ID: post_ID,
	}

	// insert comment
	err = h.Repository.CreateComment(commen)
	if err != nil {
		log.Println("4. CreateComment", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) GetAllComment(w http.ResponseWriter, r *http.Request) error {
	postID := chi.URLParam(r, "postID")

	comment, err := h.Repository.GetAllComment(postID)
	if err != nil {
		log.Println("1. GetAllComment", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, comment)
}

func (h *Handler) CreateNotification(w http.ResponseWriter, r *http.Request) error {
	not := new(NotificationType)

	if err := json.NewDecoder(r.Body).Decode(not); err != nil {
		log.Println("1. CreateNotification", err)
		return err
	}

	defer r.Body.Close()

	err := h.Repository.CreateNotification(not)
	if err != nil {
		log.Println("2. CreateNotification", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) UpdateAddFriendNotification(w http.ResponseWriter, r *http.Request) error {
	newNotif := new(UpdateNotifType)

	if err := json.NewDecoder(r.Body).Decode(newNotif); err != nil {
		log.Println("1. UpdateAddFriendNotification", err)
		return err
	}

	defer r.Body.Close()

	err := h.Repository.UpdateNotif(newNotif)
	if err != nil {
		log.Println("2. UpdateAddFriendNotification", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) UpdateNotificationRead(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	err := h.Repository.UpdatedNotifRead(userID)
	if err != nil {
		log.Println("1. UpdateNotificationRead", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) GetCountNotification(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	num, err := h.Repository.GetCountNotif(userID)
	if err != nil {
		log.Println("1. GetCountNotification", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]int{"number": num})
}

func (h *Handler) GetAllNotification(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	notif, err := h.Repository.GetAllNotif(userID)
	if err != nil {
		log.Println("1. GetAllNotification", err)
		return err
	}

	return util.WriteJSON(w, http.StatusOK, notif)
}
