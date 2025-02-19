package judge

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"src/config"
	gincontext "src/internal/gin"
	"src/internal/global"
	"src/internal/utils"
	"src/internal/utils/sql"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type Task struct {
	UID  int
	PID  int
	Name string
}

// ProcessJudgeTasks 函数用于处理判题任务
func ProcessJudgeTasks(rdb *redis.Client) {
	// 创建一个任务通道
	taskChan := make(chan Task)
	// 创建一个等待组
	var wg sync.WaitGroup

	// 创建多个worker协程
	for i := 0; i < config.MaxSandbox; i++ {
		wg.Add(1)
		go worker(taskChan, &wg)
	}

	// 无限循环，从Redis任务队列中取出任务
	for {
		// 从Redis任务队列中取出一个任务
		task, err := rdb.LPop("judgeTask").Result()
		if err == redis.Nil {
			// 如果队列为空，等待一段时间后再检查
			time.Sleep(2 * time.Second)
			continue
		} else if err != nil {
			log.Panic("[FeasOJ] Redis connect error: ", err)
		}

		// 将任务分割成用户ID和题目ID
		parts := strings.Split(task, "_")
		uid := parts[0]
		pid := strings.Split(parts[1], ".")[0]
		// 将用户ID和题目ID转换为整数
		uidInt, err := strconv.Atoi(uid)
		if err != nil {
			log.Panic(err)
		}
		pidInt, err := strconv.Atoi(pid)
		if err != nil {
			log.Panic(err)
		}

		// 将任务发送到任务通道
		taskChan <- Task{UID: uidInt, PID: pidInt, Name: task}
	}
}

// worker 使用容器池执行任务
func worker(taskChan chan Task, wg *sync.WaitGroup) {
	// 使用 defer 关键字，在函数结束时调用 wg.Done()，表示任务完成
	defer wg.Done()
	// 从任务通道中获取任务
	for task := range taskChan {
		// 从容器池中获取一个空闲容器
		containerID := AcquireContainer()
		// 将容器ID存储到全局变量中
		global.ContainerIDs.Store(task.Name, containerID)
		// 执行编译与运行
		result := CompileAndRun(task.Name, containerID)
		// 更新判题状态
		sql.ModifyJudgeStatus(task.UID, task.PID, result)

		// 发送 SSE 通知
		if client, ok := gincontext.Clients[fmt.Sprint(task.UID)]; ok {
			// 获取客户端的语言
			lang := client.Lang
			// 将语言转换为语言标签
			tag := language.Make(lang)
			// 初始化国际化
			langBundle := utils.InitI18n()
			// 创建本地化器
			localizer := i18n.NewLocalizer(langBundle, tag.String())
			// 本地化消息
			message, _ := localizer.Localize(&i18n.LocalizeConfig{
				MessageID: "problem_completed",
				TemplateData: map[string]interface{}{
					"PID": task.PID,
				},
			})
			// 将消息发送到客户端的消息通道中
			client.MessageChan <- message
		}

		// 将容器归还到池中（内部会先重置环境）
		ReleaseContainer(containerID)
	}
}
