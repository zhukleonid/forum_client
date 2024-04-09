package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
	"net/http"
)

// Функция конвертации данных при создании нового поста
func ConvertCreatePost(r *http.Request) ([]byte, error) {
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
