package src

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type UbVideoId struct {
	Id string `uri:"id" binding:"required"`
}

// =================[public functions]=====================
func GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"title": "Main Title",
	})
}

func GetVideo(c *gin.Context) {
	var videoId UbVideoId
	if err := c.ShouldBindUri(&videoId); err != nil {
		errorPage(c, 400)
	}
	if videoId.Id == "test" {
		c.HTML(http.StatusOK, "videoId", gin.H{
			"title":    "Video Page",
			"video_id": videoId.Id,
		})
	} else {
		errorPage(c, 404)
	}
}

func GetLogin(c *gin.Context) {
	if isAuthenticated(c) {
		c.Redirect(http.StatusMovedPermanently, "index")
	} else {
		c.HTML(http.StatusOK, "login", gin.H{
			"title": "Login Page",
		})
	}
}

func PostLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user User

	db := SqlConnect()
	defer db.Close()
	db.Where("name = ?", username).Find(&user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusOK, "login", gin.H{
			"title":     "Login Page",
			"error_msg": "Invalid username or password",
		})
		return
	}
	login(c, user.Uid)
	c.String(http.StatusOK, user.Uid)
}

func GetLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		panic(err)
	}
	c.String(http.StatusOK, "session Deleted")
}

func GetRegister(c *gin.Context) {
	if isAuthenticated(c) {
		c.Redirect(http.StatusMovedPermanently, "index")
	} else {
		c.HTML(http.StatusOK, "register", gin.H{
			"title": "Register Page",
		})
	}
}

func PostRegister(c *gin.Context) {
	uid := extractULID()
	username := c.PostForm("username")
	email := c.PostForm("email")
	rawPassword := c.PostForm("password")
	passwordConfirmation := c.PostForm("password_confirmation")

	if rawPassword != passwordConfirmation {
		c.HTML(http.StatusOK, "register", gin.H{
			"title":    "Register Page",
			"errorMsg": "Not Match Password",
		})
		return
	}
	encPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		errorPage(c, 400)
	}
	user := User{
		Uid:       uid,
		Name:      username,
		Email:     email,
		Password:  string(encPassword),
		LastLogin: time.Now(),
		CreatedAt: time.Now(),
	}
	err = user.Validate()
	if err != nil {
		c.HTML(http.StatusOK, "register", gin.H{
			"title":    "Register Page",
			"errorMsg": "Validation Error",
		})
		return
	}

	db := SqlConnect()
	defer db.Close()
	db.Create(&user)

	login(c, uid)
	c.Redirect(http.StatusMovedPermanently, "index")
}

func GetUpload(c *gin.Context) {
	c.HTML(http.StatusOK, "upload", gin.H{
		"title": "Upload Page",
	})
}

func PostUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		errorPage(c, 400)
	}
	fileName := header.Filename
	dir, _ := os.Getwd()
	out, err := os.Create(dir + "/resources/video/" + fileName)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	c.String(http.StatusOK, "Upload Success")
}

// =================[private functions]=====================
func errorPage(c *gin.Context, code int) {
	codeMap := map[int]string{400: "Bad Request", 404: "Not Found", 500: "Server Error"}
	if v, ok := codeMap[code]; ok {
		c.HTML(http.StatusBadRequest, "error", gin.H{
			"title":   v,
			"code":    code,
			"message": v,
		})
		return
	} else {
		panic("specified incorrect error code:" + string(code))
	}

}

func extractULID() string {
	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func isAuthenticated(c *gin.Context) bool {
	session := sessions.Default(c)
	v := session.Get("uid")
	return v != nil
}

func login(c *gin.Context, uid string) {
	session := sessions.Default(c)
	session.Set("uid", uid)
	if err := session.Save(); err != nil {
		errorPage(c, 400)
	}
}
