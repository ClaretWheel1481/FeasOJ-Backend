[English](README.md) | 简体中文
<p align="center">
    <a href="https://github.com/ClaretWheel1481/FeasOJ-Backend">
        <img src="public/logo.png" height="200"/>
    </a>
</p>

# FeasOJ-Backend
### 项目简介
FeasOJ 是一个基于 Vue 和 Golang 的在线编程练习平台，旨在为用户提供一个方便、高效的学习和练习环境。
<br>
[FeasOJ-Frontend](https://github.com/ClaretWheel1481/FeasOJ)
[ImageGuard](https://github.com/ClaretWheel1481/ImageGuard)

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
│  │     └─sql
│  ├─go.mod
│  └─main.go    # 程序主入口
└─Sandbox
```

### 环境
- Golang 1.23.1
- Redis
- MySQL 8.0+
- Docker

### 如何运行
1. 克隆此库。
2. 安装 Docker。
3. 运行 `cd src` 和 `go mod tidy` 下载依赖.
4. 运行 `go run main.go` 启动后端服务器.

### 注意
这是我第一次用Vue + Golang写大项目，所以代码会很糟糕，不过我会一直去改进它！
如果你找到任何Bug请发布Issue告诉我。

### 致谢
- [Go](https://github.com/golang/go)
- [Gin](https://github.com/gin-gonic/gin)
- [gorm](https://github.com/go-gorm/gorm)
- [jwt-go](https://github.com/golang-jwt/jwt)
- [docker](https://github.com/moby/moby)
- [gomail](https://github.com/go-gomail/gomail)
- [go-redis](https://github.com/redis/go-redis)