package main

import (
	"html/template"
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	html := template.Must(template.ParseFiles("templates/base.gohtml", "templates/index.gohtml"))
	r.Add("index", html)
	html = template.Must(template.ParseFiles("templates/base.gohtml", "templates/video.gohtml"))
	r.Add("video", html)
	html = template.Must(template.ParseFiles("templates/base.gohtml", "templates/error.gohtml"))
	r.Add("error", html)
	return r
}

type Video struct {
	Id string `uri:"id" binding:"required"`
}

func main() {
	router := gin.Default()
	router.Static("/static", "./static")
	router.Static("/resources", "./resources")
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.HTMLRender = createMyRender()

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{
			"title": "Main Theme",
		})
	})
	router.GET("/video/:id", func(c *gin.Context) {
		var video Video
		if err := c.ShouldBindUri(&video); err != nil {
			c.HTML(http.StatusInternalServerError, "error", gin.H{
				"title":   "Server Error",
				"code":    500,
				"message": "Server Error",
			})
		}
		if video.Id == "test" {
			c.HTML(http.StatusOK, "video", gin.H{
				"title":    "Video Page",
				"video_id": video.Id,
			})
		} else {
			c.HTML(http.StatusNotFound, "error", gin.H{
				"title":   "Page Not Found",
				"code":    404,
				"message": "Page Not Found.",
			})
		}
	})

	router.Run(":3000")
}
