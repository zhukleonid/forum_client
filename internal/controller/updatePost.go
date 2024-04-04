package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"lzhuk/clients/internal/convertor"
	"net/http"
	"strings"
)

func updatePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		path := r.URL.Path
		parts := strings.Split(path, "/")
		id := parts[len(parts)-1]
		getUserPostId := fmt.Sprintf(getUserPost+"%s", id)

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
			t, err := template.ParseFiles("./ui/html/update_post.html")
			if err != nil {
				http.Error(w, "Error parsing template", http.StatusInternalServerError)
				return
			}
			err = t.ExecuteTemplate(w, "update_post.html", result)
			if err != nil {
				http.Error(w, "Error executing template", http.StatusInternalServerError)
				return
			}
		}
	case http.MethodPost:
		jsonData, err := convertor.NewConvertUpdatePost(r)
		if err != nil {
			http.Error(w, "Marshal CreateComment error", http.StatusInternalServerError)
			return
		}

		updatePostsId := fmt.Sprintf(updatePosts+"%s", r.FormValue("postId"))
		req, err := http.NewRequest("PUT", updatePostsId, bytes.NewBuffer(jsonData))
		if err != nil {
			http.Error(w, "Request CreateComment error", http.StatusInternalServerError)
			return
		}
		req.AddCookie(r.Cookies()[0])
		req.Header.Set("Content-Type", "application/json")
		fmt.Println(req.Body)
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Request client CreateComment error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			path := r.URL.Path
			parts := strings.Split(path, "/")
			id := parts[len(parts)-1]
			getUserPostId := fmt.Sprintf(getUserPost+"%s", id)

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
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
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
				}
			}
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
