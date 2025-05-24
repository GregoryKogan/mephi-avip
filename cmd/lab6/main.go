package main

import (
	"fmt"
	"os"

	"github.com/GregoryKogan/mephi-av-processing/pkg/imgproc"
	"github.com/vicanso/go-charts/v2"
)

func main() {
	raw, _ := imgproc.OpenPNG("assets/iloveyou.png")
	img := imgproc.GetHalfTone(imgproc.InvertColors(raw))

	horProfile := imgproc.HorizontalProfile(img)
	verProfile := imgproc.VerticalProfile(img)

	saveHorizontalProfile(horProfile)
	saveVerticalProfile(verProfile)

	rectangles := imgproc.SegmentLetters(img)
	segmented := imgproc.DrawRectangles(img, rectangles)
	imgproc.SavePNG(segmented, "output/lab6/segmented.png")
}

func saveHorizontalProfile(profile []float64) {
	labels := make([]string, len(profile))
	for i := range profile {
		if i%20 == 0 {
			labels[i] = fmt.Sprint(i) + "-"
		} else {
			labels[i] = "|"
		}
	}
	p, _ := charts.HorizontalBarRender(
		[][]float64{profile},
		charts.TitleTextOptionFunc("Horizontal profile"),
		charts.YAxisDataOptionFunc(labels),
	)
	buf, _ := p.Bytes()
	os.WriteFile("output/lab6/horizontal.png", buf, 0600)
}

func saveVerticalProfile(profile []float64) {
	labels := make([]string, len(profile))
	for i := range profile {
		labels[i] = fmt.Sprint(i)
	}
	p, _ := charts.BarRender(
		[][]float64{profile},
		charts.TitleTextOptionFunc("Vertical profile"),
		charts.XAxisDataOptionFunc(labels),
	)
	buf, _ := p.Bytes()
	os.WriteFile("output/lab6/vertical.png", buf, 0600)
}
