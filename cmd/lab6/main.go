package main

import (
	"fmt"
	"image"
	"os"

	"github.com/GregoryKogan/mephi-av-processing/pkg/imgproc"
	"github.com/vicanso/go-charts/v2"
)

func main() {
	raw, _ := imgproc.OpenPNG("assets/iloveyou.png")
	img := imgproc.GetHalfTone(imgproc.InvertColors(raw))

	horProfile := imgproc.HorizontalProfile(img)
	verProfile := imgproc.VerticalProfile(img)

	saveHorizontalProfile(horProfile, "line")
	saveVerticalProfile(verProfile, "line")

	rectangles := imgproc.SegmentLetters(img)
	segmented := imgproc.DrawRectangles(img, rectangles)
	imgproc.SavePNG(segmented, "output/lab6/segmented.png")

	for i, r := range rectangles {
		letterImg := raw.(interface {
			SubImage(image.Rectangle) image.Image
		}).SubImage(r)
		name := "letter-" + fmt.Sprint(i+1)
		imgproc.SavePNG(letterImg, "output/lab6/"+name+".png")

		letterImgPrep := img.SubImage(r).(*image.Gray)
		horProfile := imgproc.HorizontalProfile(letterImgPrep)
		verProfile := imgproc.VerticalProfile(letterImgPrep)
		saveHorizontalProfile(horProfile, name)
		saveVerticalProfile(verProfile, name)
	}
}

func saveHorizontalProfile(profile []float64, name string) {
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
		charts.TitleTextOptionFunc("Horizontal profile, "+name),
		charts.YAxisDataOptionFunc(labels),
	)
	buf, _ := p.Bytes()
	os.WriteFile("output/lab6/horizontal-"+name+".png", buf, 0600)
}

func saveVerticalProfile(profile []float64, name string) {
	labels := make([]string, len(profile))
	for i := range profile {
		labels[i] = fmt.Sprint(i)
	}
	p, _ := charts.BarRender(
		[][]float64{profile},
		charts.TitleTextOptionFunc("Vertical profile, "+name),
		charts.XAxisDataOptionFunc(labels),
	)
	buf, _ := p.Bytes()
	os.WriteFile("output/lab6/vertical-"+name+".png", buf, 0600)
}
