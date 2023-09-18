package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/gin-gonic/gin"
)

//go:embed web
var embededFiles embed.FS

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func main() {
	webfs, err := fs.Sub(embededFiles, "web")
	if err != nil {
		panic(err)
	}
	fsys, err := fs.Sub(embededFiles, "web/assets")
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	api := r.Group("api")
	{
		api.GET("/hello", func(context *gin.Context) {
			context.String(http.StatusOK, "Hello World!")
		})
	}
	templ := template.Must(template.New("").ParseFS(webfs, "*.html"))
	r.SetHTMLTemplate(templ)

	r.StaticFileFS("/vite.svg", "/vite.svg", http.FS(webfs))
	r.StaticFileFS("/index.html", "/index.html", http.FS(webfs))
	r.StaticFS("/assets", http.FS(fsys))
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	go open("http://localhost:8080")
	r.Run(":8080")
}
