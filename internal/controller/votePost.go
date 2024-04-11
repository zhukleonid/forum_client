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

func votePost(w http.ResponseWriter, r *http.Request) {
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
		// Конвертация данных при постановке реакции на пост
		jsonData, err := convertor.ConvertVotePost(r)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при конвертации данных при постановке реакции на пост. Ошибка: %v", err)
			return
		}
		// Формируем запрос
		req, err := http.NewRequest("POST", votePosts, bytes.NewBuffer(jsonData))
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при формировании запроса при постановке реакции на пост. Ошибка: %v", err)
			return
		}
		// Записываем куки из бразура в запрос на сервис сервера
		req.AddCookie(r.Cookies()[0])
		req.Header.Set("Content-Type", "application/json")
		// Формируем структуру нового клиента
		client := http.Client{}
		// Отправляем запрос на сервис сервера
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при отправке запроса при постановке реакции на пост. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()
		switch resp.StatusCode {
		case http.StatusOK:
			newReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8082/userd3/post/%s", r.FormValue("postId")), nil)
			if err != nil {
				http.Error(w, "Request client registry error", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, newReq.URL.String(), 302)
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
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на получение старницы с данными о посте который пользователь собирается редактировать")
				return
			}

		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на получение старницы с данными о посте который пользователь собирается редактировать")
			return
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("Не верный метод запроса при постановке голоса на пост")
		return

	}
}
