package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type data struct {
	Ext     string `json:"ext"`
	Tim     int    `json:"tim"`
	Source  string `json:"-"`
	IsImage bool   `json:"-"`
}
type list struct {
	Posts []data `json:"posts"`
}

// for prod
func H(w http.ResponseWriter, r *http.Request) {
	router().ServeHTTP(w, r)
}

// for dev
func main() {
	router().Run()
}
func router() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("./templates/*")
	r.GET("/:board/thread/:thread/:slug", serveThread)
	r.GET("/:board/thread/:thread", serveThread)
	r.NoRoute(serveNotFound)
	return r
}
func serveNotFound(c *gin.Context) {
	c.HTML(http.StatusOK, "main.html", gin.H{})
}
func serveThread(c *gin.Context) {
	board := c.Param("board")
	thread := c.Param("thread")
	url := fmt.Sprintf("https://a.4cdn.org/%s/thread/%s.json", board, thread)
	res, err := http.Get(url)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	var all list
	err = json.Unmarshal(b, &all)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	var media list
	for _, val := range all.Posts {
		img := isImage(val.Ext)
		vid := isVideo(val.Ext)
		if img || vid {
			source := fmt.Sprintf("http://is2.4chan.org/%s/%d%s", board, val.Tim, val.Ext)
			media.Posts = append(media.Posts, data{Source: source, IsImage: img})
		}
	}
	id := fmt.Sprintf("ID-%s-%s", board, thread)
	c.HTML(http.StatusOK, "thread.html", gin.H{"List": media, "ID": id})
}
func isImage(ext string) bool {
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}
func isVideo(ext string) bool {
	return ext == ".webm"
}
