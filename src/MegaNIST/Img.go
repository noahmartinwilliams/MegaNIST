package main

import . "gopkg.in/gographics/imagick.v3/imagick"
import "sync"
import "strconv"

func abs(input int) int {
	if input >= 0 {
		return input
	} else {
		return -input
	}
}
type Img struct {
	Image *MagickWand
	Number int
}

type Merged struct {
	Coordinate Coord
	Angle float64
	Fonts *DrawingWand
	Num int
}

func GetImages(coords chan Coord, angles chan float64, fonts chan *DrawingWand, numc chan int) chan Img {
	retc := make(chan Img, 1024)
	inputc := make(chan Merged, 1024)
	go func() {
		defer close(inputc)
		for {
			coord, ok := <-coords
			if !ok {
				return
			}

			angle, ok := <- angles
			if !ok {
				return
			}

			font, ok := <- fonts
			if !ok {
				return
			}

			num, ok := <- numc
			if !ok {
				return
			}
			inputc <- Merged{Coordinate:coord, Angle:angle, Fonts:font, Num:abs(num % 10)}
		}
	} ()
	go func() {
		var wg sync.WaitGroup
		defer close(retc)
		for input := range inputc {
			wg.Add(1)
			launchImager(&wg, retc, input)
		}
		wg.Wait()
	} ()
	return retc
}

func launchImager(wg *sync.WaitGroup, retc chan Img, input Merged) {
	go func() {
		defer wg.Done()
		pxwand := NewPixelWand()
		pxwand.SetColor("#FFFFFF")
		input.Fonts.SetFillColor(pxwand)
		input.Fonts.SetFontSize(12)
		input.Fonts.Annotation(float64(input.Coordinate.X), float64(input.Coordinate.Y), strconv.Itoa(input.Num))
		input.Fonts.Rotate(input.Angle)
		defer pxwand.Destroy()

		blackWand := NewDrawingWand()
		pxwand2 := NewPixelWand()
		pxwand2.SetColor("#000000")

		defer pxwand2.Destroy()

		blackWand.SetFillColor(pxwand2)
		blackWand.Color(0, 0, PAINT_METHOD_RESET )

		img := NewMagickWand()
		img.SetSize(28, 28)
		img.DrawImage(blackWand)
		img.DrawImage(input.Fonts)
		retc <- Img{Image:img, Number:input.Num}

	}()
}
