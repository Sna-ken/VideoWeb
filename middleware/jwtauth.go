package middleware

import (
	"context"
	"log"
	"time"

	"github.com/Sna-ken/videoweb/biz/model/user"
	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/pkg/jwt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func JWTAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		accesstoken := c.Request.Header.Get("access_token")
		refreshtoken := c.Request.Header.Get("refresh_token")
		if accesstoken == "" {
			c.JSON(consts.StatusUnauthorized, &user.Base{
				Code: consts.StatusUnauthorized, Msg: "token is empty",
			})
			c.Abort()
			return
		}

		acClaims, err := jwt.ValidateAccessToken(accesstoken)
		if err == nil {
			c.Set("user_id", acClaims.UserID)
			c.Next(ctx)
			return
		}

		if refreshtoken == "" {
			c.JSON(consts.StatusUnauthorized, &user.Base{
				Code: consts.StatusUnauthorized, Msg: "accesstoken expried and lack refresh token",
			})
			c.Abort()
			return
		}

		rfClaims, err := jwt.ValidateRefreshToken(refreshtoken)
		if err != nil {
			c.JSON(consts.StatusUnauthorized, &user.Base{
				Code: consts.StatusUnauthorized, Msg: "refreshtoken expried",
			})
			c.Abort()
			return
		}

		val, err := config.REDISDB.Get(ctx, "user_rftoken:"+rfClaims.UserID).Result()
		if err != nil {
			c.JSON(consts.StatusUnauthorized, &user.Base{
				Code: consts.StatusUnauthorized, Msg: "session invalid:" + err.Error(),
			})
			c.Abort()
			return
		}

		if val != refreshtoken {
			c.JSON(consts.StatusUnauthorized, &user.Base{
				Code: consts.StatusUnauthorized, Msg: "session invalid: token mismatch",
			})
			c.Abort()
			return
		}

		newAccesstoken, errA := jwt.GenerateAccessToken(rfClaims.UserID)
		newRefreshtoken, errR := jwt.GenerateRefreshToken(rfClaims.UserID)
		if errA != nil || errR != nil {
			c.JSON(consts.StatusInternalServerError, &user.Base{
				Code: consts.StatusInternalServerError, Msg: "genrate token failed" + errA.Error() + " " + errR.Error(),
			})
			c.Abort()
			return
		}

		duration := time.Duration(config.JWTConfig.RefreshTokenExpiry) * time.Second
		err = config.REDISDB.Set(ctx, "user_rftoken:"+rfClaims.UserID, newRefreshtoken, duration).Err()
		if err != nil {
			log.Println("Failed to save refresh token to Redis:", err)
		} else {
			log.Println("Successfully saved Redis key:", "user_rftoken:"+rfClaims.UserID)
		}

		c.Header("new_access_token", newAccesstoken)
		c.Header("new_refresh_token", newRefreshtoken)

		c.Set("user_id", rfClaims.UserID)
		c.Next(ctx)
	}
}
