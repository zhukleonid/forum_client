package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/internal/validation"
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
		form := validation.New()
		form.CheckField(form.EmailValid(r.FormValue("email")), "Не корректная почта", "")
		form.CheckField(form.EmptyFieldValid(r.FormValue("email")), "Пустой ввод почты не допускается", "")
		form.CheckField(form.EmptyFieldValid(r.FormValue("name")), "Пустой ввод имени не допускается", "")
		form.CheckField(form.NameValid(r.FormValue("name")), "Не корректное имя", "")
		form.CheckField(form.PasswordValid(r.FormValue("password")), "Не корректный пароль", "")
		if !form.Valid() {
			fmt.Println(form.Errors)
			t, err := template.ParseFiles("./ui/html/register.html")
			if err != nil {
				http.Error(w, "Error parsing template", http.StatusInternalServerError)
				return
			}
			err = t.ExecuteTemplate(w, "register.html", form.Errors)
			if err != nil {
				http.Error(w, "Error executing template", http.StatusInternalServerError)
				return
			}
		}
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

		if resp.StatusCode == http.StatusCreated {
			http.Redirect(w, r, "http://localhost:8082/login", 300)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
