English | [简体中文](README_CN.md)
<p align="center">
    <a href="https://github.com/ClaretWheel1481/FeasOJ-Backend">
        <img src="public/logo.png" height="200"/>
    </a>
</p>

# FeasOJ
### Project Description
FeasOJ is an online programming practice platform based on Vue and Golang, supporting multi-languages, discussion forums, contests and other features, aiming to provide users with a convenient and efficient learning and practice environment.
<br>
[FeasOJ-Frontend](https://github.com/ClaretWheel1481/FeasOJ)
[ImageGuard](https://github.com/ClaretWheel1481/ImageGuard)
[Profanity Detector](https://github.com/ClaretWheel1481/ProfanityDetector)
[API Document(CN)](https://claret-feasoj.apifox.cn)
[API Document(EN)](https://claret-feasoj.apifox.cn/en/)

### Notice
The project is expanding as well as being microserviced, so watch out for the Readme.
If you find any bugs, please open an issue.

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
│  │    ├─locales   # i18n files
│  │    ├─scheduler
│  │    └─sql
│  ├─go.mod
│  └─main.go    # Main entry file
└─Sandbox   # After you build the executable file, copy the folder to the same directory as the executable file
```

### Environment
- Golang 1.24.0
- Redis
- MySQL 8.0+
- Docker
- RabbitMQ

### How to run
1. Clone this repository and [ImageGuard](https://github.com/ClaretWheel1481/ImageGuard), [Profanity Detector](https://github.com/ClaretWheel1481/ProfanityDetector)
2. Install Docker, RabbitMQ, Redis and MySQL
3. Start `MySQL`, `Redis`, `Docker`, `RabbitMQ`, `ImageGuard`, `ProfanityDetector` Services
4. Run `cd src` and `go mod tidy` Install dependencies
5. Config `src/config/global.go` (Check it for more details)
6. Run `go run main.go` to start the back-end server
7. Enter information that system show in the terminal

### Thanks
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
