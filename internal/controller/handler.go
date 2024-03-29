package controller

import (
	"bytes"
	"html/template"
	"lzhuk/clients/internal/convertor"
	"net/http"
)

const (
	allPost  = "http://localhost:8083/userd3"
	registry = "http://localhost:8083/register"
	login    = "http://localhost:8083/login"
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
		t, err := template.ParseFiles("./ui/html/register.html")
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		jsonData, err := convertor.NewConvertRegister(r)
		if err != nil {
			http.Error(w, "Marshal registry error", http.StatusInternalServerError)
			return
		}
		req, err := http.NewRequest("POST", registry, bytes.NewBuffer(jsonData))
		if err != nil {
			http.Error(w, "Request registry error", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Request client registry error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t, err := template.ParseFiles("./ui/html/login.html")
			if err != nil {
				http.Error(w, "Error parsing template", http.StatusInternalServerError)
				return
			}
			err = t.ExecuteTemplate(w, "login.html", nil)
			if err != nil {
				http.Error(w, "Error executing template", http.StatusInternalServerError)
				return
			}
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("./ui/html/login.html")
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		jsonData, err := convertor.NewConvertLogin(r)
		if err != nil {
			http.Error(w, "Marshal login error", http.StatusInternalServerError)
			return
		}
		req, err := http.NewRequest("POST", login, bytes.NewBuffer(jsonData))
		if err != nil {
			http.Error(w, "Request login error", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Request login registry error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		
		if resp.StatusCode == http.StatusOK {
			cookie, err := convertor.NewConvertCookie(resp)
			if err != nil {
				http.Error(w, "Error decode cookie", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, cookie)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
