# 使用 Go 语言调用 Azure Vision 图像分析免费服务
Note: 这不是 [AzureVision](https://portal.vision.cognitive.azure.com/gallery/imageanalysis). 提供的普通商业 API。 
## 起因
手机上太多图片了,想给他们重命名一下  
不想传到电脑,因为比较大,还要传来传去,想在手机上直接跑,所以选择golang直接编译为二进制包  
最开始的方案是开源的yolov5,测试之后发现他只能识别一张图片中的有多少个物体即tags功能,而且模型很大  
然后看到了yolo-mini的模型,只有88个分类,而且效果并不好.  
比如很多照片上全部是人,他就只能识别出很多person.... 我知道这些有什么用 = =  
本来想用离线的,方案是 CLIP + GPT2,但是模型太大了,完全不适合手机运行,想想在手机上装 PyTorch 或 TensorFlow....
最终放弃了离线

## CURL
```curl
curl 'https://portal.vision.cognitive.azure.com/api/demo/analyze?features=caption&language=en' \
  -F 'file=@/Users/parapeng/Documents/WallPaper/azure-vision-free/output.jpg'
```
返回
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
所有的功能(Feature)你可以在这里找到,就一个tag的区别:
[官网](https://portal.vision.cognitive.azure.com/gallery/imageanalysis)

### 题外话
在线方案有google的和azure的两种，阿里和百度居然没有.... 应该是看不上吧   
google要api,azure商业版也要,但是这里用的是他的demo接口  
azure Vision的价格如下:
5,000 free transactions per month  
20 transactions per minute  
超出: 0-1M transactions - $1 per 1,000 transactions  
但是要先绑卡(坑啊之前在amazon没注意超了扣了好多钱)
azure 并没有做任何加密处理之类的,点名批评国内的云服务商,很多都是没有demo的
