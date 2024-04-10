package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"lzhuk/clients/internal/convertor"
	"lzhuk/clients/pkg/config/errors"
	"net/http"
	"strings"
)

func updatePost(w http.ResponseWriter, r *http.Request) {
	// Проверяем что в запросе присутствуют куки с валидным имененем
	switch {
	case len(r.Cookies()) < 1:
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	case !strings.HasPrefix(r.Cookies()[0].Name, "CookieUUID"):
		fmt.Println(strings.HasPrefix(r.Cookies()[0].Name, "CookieUUID"))
		http.Redirect(w, r, "http://localhost:8082/login", 302)
		return
	}
	// Создаем шаблон страницы с постом для его изменения
	t, err := template.ParseFiles("./ui/html/update_post.html")
	if err != nil {
		errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
		log.Printf("Произошла ошибка при создании шаблона страницы с постом для его изменения. Ошибка: %v", err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		// Разбиваем путь URL на срез по признаку слеша
		parts := strings.Split(r.URL.Path, "/")
		// Формируем URL запроса на сервис сервера с конкретным id поста
		getUserPostId := fmt.Sprintf(getUserPost+"%s", parts[len(parts)-1])
		// Формируем запрос на сервис сервера
		req, err := http.NewRequest("GET", getUserPostId, nil)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при создании запроса на сервер для изменения поста. Ошибка: %v", err)
			return
		}
		// Записываем куки из браузера в запрос к серверу
		req.AddCookie(r.Cookies()[0])
		// Создаем структуру нового клиента
		client := http.Client{}
		// Передаем запрос на сервис сервера
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
			log.Printf("Произошла ошибка при передаче запроса на сервер для изменения поста. Ошибка: %v", err)
			return
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			// Конвертируем данные о посте из БД который пользователь будет редактировать
			result, err := convertor.ConvertGetPosts(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при конвертации данных о посте из БД который пользователь будет редактировать. Ошибка: %v", err)
				return
			}

			err = t.ExecuteTemplate(w, "update_post.html", result)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при парсинге страницы с данными о посте который пользователь будет редактировать. Ошибка: %v", err)
				return
			}
		case http.StatusInternalServerError:
			discriptionMsg, err := convertor.DecodeErrorResponse(resp)
			if err != nil {
				errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при декодировании ответа ошибки и описания от сервера на запрос об данных о посте который пользователь будет редактировать")
				return
			}
			switch {
			// Получена ошибка что почта уже используется
			case discriptionMsg.Discription == "Email already exist":
				errorPage(w, errors.EmailAlreadyExists, http.StatusConflict)
				log.Printf("Пользователь пытается зарегестировать почту которая используется под другим аккаунтом")
				return
				// Получена ошибка что введены неверные учетные данные
			case discriptionMsg.Discription == "Invalid Credentials":
				errorPage(w, errors.InvalidCredentials, http.StatusBadRequest)
				log.Printf("Не валидные данные")
				return
			case discriptionMsg.Discription == "Not Found Any Data":
				errorPage(w, errors.NotFoundAnyDate, http.StatusBadRequest)
				log.Printf("Не найдено")
				return
			default:
				errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
				log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на получение старницы с данными о посте который пользователь собирается редактировать")
				return
			}
		default:
			errorPage(w, errors.ErrorNotMethod, http.StatusMethodNotAllowed)
			log.Printf("Получена ошибка сервера от сервиса сервера при передаче запроса на получение старницы с данными о посте который пользователь собирается редактировать")
			return
		}
	case http.MethodPost:
		// Конвертируем данные с изменениями поста от пользователя для передачи на сервер
		jsonData, err := convertor.ConvertUpdatePost(r)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при конвертации данных c изменениями внесенными пользователем о посте для передачи на сервер. Ошибка: %v", err)
				return
		}
		// Формирование URL  c конкретным id поста для редактирования
		updatePostsId := fmt.Sprintf(updatePosts+"%s", r.FormValue("postId"))
		// Формируем запрос на измение данных в посте
		req, err := http.NewRequest("PUT", updatePostsId, bytes.NewBuffer(jsonData))
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при подготовке запроса об изменениях данных в посте. Ошибка: %v", err)
				return
		}
		// Записываем куки из браузера в запрос на сервер
		req.AddCookie(r.Cookies()[0])
		// Указываем формат передачи данных 
		req.Header.Set("Content-Type", "application/json")
		// Создаем структуру нового клиента
		client := http.Client{}
		// Передаем запрос на сервис сервера
		resp, err := client.Do(req)
		if err != nil {
			errorPage(w, errors.ErrorServer, http.StatusInternalServerError)
				log.Printf("Произошла ошибка при передаче запроса на сервис сервера об изменениях данных в посте. Ошибка: %v", err)
				return
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusAccepted:
			
		case http.StatusInternalServerError:
		case default:
		}
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
					userPosts, err := convertor.ConvertAllPosts(resp)
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
