package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"src/internal/config"

	"github.com/nfnt/resize"
)

// 压缩图像为256*256
func CompressImage(inputPath, outputPath string) error {
	// 打开图像文件
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}
	if format != "png" {
		return errors.New("unsupported image format")
	}

	newImage := resize.Resize(256, 256, img, resize.Lanczos3)

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	err = png.Encode(out, newImage)
	return err
}

// ImageGuardPing 检查服务是否可用
func ImageGuardPing() bool {
	resp, err := http.Get(config.ImageGuardAddress)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// PredictImage 判断图片是否违规
func PredictImage(imagePath string) bool {
	// 打开图片文件
	file, err := os.Open(imagePath)
	if err != nil {
		return false
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", imagePath)
	if err != nil {
		return false
	}
	if _, err = io.Copy(part, file); err != nil {
		return false
	}
	if err = writer.Close(); err != nil {
		return false
	}

	req, err := http.NewRequest("POST", config.ImageGuardAddress, &buf)
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Prediction string `json:"prediction"`
		Error      string `json:"error,omitempty"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}

	// 返回 true 表示预测结果为 "neutral"
	return result.Prediction == "neutral"
}
