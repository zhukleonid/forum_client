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
	getUserPost = fmt.Sprintf(getUserPost, id)

	req, err := http.NewRequest("POST", getUserPost, nil)
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

		cookies := r.Cookies()

		result, err := convertor.NewConvertAllPosts(resp)
		if err != nil {
			http.Error(w, "Error request all posts", http.StatusInternalServerError)
			return
		}
		t, err := template.ParseFiles("./ui/html/post.html")
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Posts":  result,
			"Cookie": len(cookies) > 0, // Передаем true, если есть куки, иначе false
		}

		err = t.ExecuteTemplate(w, "post.html", data)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	}
}
