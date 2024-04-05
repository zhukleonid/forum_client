package controller

import (
	"bytes"
	"lzhuk/clients/internal/convertor"
	"net/http"
)

func logoutUser(w http.ResponseWriter, r *http.Request) {
	jsonData, err := convertor.NewConvertLogout(r)
	if err != nil {
		http.Error(w, "Еггог convert logout", http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequest("POST", logoutUsers, bytes.NewBuffer(jsonData))
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

	if resp.StatusCode == http.StatusOK {
		expiredCookie := &http.Cookie{
			Name:   "CookieUUID",
			Value:  "",
			MaxAge: -1,
			Path:   "/",
		}
		http.SetCookie(w, expiredCookie)
		http.Redirect(w, r, "http://localhost:8082", 300)
	}
}
