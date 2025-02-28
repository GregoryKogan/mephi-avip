package main

import (
	"github.com/GregoryKogan/mephi-av-processing/pkg/imgproc"
)

func main() {
	files := []string{
		"cartoon.png",
		"fingerprint.png",
		"khachapuri.png",
		"map.png",
		"page.png",
		"xray.png",
	}

	for _, file := range files {
		img, _ := imgproc.OpenPNG("assets/" + file)

		// 1. Приведение полноцветного изображения к полутоновому.
		halfTone := imgproc.GetLightness(img)
		imgproc.SavePNG(halfTone, "output/lab2/1-"+file)

		// 2. Приведение полутонового изображения к монохромному
		// методом пороговой обработки
		bin := imgproc.GetNiblackThresholding(halfTone, 15, 0.2)
		imgproc.SavePNG(bin, "output/lab2/2-"+file)
	}
}
