package main

import (
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "templates/index.gohtml")
	r.AddFromFiles("video", "templates/view_video.gohtml")
	return r
}

type Video struct {
	Id string `uri:"id" binding:"required"`
}

func main() {
	router := gin.Default()
	router.Static("/static", "./static")
	router.HTMLRender = createMyRender()

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{
			"title": "Main Theme",
		})
	})
	router.GET("/video/:id", func(c *gin.Context) {
		var video Video
		if err := c.ShouldBindUri(&video); err != nil {
			c.String(http.StatusInternalServerError, "Server Error")
			return
		}
		if video.Id == "test" {
			c.HTML(http.StatusOK, "video", gin.H{
				"title":    "Video Page",
				"video_id": video.Id,
			})
		} else {
			c.String(http.StatusNotFound, "Not Found")
		}
	})

	router.Run(":3000")
}
