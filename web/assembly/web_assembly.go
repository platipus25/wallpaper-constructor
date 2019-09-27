package main

import (
	"fmt"
	"syscall/js"

	"github.com/platipus25/wallpaperconstructor/web"
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
	outStr := web.Process(i[0].Int(), i[1].Int(), i[2].String(), i[3].Int())
	fmt.Println("Done.")
	return js.ValueOf(outStr)
}
