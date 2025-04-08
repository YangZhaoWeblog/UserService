

# 架构设计
总共分为 6 个 server:

+ Gatewaysvr: 处理前端请求, 聚合后台各个微服务 server 的处理结果, jwt 鉴权, router 路由
+ Commentsvr: 评论的创建、删除、列表查询等评论相关功能
+ Favoritesvr: 视频点赞、取消点赞、获取点赞列表等交互功能
+ Relationsvr: 用户关注、取消关注、粉丝列表等社交关系功能
+ Usersvr: 用户注册、登录、信息管理等账号相关功能
+ Videosvr: 视频上传、发布、流推送、视频信息管理等核心功能
  所有客户端的请求都是打到 gatewaysvr, 由 gatewaysvr 负责路由到下游具体的业务 server, 再由下游的对应的业务server处理完之后将数据返回到 gatewaysvr,
  gatesvr再根据需求，看是否需要做下游数据聚合，最终将数据处理成前端所需要的格式返回给抖音客户端

所有客户端请求统一由 Gatewaysvr 处理, 实现了:
1. 请求的统一入口
2. 服务端 API 聚合
3. 用户认证与授权
4. 请求的路由分发


project-root/
├── api                     # API 定义目录
│   ├── user               # 用户服务接口
│   │   ├── v1
│   │   │   └── user.proto
│   ├── video              # 视频服务接口
│   │   ├── v1
│   │   │   └── video.proto
│   ├── comment            # 评论服务接口
│   │   └── v1
│   └── gateway            # gateway接口
│       └── v1
├── cmd                    # 服务启动目录
│   ├── user              # 用户服务
│   │   └── main.go      
│   ├── video             # 视频服务
│   │   └── main.go
│   ├── comment           # 评论服务
│   │   └── main.go
│   └── gateway           # gateway服务
│       └── main.go
├── configs               # 配置文件目录
│   ├── user.yaml
│   ├── video.yaml
│   ├── comment.yaml
│   └── gateway.yaml
├── internal             # 内部代码
│   ├── user            # 用户服务实现
│   │   ├── biz
│   │   ├── data
│   │   ├── server
│   │   └── service
│   ├── video           # 视频服务实现
│   ├── comment         # 评论服务实现
│   └── gateway         # gateway实现
├── pkg                 # 公共代码包
│   ├── jwt            # JWT工具
│   ├── middleware     # 中间件
│   └── utils          # 通用工具
└── scripts            # 脚本文件
├── build.sh       # 构建脚本
└── run.sh         # 运行脚本