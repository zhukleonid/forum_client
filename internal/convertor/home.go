package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
	"net/http"
	"time"
)

// Функция конвертации данных полученных при запросе всех постов
func ConvertAllPosts(resp *http.Response) (model.AllPosts, error) {
	posts := model.AllPosts{}
	err := json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		date := posts[i].CreatedAt
		dateString := date.String()
		layout := "2006-01-02 15:04:05.999999999 -0700 -07"
		dateTime, err := time.Parse(layout, dateString)
		if err != nil {
			// errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			// log.Printf("Произошла ошибка при преобразовании формата времени в строку в постах. Ошибка: %v", err)
			return nil, err
		}
		formattedStr := dateTime.Format("2006-01-02 15:04:05")
		newT, err := time.Parse("2006-01-02 15:04:05", formattedStr)
		if err != nil {
			// errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			// log.Printf("Произошла ошибка при преобразовании формата времени из строки в постах. Ошибка: %v", err)
			return nil, err
		}
		posts[i].CreatedAt = newT
		// fmt.Println(v.CreatedAt)
	}

	return posts, nil
}
