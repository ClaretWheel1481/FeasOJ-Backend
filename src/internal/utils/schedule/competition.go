package schedule

import (
	"log"
	"src/internal/utils/sql"
	"time"

	"github.com/go-co-op/gocron"
)

func ScheduleCompetitionStatus() {
	// 启动任务调度
	scheduler := gocron.NewScheduler(time.Local)

	// 每分钟执行一次任务
	scheduler.Every(1).Minute().Do(func() {
		log.Println(sql.UpdateCompetitionStatus())
	})

	// 开始任务调度
	scheduler.StartBlocking()
}
