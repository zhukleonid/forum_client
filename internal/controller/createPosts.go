package controller

import (
	"bytes"
	"html/template"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/internal/helpers"
	"lzhuk/clients/internal/validation"
	"lzhuk/clients/pkg/config/errors"
	"net/http"
	"strings"
)

func createPost(w http.ResponseWriter, r *http.Request) {
	// Проверяем что в запросе присутствуют куки с валидным имененем
	switch {
	case len(r.Cookies()) < 1:
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	case !strings.HasPrefix(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())].Name, "CookieUUID"):
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	}
	// Создание шаблона для страницы создания поста
	t, err := template.ParseFiles("./ui/html/create_post.html")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка создания шаблона страницы для создания поста. Ошибка: %v", err)
		return
	}
	// Проверка метода запроса
	switch r.Method {
	case http.MethodGet:

		err = t.ExecuteTemplate(w, "create_post.html", nil)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при рендеринге страницы создания поста. Ошибка: %v", err)
			return
		}
	case http.MethodPost:
		// Проверка на валидность пользовательских данных
		validDatePost, _ := validation.ValidDatePost(r)
		if validDatePost == false {
			errorPage(w, errors.EmptyDatePost, http.StatusBadRequest)
			log.Printf("Произошла ошибка при рендеринге шаблона страницы создания поста пользователем при проверке на валидность данных. Ошибка: %v", err)
			return

		} else {
			// Конвертация данных при создании нового поста
			jsonData, err := convertor.ConvertCreatePost(r)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при конвертации данных для создания нового поста в JSON. Ошибка: %v", err)
				return
			}
			// Создание POST запроса на внесение информации о новом посте
			req, err := http.NewRequest("POST", createPosts, bytes.NewBuffer(jsonData))
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при создании запроса о новом посте на сервис сервера. Ошибка: %v", err)
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
			switch resp.StatusCode {
			case http.StatusCreated:
				http.Redirect(w, r, "http://localhost:8082/userd3", 302)
				return
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
					log.Printf("Получена ошибка сервера от сервиса сервера")
					return
				}
			default:
				errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
				log.Printf("Получена ошибка сервера от сервиса сервера")
				return
			}
		}
		// Метод запроса с браузера не POST и не GET
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса на домашнюю страницу не верный метод запроса")
		return
	}
}
