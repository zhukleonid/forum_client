package controller

import (
	"bytes"
	"html/template"
	"lzhuk/clients/internal/convertor"
	"net/http"
)

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
			cookies, err := convertor.NewConvertCookie(resp)
			if err != nil {
				http.Error(w, "Error decode cookie", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, cookies[0])
			http.Redirect(w, r, "http://localhost:8082/userd3", 300)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
