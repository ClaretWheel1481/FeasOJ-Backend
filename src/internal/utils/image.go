package utils

import (
	"errors"
	"image"
	"image/png"
	"os"

	"github.com/nfnt/resize"
)

// 压缩图像为128*128
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

	newImage := resize.Resize(128, 128, img, resize.Lanczos3)

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	err = png.Encode(out, newImage)
	return err
}

// TODO: 接入本地模型判断图片是否违规
