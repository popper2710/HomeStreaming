package src

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/my-repo/home_streaming/config"
	"html/template"
)

func parsedHtmlWithBase(name string) *template.Template {
	html := template.Must(template.ParseFiles("templates/base.tmpl", "templates/"+name+".tmpl"))
	return html
}
func createMyRender() multitemplate.Renderer {
	renderer := multitemplate.NewRenderer()
	pathList := []string{"index", "video", "upload", "login", "register", "error"}
	for _, v := range pathList {
		renderer.Add(v, parsedHtmlWithBase(v))
	}
	return renderer
}

func Router() *gin.Engine {
	router := gin.Default()

	// setting static files
	router.Static("/static", "./static")
	router.Static("/resources", "./resources")
	router.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	router.HTMLRender = createMyRender()

	// setting session
	var secret config.Secret
	secret.Init()
	store, err := redis.NewStore(
		secret.Redis.Size,
		secret.Redis.Network,
		secret.Redis.Host+":"+secret.Redis.Port,
		secret.Redis.Password,
		[]byte(secret.Redis.KeyPair))
	if err != nil {
		panic(err)
	}
	router.Use(sessions.Sessions("hs_sid", store))

	// setting urls
	router.GET("/index", GetIndex)
	router.GET("/video/:id", GetVideo)
	router.GET("/login", GetLogin)
	router.POST("/login", PostLogin)
	router.GET("/logout", GetLogout)
	router.GET("/register", GetRegister)
	router.POST("/register", PostRegister)
	router.GET("/upload", GetUpload)
	router.POST("/upload", PostUpload)
	return router
}
