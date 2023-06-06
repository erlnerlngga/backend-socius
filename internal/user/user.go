package user

import "time"

type UserType struct {
	User_ID       string `json:"user_id"`
	User_Name     string `json:"user_name"`
	Email         string `json:"email"`
	Photo_Profile string `json:"photo_profile"`
}

type UserFriendType struct {
	User_ID        string `json:"user_id"`
	User_Name      string `json:"user_name"`
	Email          string `json:"email"`
	Photo_Profile  string `json:"photo_profile"`
	User_Friend_ID string `json:"user_friend_id"`
}

type SignInType struct {
	Email string `json:"email"`
}

type VerifyResType struct {
	Status string    `json:"status"`
	Token  string    `json:"token"`
	User   *UserType `json:"user"`
}

type User_FriendType struct {
	User_Friend_ID string `json:"user_friend_id"`
	User_ID        string `json:"user_id"`
	Friend_ID      string `json:"friend_id"`
}

type PostType struct {
	Post_ID    string    `json:"post_id"`
	User_ID    string    `json:"user_id"`
	Content    string    `json:"content"`
	Type       string    `json:"type"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

type PostReqType struct {
	User_ID string   `json:"user_id"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

type Image_PostType struct {
	Image_Post_ID string    `json:"image_post_id"`
	Post_ID       string    `json:"post_id"`
	User_ID       string    `json:"user_id"`
	Image         string    `json:"image"`
	Created_At    time.Time `json:"created_at"`
	Updated_At    time.Time `json:"updated_at"`
}

type GetPostType struct {
	Post_ID           string            `json:"post_id"`
	User_ID           string            `json:"user_id"`
	Content           string            `json:"content"`
	Type              string            `json:"type"`
	Created_At        time.Time         `json:"created_at"`
	Updated_At        time.Time         `json:"updated_at"`
	Images            []*Image_PostType `json:"images"`
	Number_Of_Comment int               `json:"number_of_comment"`
}

type GetPostResType struct {
	Post_ID           string            `json:"post_id"`
	User_ID           string            `json:"user_id"`
	Content           string            `json:"content"`
	Type              string            `json:"type"`
	Created_At        time.Time         `json:"created_at"`
	Updated_At        time.Time         `json:"updated_at"`
	Images            []*Image_PostType `json:"images"`
	Number_Of_Comment int               `json:"number_of_comment"`
	User_Name         string            `json:"user_name"`
	Email             string            `json:"email"`
	Photo_Profile     string            `json:"photo_profile"`
}

type CommentType struct {
	Comment_ID      string    `json:"comment_id"`
	Post_ID         string    `json:"post_id"`
	Comment_Post_ID string    `json:"comment_post_id"`
	Created_At      time.Time `json:"created_at"`
	Updated_At      time.Time `json:"updated_at"`
}

type CommentReqType struct {
	Post_ID string   `json:"post_id"`
	User_ID string   `json:"user_id"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

type NotificationType struct {
	Notification_ID string    `json:"notification_id"`
	Issuer          string    `json:"issuer"`
	Issuer_Name     string    `json:"issuer_name"`
	Notifier        string    `json:"notifier"`
	Notifier_Name   string    `json:"notifier_name"`
	Status          string    `json:"status"`
	Accept          string    `json:"accept"`
	Post_ID         string    `json:"post_id"`
	Type            string    `json:"type"`
	Created_At      time.Time `json:"created_at"`
	Updated_At      time.Time `json:"updated_at"`
}

type UpdateNotifType struct {
	Notification_ID string    `json:"notification_id"`
	Accept          string    `json:"accept"`
	Updated_At      time.Time `json:"updated_at"`
}
