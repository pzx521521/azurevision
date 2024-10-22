[![简体中文](https://img.shields.io/badge/lang-中文-red.svg)](README.zh-CN.md)

# Azure Vision Image Analysis FREE in golang
Note: This is not normal commercial API provided by [AzureVision](https://portal.vision.cognitive.azure.com/gallery/imageanalysis).
## CURL
```curl
curl 'https://portal.vision.cognitive.azure.com/api/demo/analyze?features=caption&language=en' \
  -F 'file=@/Users/parapeng/Documents/WallPaper/azure-vision-free/output.jpg'
```
the reseponse is a json
```json
{"modelVersion":"2023-10-01","captionResult":{"text":"a cartoon of a woman with blue hair","confidence":0.80237787961959839},"metadata":{"width":1202,"height":751}}
```
## Install:
```
go get github.com/pzx521521/azurevision
```


Example usage:
```go
package main

import (
    "fmt"
    "github.com/pzx521521/azurevision"
)

func main(){
	av := azurevision.NewAzureVision()
	anlyze, _ := av.Anlyze("input.png")
	fmt.Printf("%v\n", anlyze)
	// Output:map[captionResult:map[confidence:0.8023778796195984 text:a cartoon of a woman with blue hair] metadata:map[height:751 width:1202] modelVersion:2023-10-01]
}
```
Compressing the image
To handle network issues, the image is compressed by default, with a default quality of 10. If you want to check the compression result, you can use the following code:
```go
av := azurevision.NewAzureVision()
av.Quality = 20
av.TestCompress("input.png", "output.jpg")
```
If you don't want to use compression, you can set the quality to 0 (or less than or equal to 0):
```go
av := azurevision.NewAzureVision()
av.Quality = 0
```
Using other features
If you want to use other features (e.g., generating tags), the default feature is to get a caption. You can switch features using the following code:
```golang
av := azurevision.NewAzureVision()
av.Feature = "tags"
anlyze, _ := av.Anlyze("input.png")
fmt.Printf("%v\n", anlyze)
```
All features can be found here:
[Official Website](https://portal.vision.cognitive.azure.com/gallery/imageanalysis)