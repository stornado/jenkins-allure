package jenkinsallure

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
)

func CompressImageResource(data []byte) []byte {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return data
	}
	buf := bytes.Buffer{}
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 40})
	if err != nil {
		return data
	}
	if buf.Len() > len(data) {
		return data
	}
	return buf.Bytes()
}

func CompressPNGResource(data []byte) []byte {
	imgSrc, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return data
	}
	newImg := image.NewRGBA(imgSrc.Bounds())
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(newImg, newImg.Bounds(), imgSrc, imgSrc.Bounds().Min, draw.Over)

	buf := bytes.Buffer{}
	err = jpeg.Encode(&buf, newImg, &jpeg.Options{Quality: 40})
	if err != nil {
		return data
	}
	if buf.Len() > len(data) {
		return data
	}
	return buf.Bytes()
}
