package controller

import (
	"bytes"
	"fmt"
	"html/template"
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

	if resp.StatusCode == http.StatusOK {
		getUserPostId := fmt.Sprintf(getUserPost+"%s", r.FormValue("postId"))
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
}
