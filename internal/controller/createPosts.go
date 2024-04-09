package controller

import (
	"bytes"
	"html/template"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/config/errors"
	"net/http"
)

func createPost(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		http.Redirect(w, r, "http://localhost:8082/login", 300)
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
		// Конвертация данных при создании нового поста
		jsonData, err := convertor.ConvertCreatePost(r)
		if err != nil {
			http.Error(w, "Marshal CreatePost error", http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("POST", createPosts, bytes.NewBuffer(jsonData))
		if err != nil {
			http.Error(w, "Request registry error", http.StatusInternalServerError)
			return
		}
		req.AddCookie(r.Cookies()[0])
		req.Header.Set("Content-Type", "application/json")
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Request client registry error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			http.Redirect(w, r, "http://localhost:8082/userd3/myposts", 300)
			return
		}
	}
}
