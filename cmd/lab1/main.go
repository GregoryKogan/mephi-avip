package main

import (
	"github.com/GregoryKogan/mephi-av-processing/pkg/imgproc"
)

func main() {
	img, _ := imgproc.OpenPNG("assets/khachapuri.png")

	// 1. Цветовые модели

	// 1) Выделить компоненты R, G, B и сохранить как отдельные изображения.
	rImg, gImg, bImg := imgproc.GetColorComponents(img)
	imgproc.SavePNG(rImg, "output/lab1/1.1-red.png")
	imgproc.SavePNG(gImg, "output/lab1/1.1-green.png")
	imgproc.SavePNG(bImg, "output/lab1/1.1-blue.png")

	// 2) Привести изображение к цветовой модели HSI, сохранить яркостную
	// компоненту как отдельное изображение.
	lImg := imgproc.GetLightness(img)
	imgproc.SavePNG(lImg, "output/lab1/1.2-lightness.png")

	// 3) Инвертировать яркостную компоненту в исходном изображении, сохранить
	// производное изображение.
	invImg := imgproc.InvertColors(lImg)
	imgproc.SavePNG(invImg, "output/lab1/1.3-inverted.png")

	// 2. Передискретизация
	m, n := 3, 7

	// 1) Растяжение (интерполяция) изображения в M раз;
	interpolated := imgproc.Resize(img, imgproc.NewResizeRatio(m, 1))
	imgproc.SavePNG(interpolated, "output/lab1/2.1-interpolated.png")

	// 2) Сжатие (децимация) изображения в N раз;
	decimated := imgproc.Resize(img, imgproc.NewResizeRatio(1, n))
	imgproc.SavePNG(decimated, "output/lab1/2.2-decimated.png")

	// 3) Передискретизация изображения в K=M/N раз путём растяжения и
	// последующего сжатия (в два прохода);
	resized2steps := imgproc.Resize(interpolated, imgproc.NewResizeRatio(1, n))
	imgproc.SavePNG(resized2steps, "output/lab1/2.3-resized-2-steps.png")

	// 4) Передискретизация изображения в K раз за один проход.
	resized := imgproc.Resize(img, imgproc.NewResizeRatio(m, n))
	imgproc.SavePNG(resized, "output/lab1/2.4-resized.png")
}
