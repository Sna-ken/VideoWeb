# 项目介绍
一项基于 Go 语言构建的视频社交类后端服务项目，涵盖用户管理，视频操作，社交互动，行为交互，聊天等功能，同时集成JWT认证等安全机制。
# 技术栈
- **核心语言：** Go
- **接口定义：** Thrift
- **web框架：** Hertz
- **数据库：** MySQL
- **缓存：** Redis
# 项目结构
```txt
VideoWeb/
├── biz/                // 核心业务层：处理、服务、路由、模型定义
│   ├── handler/        // 请求处理：按业务模块拆分（user/video/social/interact）
│   │   ├── interact/   
│   │   ├── social/     
│   │   ├── user/       
│   │   └── video/      
│   ├── model/          // 业务模型：各模块数据结构定义
│   ├── router/         // 路由注册：按模块配置路由与中间件
│   │   ├── interact/   
│   │   ├── social/     
│   │   ├── user/       
│   │   └── video/      
│   └── service/        // 业务服务：核心逻辑实现
│       ├── interact/   
│       ├── social/     
│       ├── user/       
│       └── video/      
├── cmd/                // 程序入口
├── config/             // 配置管理（config.go/config.yml）
├── idl/                // 接口定义
├── internal/dao/       // 数据访问层：对数据库进行操作
├── middleware/         // 中间件：JWT 认证
├── pkg/                // 公共工具
│   ├── e/              // 错误码定义
│   ├── jwt/            // JWT 工具包
│   └── utils/          // 工具类
├── static/             // 静态资源：视频、头像等文件存储
│   ├── avatar/         // 头像存储目录
│   └── video/          // 视频存储目录
├── docker-compose.yml  // 容器化部署配置
├── dockerfile          // Docker 构建配置
├── go.mod/go.sum       // Go 依赖管理
└── README.md           // 项目文档（含 API 文档地址）
```
# 文档地址
**apifox接口文档：** https://7dveb16f8y.apifox.cn/418681111e0