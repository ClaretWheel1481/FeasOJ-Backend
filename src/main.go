package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"src/config"
	"src/internal/global"
	"src/internal/judge"
	"src/internal/router"
	"src/internal/utils"
	"src/internal/utils/sql"

	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("[FeasOJ] The server is starting...")
	global.CurrentDir, _ = os.Getwd()
	global.ParentDir = filepath.Dir(global.CurrentDir)

	// 定义目录映射
	dirs := map[string]*string{
		"config":      &global.ConfigDir,
		"certificate": &global.CertDir,
		"avatars":     &global.AvatarsDir,
		"codefiles":   &global.CodeDir,
		"logs":        &global.LogDir,
	}

	// 遍历map，设置路径并创建不存在的目录
	for name, dir := range dirs {
		var mainDir string
		if global.DebugMode {
			mainDir = global.ParentDir
		} else {
			mainDir = global.CurrentDir
		}
		*dir = filepath.Join(mainDir, name)
		if _, err := os.Stat(*dir); os.IsNotExist(err) {
			os.Mkdir(*dir, os.ModePerm)
		}
	}

	// 初始化Logger
	logFile, err := utils.InitializeLogger()
	if err != nil {
		log.Fatalf("[FeasOJ] Failed to initialize logger: %v", err)
	}
	defer utils.CloseLogger(logFile)

	// 初始化配置文件
	config.InitConfig()

	// 初始化数据库
	if utils.ConnectSql() == nil {
		return
	}
	utils.InitTable()

	// 初始化管理员账户
	if sql.SelectAdminUser(1) {
		log.Println("[FeasOJ] The administrator account already exists and will continue.")
	} else {
		sql.Register(utils.InitAdminAccount())
	}
	log.Println("[FeasOJ] MySQL initialization complete.")

	// 测试邮箱模块是否正常
	if !utils.TestSend(config.InitEmailConfig()) {
		log.Println("[FeasOJ] Email service startup failed, please check the configuration.")
		return
	} else {
		log.Println("[FeasOJ] Email service initialization complete.")
	}

	// 测试Redis连接
	if utils.ConnectRedis() == nil {
		log.Println("[FeasOJ] Redis connection failed, please check the configuration.")
		return
	} else {
		log.Println("[FeasOJ] Redis connection successful.")
	}

	// 构建沙盒镜像
	if judge.BuildImage() {
		log.Println("[FeasOJ] SandBox builds successfully.")
	} else {
		log.Println("[FeasOJ] SandBox builds fail, please make sure Docker is running and up to date!")
		return
	}

	// 启动服务器
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	router.LoadRouter(r)

	// 挂载头像文件夹
	r.StaticFS("/avatar", http.Dir(global.AvatarsDir))

	// 实时检测Redis JudgeTask中是否有任务
	rdb := utils.ConnectRedis()
	go judge.ProcessJudgeTasks(rdb)

	startServer := func(protocol, address, certFile, keyFile string) {
		for {
			var err error
			if protocol == "http" {
				err = r.Run(address)
			} else {
				err = r.RunTLS(address, certFile, keyFile)
			}
			if err != nil {
				log.Printf("[FeasOJ] Server start error: %v\n", err)
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}
	}
	// TODO: 根据需求调整
	// go startServer("https", "127.0.0.1:37881", "./certificate/fullchain.pem", "./certificate/privkey.key")
	go startServer("http", "127.0.0.1:37881", "", "")

	log.Println("[FeasOJ] Server activated.")

	// 监听终端输入
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if scanner.Text() == "quit" {
				log.Println("[FeasOJ] The server is being shut down....")
				os.Exit(0)
			}
		}
	}()

	// 等待中断信号关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	log.Println("[FeasOJ] Input 'quit' or Ctrl+C to stop the server.")
	<-quit

	// 关闭服务器前的清理工作
	log.Println("[FeasOJ] The server is shutting down...")
	utils.CloseLogger(logFile)
}
