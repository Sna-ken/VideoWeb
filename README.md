文档地址：https://7dveb16f8y.apifox.cn/418681111e0

```
目录树
│
├───biz
│   ├───handler
│   │   │   ping.go
│   │   │
│   │   ├───interact
│   │   │       interact_service.go
│   │   │
│   │   ├───social
│   │   │       social_service.go
│   │   │
│   │   ├───user
│   │   │       user_service.go
│   │   │
│   │   └───video
│   │           video_service.go
│   │
│   ├───model
│   │   ├───interact
│   │   │       interact.go
│   │   │
│   │   ├───social
│   │   │       social.go
│   │   │
│   │   ├───user
│   │   │       user.go
│   │   │
│   │   └───video
│   │           video.go
│   │
│   ├───router
│   │   │   register.go
│   │   │
│   │   ├───interact
│   │   │       interact.go
│   │   │       middleware.go
│   │   │
│   │   ├───social
│   │   │       middleware.go
│   │   │       social.go
│   │   │
│   │   ├───user
│   │   │       middleware.go
│   │   │       user.go
│   │   │
│   │   └───video
│   │           middleware.go
│   │           video.go
│   │
│   └───service
│       ├───interact
│       │       interact.go
│       │
│       ├───social
│       │       social.go
│       │
│       ├───user
│       │       user.go
│       │
│       └───video
│               video.go
│
├───cmd
│       main.go
│
├───config
│       cofig.go
│       db.go
│
├───idl
│       interact.thrift
│       social.thrift
│       user.thrift
│       video.thrift
│
├───internal
│   └───dao
│           interact.go
│           model.go
│           social.go
│           user.go
│           video.go
│
├───middleware
│       jwtauth.go
│
├───pkg
│   ├───e
│   │       e.go
│   │
│   ├───jwt
│   │       jwt.go
│   │
│   └───utils
│           file.go
│           hashpassword.go
│           store.go
│
└───static
    │   default-avatar.png
    │
    ├───avatar
    │  
    │   
    │   
    │ 
    │
    └───video
```
