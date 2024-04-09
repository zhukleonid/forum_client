package controller

import (
	"html/template"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/config/errors"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	// Создание шаблона для домашней
	t, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка создании шаблона страницы для входа пользователя. Ошибка: %v", err)
		return
	}
	// Проверка метода запроса
	switch r.Method {
	case http.MethodGet:
		// Отправка GET запроса на получение всех постов из БД сервиса сервера
		resp, err := http.Get(allPost)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при конвертации данных о входе пользователя в JSON. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			// Получение данных обо всех имеющихся постах
			result, err := convertor.ConvertAllPosts(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при конвертации данных обо всех постах из JSON. Ошибка: %v", err)
				return
			}

			var (
				nicname string // Хранит имя пользователя
				cookie  bool   // Хранит наличие куки
			)
			// Проверка на наличие куки
			switch {
			// Если куки не получены передаем пустое имя и отсутствие куки
			case len(r.Cookies()) < 1:
				nicname = ""
				cookie = false
				// При наличии куки проверяем их валидность
			default:
				// Проверяем что куки сгенерированы нашим сервисом сервера
				if r.Cookies()[0].Name == "CookieUUID" {
					cookie = true
					// Получаем по UUID имя пользователя
					value, ok := Username[r.Cookies()[0].Value]
					if ok {
						nicname = value
					} else {
						nicname = ""
						http.Redirect(w, r, "http://localhost:8082/login", 303)
						return
					}
				} else {
					nicname = ""
					cookie = false
				}
			}

			// Данные для рендеринга страницы
			data := map[string]interface{}{
				"Username": nicname, // Глобальное имя пользователя
				"Posts":    result,  // Все посты из БД
				"Cookie":   cookie,  // Передаем true, если есть куки, иначе false
			}

			// Рендеринг домашней страницы
			err = t.ExecuteTemplate(w, "home.html", data)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при рендеринге шаблона домашней страницы пользователя. Ошибка: %v", err)
				return
			}

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
				log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на создание поста")
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
