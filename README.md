English | [简体中文](README_CN.md)
<p align="center">
    <a href="https://github.com/ClaretWheel1481/FeasOJ-Backend">
        <img src="public/logo.png" height="200"/>
    </a>
</p>

# FeasOJ
### Project Description
FeasOJ is an online programming practice platform based on Vue and Golang, aiming to provide users with a convenient and efficient learning and practice environment.
<br>
[FeasOJ-Frontend](https://github.com/ClaretWheel1481/FeasOJ)
[ImageGuard](https://github.com/ClaretWheel1481/ImageGuard)

### Project Structure
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
│  └─main.go    # Main entry file
└─Sandbox
```

### Environment
- Golang 1.23.1
- Redis
- MySQL 8.0+
- Docker

### How to run
1. Clone repository.
2. Install Docker.
3. Run `cd src` and `go mod tidy` Install dependencies.
4. Run `go run main.go` to start the back-end server.

### Notice
This is the first time I've written a big project with Vue + Golang, so the code is going to be terrible, but I'll keep going to improve it!
If you find any bugs, please open an issue.

### Thanks
- [Go](https://github.com/golang/go)
- [Gin](https://github.com/gin-gonic/gin)
- [gorm](https://github.com/go-gorm/gorm)
- [jwt-go](https://github.com/golang-jwt/jwt)
- [docker](https://github.com/moby/moby)
- [gomail](https://github.com/go-gomail/gomail)
- [go-redis](https://github.com/redis/go-redis)
