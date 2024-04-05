package controller

const (
	allPost     = "http://localhost:8083/d3"
	registry    = "http://localhost:8083/register"
	login       = "http://localhost:8083/login"
	createPosts = "http://localhost:8083/d3/post-create"
	userPost    = "http://localhost:8083/d3/user-posts"
	getUserPost     = "http://localhost:8083/d3/post?id="
	createComments = "http://localhost:8083/d3/comment-create"
	updatePosts = "http://localhost:8083/d3/post-update?id="
	deletePosts = "http://localhost:8083/d3/post-delete?id="
	votePosts = "http://localhost:8083/d3/post-like"
	voteComments = "http://localhost:8083/d3/comment-like"
	likePosts = "http://localhost:8083/d3/user-likes"
	updateComments = "http://localhost:8083/d3/comment-update"
	deleteComments = "http://localhost:8083/d3/comment-delete"
	logoutUsers = "http://localhost:8083/logout"
)
