package main

import "github.com/DimKa163/gophermart/app/gophermart"

func main() {
	var conf gophermart.Config
	ParseFlags(&conf)
	server, err := gophermart.New(conf)
	if err != nil {
		panic(err)
	}
	server.Map()
	if err := server.Run(); err != nil {
		panic(err)
	}
}
