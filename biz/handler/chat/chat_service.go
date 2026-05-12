package chat

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Sna-ken/videoweb/biz/service/chat"
	"github.com/Sna-ken/videoweb/internal/dao"
	"github.com/Sna-ken/videoweb/ws"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/websocket"
)

var upgrade = websocket.HertzUpgrader{
	CheckOrigin: func(ctx *app.RequestContext) bool {
		return true
	},
}

func Chat(ctx context.Context, c *app.RequestContext) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.AbortWithStatus(consts.StatusUnauthorized)
		return
	}

	username, err := dao.FindUsernameByID(ctx, userID)
	if err != nil {
		c.AbortWithStatus(consts.StatusInternalServerError)
		return
	}

	err = upgrade.Upgrade(c, func(conn *websocket.Conn) {
		defer conn.Close()

		client := &ws.Client{
			UserID:   userID,
			Username: username,
			Conn:     conn,
			Send:     make(chan []byte, 256),
		}

		ws.ManagerInstance.AddClient(client)
		go func() {
			messages, err := chat.GetRecentMessages(ctx)
			if err != nil {
				log.Printf("Get reccent masseges failed:%v", err)
				return
			}
			for _, msg := range messages {
				data, _ := json.Marshal(msg)
				client.Send <- data
			}
		}()

		go client.WriteMessage()
		client.ReadMessage()
	})

	if err != nil {
		log.Printf("upgrade failed:%v", err)
		return
	}
}
