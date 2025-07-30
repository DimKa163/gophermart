package main

import (
	"errors"
	"github.com/DimKa163/gophermart/app/gophermart"
	"github.com/DimKa163/gophermart/internal/shared/logging"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	var conf gophermart.Config
	ParseFlags(&conf)
	server, err := gophermart.New(conf)
	if err != nil {
		panic(err)
	}
	err = server.AddLogging()
	if err != nil {
		panic(err)
	}
	server.Map()
	if err = server.Run(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logging.Log.Fatal("Failed to run server", zap.Error(err))
		}
	}
}
