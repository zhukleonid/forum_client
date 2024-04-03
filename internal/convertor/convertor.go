package convertor

import (
	"encoding/json"
	"errors"
	"fmt"
	"lzhuk/clients/model"
	"net/http"
	"strconv"
)

func NewConvertAllPosts(resp *http.Response) (model.AllPosts, error) {
	posts := model.AllPosts{}
	fmt.Println(resp.Body)
	err := json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func NewConvertRegister(r *http.Request) ([]byte, error) {
	register := model.Register{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	jsonData, err := json.Marshal(register)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func NewConvertLogin(r *http.Request) ([]byte, error) {
	login := model.Login{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	jsonData, err := json.Marshal(login)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func NewConvertCookie(resp *http.Response) ([]*http.Cookie, error) {
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		return nil, errors.New("no cookies found in the response")
	}
	return cookies, nil
}

func NewConvertCreatePost(r *http.Request) ([]byte, error) {
	createPost := model.CreatePost{
		CategoryName: r.FormValue("category"),
		Title:        r.FormValue("title"),
		Description:  r.FormValue("description"),
	}

	jsonData, err := json.Marshal(createPost)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func NewConvertGetPosts(resp *http.Response) (*model.GetPost, error) {
	getPosts := &model.GetPost{}
	err := json.NewDecoder(resp.Body).Decode(getPosts)
	if err != nil {
		return nil, err
	}
	return getPosts, nil
}

func NewConvertCreateComment(r *http.Request) ([]byte, error) {
	postId, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		return nil, err
	}
	createComment := model.CreateComment{
		Post:        postId,
		Description: r.FormValue("comment"),
	}
	jsonData, err := json.Marshal(createComment)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func NewConvertUpdatePost(r *http.Request) ([]byte, error) {
	updatePost := model.Post{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}
	jsonData, err := json.Marshal(updatePost)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func NewConvertVotePost(r *http.Request) ([]byte, error) {
	postId, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		return nil, err
	}
	var status bool
	switch r.FormValue("vote") {
	case "true":
		status = true
	case "false":
		status = false
	}
	votePost := model.VotePost{
		PostId:     postId,
		LikeStatus: status,
	}
	jsonData, err := json.Marshal(votePost)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
