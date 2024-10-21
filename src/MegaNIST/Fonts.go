package main

import . "gopkg.in/gographics/imagick.v3/imagick"
import "regexp"

func isFontFile(fname string) bool {
	matched, _ := regexp.MatchString(".*\\.otf", fname)
	return matched
}

func GetFontFiles(dirname string) []*DrawingWand {
	inputc := FindFiles(dirname, false)
	hash := make([]*DrawingWand, 0)
	for input := range inputc {
		if isFontFile(input) {
			dw := NewDrawingWand()
			dw.SetFont(input)
			hash = append(hash, dw)
		}
	}
	return hash
}
func GetFonts(dirname string, rands chan int) chan *DrawingWand {
	retc := make(chan *DrawingWand, 1024)
	hash := GetFontFiles(dirname)
	go func() {
		defer close(retc)
		length := len(hash)
		for rand := range rands {
			retc <- hash[abs(rand) % length]
		}

	} ()
	return retc

}
