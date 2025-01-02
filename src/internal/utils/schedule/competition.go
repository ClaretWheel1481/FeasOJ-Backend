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

	now := time.Now()
	delay := time.Minute - time.Duration(now.Second())*time.Second

	// 调整至每分钟的0秒
	time.AfterFunc(delay, func() {
		scheduler.Every(1).Minute().Do(func() {
			log.Println("[FeasOJ] Competition Status Update:", time.Now(), "Status:", sql.UpdateCompetitionStatus())
		})

		// 启动任务调度
		scheduler.StartBlocking()
	})
}
