package model

import "time"

type AllPosts []struct {
	PostID       int       `json:"post_id"`
	UserID       int       `json:"user_id"`
	CategoryName string    `json:"category_name"`
	Title        string    `json:"title"`
	Discription  string    `json:"discription"`
	CreateAt     time.Time `json:"create_at"`
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
