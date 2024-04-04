package controller

import (
	"bytes"
	"fmt"
	"lzhuk/clients/internal/convertor"
	"net/http"
)

func voteComment(w http.ResponseWriter, r *http.Request) {
	jsonData, err := convertor.NewConvertVoteComment(r)
	if err != nil {
		http.Error(w, "Marshal VoteComment error", http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequest("POST", voteComments, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Request registry error", http.StatusInternalServerError)
		return
	}
	req.AddCookie(r.Cookies()[0])
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Request voteComments error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		link := fmt.Sprintf("http://localhost:8082/userd3/post/%s", r.FormValue("postId"))
		http.Redirect(w, r, link, 300)
	}
}
