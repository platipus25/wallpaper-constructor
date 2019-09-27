package main

import (
	"github.com/micro/go-micro"
	"github.com/platipus25/wallpaperconstructor/web"
)

func main(){

	servie = micro.NewService(
		micro.Name("hellowworld")
	)

	fmt.Println(web)
	service.Run()
}