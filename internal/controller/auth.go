package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"lzhuk/clients/internal/cahe"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/model"
	"lzhuk/clients/pkg/config/errors"
	"net/http"
	"net/url"
	"strings"
)

type googleUserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Sub   string `json:"sub"`
}

type githubUserInfo struct {
	Name     string `json:"name"`
	Email    string `json:"login"`
	Password string `json:"node_id"`
}

func Google(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=profile email", googleAuthEndPoint, "16807684949-bp5bhvp85ar5sfj2iuasmfsf4l6bj4up.apps.googleusercontent.com", "http://localhost:8082/google/callback")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		return
	}
	s := strings.NewReader(fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code", code, "16807684949-bp5bhvp85ar5sfj2iuasmfsf4l6bj4up.apps.googleusercontent.com", "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", "http://localhost:8082/google/callback"))
	response, err := http.Post(googleAuthEndPointAccessToken, "application/x-www-form-urlencoded", s)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()

	resp, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	token := ExtractValueFromBody(resp, "access_token")

	if token == "" {
		log.Println("Error empty token:")
		return
	}
	info, err := getUserInfo(token, googleUserInfoURL)
	if err != nil {
		log.Println(err)
		return
	}

	var googleUserInfo googleUserInfo
	if err = json.Unmarshal(info, &googleUserInfo); err != nil {
		log.Printf("Failed to unmarshal user info: %v", err)
		return
	}
	/// NEW
	us := model.UserReq{Name: googleUserInfo.Name, Email: googleUserInfo.Email, Password: googleUserInfo.Sub}
	us1, err := json.Marshal(us)
	fmt.Println(string(us1))
	if err != nil {
		log.Println(err)
		return
	}

		req1, err := http.NewRequest("POST", auth, bytes.NewBuffer(us1))
		if err != nil {
			log.Println(err)
			return
		}
		req1.Header.Set("Content-Type", "application/json")
		client := http.Client{}

		resp1, err := client.Do(req1)
		if err != nil {
			log.Println(err)
			return
		}

	defer resp1.Body.Close()
	switch resp1.StatusCode {
	// Получен статус код 201 об успешной регистрации пользователя в системе
	case http.StatusOK:
		var clientName string
		// Получение сгенерированных сервером куки
		cookie, err := convertor.ConvertFirstCookie(resp1)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при конвертации кук из ответа. Ошибка: %v", err)
			return
		}
		// Получение в глобальную переменную имени вошедшего пользователя
		clientName, err = convertor.DecodeClientName(resp1)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при получении имени пользователя")
			return
		}
		// Записываем клиента в хеш-таблицу
		cahe.Username[cookie.Value] = clientName
		// Записываем в ответ браузеру полученный экземпляр куки от сервера
		http.SetCookie(w, cookie)
		// Переход на домашнюю страницу пользователя
		http.Redirect(w, r, "http://localhost:8082/userd3", 302)
		return
	}
}

func GitHub(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=user:email", githubAuthEndPoint, "21c2671efe47648ceedd", "http://localhost:8082/github/callback")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Println("Error")
		return
	}
	req, err := http.NewRequest("POST", githubAuthEndAccessToken, strings.NewReader(fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code", code, "21c2671efe47648ceedd", "acb9aa72c6f829fa1262760243287dfd71566859", "http://localhost:8082/github/callback")))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(body))
	token, err := ExtractAccessTokenFromResponse(string(body))

	if token == "" {
		log.Println("Bad Credentials")
		return
	}

	info, err := getUserInfo(token, githubUserInfoURL)
	if err != nil {
		log.Println(err)
		return
	}

	var github githubUserInfo
	if err := json.Unmarshal(info, &github); err != nil {
		log.Println(err)
		return
	}

	user := model.UserReq{Name: github.Name, Email: github.Email, Password: github.Password}
	userMarshal, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		return
	}

	req1, err := http.NewRequest("POST", auth, bytes.NewBuffer(userMarshal))
	if err != nil {

		log.Println(err)
		return
	}
	req1.Header.Set("Content-Type", "application/json")
	client1 := &http.Client{}

	re, err1 := client1.Do(req1)
	if err1 != nil {
		log.Println(err)
		return
	}
	defer re.Body.Close()

	switch re.StatusCode {
	case http.StatusOK:
		var clientName string
		cookie, err := convertor.ConvertFirstCookie(re)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при конвертации кук из ответа. Ошибка: %v", err)
			return
		}
		clientName, err = convertor.DecodeClientName(re)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при получении имени пользователя")
			return
		}
		cahe.Username[cookie.Value] = clientName
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "http://localhost:8082/userd3", 302)
		return
	}
}

func ExtractValueFromBody(body []byte, key string) string {
	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		return ""
	}

	value, ok := response[key].(string)
	if !ok {
		return ""
	}
	return value
}

func ExtractAccessTokenFromResponse(response string) (string, error) {
	params, err := url.ParseQuery(response)
	if err != nil {
		return "", err
	}
	accessToken := params.Get("access_token")
	return accessToken, nil
}

func getUserInfo(accessToken string, userInfoURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
