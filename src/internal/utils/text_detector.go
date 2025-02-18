package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"src/config"
)

// PredictionResponse 定义了检测服务返回的 JSON 结构
type PredictionResponse struct {
	Text       string `json:"text"`
	Prediction string `json:"prediction"`
}

// DetectText 检测文字是否违规
// 返回 true 表示合规（预测结果为 "normal"），返回 true 表示违规或检测过程中出错
func DetectText(text string) bool {
	payload := map[string]string{"text": text}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return true
	}

	req, err := http.NewRequest("POST", config.ProfanityDetectorAddress, bytes.NewBuffer(payloadBytes))
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

	bodyBytes, err := ioutil.ReadAll(resp.Body)
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
