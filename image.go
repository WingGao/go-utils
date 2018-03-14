package utils

import (
	"image"
	"path/filepath"
	"strings"
	"image/png"
	"os"
	"image/jpeg"
)

func GetImage(fp string) (image.Image, error) {
	ext := strings.ToLower(filepath.Ext(fp))
	file, _ := os.Open("test.jpg")
	defer file.Close()
	if ext == ".png" {
		return png.Decode(file)
	}
	if ext == ".jpg" {
		return jpeg.Decode(file)
	}
	return nil, nil
}
