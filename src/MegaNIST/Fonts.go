package main

import "gopkg.in/gographics/imagick.v3/imagick"
import "regexp"

func isFontFile(fname string) bool {
	matched, _ := regexp.MatchString(".*\\.otf", fname)
	return matched
}

func GetFonts(dirname string) chan DrawingWand {
	retc := make(chan DrawingWand, 1024)
	go func() {
		defer close(retc)
		inputc := FindFiles(dirname, false)
		for input := range inputc {
			dw := NewDrawingWand()
			if isFontFile(input) {
				dw.SetFont(input)
				retc <- dw
			}
		}
	} ()
	return retc

}
