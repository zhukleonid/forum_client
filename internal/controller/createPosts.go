package controller

import (
	"bytes"
	"html/template"
	"lzhuk/clients/internal/convertor"
	"net/http"
)

func createPost(w http.ResponseWriter, r *http.Request) {
	if len(r.Cookies()) == 0 {
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
	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("./ui/html/create_post.html")
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "create_post.html", nil)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		jsonData, err := convertor.NewConvertCreatePost(r)
		if err != nil {
			http.Error(w, "Marshal CreatePost error", http.StatusInternalServerError)
			return
		}
		req, err := http.NewRequest("POST", createPosts, bytes.NewBuffer(jsonData))
		if err != nil {
			http.Error(w, "Request registry error", http.StatusInternalServerError)
			return
		}
		req.AddCookie(r.Cookies()[0])
		req.Header.Set("Content-Type", "application/json")
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Request client registry error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			response, err := http.Get(allPost)
			if err != nil {
				http.Error(w, "Error request all posts", http.StatusInternalServerError)
				return
			}
			defer response.Body.Close()

			cookies := response.Cookies()

			result, err := convertor.NewConvertAllPosts(response)
			if err != nil {
				http.Error(w, "Error request all posts", http.StatusInternalServerError)
				return
			}
			t, err := template.ParseFiles("./ui/html/home.html")
			if err != nil {
				http.Error(w, "Error parsing template", http.StatusInternalServerError)
				return
			}

			data := map[string]interface{}{
				"Posts":  result,
				"Cookie": len(cookies) > 0, // Передаем true, если есть куки, иначе false
			}

			err = t.ExecuteTemplate(w, "home.html", data)
			if err != nil {
				http.Error(w, "Error executing template", http.StatusInternalServerError)
				return
			}

		}
	}
}
