package main

import (
	"errors"
	"github.com/DimKa163/gophermart/app/gophermart"
	"net/http"
)

func main() {
	var conf gophermart.Config
	ParseFlags(&conf)
	server, err := gophermart.New(conf)
	if err != nil {
		panic(err)
	}
	server.Map()
	if err := server.Run(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}

	}
}
