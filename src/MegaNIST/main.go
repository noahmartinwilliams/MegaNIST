package main

import "os"
import "encoding/binary"
//import "fmt"

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
		fileImg.Write(flattened)
		//bs := make([]byte, 1)
		//binary.LittleEndian.PutUint8(bs, uint8(input.Number))
		fileLabel.Write([]byte{uint8(input.Number)})
	}
}

func flattenImg(image Img) []byte {
	ret := make([]byte, 28*28)
	for y := 0 ; y < 28 ; y = y + 1 {
		for x := 0 ; x < 28 ; x=x+1 {
			colorWand, _ := image.Image.GetImagePixelColor(x, y)
			defer colorWand.Destroy()

			intensity := byte(uint8(colorWand.GetRed()*255.0))
			ret[y*28+x]=intensity
		}
	}
	image.Image.Destroy()
	return ret
}

func main() {
	numc := GetRandsLimited(60000)
	numc2 := GetRands()
	numc3 := GetRands()

	anglesc := GetRandAngles(numc2, 15*3.141592/180.0)
	fontsc := GetFonts("/usr/share/fonts")
	coordsc := GetRandCoords(numc3, 5, 5)

	imgc := GetImages(coordsc, anglesc, fontsc, numc)
	SaveImgs(60000, imgc, "images.img", "images.label")
}
