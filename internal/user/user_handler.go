package user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/erlnerlngga/backend-socius/util"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Repository *Repository
}

func NewUserHandler(r *Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) CheckEmail(w http.ResponseWriter, r *http.Request) error {
	acc := new(SignInType)
	if err := json.NewDecoder(r.Body).Decode(acc); err != nil {
		return err
	}

	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()
	email, err := h.Repository.CheckEmail(ctx, acc.Email)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, email)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) error {
	acc := new(UserType)

	if err := json.NewDecoder(r.Body).Decode(acc); err != nil {
		return err
	}

	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()
	newAcc, err := h.Repository.SignUp(ctx, acc)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, newAcc)
}

func (h *Handler) AddNewFriend(w http.ResponseWriter, r *http.Request) error {
	newFr := new(User_FriendType)

	if err := json.NewDecoder(r.Body).Decode(newFr); err != nil {
		return err
	}

	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()
	fr, err := h.Repository.AddFriend(ctx, newFr)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, fr)
}

func (h *Handler) RemoveFriend(w http.ResponseWriter, r *http.Request) error {
	userFriendId := chi.URLParam(r, "userFriendID")

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()
	err := h.Repository.RemoveFriend(ctx, userFriendId)

	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]string{"status": "succes"})
}

func (h *Handler) GetAllFriend(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()

	friends, err := h.Repository.GetAllFriend(ctx, userID)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, friends)
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) error {
	newPost := new(PostReqType)

	if err := json.NewDecoder(r.Body).Decode(newPost); err != nil {
		return err
	}

	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()

	p := &PostType{
		User_ID: newPost.User_ID,
		Content: newPost.Content,
	}

	post, err := h.Repository.CreatePost(ctx, p)
	if err != nil {
		return err
	}

	if len(newPost.Images) > 0 {
		for _, val := range newPost.Images {
			im := &Image_PostType{
				Post_ID: post.Post_ID,
				User_ID: p.User_ID,
				Image:   val,
			}

			_, err := h.Repository.CreateImagePost(ctx, im)

			if err != nil {
				return err
			}
		}
	}

	return util.WriteJSON(w, http.StatusOK, post)
}

func (h *Handler) GetAllPost(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(3)*time.Second)
	defer cancel()

	post, err := h.Repository.GetAllPost(ctx, userID)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, post)
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) error {
	postID := chi.URLParam(r, "postID")

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()

	post, err := h.Repository.GetPost(ctx, postID)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, post)
}

func (h *Handler) GetAllImage(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()

	images, err := h.Repository.GetAllImage(ctx, userID)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, images)
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) error {
	newPost := new(CommentReqType)

	if err := json.NewDecoder(r.Body).Decode(newPost); err != nil {
		return err
	}

	defer r.Body.Close()
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()

	p := &PostType{
		User_ID: newPost.User_ID,
		Content: newPost.Content,
	}

	post, err := h.Repository.CreatePost(ctx, p)
	if err != nil {
		return err
	}

	if len(newPost.Images) > 0 {
		for _, val := range newPost.Images {
			im := &Image_PostType{
				Post_ID: post.Post_ID,
				User_ID: p.User_ID,
				Image:   val,
			}

			_, err := h.Repository.CreateImagePost(ctx, im)

			if err != nil {
				return err
			}
		}
	}

	commen := &CommentType{
		Post_ID:         newPost.Post_ID,
		Comment_Post_ID: post.Post_ID,
	}

	// insert comment
	com, err := h.Repository.CreateComment(ctx, commen)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, com)
}

func (h *Handler) GetAllComment(w http.ResponseWriter, r *http.Request) error {
	postID := chi.URLParam(r, "postID")

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()

	comment, err := h.Repository.GetAllComment(ctx, postID)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, comment)
}

func (h *Handler) CreateNotification(w http.ResponseWriter, r *http.Request) error {
	not := new(NotificationType)

	if err := json.NewDecoder(r.Body).Decode(not); err != nil {
		return err
	}

	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)

	defer cancel()

	notif, err := h.Repository.CreateNotification(ctx, not)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, notif)
}

func (h *Handler) UpdateAddFriendNotification(w http.ResponseWriter, r *http.Request) error {
	newNotif := new(UpdateNotifType)

	if err := json.NewDecoder(r.Body).Decode(newNotif); err != nil {
		return err
	}

	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()

	notif, err := h.Repository.UpdateNotif(ctx, newNotif)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, notif)
}

func (h *Handler) UpdateNotificationRead(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()

	notif, err := h.Repository.UpdatedNotifRead(ctx, userID)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, notif)
}

func (h *Handler) GetCountNotification(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()

	num, err := h.Repository.GetCountNotif(ctx, userID)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]int{"number": num})
}

func (h *Handler) GetAllNotification(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(2)*time.Second)
	defer cancel()

	notif, err := h.Repository.GetAllNotif(ctx, userID)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, notif)
}
