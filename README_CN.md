[English](README.md) | 简体中文
<p align="center">
    <a href="https://github.com/ClaretWheel1481/FeasOJ-Backend">
        <img src="public/logo.png" height="200"/>
    </a>
</p>

# FeasOJ-Backend
### 项目简介
FeasOJ 是一个基于 Vue 和 Golang 的在线编程练习平台，支持多国语言、讨论区、竞赛等功能，旨在为用户提供一个方便、高效的学习和练习环境。
<br>
[FeasOJ-Frontend](https://github.com/ClaretWheel1481/FeasOJ)
[ImageGuard](https://github.com/ClaretWheel1481/ImageGuard)
[Profanity Detector](https://github.com/ClaretWheel1481/ProfanityDetector)
[API Document(CN)](https://claret-feasoj.apifox.cn)
[API Document(EN)](https://claret-feasoj.apifox.cn/en/)

### 注意
该项目正在不断扩大以及微服务化，请注意阅读Readme。
如果你找到任何Bug请发布Issue告诉我。

### 项目结构
```
FeasOJ-Backend
│ 
├─src
│  ├─config
│  ├─internal
│  │  ├─gin
│  │  ├─global
│  │  ├─judge
│  │  ├─middlewares
│  │  ├─router
│  │  └─utils
│  │    ├─locales   # i18n 翻译文件
│  │    ├─scheduler
│  │    └─sql
│  ├─go.mod
│  └─main.go    # 程序主入口
└─Sandbox   # 当你构建出可执行文件后，请将该文件夹复制到可执行文件同一目录下
```

### 环境
- Golang 1.24.0
- Redis
- MySQL 8.0+
- Docker
- RabbitMQ

### 如何运行
1. 克隆此库以及[ImageGuard](https://github.com/ClaretWheel1481/ImageGuard)和[Profanity Detector](https://github.com/ClaretWheel1481/ProfanityDetector)
2. 安装 Docker、MySQL、Redis和RabbitMQ
3. 运行 `MySQL`、`Redis`、`Docker`、`RabbitMQ`、`ImageGuard` 和 `Profanity Detector`
4. 运行 `cd src` 和 `go mod tidy` 下载依赖
5. 配置 `src/config/global.go` (打开文件查看以获取详细配置)
6. 运行 `go run main.go` 启动后端服务器
7. 输入控制台提示的连接信息即可

### 致谢
- [Go](https://github.com/golang/go)
- [Gin](https://github.com/gin-gonic/gin)
- [gorm](https://github.com/go-gorm/gorm)
- [jwt-go](https://github.com/golang-jwt/jwt)
- [docker](https://github.com/moby/moby)
- [gomail](https://github.com/go-gomail/gomail)
- [go-redis](https://github.com/redis/go-redis)
- [go-i18n](https://github.com/nicksnyder/go-i18n)
- [gocron](https://github.com/go-co-op/gocron)
- [amqp091-go](https://github.com/rabbitmq/amqp091-go)