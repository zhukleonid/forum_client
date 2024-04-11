package controller

import (
	"fmt"
	"log"
	"lzhuk/clients/internal/cahe"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/internal/helpers"
	"lzhuk/clients/pkg/config/errors"
	"net/http"
)

func categoryPosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Передаем на сервис сервера запрос на получение постов по конкретной категории
		newCategoryGet := fmt.Sprintf(categoryGet+"%s", r.FormValue("name"))
		req, err := http.NewRequest("GET", newCategoryGet, nil)
		if err != nil {
			return
		}
		// Добавление из браузера куки в запрос на сервер
		req.AddCookie(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())])
		req.Header.Set("Content-Type", "application/json")
		// Создаем структуру нового клиента
		client := http.Client{}
		// Отправляем запрос на сервер
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при передаче запроса об создании нового поста на сервис сервера. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()
		result, err := convertor.ConvertCategoryPosts(resp)
		if err != nil {
			return
		}
		cahe.CategoryPosts = result
		switch resp.StatusCode {
		case http.StatusOK:
			http.Redirect(w, r, "http://localhost:8082/userd3", 302)
		case http.StatusInternalServerError:
			return
		default:
			return
		}
	default:
		return
	}
}
