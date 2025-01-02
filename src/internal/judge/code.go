package judge

import (
	"fmt"
	"log"
	gincontext "src/internal/gin"
	"src/internal/utils/sql"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

// 实时处理Redis任务队列中的任务
func ProcessJudgeTasks(rdb *redis.Client) {
	for {
		// 从Redis任务队列中取出一个任务
		task, err := rdb.LPop("judgeTask").Result()
		if err == redis.Nil {
			// 如果队列为空，等待一段时间后再检查
			time.Sleep(2 * time.Second)
			continue
		} else if err != nil {
			log.Panic(err)
			continue
		}
		// TODO: 实现当任务较多时，Sandbox的并发处理
		// 执行任务
		StartContainer()
		str := CompileAndRun(task)

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

		// 将结果从数据库中修改
		sql.ModifyJudgeStatus(uidInt, pidInt, str)

		// 发送SSE通知
		if ch, ok := gincontext.Clients[uid]; ok {
			ch <- fmt.Sprintf("Problem %d 运行完毕。", pidInt)
		}
	}
}
