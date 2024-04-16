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
	categoryGet = "http://localhost:8083/d3/category?name="
	auth = "http://localhost:8083/auth"
)

const (
	googleAuthEndPoint            = "https://accounts.google.com/o/oauth2/auth"
	googleAuthEndPointAccessToken = "https://accounts.google.com/o/oauth2/token"
	googleUserInfoURL             = "https://www.googleapis.com/oauth2/v3/userinfo"

	githubAuthEndPoint       = "https://github.com/login/oauth/authorize"
	githubAuthEndAccessToken = "https://github.com/login/oauth/access_token"
	githubUserInfoURL        = "https://api.github.com/user"

	client_idGIT     = "21c2671efe47648ceedd"
	client_secretGIT = "fbf46e505b7583bd24c5309bd342379f80591e68"
	callbackGIT      = "http://localhost:8082/github/callback"
)