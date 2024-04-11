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
		newCategoryGet := fmt.Sprintf(categoryGet+"%s", r.FormValue("category"))
		fmt.Println(newCategoryGet)
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
			discriptionMsg, err := convertor.DecodeErrorResponse(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при декодировании ответа ошибки и описания от сервера на запрос об регистрации пользователя")
				return
			}
			switch {
			// Получена ошибка что почта уже используется
			case discriptionMsg.Discription == "Email already exist":
				errorPage(w, errors.EmailAlreadyExists, http.StatusConflict)
				log.Printf("Пользователь пытается зарегестировать почту которая используется под другим аккаунтом")
				return
				// Получена ошибка что введены неверные учетные данные
			case discriptionMsg.Discription == "Invalid Credentials":
				errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
				log.Printf("Не валидные данные")
				return
			case discriptionMsg.Discription == "Not Found Any Data":
				errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
				log.Printf("Не найдено")
				return
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на получение постов под конкретной категорией")
				return
			}
		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на получение постов под конкретной категорией")
			return
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("Не верный метод запроса при передаче запроса на категорию")
		return
	}
}
