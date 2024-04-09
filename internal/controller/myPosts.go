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

func myPosts(w http.ResponseWriter, r *http.Request) {
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
	// Создаем шаблон страницы с постами пользователя
	t, err := template.ParseFiles("./ui/html/user_posts.html")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка создания шаблона страницы постов пользователя. Ошибка: %v", err)
		return
	}
	// Создаем новый запрос на сервис сервера для получения всех постов пользователя
	req, err := http.NewRequest("GET", userPost, nil)
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при создании запроса для получения с сервера постов пользователя. Ошибка: %v", err)
		return
	}
	// Добавление из браузера куки в запрос на сервер
	req.AddCookie(r.Cookies()[0])
	// Создаем структуру нового клиента
	client := http.Client{}
	
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Request client registry error", http.StatusInternalServerError)
		return
	}
	userPosts, err := convertor.ConvertAllPosts(resp)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "user_posts.html", userPosts)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		http.Redirect(w, r, "http://localhost:8082/userd3", 300)
	}
}
