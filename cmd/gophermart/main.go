package main

import (
	"errors"
	"fmt"
	"github.com/DimKa163/gophermart/app/gophermart"
	"github.com/DimKa163/gophermart/internal/shared/logging"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	var conf gophermart.Config
	ParseFlags(&conf)
	server := gophermart.New(conf)

	if err := server.AddServices(); err != nil {
		fmt.Printf("Error adding services: %v", err)
		return
	}

	if err := server.AddLogging(); err != nil {
		fmt.Printf("Error adding logging: %v", err)
		return
	}

	server.Map()
	if err := server.Run(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logging.Log.Fatal("Failed to run server", zap.Error(err))
		}
	}
}
