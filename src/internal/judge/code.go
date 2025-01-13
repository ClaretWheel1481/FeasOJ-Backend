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

func ProcessJudgeTasks(rdb *redis.Client) {
	taskChan := make(chan Task)
	var wg sync.WaitGroup

	for i := 0; i < config.MaxWorkers; i++ {
		wg.Add(1)
		go worker(taskChan, &wg)
	}

	for {
		// 从Redis任务队列中取出一个任务
		task, err := rdb.LPop("judgeTask").Result()
		if err == redis.Nil {
			// 如果队列为空，等待一段时间后再检查
			time.Sleep(2 * time.Second)
			continue
		} else if err != nil {
			log.Panic(err)
		}

		parts := strings.Split(task, "_")
		uid := parts[0]
		pid := strings.Split(parts[1], ".")[0]
		uidInt, err := strconv.Atoi(uid)
		if err != nil {
			log.Panic(err)
		}
		pidInt, err := strconv.Atoi(pid)
		if err != nil {
			log.Panic(err)
		}

		taskChan <- Task{UID: uidInt, PID: pidInt, Name: task}
	}
}

// 并发处理
func worker(taskChan chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChan {
		containerID, err := StartContainer()
		if err != nil {
			log.Println(err)
			continue
		}
		// 存储任务的容器ID
		global.ContainerIDs.Store(task.Name, containerID)
		str := CompileAndRun(task.Name, containerID)

		// 终止任务对应的容器
		TerminateContainer(containerID)

		// 将结果从数据库中修改
		sql.ModifyJudgeStatus(task.UID, task.PID, str)

		// 发送SSE通知
		if client, ok := gincontext.Clients[fmt.Sprint(task.UID)]; ok {
			lang := client.Lang
			tag := language.Make(lang)
			langBundle := utils.InitI18n()
			localizer := i18n.NewLocalizer(langBundle, tag.String())
			message, _ := localizer.Localize(&i18n.LocalizeConfig{
				MessageID: "problem_completed",
				TemplateData: map[string]interface{}{
					"PID": task.PID,
				},
			})
			client.MessageChan <- message
		}
	}
}
