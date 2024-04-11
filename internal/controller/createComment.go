package controller

import (
	"bytes"
	"fmt"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/internal/helpers"
	"lzhuk/clients/internal/validation"
	"lzhuk/clients/pkg/config/errors"
	"net/http"
	"strings"
)

func createComment(w http.ResponseWriter, r *http.Request) {
	// Проверяем что в запросе присутствуют куки с валидным имененем
	switch {
	case len(r.Cookies()) < 1:
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	case !strings.HasPrefix(r.Cookies()[0].Name, "CookieUUID"):
		fmt.Println(strings.HasPrefix(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())].Name, "CookieUUID"))
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	}
	switch r.Method {
	case http.MethodPost:
		// Проверка на валидность данных комментария
		validDateComment, _ := validation.ValidDateComment(r)
		if validDateComment == false {
			errorPage(w, errors.EmptyComments, http.StatusBadRequest)
			log.Printf("Пустой комментарий")
		} else {
			// Конвертируем данные при создании нового комментария
			jsonData, err := convertor.ConvertCreateComment(r)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при конвертации данных нового комментария. Ошибка: %v", err)
				return
			}
			// Формирование нового запроса
			req, err := http.NewRequest("POST", createComments, bytes.NewBuffer(jsonData))
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при формирование запроса нового комментария. Ошибка: %v", err)
				return
			}
			// Запись куки из браузера в запрос для сервера
			req.AddCookie(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())])
			req.Header.Set("Content-Type", "application/json")
			client := http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при отправке запроса нового комментария на сервер. Ошибка: %v", err)
				return
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusCreated:
				newReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8082/userd3/post/%s", r.FormValue("postId")), nil)
				if err != nil {
					errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
					log.Printf("Произошла ошибка при формирования запроса на редирект. Ошибка: %v", err)
					return
				}
				http.Redirect(w, r, newReq.URL.String(), 302)
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
					errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
					log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на создание комменатрия")
					return
				}
			default:
				errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
				log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на создание комменатрия")
				return
			}
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на создание комменатрия")
		return
	}
}
