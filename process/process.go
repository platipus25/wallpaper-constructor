package wallpaperconstructor

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/disintegration/imaging"
)

// ProcessImg layers the input and resizes it to make an image that will fit within the width and height parameters
/*
	Layers:
		1. Bottom layer is semi-opaque solid color background using prominent color
		2. Next layer is a blurred and stretched copy of the image
		3. The foreground is fitted to the frame and the aspect ratio is preserved, then is streched a small ammount
*/
func ProcessImg(width, height int, img image.Image, blurRadius int) image.Image {
	blur := make(chan image.Image)
	foreground := make(chan image.Image)
	background := make(chan image.Image)

	go func(image image.Image) { // blur
		blurred := imaging.Fill(image, width, height, imaging.Center, imaging.Linear)
		blurred = imaging.Blur(blurred, float64(blurRadius))
		fmt.Println("Image Blurred")
		blur <- blurred
	}(img)

	go func(image image.Image) { // foreground
		resized := fit(image, width, height, imaging.Lanczos)
		resized = resizeToIsh(resized, width, height, imaging.Lanczos)
		fmt.Println("Image resized")
		foreground <- resized
	}(img)

	go func(image image.Image) { // background
		promColors, err := prominentcolor.Kmeans(img)
		promColor := promColors[0].Color
		if err != nil {
			log.Fatalf("prominentcolor failed: %v", err)
		}
		fmt.Printf("Prominent colors: %v\n", promColors)
		background <- imaging.New(width, height, color.NRGBA{uint8(promColor.R), uint8(promColor.G), uint8(promColor.B), 200})
	}(img)

	// put it all together
	out := <-background
	out = imaging.OverlayCenter(out, <-blur, 0.5)
	out = imaging.PasteCenter(out, <-foreground)
	fmt.Println("Composite compiled")

	return out
}

func resizeToIsh(img image.Image, targetWidth int, targetHeight int, filter imaging.ResampleFilter) *image.NRGBA {
	bounds := img.Bounds().Size()
	width, height := bounds.X, bounds.Y
	aspect := float64(width) / float64(height) // width:height aspect:1

	fmt.Printf("\tOriginal Aspect Ratio:\t%v:1\n", aspect)

	deltaWidth := float64(targetWidth-width) / 6
	deltaHeight := float64(targetHeight-height) / 5

	finalWidth := width + int(deltaWidth)
	finalHeight := height + int(deltaHeight)

	aspect = float64(finalWidth) / float64(finalHeight)
	fmt.Printf("\tResult Aspect Ratio:\t%v:1\n", aspect)

	//fmt.Printf("deltaW:%v\n deltaH:%v\n finalW:%v\n finalH:%v\n origW:%v\n origH:%v\n targetW%v\n targetH%v\n", deltaWidth, deltaHeight, finalWidth, finalHeight, width, height, targetWidth, targetHeight)

	return imaging.Resize(img, finalWidth, finalHeight, filter)
}

/*
	if original aspect > target aspect then width should be smaller
	if orignial aspect < target aspects then height should be smaller

	if width should be smaller then width should be set to target width
	if height should be smaller height should be set to target height

	the other will be changed in order to preserve aspect ratio
*/
func fit(img image.Image, targetWidth int, targetHeight int, filter imaging.ResampleFilter) *image.NRGBA {
	bounds := img.Bounds().Size()
	width, height := bounds.X, bounds.Y
	aspect := float64(width) / float64(height) // width:height aspect:1

	targetAspect := float64(targetWidth) / float64(targetHeight) // width:height aspect:1

	if aspect > targetAspect {
		targetHeight = 0 // make height change automatically to preserve aspect ratio
	}

	if aspect < targetAspect {
		targetWidth = 0 // make width change automatically to preserve aspect ratio
	}

	return imaging.Resize(img, targetWidth, targetHeight, filter)
}
