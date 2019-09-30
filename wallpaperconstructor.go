package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"strings"

	"github.com/disintegration/imaging"
	wallpaperconstructor "github.com/platipus25/wallpaperconstructor/process"
)

func main() {
	fileName := flag.String("file", "img.jpeg", "image to turn into wallpaper")
	fileNameOut := flag.String("out", "wallpaper-output.png", "file for output")
	blurRadius := flag.Int("blur", 10, "blur radius for background; suggested: 10 - 50 \n\t(the higher the blur radius, the longer it takes to run)")
	w := flag.Int("w", 1440, "The width of the output image")
	h := flag.Int("h", 900, "The height of the output image")
	width := flag.Int("width", 1440, "The width of the output image")
	height := flag.Int("height", 900, "The height of the output image")
	flag.Parse()
	tail := flag.Args()
	if len(tail) > 0 && *fileName == "img.jpeg" {
		*fileName = tail[0]
	}
	if len(tail) > 1 && *fileNameOut == "wallpaper-output.png" {
		*fileNameOut = tail[1]
	} else if *fileNameOut == "wallpaper-output.png" {
		*fileNameOut = strings.Split(*fileName, ".")[0] + "-wallpaper.png"
	}
	if *w != 1440 {
		width = w
	}
	if *h != 900 {
		height = h
	}
	fmt.Printf("\tInput File: %s\n\tOutput File: %s\n\tBlur Radius: %v\n\tWidth: %v\n\tHeight: %v\n", *fileName, *fileNameOut, *blurRadius, *width, *height)
	// Decode the JPEG data. If reading from file, create a reader with
	img, err := imaging.Open(*fileName)
	imgNRGBA, ok := img.(*image.NRGBA)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	if ok != true {
		log.Fatalf("failed to convert image to NRGBA")
	}
	fmt.Println("Image loaded")

	out := wallpaperconstructor.ProcessImg(*width, *height, imgNRGBA, *blurRadius)

	err = imaging.Save(out, *fileNameOut)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
	fmt.Println("Image saved")
}
