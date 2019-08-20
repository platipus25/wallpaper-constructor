package wallpaperconstructor

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"

	// import for side effects to register jpeg image format
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"strings"

	wallpaperconstructor "github.com/platipus25/wallpaperconstructor/process"
)

func main() {
	fmt.Println("start")
	data := ""
	out := Process(1440, 900, data, 20)

	err := ioutil.WriteFile("out.txt", []byte(out), 0644)
	if err != nil {
		log.Fatalf("Oops:", err)
	}

	//fmt.Println(out)
}

// Process takes in a base64 uses wallpaperconstructor.ProcessImg(width, height, img, blurRadius
func Process(width, height int, imgStr string, blurRadius int) string {

	reader := strings.NewReader(imgStr)
	dat := base64.NewDecoder(base64.StdEncoding, reader)
	img, _, err := image.Decode(dat)
	if err != nil {
		log.Fatalf("Error loading image:", err)
	}

	out := wallpaperconstructor.ProcessImg(width, height, img, blurRadius)

	fmt.Println(len(out.Pix))
	var b = new(bytes.Buffer)
	encoder := png.Encoder{png.BestCompression, nil}
	err = encoder.Encode(b, out)
	if err != nil {
		log.Fatalf("Error making png:", err)
	}

	outStr := base64.URLEncoding.EncodeToString(b.Bytes())
	return outStr
}
