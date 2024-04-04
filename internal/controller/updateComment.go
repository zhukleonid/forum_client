package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"lzhuk/clients/internal/convertor"
	"net/http"
)

func updateComment(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:

		updateComment, err := convertor.NewConvertUpdateComment(r)
		if err != nil {
			http.Error(w, "error update comment", http.StatusInternalServerError)
			return
		}
		t, err := template.ParseFiles("./ui/html/update_comment.html")
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}

		err = t.ExecuteTemplate(w, "update_comment.html", updateComment)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		jsonData, err := convertor.NewConvertUpdateCommentUser(r)
		if err != nil {
			http.Error(w, "Marshal UpdateComment error", http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("PUT", updateComments, bytes.NewBuffer(jsonData))
		if err != nil {
			http.Error(w, "Request UpdateComment error", http.StatusInternalServerError)
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
			link := fmt.Sprintf("http://localhost:8082/userd3/post/%s", r.FormValue("postId"))
			http.Redirect(w, r, link, 300)
		}
	}
}
