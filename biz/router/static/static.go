package static

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterStaticFS(r *server.Hertz) {
	root := r.Group("/")
	{
		root.StaticFS("/avatar/default", &app.FS{Root: "./static/", GenerateIndexPages: false})
		root.StaticFS("/avatar", &app.FS{Root: "./static/", GenerateIndexPages: false})
		root.StaticFS("/video", &app.FS{Root: "./static/", GenerateIndexPages: false})
	}
}

func GeneratedRegisterStaticFS(r *server.Hertz) {
	RegisterStaticFS(r)
}
