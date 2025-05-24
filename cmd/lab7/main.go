package main

import (
	"bufio"
	"fmt"
	"image"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/GregoryKogan/mephi-av-processing/pkg/imgproc"
)

type Features struct {
	letter                                                                 string
	relativeWeightQ1, relativeWeightQ2, relativeWeightQ3, relativeWeightQ4 float64
	relativeCenterOfMassX, relativeCenterOfMassY                           float64
	relativeInertiaX, relativeInertiaY                                     float64
}

func main() {
	os.Remove("output/lab7/guess.txt")
	file, _ := os.Create("output/lab7/guess.txt")
	out := bufio.NewWriter(file)
	defer out.Flush()

	saveRecognizedLetters()
	alphabet := GetFeatures("assets/osmanya/")
	recognized := GetFeatures("temp/")
	os.RemoveAll("temp")

	for lind, letter := range recognized {
		sort.Slice(alphabet, func(i, j int) bool {
			return letter.Similarity(alphabet[i]) > letter.Similarity(alphabet[j])
		})

		fmt.Printf("%c ", GetUnicode(alphabet[0].letter))

		fmt.Fprintf(out, "%d: ", lind+1)
		for _, g := range alphabet {
			fmt.Fprintf(out, "('%c', %.2f) ", GetUnicode(g.letter), letter.Similarity(g))
		}
		fmt.Fprintln(out)
	}
	fmt.Println()
}

func (f *Features) Similarity(other Features) float64 {
	d := f.dist(other)
	return 1 / (1 + d)
}

func (f *Features) dist(other Features) float64 {
	sum := 0.0
	sum += math.Pow(f.relativeWeightQ1-other.relativeWeightQ1, 2.0)
	sum += math.Pow(f.relativeWeightQ2-other.relativeWeightQ2, 2.0)
	sum += math.Pow(f.relativeWeightQ3-other.relativeWeightQ3, 2.0)
	sum += math.Pow(f.relativeWeightQ4-other.relativeWeightQ4, 2.0)
	sum += math.Pow(f.relativeCenterOfMassX-other.relativeCenterOfMassX, 2.0)
	sum += math.Pow(f.relativeCenterOfMassY-other.relativeCenterOfMassY, 2.0)
	sum += math.Pow(f.relativeInertiaX-other.relativeInertiaX, 2.0)
	sum += math.Pow(f.relativeInertiaY-other.relativeInertiaY, 2.0)
	return math.Sqrt(sum)
}

func GetFeatures(path string) []Features {
	entries, _ := os.ReadDir(path)
	alphabet := make([]Features, 0, len(entries))
	for _, e := range entries {
		img, _ := imgproc.OpenPNG(path + e.Name())
		ht := imgproc.GetHalfTone(imgproc.InvertColors(img))
		f := Features{letter: strings.Split(e.Name(), ".")[0]}
		f.relativeWeightQ1, f.relativeWeightQ2, f.relativeWeightQ3, f.relativeWeightQ4 = imgproc.QuartersRelativeBlackWeight(ht)
		f.relativeCenterOfMassX, f.relativeCenterOfMassY = imgproc.RelativeCenterOfMass(ht)
		f.relativeInertiaX, f.relativeInertiaY = imgproc.RelativeInertia(ht)
		alphabet = append(alphabet, f)
	}

	return alphabet
}

func saveRecognizedLetters() {
	raw, _ := imgproc.OpenPNG("assets/iloveyou3.png")
	img := imgproc.GetThresholding(imgproc.GetHalfTone(imgproc.InvertColors(raw)), 100)
	rectangles := imgproc.SegmentLetters(img)
	originalColors := imgproc.InvertColors(img)
	for i, r := range rectangles {
		name := "letter-" + fmt.Sprint(i+1)
		letterImg := originalColors.(interface {
			SubImage(image.Rectangle) image.Image
		}).SubImage(r)
		imgproc.SavePNG(letterImg, "temp/"+name+".png")
	}
}

func GetUnicode(name string) rune {
	dict := map[string]rune{
		"A":     '𐒖',
		"AA":    '𐒛',
		"ALEF":  '𐒀',
		"BA":    '𐒁',
		"CAYN":  '𐒋',
		"DEEL":  '𐒆',
		"DHA":   '𐒊',
		"E":     '𐒗',
		"EE":    '𐒜',
		"FA":    '𐒍',
		"GA":    '𐒌',
		"HA":    '𐒔',
		"I":     '𐒘',
		"JA":    '𐒃',
		"KAAF":  '𐒏',
		"KHA":   '𐒅',
		"LAAN":  '𐒐',
		"MIIN":  '𐒑',
		"NUUN":  '𐒒',
		"O":     '𐒙',
		"OO":    '𐒝',
		"QAAF":  '𐒎',
		"RA":    '𐒇',
		"SA":    '𐒈',
		"SHIIN": '𐒉',
		"TA":    '𐒂',
		"U":     '𐒚',
		"WAW":   '𐒓',
		"XA":    '𐒄',
		"YA":    '𐒕',
	}
	return dict[name]
}

// 𐒖 𐒛 𐒀 𐒁 𐒋 𐒆 𐒊 𐒗 𐒜 𐒍 𐒌 𐒔 𐒘 𐒃 𐒏 𐒅 𐒐 𐒑 𐒒 𐒙 𐒝 𐒎 𐒇 𐒈 𐒉 𐒂 𐒚 𐒓 𐒄
// 𐒆𐒖𐒕𐒐 𐒄𐒆𐒜𐒖𐒕𐒜
