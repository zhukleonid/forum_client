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
	return postsLikeUser, nil
}
