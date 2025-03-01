package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"src/internal/config"
	"src/internal/global"

	"fmt"

	"github.com/hashicorp/consul/api"
)

// PredictionResponse 定义检测服务返回的 JSON 结构
type PredictionResponse struct {
	Text       string `json:"text"`
	Prediction string `json:"prediction"`
}

// 获取Profanity Detector服务地址
func GetProfanityDetectorAddress() string {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = config.ConsulAddress
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Println("[FeasOJ] Error connecting to Consul:", err)
		return ""
	}

	services, _, err := consulClient.Catalog().Service("ProfanityDetector", "", nil)
	if err != nil {
		log.Println("[FeasOJ] Error querying Consul:", err)
		return ""
	}

	if len(services) == 0 {
		log.Println("[FeasOJ] ProfanityDetector service not found in Consul")
		return ""
	}

	service := services[0]
	global.ProfanityDetectorAddr = "http://" + service.ServiceAddress + ":" + fmt.Sprint(service.ServicePort)
	return "http://" + service.ServiceAddress + ":" + fmt.Sprint(service.ServicePort) + "/api/v1/text/predict"
}

// DetectText 检测文本是否包含敏感词汇
func DetectText(text string) bool {
	payload := map[string]string{"text": text}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return true
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/text/predict", global.ProfanityDetectorAddr), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return true
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return true
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return true
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return true
	}

	var result PredictionResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return true
	}

	// "normal" 表示合规，其他情况均视为违规
	if result.Prediction == "normal" {
		return false
	}
	return true
}

// ProfanityDetectorPing 检查服务是否可用
func ProfanityDetectorPing() bool {
	profanityDetectorAddress := GetProfanityDetectorAddress()
	if profanityDetectorAddress == "" {
		log.Println("[FeasOJ] Unable to get ProfanityDetector address from Consul")
		return false
	}

	resp, err := http.Get(profanityDetectorAddress)
	if err != nil {
		log.Println("[FeasOJ] Error requesting ProfanityDetector:", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
