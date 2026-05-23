package agent

import (
	"context"

	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/ws"
)

// 不注册用户,直接发送消息
func SendReply(ctx context.Context, reply string) {
	aiClient := ws.Client{
		UserID:   config.Agent.UserID,
		Username: config.Agent.Name,
	}

	ws.ManagerInstance.BroadcastMessage([]byte(reply), &aiClient)
}
