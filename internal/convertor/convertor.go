package convertor

import (
	"encoding/json"
	"lzhuk/clients/model"
	"net/http"
)

func NewConvertAllPosts(resp *http.Response) (model.AllPosts, error) {
	posts := model.AllPosts{}
	err := json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func NewConvertRegister(r *http.Request) ([]byte, error) {
	register := model.Register{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	jsonData, err := json.Marshal(register)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func NewConvertLogin(r *http.Request) ([]byte, error) {
	login := model.Login{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	jsonData, err := json.Marshal(login)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func NewConvertCookie(resp *http.Response) (*http.Cookie, error) {
	cookie := &http.Cookie{}
	err := json.NewDecoder(resp.Body).Decode(cookie)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}
