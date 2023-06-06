package user

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

func NewUserRepository(db DBTX) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CheckEmail(email string) (*UserType, error) {
	acc := new(UserType)

	query := `select * from user where email = ?;`
	err := r.db.QueryRow(query, email).Scan(&acc.User_ID, &acc.User_Name, &acc.Email, &acc.Photo_Profile)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}

	if err != nil {
		return nil, err
	}

	return acc, nil
}

// Sign Up
func (r *Repository) SignUp(acc *UserType) (*UserType, error) {
	account := new(UserType)

	acc.User_ID = uuid.New().String()

	query := `insert into user(user_id, user_name, email, photo_profile) values (?, ?, ?, ?);`
	_, err := r.db.Exec(query, acc.User_ID, acc.User_Name, acc.Email, "")

	if err != nil {
		return nil, err
	}

	queryRes := `select * from user where email = ?;`
	err = r.db.QueryRow(queryRes, acc.Email).Scan(&account.User_ID, &account.User_Name, &account.Email, &account.Photo_Profile)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// get user
func (r *Repository) GetUser(user_id string) (*UserType, error) {
	u := new(UserType)

	query := `select user_id ,user_name, email, photo_profile from user where user_id = ?;`

	err := r.db.QueryRow(query, user_id).Scan(&u.User_ID, &u.User_Name, &u.Email, &u.Photo_Profile)
	if err == sql.ErrNoRows {
		log.Println("1. GetUser", err)
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		log.Println("2. GetUser", err)
		return nil, err
	}

	return u, nil
}

// update USER
func (r *Repository) UpdateUser(user *UserType) error {

	query := `update user set user_name = ?, email = ?, photo_profile = ? where user_id = ?;`

	_, err := r.db.Exec(query, user.User_Name, user.Email, user.Photo_Profile, user.User_ID)
	if err != nil {
		log.Println("1. UpdateUser", err)
		return err
	}

	return nil
}

// get user BY EMAIL
func (r *Repository) GetUserbyEmail(email string) (*UserType, error) {
	u := new(UserType)

	query := `select user_id ,user_name, email, photo_profile from user where email = ?;`

	err := r.db.QueryRow(query, email).Scan(&u.User_ID, &u.User_Name, &u.Email, &u.Photo_Profile)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}

// add friend
func (r *Repository) AddFriend(acc *User_FriendType) error {

	acc.User_Friend_ID = uuid.New().String()
	query := `insert into user_friend(user_friend_id, user_id, friend_id) values(?, ?, ?);`
	_, err := r.db.Exec(query, acc.User_Friend_ID, acc.User_ID, acc.Friend_ID)

	if err != nil {
		log.Println("1. AddFriend", err)
		return err
	}

	return nil
}

// check friend'
func (r *Repository) CheckFriend(user_id string) (int, error) {
	number := new(int)
	query := "select count(*) as `number` from user_friend where user_id = ?;"

	err := r.db.QueryRow(query, user_id).Scan(&number)
	if err != nil {
		log.Println("1. CheckFriend", err)
		return -1, err
	}

	return *number, nil
}

// remove friend
func (r *Repository) RemoveFriend(user_friend_id string) error {
	query := "delete from user_friend where user_friend_id = ?;"
	_, err := r.db.Exec(query, user_friend_id)

	if err != nil {
		log.Println("1. RemoveFriend", err)
		return err
	}

	return nil
}

// remove friend
func (r *Repository) RemoveFriendByUserID(user_id, friend_id string) error {
	query := "delete from user_friend where user_id = ? and friend_id = ?;"
	_, err := r.db.Exec(query, user_id, friend_id)

	if err != nil {
		log.Println("1. RemoveFriendByUserID", err)
		return err
	}

	return nil
}

// get all my friend
func (r *Repository) GetAllFriend(user_id string) ([]*UserFriendType, error) {
	query := "select user.user_id as `user_id`, user.user_name as `user_name`, user.email as `email`, user.photo_profile as `photo_profile`, user_friend.user_friend_id as `user_friend_id` from user_friend inner join user on user_friend.friend_id = user.user_id where user_friend.user_id = ?;"

	rows, err := r.db.Query(query, user_id)

	if err != nil {
		log.Println("1. GetAllFriend", err)
		return nil, err
	}

	defer rows.Close()

	friends := []*UserFriendType{}
	for rows.Next() {
		f := new(UserFriendType)

		if err := rows.Scan(&f.User_ID, &f.User_Name, &f.Email, &f.Photo_Profile, &f.User_Friend_ID); err != nil {
			log.Println("2. GetAllFriend", err)
			return nil, err
		}

		friends = append(friends, f)
	}

	if err := rows.Err(); err != nil {
		log.Println("3. GetAllFriend", err)
		return nil, err
	}

	return friends, nil
}

// create post
func (r *Repository) CreatePost(post *PostType) (string, error) {
	post.Post_ID = uuid.New().String()
	post.Created_At = time.Now().UTC()
	post.Updated_At = time.Now().UTC()

	query := `insert into post(post_id, user_id, content, type, created_at, updated_at) values (?, ?, ?, ?, ?, ?);`
	_, err := r.db.Exec(query, post.Post_ID, post.User_ID, post.Content, post.Type, post.Created_At, post.Updated_At)

	if err != nil {
		log.Println("1. CreatePost", err)
		return "", err
	}

	return post.Post_ID, nil
}

// get user
func (r *Repository) GetUserPost(u *GetPostResType) (*GetPostResType, error) {
	query := `select user_name, email, photo_profile from user where user_id = ?;`

	err := r.db.QueryRow(query, u.User_ID).Scan(&u.User_Name, &u.Email, &u.Photo_Profile)
	if err == sql.ErrNoRows {
		log.Println("1. GetUserPost", err)
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		log.Println("2. GetUserPost", err)
		return nil, err
	}

	return u, nil
}

// get ALl post
func (r *Repository) GetAllPost(user_id string) ([]*GetPostResType, error) {
	query := "select post.post_id as `post_id`, post.user_id as `user_id`, post.content as `content`, post.type as `type`, post.created_at as `created_at`, post.updated_at as `updated_at` from user_friend inner join post on user_friend.friend_id = post.user_id or post.user_id = ? where user_friend.user_id = ? and post.type = 'main';"

	rows, err := r.db.Query(query, user_id, user_id)

	if err != nil {
		log.Println("1. GetAllPost", err)
		return nil, err
	}

	defer rows.Close()

	allPost := []*GetPostResType{}
	for rows.Next() {
		p := new(GetPostResType)

		if err := rows.Scan(&p.Post_ID, &p.User_ID, &p.Content, &p.Type, &p.Created_At, &p.Updated_At); err != nil {
			log.Println("2. GetAllPost", err)
			return nil, err
		}

		p, err = r.GetImagePost(p)

		if err != nil {
			log.Println("3. GetAllPost", err)
			return nil, err
		}

		p, err = r.GetCountPost(p)
		if err != nil {
			log.Println("4. GetAllPost", err)
			return nil, err
		}

		p, err = r.GetUserPost(p)
		if err != nil {
			log.Println("5. GetAllPost", err)
			return nil, err
		}

		allPost = append(allPost, p)
	}

	if err := rows.Err(); err != nil {
		log.Println("6. GetAllPost", err)
		return nil, err
	}

	return allPost, nil
}

// get ALl OWN post
func (r *Repository) GetAllOwnPost(user_id string) ([]*GetPostResType, error) {
	query := "select * from post where user_id = ? and type = 'main';"

	rows, err := r.db.Query(query, user_id)

	if err != nil {
		log.Println("1. GetAllOwnPost", err)
		return nil, err
	}

	defer rows.Close()

	allPost := []*GetPostResType{}
	for rows.Next() {
		p := new(GetPostResType)

		if err := rows.Scan(&p.Post_ID, &p.User_ID, &p.Content, &p.Type, &p.Created_At, &p.Updated_At); err != nil {
			log.Println("2. GetAllOwnPost", err)
			return nil, err
		}

		p, err = r.GetImagePost(p)

		if err != nil {
			log.Println("3. GetAllOwnPost", err)
			return nil, err
		}

		p, err = r.GetCountPost(p)
		if err != nil {
			log.Println("4. GetAllOwnPost", err)
			return nil, err
		}

		p, err = r.GetUserPost(p)
		if err != nil {
			log.Println("5. GetAllOwnPost", err)
			return nil, err
		}

		allPost = append(allPost, p)
	}

	if err := rows.Err(); err != nil {
		log.Println("6. GetAllOwnPost", err)
		return nil, err
	}

	return allPost, nil
}

// get Single post
func (r *Repository) GetPost(post_id string) (*GetPostResType, error) {
	res := new(GetPostResType)

	query := `select * from post where post_id = ?;`
	err := r.db.QueryRow(query, post_id).Scan(&res.Post_ID, &res.User_ID, &res.Content, &res.Type, &res.Created_At, &res.Updated_At)

	if err == sql.ErrNoRows {
		log.Println("1. GetPost", err)
		return nil, fmt.Errorf("post not found")
	}

	if err != nil {
		log.Println("2. GetPost", err)
		return nil, err
	}

	res, err = r.GetImagePost(res)
	if err != nil {
		log.Println("3. GetPost", err)
		return nil, err
	}

	res, err = r.GetCountPost(res)
	if err != nil {
		log.Println("4. GetPost", err)
		return nil, err
	}

	res, err = r.GetUserPost(res)
	if err != nil {
		log.Println("5. GetPost", err)
		return nil, err
	}

	return res, nil
}

// get count comment per Post
func (r *Repository) GetCountPost(post *GetPostResType) (*GetPostResType, error) {
	query := "select count(*) as `number_of_comment` from comment where post_id = ?;"

	err := r.db.QueryRow(query, post.Post_ID).Scan(&post.Number_Of_Comment)
	if err != nil {
		log.Println("1. GetCountPost", err)
		return nil, err
	}

	return post, nil
}

// get imageforPost
func (r *Repository) GetImagePost(img *GetPostResType) (*GetPostResType, error) {
	query := `select * from image_post where post_id = ?;`

	rows, err := r.db.Query(query, img.Post_ID)

	if err != nil {
		log.Println("1. GetImagePost", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		i := new(Image_PostType)

		if err := rows.Scan(&i.Image_Post_ID, &i.Post_ID, &i.User_ID, &i.Image, &i.Created_At, &i.Updated_At); err != nil {
			log.Println("2. GetImagePost", err)
			return nil, err
		}

		img.Images = append(img.Images, i)
	}

	if err := rows.Err(); err != nil {
		log.Println("3. GetImagePost", err)
		return nil, err
	}

	return img, nil
}

// create Image post
func (r *Repository) CreateImagePost(img *Image_PostType) error {

	img.Image_Post_ID = uuid.New().String()
	img.Created_At = time.Now().UTC()
	img.Updated_At = time.Now().UTC()

	query := `insert into image_post(image_post_id, post_id, user_id, image, created_at, updated_at) values (?, ?, ?, ?, ?, ?);`
	_, err := r.db.Exec(query, img.Image_Post_ID, img.Post_ID, img.User_ID, img.Image, img.Created_At, img.Updated_At)
	if err != nil {
		log.Println("1. CreateImagePost", err)
		return err
	}

	return nil
}

// Get all image
func (r *Repository) GetAllImage(user_id string) ([]*Image_PostType, error) {
	query := `select * from image_post where user_id = ?;`

	rows, err := r.db.Query(query, user_id)
	if err != nil {
		log.Println("1. GetAllImage", err)
		return nil, err
	}

	defer rows.Close()

	images := []*Image_PostType{}
	for rows.Next() {
		i := new(Image_PostType)

		if err := rows.Scan(&i.Image_Post_ID, &i.Post_ID, &i.User_ID, &i.Image, &i.Created_At, &i.Updated_At); err != nil {
			log.Println("2. GetAllImage", err)
			return nil, err
		}

		images = append(images, i)
	}

	if err := rows.Err(); err != nil {
		log.Println("3. GetAllImage", err)
		return nil, err
	}

	return images, nil
}

// create comment
func (r *Repository) CreateComment(comment *CommentType) error {
	comment.Comment_ID = uuid.New().String()
	comment.Created_At = time.Now().UTC()
	comment.Updated_At = time.Now().UTC()

	query := `insert into comment(comment_id, post_id, comment_post_id, created_at, updated_at) values (?, ?, ?, ?, ?);`
	_, err := r.db.Exec(query, comment.Comment_ID, comment.Post_ID, comment.Comment_Post_ID, comment.Created_At, comment.Updated_At)
	if err != nil {
		log.Println("1. CreateComment", err)
		return err
	}

	return nil
}

// get All Comment
func (r *Repository) GetAllComment(post_id string) ([]*GetPostResType, error) {
	query := "select post.post_id as `post_id`, post.user_id as `user_id`, post.content as `content`, post.type as `type`, post.created_at as `created_at`, post.updated_at as `updated_at` from comment inner join post on comment.comment_post_id = post.post_id where comment.post_id = ?;"

	rows, err := r.db.Query(query, post_id)
	if err != nil {
		log.Println("1. GetAllComment", err)
		return nil, err
	}

	defer rows.Close()

	allComment := []*GetPostResType{}
	for rows.Next() {
		c := new(GetPostResType)

		if err := rows.Scan(&c.Post_ID, &c.User_ID, &c.Content, &c.Type, &c.Created_At, &c.Updated_At); err != nil {
			log.Println("2. GetAllComment", err)
			return nil, err
		}

		c, err = r.GetImagePost(c)
		if err != nil {
			log.Println("3. GetAllComment", err)
			return nil, err
		}

		c, err = r.GetCountPost(c)
		if err != nil {
			log.Println("4. GetAllComment", err)
			return nil, err
		}

		c, err = r.GetUserPost(c)
		if err != nil {
			log.Println("5. GetAllComment", err)
			return nil, err
		}

		allComment = append(allComment, c)
	}

	if err := rows.Err(); err != nil {
		log.Println("6. GetAllComment", err)
		return nil, err
	}

	return allComment, nil
}

// create notification
func (r *Repository) CreateNotification(notif *NotificationType) error {
	notif.Notification_ID = uuid.New().String()
	notif.Created_At = time.Now().UTC()
	notif.Updated_At = time.Now().UTC()

	query := `insert into notification(notification_id, issuer, issuer_name, notifier, notifier_name, status, accept, post_id, type, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	_, err := r.db.Exec(query, notif.Notification_ID, notif.Issuer, notif.Issuer_Name, notif.Notifier, notif.Notifier_Name, notif.Status, notif.Accept, notif.Post_ID, notif.Type, notif.Created_At, notif.Updated_At)
	if err != nil {
		log.Println("1. CreateNotification", err)
		return err
	}

	return nil
}

// update friend notif
func (r *Repository) UpdateNotif(notif *UpdateNotifType) error {
	notif.Updated_At = time.Now().UTC()

	query := `update notification set accept = ?, updated_at = ? where notification_id = ?;`
	_, err := r.db.Exec(query, notif.Accept, notif.Updated_At, notif.Notification_ID)
	if err != nil {
		log.Println("1. UpdateNotif", err)
		return err
	}

	return nil
}

// updated become read
func (r *Repository) UpdatedNotifRead(user_id string) error {

	query := `update notification set status = "read", updated_at = ? where notifier = ? and status = "not_read";`
	_, err := r.db.Exec(query, time.Now().UTC(), user_id)
	if err != nil {
		log.Println("1. UpdatedNotifRead", err)
		return err
	}

	return nil
}

// get count Notif
func (r *Repository) GetCountNotif(user_id string) (int, error) {
	var number int
	query := "select count(*) as `number` from notification where notifier = ? and status = 'not_read';"

	err := r.db.QueryRow(query, user_id).Scan(&number)

	if err != nil {
		log.Println("1. GetCountNotif", err)
		return -1, err
	}

	return number, nil
}

// get All notif
func (r *Repository) GetAllNotif(user_id string) ([]*NotificationType, error) {
	query := `select * from notification where notifier = ? order by created_at desc;`

	rows, err := r.db.Query(query, user_id)

	if err != nil {
		log.Println("1. GetAllNotif", err)
		return nil, err
	}

	defer rows.Close()

	notifs := []*NotificationType{}
	for rows.Next() {
		newNotif := new(NotificationType)

		if err := rows.Scan(&newNotif.Notification_ID, &newNotif.Issuer, &newNotif.Issuer_Name, &newNotif.Notifier, &newNotif.Notifier_Name, &newNotif.Status, &newNotif.Accept, &newNotif.Post_ID, &newNotif.Type, &newNotif.Created_At, &newNotif.Updated_At); err != nil {
			log.Println("2. GetAllNotif", err)
			return nil, err
		}

		notifs = append(notifs, newNotif)
	}

	if err := rows.Err(); err != nil {
		log.Println("3. GetAllNotif", err)
		return nil, err
	}

	return notifs, nil
}
