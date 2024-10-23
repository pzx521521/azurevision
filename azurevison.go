package azurevision

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // for PNG decoding
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const analyzeURL = "https://portal.vision.cognitive.azure.com/api/demo/analyze"

var FEATURES = []string{"caption", "tags", "denseCaptions"}

type AzureVision struct {
	Feature string
	Quality int
	Width   int
}

func NewAzureVision() *AzureVision {
	return &AzureVision{Feature: FEATURES[0], Quality: 50, Width: 500}
}

func (v *AzureVision) TestCompress(inputPath, outputPath string) {
	buffer, err := compressImage(inputPath, v.Quality, v.Width)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("Compressed image size: %.2f KB\n", float64(buffer.Len())/1024)
	os.WriteFile(outputPath, buffer.Bytes(), os.ModePerm)
}

func (v *AzureVision) Anlyze(inputPath string) (ret map[string]interface{}, err error) {
	var imageData *bytes.Buffer
	// Compress image
	if v.Quality > 0 && v.Quality <= 100 {
		imageData, err = compressImage(inputPath, v.Quality, v.Width)
		if err != nil {
			return ret, err
		}
	} else {
		file, err := os.Open(inputPath)
		if err != nil {
			return ret, err
		}
		defer file.Close()
		buffer := new(bytes.Buffer)
		_, err = io.Copy(buffer, file)
		if err != nil {
			return ret, err
		}
		imageData = buffer
	}
	// Upload compressed image and get info
	ret, err = analyze(imageData, v.Feature)
	if err != nil {
		fmt.Println("Error analyzing image:", err)
	}
	return ret, nil
}

func compressImage(inputPath string, quality int, width int) (*bytes.Buffer, error) {
	// Open the image
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	if width > 0 && width < img.Bounds().Dx() {
		img = resizeWithAspectRatio(img, width, 0)
	}
	// Save the image to bytes.Buffer with compression
	imgBytes := new(bytes.Buffer)
	err = jpeg.Encode(imgBytes, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}
	return imgBytes, nil
}

func analyze(imageData *bytes.Buffer, features string) (ret map[string]interface{}, err error) {
	url := fmt.Sprintf("%s?features=%s&language=en", analyzeURL, features)
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	// 将图像数据作为普通的字段直接添加到 multipart 数据中
	part, err := writer.CreateFormFile("file", "output.jpg")
	if err != nil {
		return nil, err
	}
	// 将 imageData 写入该部分
	io.Copy(part, imageData)
	// 结束 multipart 写入
	writer.Close()
	// Prepare the file upload
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return ret, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ret, fmt.Errorf("Failed to get response: %v", resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

// resizeWithAspectRatio 根据宽高比等比例缩放图像
func resizeWithAspectRatio(img image.Image, newWidth, newHeight int) *image.RGBA {
	// 获取原始图像的宽度和高度
	srcBounds := img.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	// 如果宽度或高度为 0，则按比例计算相应的值
	if newWidth == 0 {
		newWidth = (newHeight * srcWidth) / srcHeight
	}
	if newHeight == 0 {
		newHeight = (newWidth * srcHeight) / srcWidth
	}

	// 创建新的目标图像
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// 缩放图像，使用最近邻插值法
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// 计算原图中对应的像素
			srcX := x * srcWidth / newWidth
			srcY := y * srcHeight / newHeight
			// 取原图中的像素并设置到目标图像中
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}

	return dst
}
