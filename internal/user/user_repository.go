package user

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

func NewUserRepository(db DBTX) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CheckEmail(ctx context.Context, email string) (*UserType, error) {
	acc := new(UserType)

	query := `select * from user where email = ?;`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&acc.User_ID, &acc.User_Name, &acc.Email, &acc.Photo_Profile)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}

	if err != nil {
		return nil, err
	}

	return acc, nil
}

// Sign Up
func (r *Repository) SignUp(ctx context.Context, acc *UserType) (*UserType, error) {
	account := new(UserType)

	acc.User_ID = uuid.New().String()

	query := `insert into user(user_id, user_name, email) values (?, ?, ?);`
	err := r.db.QueryRowContext(ctx, query, acc.User_ID, acc.User_Name, acc.Email).Scan(&account.User_ID, &account.User_Name, &acc.Email, &acc.Photo_Profile)

	if err != nil {
		return nil, err
	}

	return account, nil
}

// add friend
func (r *Repository) AddFriend(ctx context.Context, acc *User_FriendType) (*User_FriendType, error) {
	newFriend := new(User_FriendType)

	acc.User_Friend_ID = uuid.New().String()
	query := `insert into user_friend(user_friend_id, user_id, friend_id) values(?, ?, ?);`
	err := r.db.QueryRowContext(ctx, query, acc.User_Friend_ID, acc.User_ID, acc.Friend_ID).Scan(&newFriend.User_Friend_ID, &newFriend.User_ID, &newFriend.Friend_ID)

	if err != nil {
		return nil, err
	}

	return newFriend, nil
}

// remove friend
func (r *Repository) RemoveFriend(ctx context.Context, user_friend_id string) error {
	_, err := r.db.ExecContext(ctx, `delete from user_friend where user_friend_id = ?;`, user_friend_id)

	if err != nil {
		return err
	}

	return nil
}

// get all my friend
func (r *Repository) GetAllFriend(ctx context.Context, user_id string) ([]*UserFriendType, error) {
	query := "select user.user_id as `user_id`, user.user_name as `user_name`, user.email as `email`, user.photo_profile as `photo_profile`, user_friend.user_friend_id as `user_friend_id` from user_friend inner join user on user_friend.friend_id = user.user_id where user_friend.user_id = ?;"

	rows, err := r.db.QueryContext(ctx, query, user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	friends := []*UserFriendType{}
	for rows.Next() {
		f := new(UserFriendType)

		if err := rows.Scan(&f.User_ID, &f.User_Name, &f.Email, &f.Photo_Profile, &f.User_Friend_ID); err != nil {
			return nil, err
		}

		friends = append(friends, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return friends, nil
}

// create post
func (r *Repository) CreatePost(ctx context.Context, post *PostType) (*PostType, error) {
	newPost := new(PostType)
	post.Post_ID = uuid.New().String()
	post.Created_At = time.Now().UTC()
	post.Updated_At = time.Now().UTC()

	query := `insert into post(post_id, user_id, content, created_at, updated_at) values (?, ?, ?, ?, ?);`
	err := r.db.QueryRowContext(ctx, query, post.Post_ID, post.User_ID, post.Content, post.Created_At, post.Updated_At).Scan(&newPost.Post_ID, &newPost.User_ID, &newPost.Content, &newPost.Created_At, &newPost.Updated_At)
	if err != nil {
		return nil, err
	}

	return newPost, err
}

// get ALl post
func (r *Repository) GetAllPost(ctx context.Context, user_id string) ([]*GetPostType, error) {
	query := "select post.post_id as `post_id`, post.user_id as `user_id`, post.content as `content`, post.created_at as `created_at`, post.updated_at as `updated_at` from user_friend inner join post on user_friend.friend_id = post.user_id where user_friend.user_id = ?;"

	rows, err := r.db.QueryContext(ctx, query, user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	allPost := []*GetPostType{}
	for rows.Next() {
		p := new(GetPostType)

		if err := rows.Scan(&p.Post_ID, &p.User_ID, &p.Content, &p.Created_At, &p.Updated_At); err != nil {
			return nil, err
		}

		p, err = r.GetImagePost(ctx, p)

		if err != nil {
			return nil, err
		}

		p, err = r.GetCountPost(ctx, p)
		if err != nil {
			return nil, err
		}

		allPost = append(allPost, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return allPost, nil
}

// get Single post
func (r *Repository) GetPost(ctx context.Context, post_id string) (*GetPostType, error) {
	res := new(GetPostType)

	query := `select * from post where post_id = ?;`
	err := r.db.QueryRowContext(ctx, query, post_id).Scan(&res.Post_ID, &res.User_ID, &res.Content, &res.Created_At, &res.Updated_At)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found")
	}

	if err != nil {
		return nil, err
	}

	res, err = r.GetImagePost(ctx, res)
	if err != nil {
		return nil, err
	}

	res, err = r.GetCountPost(ctx, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// get count comment per Post
func (r *Repository) GetCountPost(ctx context.Context, post *GetPostType) (*GetPostType, error) {
	query := "select count(*) as `number_of_comment` from comment where post_id = ?;"

	err := r.db.QueryRowContext(ctx, query, post.Post_ID).Scan(&post.Number_Of_Comment)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// get imageforPost
func (r *Repository) GetImagePost(ctx context.Context, img *GetPostType) (*GetPostType, error) {
	query := `select * from image_post where post_id = ?;`

	rows, err := r.db.QueryContext(ctx, query, img.Post_ID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		i := new(Image_PostType)

		if err := rows.Scan(&i.Image_Post_ID, &i.Post_ID, &i.User_ID, &i.Image, &i.Created_At, &i.Updated_At); err != nil {
			return nil, err
		}

		img.Images = append(img.Images, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return img, nil
}

// create Image post
func (r *Repository) CreateImagePost(ctx context.Context, img *Image_PostType) (*Image_PostType, error) {
	newImg := new(Image_PostType)
	img.Image_Post_ID = uuid.New().String()
	img.Created_At = time.Now().UTC()
	img.Updated_At = time.Now().UTC()

	query := `insert into image_post(image_post_id, post_id, user_id, image, created_at, updated_at) values (?, ?, ?, ?, ?, ?);`
	err := r.db.QueryRowContext(ctx, query, img.Image_Post_ID, img.Post_ID, img.User_ID, img.Image, img.Created_At, img.Updated_At).Scan(&newImg.Image_Post_ID, &newImg.Post_ID, &newImg.User_ID, &newImg.Image, &newImg.Created_At, &newImg.Updated_At)
	if err != nil {
		return nil, err
	}

	return newImg, err
}

// Get all image
func (r *Repository) GetAllImage(ctx context.Context, user_id string) ([]*Image_PostType, error) {
	query := `select * from image_post where user_id = ?;`

	rows, err := r.db.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	images := []*Image_PostType{}
	for rows.Next() {
		i := new(Image_PostType)

		if err := rows.Scan(&i.Image_Post_ID, &i.Post_ID, &i.User_ID, &i.Image, &i.Created_At, &i.Updated_At); err != nil {
			return nil, err
		}

		images = append(images, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

// create comment
func (r *Repository) CreateComment(ctx context.Context, comment *CommentType) (*CommentType, error) {
	newCom := new(CommentType)
	comment.Comment_ID = uuid.New().String()
	comment.Created_At = time.Now().UTC()
	comment.Updated_At = time.Now().UTC()

	query := `insert into comment(comment_id, post_id, comment_post_id, created_at, updated_at) values (?, ?, ?, ?, ?);`
	err := r.db.QueryRowContext(ctx, query, comment.Comment_ID, comment.Post_ID, comment.Comment_Post_ID, comment.Created_At, comment.Updated_At).Scan(&newCom.Comment_ID, &newCom.Post_ID, &newCom.Comment_Post_ID, &newCom.Created_At, &newCom.Updated_At)
	if err != nil {
		return nil, err
	}

	return newCom, err
}

// get All Comment
func (r *Repository) GetAllComment(ctx context.Context, post_id string) ([]*GetPostType, error) {
	query := "select post.post_id as `post_id`, post.user_id as `user_id`, post.content as `content`, post.created_at as `created_at`, post.updated_at as `updated_at` from comment inner join post on comment.comment_post_id = post.post_id where comment.post_id = ?;"

	rows, err := r.db.QueryContext(ctx, query, post_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	allComment := []*GetPostType{}
	for rows.Next() {
		c := new(GetPostType)

		if err := rows.Scan(&c.Post_ID, &c.User_ID, &c.Content, &c.Created_At, &c.Updated_At); err != nil {
			return nil, err
		}

		c, err = r.GetImagePost(ctx, c)
		if err != nil {
			return nil, err
		}

		c, err = r.GetCountPost(ctx, c)
		if err != nil {
			return nil, err
		}

		allComment = append(allComment, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return allComment, nil
}

// create notification
func (r *Repository) CreateNotification(ctx context.Context, notif *NotificationType) (*NotificationType, error) {
	newNotif := new(NotificationType)
	notif.Notification_ID = uuid.New().String()
	notif.Created_At = time.Now().UTC()
	notif.Updated_At = time.Now().UTC()

	query := `insert into notification(notification_id, issuer, notifier, notifier_name, status, accept, post_id, type, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	err := r.db.QueryRowContext(ctx, query, notif.Notification_ID, notif.Issuer, notif.Notifier, notif.Notifier_Name, notif.Status, notif.Accept, notif.Post_ID, notif.Type, notif.Created_At, notif.Updated_At).Scan(&newNotif.Notification_ID, &newNotif.Issuer, &newNotif.Notifier, &newNotif.Notifier_Name, &newNotif.Status, &newNotif.Accept, &newNotif.Post_ID, &newNotif.Type, &newNotif.Created_At, &newNotif.Updated_At)
	if err != nil {
		return nil, err
	}

	return newNotif, nil
}

// update friend notif
func (r *Repository) UpdateNotif(ctx context.Context, notif *UpdateNotifType) (*NotificationType, error) {
	newNotif := new(NotificationType)
	notif.Updated_At = time.Now().UTC()

	query := `update notification set accept = ?, updated_at = ? where notification_id = ?;`
	err := r.db.QueryRowContext(ctx, query, notif.Accept, notif.Updated_At, notif.Notification_ID).Scan(&newNotif.Notification_ID, &newNotif.Issuer, &newNotif.Notifier, &newNotif.Status, &newNotif.Accept, &newNotif.Post_ID, &newNotif.Type, &newNotif.Created_At, &newNotif.Updated_At)
	if err != nil {
		return nil, err
	}

	return newNotif, nil
}

// updated become read
func (r *Repository) UpdatedNotifRead(ctx context.Context, user_id string) (*NotificationType, error) {
	newNotif := new(NotificationType)

	query := `update notification set status = "read", updated_at = ? where notifier = ? and status = "not_read";`
	err := r.db.QueryRowContext(ctx, query, time.Now().UTC(), user_id).Scan(&newNotif.Notification_ID, &newNotif.Issuer, &newNotif.Notifier, &newNotif.Status, &newNotif.Accept, &newNotif.Post_ID, &newNotif.Type, &newNotif.Created_At, &newNotif.Updated_At)
	if err != nil {
		return nil, err
	}

	return newNotif, nil
}

// get count Notif
func (r *Repository) GetCountNotif(ctx context.Context, user_id string) (int, error) {
	var number int
	query := "select count(*) as `number` from notification where notifier = ? and status = 'not_read';"

	err := r.db.QueryRowContext(ctx, query, user_id).Scan(&number)

	if err != nil {
		return -1, err
	}

	return number, nil
}

// get All notif
func (r *Repository) GetAllNotif(ctx context.Context, user_id string) ([]*NotificationType, error) {
	query := `select * from notification where notifier = ? order by created_at desc;`

	rows, err := r.db.QueryContext(ctx, query, user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	notifs := []*NotificationType{}
	for rows.Next() {
		newNotif := new(NotificationType)

		if err := rows.Scan(&newNotif.Notification_ID, &newNotif.Issuer, &newNotif.Notifier, &newNotif.Status, &newNotif.Accept, &newNotif.Post_ID, &newNotif.Type, &newNotif.Created_At, &newNotif.Updated_At); err != nil {
			return nil, err
		}

		notifs = append(notifs, newNotif)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifs, nil
}
