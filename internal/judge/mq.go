package judge

import (
	"encoding/json"
	"fmt"
	"log"
	gincontext "src/internal/gin"
	"src/internal/global"
	"src/internal/utils"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// 发送评测结果到消息队列
func ConsumeJudgeResults() {
	conn, ch, err := utils.ConnectRabbitMQ()
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()
	defer ch.Close()

	// 声明结果队列
	q, err := ch.QueueDeclare(
		"judgeResults", // 队列名称
		true,           // 持久化
		false,          // 自动删除
		false,          // 排他性
		false,          // 不等待
		nil,            // 参数
	)
	if err != nil {
		log.Fatal("Failed to declare a queue:", err)
	}

	msgs, err := ch.Consume(
		q.Name, // 队列名称
		"",     // 消费者标签
		true,   // 自动应答
		false,  // 排他性
		false,  // 不等待
		false,  // 参数
		nil,
	)
	if err != nil {
		log.Fatal("Failed to register consumer:", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var result global.JudgeResultMessage
			if err := json.Unmarshal(d.Body, &result); err != nil {
				log.Printf("Error decoding result: %v", err)
				continue
			}
			// 发送 SSE 通知
			if client, ok := gincontext.Clients[fmt.Sprint(result.UserID)]; ok {
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
					TemplateData: map[string]any{
						"PID": result.ProblemID,
					},
				})
				// 将消息发送到客户端的消息通道中
				client.MessageChan <- message
			}

			log.Printf("Processed result for user %d", result.UserID)
		}
	}()

	<-forever
}
