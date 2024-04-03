package controller

import (
	"fmt"
	"net/http"
)

func deletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	deletePostsId := fmt.Sprintf(deletePosts+"%s", r.FormValue("postId"))
	fmt.Println(deletePostsId)
	req, err := http.NewRequest("DELETE", deletePostsId, nil)
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
		http.Redirect(w, r, "http://localhost:8082/userd3/myposts", 300)
	}
}
