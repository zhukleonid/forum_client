package controller

import (
	"lzhuk/clients/pkg/config"
	"net/http"
	"time"
)

func Router(cfg config.Config) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", startPage)
	mux.HandleFunc("/userd3", homePage)
	mux.HandleFunc("/register", registerPage)
	mux.HandleFunc("/login", loginPage)
	mux.HandleFunc("/userd3/posts", createPost)
	mux.HandleFunc("/userd3/myposts", myPosts)
	mux.HandleFunc("/userd3/post/", getPost)
	mux.HandleFunc("/userd3/createcomment", createComment)
	mux.HandleFunc("/userd3/updatepost/", updatePost)
	mux.HandleFunc("/userd3/deletepost/", deletePost)
	mux.HandleFunc("/userd3/votepost", votePost)
	mux.HandleFunc("/userd3/votecomment", voteComment)
	mux.HandleFunc("/userd3/likeposts", likePost)
	mux.HandleFunc("/userd3/updatecomment", updateComment)

	fileServer := http.FileServer(http.Dir("ui/css"))
	mux.Handle("/ui/css/", http.StripPrefix("/ui/css/", fileServer))
	s := &http.Server{
		Addr:         cfg.Port,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,  // время ожидания для чтения данных
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second, // время ожидания для записи данных
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,  // время простоя
		Handler:      mux,
	}
	return s
}
