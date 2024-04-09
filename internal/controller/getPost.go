package controller

import (
	"fmt"
	"html/template"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/config/errors"
	"net/http"
	"strings"
)

func getPost(w http.ResponseWriter, r *http.Request) {
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
	// Создаем шаблон страницы с конкретным постом
	t, err := template.ParseFiles("./ui/html/post.html")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при создании шаблона страницы с конкретным постом. Ошибка: %v", err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Разбиваем путь URL на срез по признаку слеша
		parts := strings.Split(r.URL.Path, "/")
		// Формируем URL запроса на сервис сервера с конкретным id поста
		getUserPostId := fmt.Sprintf(getUserPost+"%s", parts[len(parts)-1])
		// Формируем запрос
		req, err := http.NewRequest("GET", getUserPostId, nil)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла формировании запроса об конкретном посте. Ошибка: %v", err)
			return
		}
		// Записываем куки из браузера в запрос к серверу
		req.AddCookie(r.Cookies()[0])
		// Создаем структуру нового клиента
		client := http.Client{}
		// Передаем запрос на сервис сервера
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при передаче запроса об получении конкретного поста. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()
		switch resp.StatusCode {
		case http.StatusOK:
			result, err := convertor.ConvertGetPosts(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при передаче запроса об получении конкретного поста. Ошибка: %v", err)
				return
			}
			err = t.ExecuteTemplate(w, "post.html", result)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при рендеренге страницы с конкретным постом. Ошибка: %v", err)
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
				log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на получение конкретного поста")
				return
			}
		default:
			errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
			log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на получение конкретного поста")
			return
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на получение конкретного поста")
		return
	}
}
