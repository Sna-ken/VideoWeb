package user

import (
	"strings"
	"testing"

	user "github.com/Sna-ken/videoweb/biz/model/user"
	service "github.com/Sna-ken/videoweb/biz/service/user"
	"github.com/Sna-ken/videoweb/pkg/e"
	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	type testCase struct {
		name         string
		url          string
		body         string
		mockerr      error
		expectedMsg  string
		expectedCode int
	}

	testCases := []testCase{
		{
			name:         "success",
			url:          "/user/register",
			body:         `{"username":"test","password":"test123"}`,
			mockerr:      nil,
			expectedMsg:  `User register successfully`,
			expectedCode: 200,
		},
		{
			name:         "username already exists",
			url:          "/user/register",
			body:         `{"username":"exist","password":"test123"}`,
			mockerr:      e.New(consts.StatusConflict, "User has existed", nil),
			expectedMsg:  `User has existed`,
			expectedCode: 409,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/user/register", Register)

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock((*service.UserService).RegisterService).To(
				func(s *service.UserService, req *user.RegisterReq) error {
					if req.Username == "exist" {
						return tc.mockerr
					}
					return nil
				}).Build()

			// 构造请求体
			var body *ut.Body
			if tc.body != "" {
				body = &ut.Body{
					Body: strings.NewReader(tc.body),
					Len:  len(tc.body),
				}
				// 发送模拟HTTP请求
				resp := ut.PerformRequest(router, consts.MethodPost, tc.url, body,
					ut.Header{Key: "Content-Type", Value: "application/json"})
				// 验证响应
				assert.Equal(t, tc.expectedCode, resp.Result().StatusCode())
				assert.Contains(t, string(resp.Result().Body()), tc.expectedMsg)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	type testCase struct {
		name         string
		url          string
		body         string
		mockerr      error
		mockresp     user.LoginResp
		expectedMsg  string
		expectedCode int
	}
	// bool类型没有引号
	testCases := []testCase{
		{
			name:    "success without mfa",
			url:     "/user/login",
			body:    `{"username":"test","password":"test123","mfa_enabled":false,"code":""}`,
			mockerr: nil,
			mockresp: user.LoginResp{
				Base:         &user.Base{Code: consts.StatusOK, Msg: "User login successfully"},
				Data:         &user.Data{UserID: "123456", Username: "success without mfa", AvatarURL: "avatar_url"},
				AccessToken:  "accesstoken",
				RefreshToken: "refreshtoken",
			},
			expectedMsg:  `User login successfully`,
			expectedCode: 200,
		},
		{
			name:    "success with mfa",
			url:     "/user/login",
			body:    `{"username":"test","password":"test123","mfa_enabled":true,"code":"123456"}`,
			mockerr: nil,
			mockresp: user.LoginResp{
				Base:         &user.Base{Code: consts.StatusOK, Msg: "User login successfully"},
				Data:         &user.Data{UserID: "123456", Username: "success with mfa", AvatarURL: "avatar_url"},
				AccessToken:  "accesstoken",
				RefreshToken: "refreshtoken",
			},
			expectedMsg:  `User login successfully`,
			expectedCode: 200,
		},
		{
			name:         "wrong password",
			url:          "/user/login",
			body:         `{"username":"test","password":"worng123","mfa_enabled":false,"code":""}`,
			mockerr:      e.New(consts.StatusUnauthorized, "Wrong password", nil),
			expectedMsg:  `Wrong password`,
			expectedCode: 401,
		},
		{
			name:         "empty mfa",
			url:          "/user/login",
			body:         `{"username":"test","password":"test123","mfa_enabled":true,"code":""}`,
			mockerr:      e.New(consts.StatusUnauthorized, "MFA required", nil),
			expectedMsg:  `MFA required`,
			expectedCode: 401,
		},
		{
			name:         "wrong code",
			url:          "/user/login",
			body:         `{"username":"test","password":"test123","mfa_enabled":true,"code":"456789"}`,
			mockerr:      e.New(consts.StatusUnauthorized, "Invalid MFA code", nil),
			expectedMsg:  `Invalid MFA code`,
			expectedCode: 401,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/user/login", Login)

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock((*service.UserService).LoginService).To(
				func(s *service.UserService, req *user.LoginReq) (error, *user.LoginResp) {
					if req.MfaEnabled {
						if req.Code == "" {
							return tc.mockerr, nil
						}
						if req.Code != "123456" {
							return tc.mockerr, nil
						}
					}

					if req.Password != "test123" {
						return tc.mockerr, nil
					}

					return nil, &tc.mockresp
				}).Build()

			var body *ut.Body
			if tc.body != "" {
				body = &ut.Body{
					Body: strings.NewReader(tc.body),
					Len:  len(tc.body),
				}

				resp := ut.PerformRequest(router, consts.MethodPost, tc.url, body,
					ut.Header{Key: "Content-Type", Value: "application/json"})

				assert.Equal(t, tc.expectedCode, resp.Result().StatusCode())
				assert.Contains(t, string(resp.Result().Body()), tc.expectedMsg)
			}
		})
	}
}

func TestUserInfo(t *testing.T) {
	type testCase struct {
		name         string
		url          string
		userid       string
		mockerr      error
		mockresp     user.UserInfoResp
		expectedMsg  string
		expectedCode int
	}

	testCases := []testCase{
		{
			name:    "success",
			url:     "/user/info",
			userid:  "123456",
			mockerr: nil,
			mockresp: user.UserInfoResp{
				Base: &user.Base{Code: consts.StatusOK, Msg: "User info fetched successfully"},
				Data: &user.Data{UserID: "123456", Username: "success", AvatarURL: "avatar_url"},
			},
			expectedMsg:  "User info fetched successfully",
			expectedCode: 200,
		},
		{
			name:         "wrong_token",
			url:          "/user/info",
			userid:       "456789",
			mockerr:      e.New(consts.StatusNotFound, "User not found", nil),
			expectedMsg:  "User not found",
			expectedCode: 404,
		},
	}
	router := route.NewEngine(&config.Options{})
	router.GET("/user/info", UserInfo)

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock((*service.UserService).UserInfoService).To(
				func(s *service.UserService, req *user.UserInfoReq, userID string) (error, *user.UserInfoResp) {
					if userID != "123456" {
						return tc.mockerr, nil
					}
					return nil, &tc.mockresp
				}).Build()

			mockey.Mock((*app.RequestContext).GetString).To(
				func(c *app.RequestContext, key string) string {
					if key == "user_id" {
						return tc.userid
					}
					return ""
				}).Build()

			resp := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)

			assert.Equal(t, tc.expectedCode, resp.Result().StatusCode())
			assert.Contains(t, string(resp.Result().Body()), tc.expectedMsg)
		})
	}

}
