package controller

const (
	allPost     = "http://localhost:8083/userd3"
	registry    = "http://localhost:8083/register"
	login       = "http://localhost:8083/login"
	createPosts = "http://localhost:8083/userd3/post-create"
	userPost    = "http://localhost:8083/userd3/myposts"
	getUserPost     = "http://localhost:8083/userd3/post?id="
	createComments = "http://localhost:8083/userd3/comment-create"
	updatePosts = "http://localhost:8083/userd3/post-update?id="
	deletePosts = "http://localhost:8083/userd3/post-delete?id="
	votePosts = "http://localhost:8083/userd3/post-like"
	voteComments = "http://localhost:8083/userd3/comment-like"
)
