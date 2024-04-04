package controller

import (
	"fmt"
	"html/template"
	"lzhuk/clients/internal/convertor"
	"net/http"
)

func likePost(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", likePosts, nil)
	if err != nil {
		http.Error(w, "Request user like post error", http.StatusInternalServerError)
		return
	}
	req.AddCookie(r.Cookies()[0])
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Request user like post error", http.StatusInternalServerError)
		return
	}
	userLikePosts, err := convertor.NewConvertUserLikePosts(resp)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles("./ui/html/user_likeposts.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	fmt.Println(userLikePosts)

	err = t.ExecuteTemplate(w, "user_likeposts.html", userLikePosts)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		http.Redirect(w, r, "http://localhost:8082/userd3", 300)
	}
}
