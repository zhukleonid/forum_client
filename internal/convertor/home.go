package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
	"net/http"
)

// Функция конвертации данных полученных при запросе всех постов
func ConvertAllPosts(resp *http.Response) (model.AllPosts, error) {
	posts := model.AllPosts{}
	err := json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
