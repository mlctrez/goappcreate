package internal

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"strconv"
	"strings"
)

func CreatePng(width, height int) (pngBytes []byte, err error) {
	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: width, Y: height},
	})
	center := image.Point{X: width / 2, Y: height / 2}

	alpha := func(x, y int) float64 {
		distance := math.Sqrt(math.Pow(float64(x-center.X), 2) + math.Pow(float64(y-center.Y), 2))
		// alpha should be 0 at edges and ff in center, assume square
		if distance > float64(width/2) {
			return 0
		}
		return (float64(width/2) - distance) / float64(width/2)

	}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2:
				img.Set(x, y, RGBAFromHex("0x4CEADB", alpha(x, y)))
			case x >= width/2 && y < height/2:
				img.Set(x, y, RGBAFromHex("0x56ACE7", alpha(x, y)))
			case x >= width/2 && y >= height/2:
				img.Set(x, y, RGBAFromHex("0x5A93EC", alpha(x, y)))
			case x < width/2 && y >= height/2:
				img.Set(x, y, RGBAFromHex("0x6260F4", alpha(x, y)))
			default:
			}
		}
	}

	b := &bytes.Buffer{}
	err = png.Encode(b, img)
	pngBytes = b.Bytes()

	return
}

func RGBAFromHex(hex string, alpha float64) color.RGBA {
	value, err := strconv.ParseInt(strings.TrimPrefix(hex, "0x"), 16, 64)
	if err != nil {
		panic(err)
	}
	red := float64((value>>16)&0xff) * alpha
	green := float64((value>>8)&0xff) * alpha
	blue := float64(value&0xff) * alpha
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(alpha * 255)}
}
