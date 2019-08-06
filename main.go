package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"strings"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/disintegration/imaging"
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
	fmt.Printf("\tInput File: %s\n\tOutput File: %s\n\tBlur Radius: %v\n", *fileName, *fileNameOut, *blurRadius)
	// Decode the JPEG data. If reading from file, create a reader with
	img, err := imaging.Open(*fileName)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	fmt.Println("Image loaded")

	blur := make(chan image.Image)
	foreground := make(chan image.Image)
	background := make(chan image.Image)

	go func(image image.Image) { // blur
		blurred := imaging.Blur(image, float64(*blurRadius))
		blurred = imaging.Fill(blurred, *width, *height, imaging.Center, imaging.Linear)
		fmt.Println("Image Blurred")
		blur <- blurred
	}(img)

	go func(image image.Image) { // foreground
		resized := imaging.Fit(image, *width, *height, imaging.Lanczos)
		resized = resizeToIsh(resized, *width, *height, imaging.Lanczos)
		fmt.Println("Image resized")
		foreground <- resized
	}(img)

	go func(image image.Image) { // foreground
		promColors, err := prominentcolor.Kmeans(img)
		promColor := promColors[0].Color
		if err != nil {
			log.Fatalf("prominentcolor failed: %v", err)
		}
		fmt.Printf("Prominent colors: %v\n", promColors)
		background <- imaging.New(*width, *height, color.NRGBA{uint8(promColor.R), uint8(promColor.G), uint8(promColor.B), 200})
	}(img)

	// put it all together
	out := <-background
	out = imaging.OverlayCenter(out, <-blur, 0.5)
	out = imaging.PasteCenter(out, <-foreground)
	fmt.Println("Composite compiled")

	err = imaging.Save(out, *fileNameOut)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
	fmt.Println("Image saved")
}

func resizeToIsh(img image.Image, targetWidth int, targetHeight int, filter imaging.ResampleFilter) *image.NRGBA {
	bounds := img.Bounds().Size()
	width, height := bounds.X, bounds.Y
	aspect := float64(width) / float64(height) // width:height 1:aspect

	fmt.Printf("\tOriginal Aspect Ratio: %v:1\n", aspect)

	deltaWidth := (targetWidth - width) / 4
	deltaHeight := (targetHeight - height) / 5

	finalWidth := width + deltaWidth
	finalHeight := height + deltaHeight

	aspect = float64(finalWidth) / float64(finalHeight)
	fmt.Printf("\tResult Aspect Ration: %v:1\n", aspect)

	return imaging.Resize(img, finalWidth, finalHeight, filter)
}
