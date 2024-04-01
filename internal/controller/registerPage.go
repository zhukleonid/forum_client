package controller

import (
	"bytes"
	"html/template"
	"lzhuk/clients/internal/convertor"
	"net/http"
)

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
