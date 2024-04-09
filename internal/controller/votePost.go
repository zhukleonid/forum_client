package controller

import (
	"bytes"
	"fmt"
	"lzhuk/clients/internal/convertor"
	"net/http"
)

func votePost(w http.ResponseWriter, r *http.Request) {
	jsonData, err := convertor.NewConvertVotePost(r)
	if err != nil {
		http.Error(w, "Marshal VotePost error", http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequest("POST", votePosts, bytes.NewBuffer(jsonData))
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
		newReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8082/userd3/post/%s", r.FormValue("postId")), nil)
		if err != nil {
			http.Error(w, "Request client registry error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, newReq.URL.String(), 302)
	}
}
