package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/internal/helpers"
	"lzhuk/clients/pkg/config/errors"
	"net/http"
	"strings"
)

func updateComment(w http.ResponseWriter, r *http.Request) {
	// Проверяем что в запросе присутствуют куки с валидным имененем
	switch {
	case len(r.Cookies()) < 1:
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	case !strings.HasPrefix(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())].Name, "CookieUUID"):
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	}

	// Создание шаблона страницы для редактирования комментария
	t, err := template.ParseFiles("./ui/html/update_comment.html")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при создании шаблона страницы для редактирования комментария. Ошибка: %v", err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Конвертация данных для редактирования комментария
		updateComment, err := convertor.ConvertUpdateComment(r)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при конвертации данных для редактирования комментария. Ошибка: %v", err)
			return
		}
		// Рендеринг страницы с изменениями для комментариев
		err = t.ExecuteTemplate(w, "update_comment.html", updateComment)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при рендеринге страницы для редактирования комментария. Ошибка: %v", err)
			return
		}
	case http.MethodPost:
		// Конвертация данных для отправки изменений в комментарии
		jsonData, err := convertor.ConvertUpdateCommentUser(r)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при конвертации данных для запроса редактирования комментария. Ошибка: %v", err)
			return
		}
		// Формируем запрос не сервис сервера
		req, err := http.NewRequest("PUT", updateComments, bytes.NewBuffer(jsonData))
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при формировании запроса редактирования комментария. Ошибка: %v", err)
			return
		}
		// Записываем куки из браузера в запрос на сервис сервера
		req.AddCookie(r.Cookies()[helpers.CheckCookieIndex(r.Cookies())])
		req.Header.Set("Content-Type", "application/json")
		// Создаем структуру нового клиента
		client := http.Client{}
		// Отправляем запрос на сервис сервера
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка отправке запроса на сервис сервера для редактирования комментария. Ошибка: %v", err)
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
				log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на получение старницы с данными о комментарии который пользователь собирается редактировать")
				return
			}
		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на получение старницы с данными о комментарии который пользователь собирается редактировать")
			return
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("Не верный метод запроса при изменении данных комментария")
		return
	}
}
