package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"syscall/js"

	// import for side effects to register jpeg image format
	_ "image/jpeg"
	"image/png"
	"log"
	"strings"

	wallpaperconstructor "github.com/platipus25/wallpaperconstructor/process"
)

func main() {
	/*fmt.Println("start")
	data := ""
	out := Process(1440, 900, data, 20)

	err := ioutil.WriteFile("out.txt", []byte(out), 0644)
	if err != nil {
		log.Fatalf("Oops:", err)
	}*/
	holdOpen := make(chan struct{}, 0)

	fmt.Println("Hello Web,\n\tGolang")
	var fun js.Func
	fun = js.FuncOf(jsProcess)

	js.Global().Set("processImg", fun)

	<-holdOpen
	//fmt.Println(out)
}

func jsProcess(this js.Value, i []js.Value) interface{} {
	fmt.Println("Processing...")
	outStr := Process(i[0].Int(), i[1].Int(), i[2].String(), i[3].Int())
	fmt.Println("Done.")
	return js.ValueOf(outStr)
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
