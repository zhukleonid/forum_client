package controller

import (
	"html/template"
	"lzhuk/clients/internal/convertor"
	"net/http"
)

func myPosts(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodGet {
	// 	http.Error(w, "Not method get", http.StatusMethodNotAllowed)
	// 	return
	// }
	req, err := http.NewRequest("GET", userPost, nil)
	if err != nil {
		http.Error(w, "Request user post error", http.StatusInternalServerError)
		return
	}
	req.AddCookie(r.Cookies()[0])
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Request client registry error", http.StatusInternalServerError)
		return
	}
	userPosts, err := convertor.NewConvertAllPosts(resp)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles("./ui/html/user_posts.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "user_posts.html", userPosts)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		http.Redirect(w, r, "http://localhost:8082/userd3", 300)
	}
}