package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
	"gocv.io/x/gocv"
)

func main() {
	//1.1开启摄像头
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	//1.2创建电脑窗口，承载图片(容器)
	windom := gocv.NewWindow("dog_face")
	defer windom.Close()

	//1.3加载图片识别的分类
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	//根据那个文件，甄别捕获图片中是否由人脸
	if !classifier.Load("haarcascade_frontalface_default.xml") {
		fmt.Println("加载分类器失败")
		return
	}

	//2 读取每一帧图像(创建图像矩阵)
	img := gocv.NewMat()
	defer img.Close()
	//开启for循环，不间断的提取图片
	for {
		//通过摄像头捕获图片,捕获到的图片放置到img中
		if ok := webcam.Read(&img); !ok {
			fmt.Println("读取图片失败")
			return
		}
		if img.Empty() {
			continue
		}
		//2.1 读取图像后，判断此图像中是否由人脸
		//2.2 将人脸所在区域捕获出来，为了存储多张人脸，返回矩形
		rects := classifier.DetectMultiScale(img)
		fmt.Println("捕获人脸的数量:", len(rects))

		//3 初始化狗头图片和摄像机捕获背景图
		//3.1加载狗头png图片,并最终将图片转换成IMage格式
		file, _ := os.Open("dog.png") //将图片保存到内存
		srcPng, _ := png.Decode(file)

		//3.2 获取和窗口等大的img图片，并获取其区域
		backgroundImg, _ := img.ToImage()
		bounds := backgroundImg.Bounds() //获取图片的矩形的区域

		//3.3将background中内容绘制到自己创建的包含(红，绿，蓝，透明)的图片中
		address := image.NewRGBA(bounds) //依据矩形区域构件图，也就是复制的图片的位置
		draw.Draw(address, bounds, backgroundImg, image.ZP, draw.Over)

		//4.比较获取react中的面积最大的图片
		var maxIndex = 0
		for i := 0; i < len(rects)-1; i++ {
			//   rectAngle:=rects[i]
			//   rectAngle.Dx()*rectAngle.Dy()

			if rects[maxIndex].Dx()*rects[maxIndex].Dy() < rects[i+1].Dx()*rects[i+1].Dy() {
				//4.1 获取每一个矩形面积，进行大小比对
				//4.2 获取面积最大的矩形的索引位置
				maxIndex = i + 1
			}
		}
		if len(rects) > 0 {
			r := rects[maxIndex]
			//修改矩形起始点的位置和终点的位置
			startPoint := image.Pt(r.Min.X, r.Min.Y-80)
			endPoint := image.Pt(r.Max.X+50, r.Max.Y+20)

			//构建替换狗头的全新的矩形区域
			//扩大了矩形的区域，在扩大的区域进行狗头替换
			r = image.Rectangle{Min: startPoint, Max: endPoint}
			srcPng = imaging.Fill(srcPng, r.Size().X, r.Size().Y, imaging.Center, imaging.Lanczos)
			//替换图片
			draw.Draw(address, r, srcPng, image.ZP, draw.Over)
		}
		//将图片转换成矩阵
		mat := gocv.NewMat()
		mat, _ = gocv.ImageToMatRGB(address)
		if mat.Empty() {
			continue
		}
		windom.IMShow(mat)

		//按键盘的时候，跳出死循环
		if windom.WaitKey(1) >= 0 {
			break
		}
	}

}
