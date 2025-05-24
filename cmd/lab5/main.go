package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/GregoryKogan/mephi-av-processing/pkg/imgproc"
	"github.com/vicanso/go-charts/v2"
)

type features struct {
	letter                                                                 string
	weightQ1, weightQ2, weightQ3, weightQ4                                 float64
	relativeWeightQ1, relativeWeightQ2, relativeWeightQ3, relativeWeightQ4 float64
	centerOfMassX, centerOfMassY                                           float64
	relativeCenterOfMassX, relativeCenterOfMassY                           float64
	inertiaX, inertiaY                                                     float64
	relativeInertiaX, relativeInertiaY                                     float64
}

func main() {
	entries, _ := os.ReadDir("assets/osmanya/")
	writer, file, _ := createCSVWriter("output/lab5/features.csv")
	defer file.Close()
	writeCSVRecord(writer, []string{
		"letter",
		"weightQ1", "weightQ2", "weightQ3", "weightQ4",
		"relativeWeightQ1", "relativeWeightQ2", "relativeWeightQ3", "relativeWeightQ4",
		"centerOfMassX", "centerOfMassY",
		"relativeCenterOfMassX", "relativeCenterOfMassY",
		"inertiaX", "inertiaY",
		"relativeInertiaX", "relativeInertiaY",
	})

	for _, e := range entries {
		img, _ := imgproc.OpenPNG("assets/osmanya/" + e.Name())
		ht := imgproc.GetHalfTone(imgproc.InvertColors(img))
		f := features{letter: strings.Split(e.Name(), ".")[0]}
		f.weightQ1, f.weightQ2, f.weightQ3, f.weightQ4 = imgproc.QuartersBlackWeight(ht)
		f.relativeWeightQ1, f.relativeWeightQ2, f.relativeWeightQ3, f.relativeWeightQ4 = imgproc.QuartersRelativeBlackWeight(ht)
		f.centerOfMassX, f.centerOfMassY = imgproc.CenterOfMass(ht)
		f.relativeCenterOfMassX, f.relativeCenterOfMassY = imgproc.RelativeCenterOfMass(ht)
		f.inertiaX, f.inertiaY = imgproc.Inertia(ht)
		f.relativeInertiaX, f.relativeInertiaY = imgproc.RelativeInertia(ht)
		f.write(writer)

		horProfile := imgproc.HorizontalProfile(ht)
		verProfile := imgproc.VerticalProfile(ht)

		saveHorizontalProfile(horProfile, f.letter)
		saveVerticalProfile(verProfile, f.letter)
	}
	writer.Flush()
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
		charts.TitleTextOptionFunc("Horizontal profile, letter "+name),
		charts.YAxisDataOptionFunc(labels),
	)
	buf, _ := p.Bytes()
	os.WriteFile("output/lab5/hor"+name+".png", buf, 0600)
}

func saveVerticalProfile(profile []float64, name string) {
	labels := make([]string, len(profile))
	for i := range profile {
		labels[i] = fmt.Sprint(i)
	}
	p, _ := charts.BarRender(
		[][]float64{profile},
		charts.TitleTextOptionFunc("Vertical profile, letter "+name),
		charts.XAxisDataOptionFunc(labels),
	)
	buf, _ := p.Bytes()
	os.WriteFile("output/lab5/ver"+name+".png", buf, 0600)
}

func createCSVWriter(filename string) (*csv.Writer, *os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}
	writer := csv.NewWriter(f)
	return writer, f, nil
}

func writeCSVRecord(writer *csv.Writer, record []string) {
	err := writer.Write(record)
	if err != nil {
		fmt.Println("Error writing record to CSV:", err)
	}
}

func (f *features) write(writer *csv.Writer) {
	writeCSVRecord(writer, []string{
		f.letter,
		fmt.Sprint(f.weightQ1),
		fmt.Sprint(f.weightQ2),
		fmt.Sprint(f.weightQ3),
		fmt.Sprint(f.weightQ4),
		fmt.Sprint(f.relativeWeightQ1),
		fmt.Sprint(f.relativeWeightQ2),
		fmt.Sprint(f.relativeWeightQ3),
		fmt.Sprint(f.relativeWeightQ4),
		fmt.Sprint(f.centerOfMassX),
		fmt.Sprint(f.centerOfMassY),
		fmt.Sprint(f.relativeCenterOfMassX),
		fmt.Sprint(f.relativeCenterOfMassY),
		fmt.Sprint(f.inertiaX),
		fmt.Sprint(f.inertiaY),
		fmt.Sprint(f.relativeInertiaX),
		fmt.Sprint(f.relativeInertiaY),
	})
}
