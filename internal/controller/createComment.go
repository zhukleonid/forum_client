package controller

import (
	"bytes"
	"lzhuk/clients/internal/convertor"
	"net/http"
)

func createComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jsonData, err := convertor.NewConvertCreateComment(r)
	if err != nil {
		http.Error(w, "Marshal CreateComment error", http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequest("POST", createComments, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Request CreateComment error", http.StatusInternalServerError)
		return
	}
	req.AddCookie(r.Cookies()[0])
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Request client CreateComment error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
}
