package middleware

import (
	"context"

	"github.com/Sna-ken/videoweb/biz/model/user"
	"github.com/Sna-ken/videoweb/config"
	"github.com/Sna-ken/videoweb/pkg/jwt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

/*
较最开始进行答辩时的版本对JWT进行了优化，相较之前解决了只需要refresh_token正确access_token无论怎么写都能登陆的问题
目前是只对access_token进行验证，若有效即放行；无效即对refresh_token进行验证，通过后方向并返回新的accesss_token
*/

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
				Msg:  "refresh token invalid",
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

		newAccess, err := jwt.GenerateAccessToken(rfClaims.UserID)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, &user.Base{
				Code: consts.StatusInternalServerError,
				Msg:  "failed to generate access token",
			})
			c.Abort()
			return

		}

		c.Header("access_token", newAccess)
		c.Set("user_id", rfClaims.UserID)
		c.Next(ctx)
	}
}
