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

			data := map[string]interface{}{
				"Username": ClientName,                          // Глобальное имя пользователя
				"Posts":    result,                              // Все посты из БД
				"Cookie":   r.Cookies()[0].Name == "CookieUUID", // Передаем true, если есть куки, иначе false
			}

			err = t.ExecuteTemplate(w, "home.html", data)
			if err != nil {
				http.Error(w, "Error executing template", http.StatusInternalServerError)
				return
			}
		}
	default:
	}
}
