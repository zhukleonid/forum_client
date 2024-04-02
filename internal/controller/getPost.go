package controller

import (
	"fmt"
	"html/template"
	"lzhuk/clients/internal/convertor"
	"net/http"
	"strings"
)

func getPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	parts := strings.Split(path, "/")
	id := parts[len(parts)-1]
	getUserPostId := fmt.Sprintf(getUserPost + "%s", id)
	fmt.Println(getUserPostId)

	req, err := http.NewRequest("GET", getUserPostId, nil)
	if err != nil {
		http.Error(w, "Request getUserPost error", http.StatusInternalServerError)
		return
	}
	req.AddCookie(r.Cookies()[0])
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Request client registry error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {


		result, err := convertor.NewConvertGetPosts(resp)
		if err != nil {
			http.Error(w, "Error request get posts", http.StatusInternalServerError)
			return
		}
		t, err := template.ParseFiles("./ui/html/post.html")
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}
		fmt.Println(result)
		err = t.ExecuteTemplate(w, "post.html", result)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	}
}
