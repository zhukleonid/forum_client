package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"lzhuk/clients/internal/convertor"
	"net/http"
)

func createPost(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
		http.Redirect(w, r, "http://localhost:8082/login", 300)
		return
	}
	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("./ui/html/create_post.html")
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "create_post.html", nil)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		
		jsonData, err := convertor.NewConvertCreatePost(r)
		if err != nil {
			http.Error(w, "Marshal CreatePost error", http.StatusInternalServerError)
			return
		}
		fmt.Println(string(jsonData))
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
