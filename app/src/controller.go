package src

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type UbVideo struct {
	Id string `uri:"id" binding:"required"`
}

// =================[public functions]=====================
func GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"title": "Main Title",
	})
}

func GetVideo(c *gin.Context) {
	var ubVideo UbVideo
	if err := c.ShouldBindUri(&ubVideo); err != nil {
		errorPage(c, 400)
	}
	var video Video
	db := SqlConnect()
	defer db.Close()
	db.Where(&Video{Uid: ubVideo.Id}).Find(&video)
	if video.Uid == "" {
		errorPage(c, 404)
	} else {
		c.HTML(http.StatusOK, "video", gin.H{
			"title":     "Video Page",
			"videoId":   video.Uid,
			"videoName": video.Name,
		})
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
	uid, err := uuid.NewRandom()
	if err != nil {
		errorPage(c, 500)
	}
	videoName := c.PostForm("filename")
	file, header, err := c.Request.FormFile("file")
	if header == nil || err != nil {
		errorPage(c, 400)
	}
	ext := extractExtension(header.Filename)
	dir, _ := os.Getwd()
	videoDir := fmt.Sprintf("%s/resources/video/%s", dir, uid.String())
	os.Mkdir(videoDir, 777)
	videoPath := videoDir + "/hs" + ext
	out, _ := os.Create(videoPath)
	defer out.Close()
	_, err = io.Copy(out, file)

	db := SqlConnect()
	defer db.Close()
	video := Video{
		Uid:      uid.String(),
		Name:     videoName,
		IsEncode: false,
	}
	if err := video.Validate(); err != nil {
		errorPage(c, 400)
	}
	db.Create(&video)
	go func() {
		if err := convertHls(videoPath); err != nil {
			println(err)
		} else {
			db := SqlConnect()
			defer db.Close()
			video.IsEncode = true
			db.Save(&video)
		}
	}()

	c.String(http.StatusOK, "Upload Success")
}

func NotFound(c *gin.Context) {
	errorPage(c, 404)
}

// =================[private functions]=====================
func extractExtension(filename string) string {
	pos := strings.LastIndex(filename, ".")
	ext := filename[pos:]
	return ext
}

func convertHls(inputPath string) error {
	fmt.Println("Start convert to Hls...")
	prev, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	defer os.Chdir(prev)

	filename := filepath.Base(inputPath)
	vdir := filepath.Dir(inputPath)
	os.Chdir(vdir)

	cmd := exec.Command("ffmpeg", "-i", filename, "-vcodec",
		"libx264", "-vprofile", "baseline", "-acodec", "copy", "-ar",
		"44100", "-ac", "1", "-f", "segment", "-segment_format", "mpegts",
		"-segment_time", "10", "-segment_list", "hs.m3u8", "hs%3d.ts")
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		println("Error occurred")
	}
	fmt.Println("end")

	return err
}

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
