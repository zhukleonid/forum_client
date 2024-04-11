package controller

import (
	"bytes"
	"fmt"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/internal/helpers"
	"lzhuk/clients/pkg/config/errors"
	"net/http"
	"strings"
)

func deleteComment(w http.ResponseWriter, r *http.Request) {
	// Проверяем что в запросе присутствуют куки с валидным имененем
	switch {
	case len(r.Cookies()) < 1:
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	case !strings.HasPrefix(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())].Name, "CookieUUID"):
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	}

	switch r.Method {
	case http.MethodPost:
		// Конвертируем данные для удаления комментария
		jsonData, err := convertor.ConvertDeleteComment(r)
		if err != nil {
			http.Error(w, "error update comment", http.StatusInternalServerError)
			return
		}
		// Формируем запрос на удаление комментария
		req, err := http.NewRequest("DELETE", deleteComments, bytes.NewBuffer(jsonData))
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при формировании запроса на удаления комментария. Ошибка: %v", err)
			return
		}
		// Записываем куки из браузера в запрос на сервис сервера
		req.AddCookie(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())])
		// Создаем структуру нового клиента
		client := http.Client{}
		// Передаем запрос на сервер
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при передаче запроса на сервис сервера для удаления комментария. Ошибка: %v", err)
			return
		}

		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusAccepted:
			link := fmt.Sprintf("http://localhost:8082/userd3/post/%s", r.FormValue("postId"))
			http.Redirect(w, r, link, 302)
		case http.StatusInternalServerError:
			discriptionMsg, err := convertor.DecodeErrorResponse(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при декодировании ответа ошибки и описания от сервера на запрос об изменении комментария")
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
				log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на удаления комментария")
				return
			}
		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на удаления комментария")
			return
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("Не верный метод запроса при удаления комментария")
		return
	}
}
