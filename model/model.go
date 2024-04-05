package model

import "time"

type AllPosts []struct {
	PostID       int       `json:"post_id"`
	UserID       int       `json:"user_id"`
	CategoryName string    `json:"category_name"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"create_at"`
	Author       string    `json:"name"`
	Like         int       `json:"likes"`
	Dislike      int       `json:"dislikes"`
}

type Post struct {
	PostId       int    `json:"post_id"`
	UserId       int    `json:"user_id"`
	CategoryName string `json:"category_name"`
	Title        string
	Description  string
	CreateDate   time.Time `json:"create_at"`

	Author  string `json:"name"`
	Like    int    `json:"likes"`
	Dislike int    `json:"dislikes"`
}

type Register struct {
	Name     string
	Email    string
	Password string
}

type Error struct {
	Status      int
	Discription string
}

type Login struct {
	Email    string
	Password string
}

type Cookie struct {
	Name   string
	Value  string
	Path   string
	MaxAge int
}

type CreatePost struct {
	CategoryName string `json:"category_name"`
	Title        string
	Description  string
}

type Comment struct {
	ID          int       `json:"id"`
	User        int       `json:"user_id"`
	Post        int       `json:"post_id"`
	Description string    `json:"description"`
	CreatedDate time.Time `json:"created_at"`
	UpdatedDate time.Time `json:"updated_at"`

	Name    string `json:"name"`
	Like    int    `json:"likes"`
	Dislike int    `json:"dislikes"`
}

type GetPost struct {
	Post     *Post
	Comments []*Comment
}

type CreateComment struct {
	Post        int    `json:"post_id"`
	Description string `json:"description"`
}

type UpdateComment struct {
	ID          int    `json:"id"`
	Post        int    `json:"post_id"`
	Description string `json:"description"`
}

type VotePost struct {
	PostId     int  `json:"post_id"`
	LikeStatus bool `json:"status"`
}

type VoteComment struct {
	CommentId  int  `json:"comment_id"`
	LikeStatus bool `json:"status"`
}
