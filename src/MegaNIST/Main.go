package main

import "os"
import "encoding/binary"
import "gopkg.in/gographics/imagick.v3/imagick"
import "runtime"
import "flag"

func SaveImgs(numImgs uint32, inputc chan Img, fnameImg string, fnameLabel string) {
	fileImg, err := os.Create(fnameImg)
	if err != nil {
		panic(err)
	}
	defer fileImg.Close()

	fileLabel, err := os.Create(fnameLabel)
	if err != nil {
		panic(err)
	}
	defer fileLabel.Close()

	var magicNum uint32
	magicNum = 2051
	var size uint32
	size = 28

	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, magicNum)
	fileImg.Write(bs)
	binary.LittleEndian.PutUint32(bs, numImgs)
	fileImg.Write(bs)
	binary.LittleEndian.PutUint32(bs, size)
	fileImg.Write(bs)
	fileImg.Write(bs)

	magicNum = 2049
	binary.LittleEndian.PutUint32(bs, magicNum)
	fileLabel.Write(bs)
	binary.LittleEndian.PutUint32(bs, numImgs)
	fileLabel.Write(bs)

	for input := range inputc {
		flattened:=flattenImg(input)
		_, err := fileImg.Write(flattened)
		if err != nil {
			panic(err)
		}
		//bs := make([]byte, 1)
		//binary.LittleEndian.PutUint8(bs, uint8(input.Number))
		_, err = fileLabel.Write([]byte{uint8(input.Number)})
		if err != nil {
			panic(err)
		}
		runtime.GC()
	}
}

func flattenImg(image Img) []byte {
	ret := make([]byte, 28*28)
	var colorWand *imagick.PixelWand
	var err error
	for y := 0 ; y < 28 ; y = y + 1 {
		for x := 0 ; x < 28 ; x=x+1 {
			colorWand, err = image.Image.GetImagePixelColor(x, y)
			if err != nil {
				panic(err)
			}

			intensity := byte(uint8(colorWand.GetRed()*255.0))
			ret[y*28+x]=intensity
			//colorWand.Destroy()
		}
	}
	image.Image.Destroy()
	return ret
}

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	numImgsPtr := flag.Int("num-images", 60000, "The number of images to generate.")
	fontsLocationPtr := flag.String("fonts", "/usr/share/fonts", "The location to get fonts from.")
	maxAnglePtr := flag.Float64("angle", 15.0, "The maximum angle (in degrees) to rotate the letters.")
	maxOffsetX := flag.Int("max-offset-x", 5, "The maximum number of pixels (X axis) to offset the digit from the top left corner.")
	maxOffsetY := flag.Int("max-offset-y", 5, "The maximum number of pixels (Y axis)  to offset the digit from the top left corner.")
	saveFileImages := flag.String("img-file", "images.img", "The file to store the MegaNIST images in.")
	saveFileLabels := flag.String("label-file", "images.label", "The file to store the MegaNIST labels in.")
	doReadImages := flag.Bool("use-files", false, "Use hand drawn image files rather than fonts. Use in combination with -use-dir.")
	dirUsed := flag.String("use-dir", "dir", "Specify what directory to use for gathering hand written digits.")

	flag.Parse()

	if *doReadImages {
		filec := FindFiles(*dirUsed, true)
		imgc := GetDrawnImages(filec, false)
		SaveImgs(uint32(*numImgsPtr), imgc, *saveFileImages, *saveFileLabels)

	} else {
		numc := GetRandsLimited(int(*numImgsPtr))
		numc2 := GetRands()

		anglesc := GetRandAngles(numc2, (*maxAnglePtr)*3.141592/180.0)
		fontsc := GetFonts(*fontsLocationPtr, numc)
		coordsc := GetRandCoords(numc2, *maxOffsetX, *maxOffsetY)

		imgc := GetFontImages(coordsc, anglesc, fontsc, numc2)
		SaveImgs(uint32(*numImgsPtr), imgc, *saveFileImages, *saveFileLabels)
	}
}
