package convertor

import (
	"encoding/json"
	"fmt"
	"lzhuk/clients/model"
	"net/http"
	"strconv"
)


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

func NewConvertVoteComment(r *http.Request) ([]byte, error) {
	commentId, err := strconv.Atoi(r.FormValue("commentId"))
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
	voteComment := model.VoteComment{
		CommentId:  commentId,
		LikeStatus: status,
	}
	jsonData, err := json.Marshal(voteComment)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func NewConvertUserLikePosts(resp *http.Response) ([]model.Post, error) {
	postsLikeUser := []model.Post{}
	err := json.NewDecoder(resp.Body).Decode(&postsLikeUser)
	if err != nil {
		return nil, err
	}
	return postsLikeUser, nil
}

func NewConvertUpdateComment(r *http.Request) (model.UpdateComment, error) {
	commentId, err := strconv.Atoi(r.FormValue("commentId"))
	if err != nil {
		return model.UpdateComment{}, err
	}
	posttId, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		return model.UpdateComment{}, err
	}

	updateComment := model.UpdateComment{
		ID:          commentId,
		Post:        posttId,
		Description: r.FormValue("description"),
	}

	return updateComment, nil
}

func NewConvertUpdateCommentUser(r *http.Request) ([]byte, error) {
	commentId, err := strconv.Atoi(r.FormValue("commentId"))
	if err != nil {
		return nil, err
	}
	posttId, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		return nil, err
	}

	updateComment := model.UpdateComment{
		ID:          commentId,
		Post:        posttId,
		Description: r.FormValue("updatedComment"),
	}
	jsonData, err := json.Marshal(updateComment)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func NewConvertDeleteComment(r *http.Request) ([]byte, error) {
	commentId, err := strconv.Atoi(r.FormValue("commentId"))
	if err != nil {
		return nil, err
	}
	posttId, err := strconv.Atoi(r.FormValue("postId"))
	if err != nil {
		return nil, err
	}

	deleteComment := model.Comment{
		ID:   commentId,
		Post: posttId,
	}
	jsonData, err := json.Marshal(deleteComment)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func NewConvertLogout(r *http.Request) ([]byte, error) {
	cookie := r.Cookies()
	uuid := cookie[0].String()
	fmt.Println(uuid)
	logoutUser := model.LogoutUser{
		UUID: uuid,
	}
	jsonData, err := json.Marshal(logoutUser)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
