package controller

import (
	"bytes"
	"html/template"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/config/errors"
	"net/http"
)

func loginPage(w http.ResponseWriter, r *http.Request) {
	// Создание шаблона для страницы входа пользователем
	t, err := template.ParseFiles("./ui/html/login.html")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка создании шаблона страницы для входа пользователя. Ошибка: %v", err)
		return
	}
	// Проверка метода запроса
	switch r.Method {
	case http.MethodGet:
		err = t.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при рендеринге шаблона страницы входа пользователя. Ошибка: %v", err)
			return
		}
	case http.MethodPost:
		// Конвертор полученных пользовательских данных о входе из формы html в формат JSON
		jsonData, err := convertor.ConvertLogin(r)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при конвертации данных о входе пользователя в JSON. Ошибка: %v", err)
			return
		}
		// Формирования POST запроса на вход пользователя на сервере
		req, err := http.NewRequest("POST", login, bytes.NewBuffer(jsonData))
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при формировании запроса на вход пользователя. Ошибка: %v", err)
			return
		}
		// Записываем тип контента в заголовок запроса
		req.Header.Set("Content-Type", "application/json")
		// Создаем структуру клиента для передачи запроса
		client := http.Client{}
		// Отправляем запрос на сервис сервера
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при передаче запроса от клиента к серверу при входе пользователя. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			// Получение сгенерированных сервером куки
			cookie, err := convertor.ConvertFirstCookie(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при конвертации кук из ответа. Ошибка: %v", err)
				return
			}
			// Получение в глобальную переменную имени вошедшего пользователя
			ClientName, err = convertor.DecodeClientName(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при получении имени пользователя")
				return
			}
			// Записываем в ответ браузеру полученный экземпляр куки от сервера
			http.SetCookie(w, cookie)
			// Переход на домашнюю страницу пользователя
			http.Redirect(w, r, "http://localhost:8082/userd3", 302)
			return
		case http.StatusInternalServerError:
			discriptionMsg, err := convertor.DecodeErrorResponse(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при декодировании ответа ошибки и описания от сервера на запрос об входе пользователя")
				return
			}
			switch {
			// Получена ошибка что введены неверные учетные данные
			case discriptionMsg.Discription == "Invalid Credentials":
				errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
				log.Printf("Пользователь пытается зарегестировать почту которая используется под другим аккаунтом")
				return
			default:
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка сервере. Ошибка: %v", err)
				return
			}
		}
		// Метод запроса с браузера не POST и не GET
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса на вход пользователя не верный метод запроса")
		return
	}
}
