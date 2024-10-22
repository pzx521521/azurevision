# Azure Vision Image Analysis FREE in golang
Note: This is not normal commercial API provided by [AzureVision](https://portal.vision.cognitive.azure.com/gallery/imageanalysis).
## CURL
```curl
curl 'https://portal.vision.cognitive.azure.com/api/demo/analyze?features=caption&language=en' \
  -F 'file=@/Users/parapeng/Documents/WallPaper/azure-vision-free/output.jpg'
```
the reseponse is a html, result in
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
	v := NewAzureVision()
	anlyze, _ := v.Anlyze("input.png")
	fmt.Printf("%v\n", anlyze)
	// Output:map[captionResult:map[confidence:0.8023778796195984 text:a cartoon of a woman with blue hair] metadata:map[height:751 width:1202] modelVersion:2023-10-01]
}
```
考虑到网络问题,对图片进行了压缩,默认质量为10。如果你想看下压缩的效果
```go
v := NewAzureVision()
v.Quality = 20
v.TestCompress("input.png", "output.jpg")
```
如果你不想使用压缩,可以设置质量为0(or<=0)
```go
v := NewAzureVision()
v.Quality = 0
```
如果你想使用其他的功能(如生成标签),默认功能是获取Caption，可以使用如下
```golang
v := NewAzureVision()
v.Feature = "tags"
anlyze, _ := v.Anlyze("input.png")
fmt.Printf("%v\n", anlyze)
```
All features can be found here:
[官网](https://portal.vision.cognitive.azure.com/gallery/imageanalysis)