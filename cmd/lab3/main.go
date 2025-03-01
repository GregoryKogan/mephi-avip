package main

import (
	"os"

	"github.com/GregoryKogan/mephi-av-processing/pkg/imgproc"
)

func main() {
	entries, _ := os.ReadDir("assets/noise/")
	for _, file := range entries {
		img, _ := imgproc.OpenPNG("assets/noise/" + file.Name())

		halfTone := imgproc.GetHalfTone(img)
		imgproc.SavePNG(halfTone, "output/lab3/original-"+file.Name())

		// 1. отфильтрованное монохромное (полутоновое) изображение;
		filtered := imgproc.MedianFilter3x3(halfTone,
			[3][3]int{
				{1, 0, 1},
				{0, 1, 0},
				{1, 0, 1},
			}, 3)

		imgproc.SavePNG(filtered, "output/lab3/filtered-"+file.Name())

		// 2. разностное изображение (монохромный xor или модуль разности для полутона).
		diff := imgproc.GetDifference(halfTone, filtered)
		imgproc.SavePNG(diff, "output/lab3/difference-"+file.Name())
	}
}
