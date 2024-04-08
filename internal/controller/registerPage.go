package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/internal/validation"
	"lzhuk/clients/pkg/config/errors"
	"net/http"
)

func registerPage(w http.ResponseWriter, r *http.Request) {
	// Создание шаблона для страницы регистрации
	t, err := template.ParseFiles("./ui/html/register.html")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка создании шаблона страницы регистрации пользователя. Ошибка: %v", err)
		return
	}
	// Проверка метода запроса
	switch r.Method {
	case http.MethodGet:
		err = t.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при рендеринге шаблона страницы регистрации пользователя. Ошибка: %v", err)
			return
		}
	case http.MethodPost:
		// Проверка на валидность пользовательских данных
		validDate, errValid := validation.ValidDate(r)
		// Рендеринг страницы при невалидных данных регистрации пользователем
		if validDate == false {
			err = t.ExecuteTemplate(w, "register.html", errValid)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при рендеринге шаблона страницы регистрации пользователя при проверку на валидность данных. Ошибка: %v", err)
				return
			}
		}
		// Конвертор полученных пользовательских данных из формы html в формат JSON
		jsonData, err := convertor.ConvertRegister(r)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при конвертации данных о регистрации пользователя в JSON. Ошибка: %v", err)
			return
		}
		// Формирования POST запроса на регистрацию нового пользователя на сервере
		req, err := http.NewRequest("POST", registry, bytes.NewBuffer(jsonData))
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при формировании запроса на регистрацию пользователя. Ошибка: %v", err)
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
			log.Printf("Произошла ошибка при передаче запроса от клиента к серверу при регистрации нового пользователя. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()
		fmt.Println(resp.Body)
		// Проверка кода статуса ответа сервера
		switch resp.StatusCode {
		// Получен статус код 201 об успешной регистрации пользователя в системе
		case http.StatusCreated:
			http.Redirect(w, r, "http://localhost:8082/login", 303)
			// Получен статус код 500 об ошибке на сервере
		case http.StatusInternalServerError:
			discriptionMsg, err := convertor.DecodeErrorResponse(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при декодировании ответа ошибки и описания от сервера на запрос об регистрации пользователя")
				return
			}
			switch {
				// Полуена ошибка что почта уже используется 
			case discriptionMsg.Discription == "Email already exist":
				errorPage(w, errors.EmailAlreadyExists, http.StatusConflict)
				log.Printf("Пользователь пытается зарегестировать почту которая используется под другим аккаунтом")
				return
			case discriptionMsg.Discription == "Invalid Credentials":
				errorPage(w, errors.InvalidCredentials, http.StatusUnauthorized)
				log.Printf("Пользователь пытается зарегестировать почту которая используется под другим аккаунтом")
				return
			case discriptionMsg.Discription == "Not Found Any Data":
				return
			}

			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка на сервере при регистрации пользователя.ss")
			return
			// Получен статус код 405 об неверном методе запроса с сервера
		case http.StatusMethodNotAllowed:
			errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
			log.Printf("При передаче запроса на регистрацию нового пользователя не верный метод запроса")
			return
		default:
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка на сервере при регистрации пользователя.")
			return
		}
	default:
		errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
		log.Printf("При передаче запроса на регистрацию нового пользователя не верный метод запроса")
		return
	}
}
