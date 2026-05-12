package chat

import (
	"github.com/Sna-ken/videoweb/biz/handler/chat"
	"github.com/Sna-ken/videoweb/middleware"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func Register(r *server.Hertz) {
	root := r.Group("/", middleware.JWTAuth())
	{
		root.GET("/ws/chat", chat.Chat)
	}
}

func GeneratedRegister(r *server.Hertz) {
	Register(r)
}
