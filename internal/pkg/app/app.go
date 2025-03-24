package app

import (
	"fmt"
	"io"
	"log"
	"myapp/internal/app/communication/bookmark"
	"myapp/internal/app/communication/commentDelete"
	"myapp/internal/app/communication/commentDislike"
	"myapp/internal/app/communication/commentLike"
	"myapp/internal/app/communication/like"
	"myapp/internal/app/communication/repost"
	"myapp/internal/app/communication/subscribe"
	"myapp/internal/app/communication/view"
	"myapp/internal/app/controllers"
	"myapp/internal/app/createComment"
	"myapp/internal/app/createImage"
	"myapp/internal/app/createPost"
	"myapp/internal/app/deleteImage"
	"myapp/internal/app/deletePost"
	"myapp/internal/app/getComments"
	"myapp/internal/app/getPost"
	"myapp/internal/app/getPosts"
	"myapp/internal/app/getUser"
	"myapp/internal/app/getUserInfo"
	getuserposts "myapp/internal/app/getUserPosts"
	"myapp/internal/app/getUsers"
	"myapp/internal/app/login"
	"myapp/internal/app/ping"
	"myapp/internal/app/randomPosts"
	"myapp/internal/app/search"
	"myapp/internal/app/service"
	"myapp/internal/app/updateUser"
	"myapp/internal/mw/checkToken"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	e              *getPosts.GetPosts
	s              *service.Service
	ping           *ping.Ping
	getUsers       *getUsers.GetUsers
	createUser     *controllers.CreateUser
	login          *login.Login
	getPost        *getPost.GetPost
	createPost     *createPost.CreatePost
	deletePost     *deletePost.DeletePost
	like           *like.Like
	repost         *repost.Repost
	subscribe      *subscribe.Subscribe
	bookmark       *bookmark.Bookmark
	view           *view.View
	getUser        *getUser.GetUser
	getUserInfo    *getUserInfo.GetUserInfo
	getUserPosts   *getuserposts.GetUserPosts
	updateUser     *updateUser.UpdateUser
	randomPosts    *randomPosts.RandomPosts
	search         *search.Search
	createImage    *createImage.CreateImage
	deleteImage    *deleteImage.DeleteImage
	createComment  *createComment.CreateComment
	getComments    *getComments.GetComments
	commentLike    *commentLike.CommentLike
	commentDislike *commentDislike.CommentDislike
	commentDelete  *commentDelete.CommentDelete

	echo *echo.Echo
}

func New() (*App, error) {
	a := &App{}
	a.echo = echo.New()

	a.echo.Use(middleware.Logger())
	a.echo.Use(middleware.Recover())
	a.echo.Use(middleware.CORS())
	a.echo.Logger.SetOutput(io.Discard)

	a.e = getPosts.New()
	a.ping = ping.New()
	a.getUsers = getUsers.New()
	a.getUser = getUser.New()
	a.getUserInfo = getUserInfo.New()
	a.getUserPosts = getuserposts.New()
	a.search = search.New()

	a.createUser = controllers.New()
	a.login = login.New()

	a.getPost = getPost.New()
	a.createPost = createPost.New()
	a.deletePost = deletePost.New()
	a.randomPosts = randomPosts.New()

	a.like = like.New()
	a.bookmark = bookmark.New()
	a.subscribe = subscribe.New()
	a.repost = repost.New()

	a.updateUser = updateUser.New()

	a.createImage = createImage.New()
	a.deleteImage = deleteImage.New()

	a.createComment = createComment.New()
	a.commentLike = commentLike.New()
	a.commentDislike = commentDislike.New()
	a.commentDelete = commentDelete.New()

	a.echo.GET("/ping", a.ping.Status)
	a.echo.GET("/getPosts", a.e.Status, checkToken.CheckToken)
	a.echo.GET("/getRandomPosts", a.randomPosts.Status, checkToken.CheckToken)
	a.echo.GET("/search", a.search.Status, checkToken.CheckToken)

	a.echo.POST("/register", a.createUser.Status)
	a.echo.POST("/login", a.login.Status)

	a.echo.GET("/getPost", a.getPost.Status, checkToken.CheckToken)
	a.echo.POST("/createPost", a.createPost.Status, checkToken.CheckToken)
	a.echo.DELETE("/deletePost", a.deletePost.Status, checkToken.CheckToken)
	a.echo.GET("/getSelf", a.getUser.Status, checkToken.CheckToken)
	a.echo.GET("/getUser", a.getUserInfo.Status, checkToken.CheckToken)
	a.echo.POST("/getUserPosts", a.getUserPosts.Status)

	a.echo.POST("/like", a.like.Status, checkToken.CheckToken)
	a.echo.POST("/bookmark", a.bookmark.Status, checkToken.CheckToken)
	a.echo.POST("/subscribe", a.subscribe.Status, checkToken.CheckToken)
	a.echo.POST("/view", a.view.Status, checkToken.CheckToken)
	a.echo.POST("/repost", a.repost.Status, checkToken.CheckToken)

	a.echo.POST("/createImage", a.createImage.Status, checkToken.CheckToken)
	a.echo.DELETE("/deleteImage", a.deleteImage.Status, checkToken.CheckToken)

	a.echo.POST("/createComment", a.createComment.Status, checkToken.CheckToken)
	a.echo.POST("/likeComment", a.commentLike.Status, checkToken.CheckToken)
	a.echo.POST("/dislikeComment", a.commentDislike.Status, checkToken.CheckToken)
	a.echo.POST("/deleteComment", a.commentDelete.Status, checkToken.CheckToken)
	a.echo.GET("/getComments", a.getComments.Status, checkToken.CheckToken)

	a.echo.PATCH("/updateUser", a.updateUser.Status, checkToken.CheckToken)

	return a, nil
}

func (a *App) Run() error {
	fmt.Println("server running")

	err := a.echo.Start(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
