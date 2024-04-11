package controller

import (
	"lzhuk/clients/pkg/config"
	"net/http"
	"time"
)



func Router(cfg config.Config) *http.Server {
	router := http.NewServeMux()

	router.HandleFunc("/", startPage)
	router.HandleFunc("/userd3", homePage)
	router.HandleFunc("/register", registerPage)
	router.HandleFunc("/login", loginPage)
	router.HandleFunc("/userd3/posts", createPost)
	router.HandleFunc("/userd3/myposts", myPosts)
	router.HandleFunc("/userd3/post/", getPost)
	router.HandleFunc("/userd3/createcomment", createComment)
	router.HandleFunc("/userd3/updatepost/", updatePost)
	router.HandleFunc("/userd3/deletepost/", deletePost)
	router.HandleFunc("/userd3/votepost", votePost)
	router.HandleFunc("/userd3/votecomment", voteComment)
	router.HandleFunc("/userd3/likeposts", likePost)
	router.HandleFunc("/userd3/updatecomment", updateComment)
	router.HandleFunc("/userd3/deletecomment", deleteComment)
	router.HandleFunc("/logout", logoutUser)
	router.HandleFunc("/userd3/category", categoryPosts)

	fileServer := http.FileServer(http.Dir("ui/css"))
	router.Handle("/ui/css/", http.StripPrefix("/ui/css/", fileServer))
	s := &http.Server{
		Addr:         cfg.Port,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,  // время ожидания для чтения данных
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second, // время ожидания для записи данных
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,  // время простоя
		Handler:      router,
	}
	return s
}
