package controller

import (
	"fmt"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/config/errors"
	"net/http"
	"strings"
)

func deletePost(w http.ResponseWriter, r *http.Request) {
	// Проверяем что в запросе присутствуют куки с валидным имененем
	switch {
	case len(r.Cookies()) < 1:
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	case !strings.HasPrefix(r.Cookies()[0].Name, "CookieUUID"):
		fmt.Println(strings.HasPrefix(r.Cookies()[0].Name, "CookieUUID"))
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	}

	switch r.Method {
	case http.MethodPost:
		// Формируем запрос на удаление конкретного поста по id
		req, err := http.NewRequest("DELETE", fmt.Sprintf(deletePosts+"%s", r.FormValue("postId")), nil)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при создании запроса на удаление конкретного поста по id. Ошибка: %v", err)
			return
		}
		// Записываем куки из браузера в запрос на сервис сервера
		req.AddCookie(r.Cookies()[0])
		// Создаем структуру нового клиента
		client := http.Client{}
		// Отправляем запрос на сервис сервера
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при отправке запроса на удаление конкретного поста по id. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusAccepted:
			http.Redirect(w, r, "http://localhost:8082/userd3/myposts", 302)
		case http.StatusInternalServerError:
			discriptionMsg, err := convertor.DecodeErrorResponse(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при декодировании ответа ошибки и описания от сервера на запрос об данных о посте который пользователь будет редактировать")
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
				errorPage(w, errors.ErrorNotMethod, http.StatusInternalServerError)
				log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на удаление конкретного поста")
				return
			}
		default:
			errorPage(w, errors.ErrorNotMethod, http.StatusInternalServerError)
			log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на удаление конкретного поста")
			return
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("Не верный метод запроса при удалении поста")
		return
	}
}
