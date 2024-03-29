package controller

import (
	"html/template"
	"lzhuk/clients/internal/convertor"
	"net/http"
)

const (
	allPost = "http://localhost:8083/userd3"
)

func startPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	t, err := template.ParseFiles("./ui/html/start.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "start.html", nil)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response, err := http.Get(allPost)
	if err != nil {
		http.Error(w, "Error request all posts", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	result, err := convertor.NewConvertAllPosts(response)
	if err != nil {
		http.Error(w, "Error request all posts", http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "home.html", result)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func registerPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

	case http.MethodPost:
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
