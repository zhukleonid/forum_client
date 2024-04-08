package controller

import (
	"html/template"
	"net/http"
)

// Функция обработчика ошибок
func errorPage(w http.ResponseWriter, statusText string, statusCode int) {
	w.WriteHeader(statusCode)

	data := struct {
		StatusText string
		StatusCode int
	}{
		statusText,
		statusCode,
	}

	t, err := template.ParseFiles("./ui/html/error.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "error.html", data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusInternalServerError)
		return
	}
}
