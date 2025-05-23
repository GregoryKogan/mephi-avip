package main

import "github.com/GregoryKogan/mephi-av-processing/pkg/imgproc"

func main() {
	files := []string{
		"cartoon.png",
		// "fingerprint.png",
		// "khachapuri.png",
		// "map.png",
		"page.png",
		// "xray.png",
	}

	for _, file := range files {
		img, _ := imgproc.OpenPNG("assets/" + file)

		halfTone := imgproc.GetHalfTone(img)
		gx, gy, g := imgproc.GetScharrEdges(halfTone)
		bin := imgproc.GetThresholding(g, 10)
		imgproc.SavePNG(gx, "output/lab4/Gx-"+file)
		imgproc.SavePNG(gy, "output/lab4/Gy-"+file)
		imgproc.SavePNG(g, "output/lab4/G-"+file)
		imgproc.SavePNG(bin, "output/lab4/Bin-"+file)
	}
}
