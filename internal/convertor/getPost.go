package convertor

import (
	"encoding/json"
	"fmt"
	"lzhuk/clients/model"
	"net/http"
)

func ConvertGetPosts(resp *http.Response) (model.GetPostDate, error) {
	getPosts := &model.GetPost{}
	err := json.NewDecoder(resp.Body).Decode(getPosts)
	if err != nil {
		return model.GetPostDate{}, err
	}
	convertDatePosts := model.GetPostDate{
		Post:    &model.PostDate{},
		Comment: make([]*model.CommentDate, len(getPosts.Comment)),
	}

	date := getPosts.Post.CreateDate
	formattedStr := date.Format("2006-01-02 15:04:05")
	convertDatePosts.Post.PostId = getPosts.Post.PostId
	convertDatePosts.Post.UserId = getPosts.Post.UserId
	convertDatePosts.Post.CategoryName = getPosts.Post.CategoryName
	convertDatePosts.Post.Title = getPosts.Post.Title
	convertDatePosts.Post.Description = getPosts.Post.Description
	convertDatePosts.Post.CreateDate = formattedStr
	convertDatePosts.Post.Author = getPosts.Post.Author
	convertDatePosts.Post.Like = getPosts.Post.Like
	convertDatePosts.Post.Dislike = getPosts.Post.Dislike

	for i := range getPosts.Comment {
		date := getPosts.Comment[i].CreatedDate
		formattedStr := date.Format("2006-01-02 15:04:05")
		convertDatePosts.Comment[i].ID = getPosts.Comment[i].ID
		convertDatePosts.Comment[i].User = getPosts.Comment[i].User
		convertDatePosts.Comment[i].Post = getPosts.Comment[i].Post
		convertDatePosts.Comment[i].Description = getPosts.Comment[i].Description
		convertDatePosts.Comment[i].CreatedDate = formattedStr
		convertDatePosts.Comment[i].UpdatedDate = formattedStr
		convertDatePosts.Comment[i].Name = getPosts.Comment[i].Name
		convertDatePosts.Comment[i].Like = getPosts.Comment[i].Like
		convertDatePosts.Comment[i].Dislike = getPosts.Comment[i].Dislike
	}
	fmt.Println(getPosts.Comment)
	return convertDatePosts, nil
}
