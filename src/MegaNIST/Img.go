package main

import . "gopkg.in/gographics/imagick.v3/imagick"
import "sync"
import "strconv"

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

func GetFontImages(coords chan Coord, angles chan float64, fonts chan *DrawingWand, numc chan int) chan Img {
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
		numInputs := 0
		for input := range inputc {
			if numInputs % 512 == 1 {
				wg.Wait()
			}
			wg.Add(1)
			launchImager(&wg, retc, input)
			numInputs = numInputs + 1
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

		pxwand2 := NewPixelWand()
		pxwand2.SetColor("#000000")



		img := NewMagickWand()
		img.NewImage(28, 28, pxwand2)
		img.DrawImage(input.Fonts)
		retc <- Img{Image:img, Number:input.Num}

	}()
}

func GetDrawnImages(inputc chan string, doPanic bool) chan Img {
	retc := make(chan Img, 1024)
	go func() {
		defer close(retc)
		var wg sync.WaitGroup
		x := 0 
		for input := range inputc {
			if x % 1024 == 1 {
				wg.Wait()
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				mw := NewMagickWand()
				err := mw.ReadImage(input) 
				if err != nil && doPanic{
					panic(err)
				} else if err != nil{
					return
				}
				num, err := strconv.Atoi(string(input[0]))
				if err != nil && doPanic {
					panic(err)
				} else if err != nil {
					return
				}
				retc <- Img{Image:mw, Number:num}
			} ()
		}
		wg.Wait()
	} ()
	return retc
}
