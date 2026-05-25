# 项目介绍
一项基于 Go 语言构建的视频社交类后端服务项目，涵盖用户管理，视频操作，社交互动，行为交互，聊天等功能，同时集成JWT认证等安全机制。
# 技术栈
- **核心语言：** Go
- **接口定义：** Thrift
- **web框架：** Hertz
- **agent框架:**Eino
- **数据库：** MySQL
- **缓存：** Redis
# 项目结构
```txt
VideoWeb/
│   docker-compose.yml
│   dockerfile
│   go.mod
│   go.sum
│   README.md
│   router.go
│   router_gen.go
│
├───ai
│   ├───agent
│   │       agent.go
│   │       sender.go
│   └───internal
│       ├───chain
│       │       decide_chain.go
│       │       generate_chain.go
│       │
│       ├───memory
│       │       memory.go
│       │
│       ├───provider
│       │       model.go
│       │
│       └───trigger
│               emotion.go
│               mention.go
│               silence.go
│               trigger.go
├───biz
│   ├───handler
│   │   │   ping.go
│   │   │
│   │   ├───chat
│   │   │       chat_service.go
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
│   ├───router
│   │   │   register.go
│   │   │
│   │   ├───chat
│   │   │       chat.go
│   │   │
│   │   ├───interact
│   │   │       interact.go
│   │   │       middleware.go
│   │   │
│   │   ├───social
│   │   │       middleware.go
│   │   │       social.go
│   │   │
│   │   ├───static
│   │   │       static.go
│   │   │
│   │   ├───user
│   │   │       middleware.go
│   │   │       user.go
│   │   │
│   │   └───video
│   │           middleware.go
│   │           video.go
│   └───service
│       ├───chat
│       │       chat.go
│       │
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
├───cmd
│       main.go
│
├───config
│       ai_config.go
│       ai_config.yml
│       config.go
│       config.yml
│       init.go
│
├───idl
│       interact.thrift
│       social.thrift
│       user.thrift
│       video.thrift
│
├───internal
│   ├───dao
│   │       chat.go
│   │       interact.go
│   │       social.go
│   │       user.go
│   │       video.go
│   │
│   └───model
│           model.go
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
│           hashpassword.go
│           mfa.go
│           search.go
│           store.go
├───static
│   ├───avatar
│   │
│   └───video
│
└───ws
        client.go
        manager.go
```
# 文档地址
**apifox接口文档：** https://7dveb16f8y.apifox.cn/418681111e0
**飞书文档地址**https://ocnz01b9ahxp.feishu.cn/wiki/UHNVwV5fJiVTRLkxSwAcMRTenbb?from=from_copylink
