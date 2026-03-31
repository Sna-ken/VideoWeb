package middleware

import (
	"context"

	"github.com/Sna-ken/videoweb/biz/model/user"
	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/pkg/jwt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func JWTAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {

		access := string(c.GetHeader("access_token"))

		if claims, err := jwt.ValidateAccessToken(access); err == nil {
			c.Set("user_id", claims.UserID)
			c.Next(ctx)
			return
		}

		refresh := string(c.GetHeader("refresh_token"))
		if refresh == "" {
			c.JSON(consts.StatusUnauthorized, &user.Base{
				Code: consts.StatusUnauthorized,
				Msg:  "missing refresh token",
			})
			c.Abort()
			return
		}

		rfClaims, err := jwt.ValidateRefreshToken(refresh)
		if err != nil {
			c.JSON(consts.StatusUnauthorized, &user.Base{
				Code: consts.StatusUnauthorized,
				Msg:  "refreshtoken invalid",
			})
			c.Abort()
			return
		}

		val, err := config.REDISDB.Get(ctx, "user_rftoken:"+rfClaims.UserID).Result()
		if err != nil || val != refresh {
			c.JSON(consts.StatusUnauthorized, &user.Base{
				Code: consts.StatusUnauthorized,
				Msg:  "session invalid",
			})
			c.Abort()
			return
		}

		newAccess, _ := jwt.GenerateAccessToken(rfClaims.UserID)

		c.Header("access_token", newAccess)
		c.Set("user_id", rfClaims.UserID)
		c.Next(ctx)
	}
}
