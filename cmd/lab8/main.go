package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/GregoryKogan/mephi-av-processing/pkg/imgproc"
	"github.com/vicanso/go-charts/v2"
)

func main() {
	files := []string{
		"texture-1.png",
		"texture-2.png",
	}

	for _, file := range files {
		img, _ := imgproc.OpenPNG("assets/" + file)
		name := strings.Split(file, ".")[0]

		halfTone := imgproc.GetHalfTone(img)
		imgproc.SavePNG(halfTone, "output/lab8/halftone-"+file)
		saveHistogram(imgproc.BrightnessHistogram(halfTone), "halftone-"+name)

		contrast := imgproc.LogLightnessAdjust(img)
		lightness := imgproc.GetLightness(contrast)
		imgproc.SavePNG(lightness, "output/lab8/lightness-"+file)
		saveHistogram(imgproc.BrightnessHistogram(lightness), "lightness-"+name)

		hogHt := imgproc.HOG(halfTone)
		saveHistogram(hogHt, "hog-halftone-"+name)

		hogLt := imgproc.HOG(lightness)
		saveHistogram(hogLt, "hog-lightness-"+name)

		saveHistogram(imgproc.HNorm(hogHt), "hog-halftone-norm-"+name)
		saveHistogram(imgproc.HNorm(hogLt), "hog-lightness-norm-"+name)
	}
}

func saveHistogram(bins []float64, name string) {
	labels := make([]string, len(bins))
	for i := range bins {
		labels[i] = fmt.Sprint(i)
	}
	p, _ := charts.BarRender(
		[][]float64{bins},
		charts.TitleTextOptionFunc("Histogram, "+name),
		charts.XAxisDataOptionFunc(labels),
	)
	buf, _ := p.Bytes()
	os.WriteFile("output/lab8/histogram-"+name+".png", buf, 0600)
}
