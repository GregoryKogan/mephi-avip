package main

import (
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/GregoryKogan/mephi-av-processing/pkg/imgproc"
)

func main() {
	entries, err := os.ReadDir("assets/osmanya/")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		img, _ := imgproc.OpenPNG("assets/osmanya/" + e.Name())
		ht := imgproc.FillBackground(img, color.White)
		newName := strings.Split(strings.Split(e.Name(), "_")[2], ".")[0] + ".png"
		imgproc.SavePNG(ht, "output/lab5/"+newName)
	}
}
