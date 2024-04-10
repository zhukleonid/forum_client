package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
	"net/http"
)

// Функция конвертации данных для получения всех понравишихся постов пользователя
func ConvertUserLikePosts(resp *http.Response) ([]model.Post, error) {
	postsLikeUser := []model.Post{}
	err := json.NewDecoder(resp.Body).Decode(&postsLikeUser)
	if err != nil {
		return nil, err
	}
	reverseLikePost(postsLikeUser)
	return postsLikeUser, nil
}

func reverseLikePost(s []model.Post) {
	for i, j := 0, len(s) - 1; i < j; i, j = i + 1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}