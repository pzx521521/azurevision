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
}

func NewAzureVision() *AzureVision {
	return &AzureVision{Feature: FEATURES[0], Quality: 10}
}

func (v *AzureVision) TestCompress(inputPath, outputPath string) {
	buffer, err := compressImage(inputPath, v.Quality)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	os.WriteFile(outputPath, buffer.Bytes(), os.ModePerm)
}

func (v *AzureVision) Anlyze(inputPath string) (ret map[string]interface{}, err error) {
	var imageData *bytes.Buffer
	// Compress image
	if v.Quality > 0 && v.Quality <= 100 {
		imageData, err = compressImage("input.png", 10)
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

func compressImage(inputPath string, quality int) (*bytes.Buffer, error) {
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
	// Save the image to bytes.Buffer with compression
	imgBytes := new(bytes.Buffer)
	err = jpeg.Encode(imgBytes, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}
	fmt.Printf("Compressed image size: %.2f KB\n", float64(imgBytes.Len())/1024)
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
