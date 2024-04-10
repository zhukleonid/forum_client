package controller

import (
	"bytes"
	"fmt"
	"lzhuk/clients/internal/convertor"
	"net/http"
	"strings"
)

func deleteComment(w http.ResponseWriter, r *http.Request) {
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

	switch r.Method {
	case http.MethodGet:
		// Конвертируем данные для удаления комментария
		jsonData, err := convertor.ConvertDeleteComment(r)
		if err != nil {
			http.Error(w, "error update comment", http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("DELETE", deleteComments, bytes.NewBuffer(jsonData))
		if err != nil {
			http.Error(w, "Request gdelete comment error", http.StatusInternalServerError)
			return
		}
		req.AddCookie(r.Cookies()[0])
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Request delete comment error", http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			link := fmt.Sprintf("http://localhost:8082/userd3/post/%s", r.FormValue("postId"))
			http.Redirect(w, r, link, 300)
		}
	default:
	}
}
