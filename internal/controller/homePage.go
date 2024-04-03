package controller

import (
	"html/template"
	"lzhuk/clients/internal/convertor"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodGet {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	response, err := http.Get(allPost)
	if err != nil {
		http.Error(w, "Error request all posts", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	cookies := r.Cookies()

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
