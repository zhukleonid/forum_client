package model

import "time"

type AllPosts []struct {
	PostID       int      `json:"post_id"`
	UserID       int      `json:"user_id"`
	CategoryName string   `json:"category_name"`
	Title        string   `json:"title"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"create_at"`
}

type Post struct {
	PostId       int
	UserId       int
	CategoryName string `json:"category_name"`
	Title        string
	Description  string
	CreateDate   time.Time
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
	CategoryName string
	Title        string
	Description  string
}

type Comment struct {
	ID          int
	User        int
	Post        int
	Description string
	CreatedDate time.Time
	UpdatedDate time.Time
}

type GetPost struct {
	Post     *Post
	Comments []*Comment
}

type CreateComment struct {
	Post        int    `json:"post_id"`
	Description string `json:"description"`
}
