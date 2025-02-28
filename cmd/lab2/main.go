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

		halfTone := imgproc.GetLightness(img)
		imgproc.SavePNG(halfTone, "output/lab2/1-"+file)

		bin := imgproc.GetNiblackThresholding(halfTone, 15, 0.2)
		imgproc.SavePNG(bin, "output/lab2/2-"+file)
	}
}
