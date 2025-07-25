package main

import (
	"flag"
	"github.com/DimKa163/gophermart/app/gophermart"
	"os"
)

func ParseFlags(config *gophermart.Config) {
	flag.StringVar(&config.Addr, "a", ":8080", "The address to listen on")
	flag.StringVar(&config.Database, "d", "", "The database to connect to")
	flag.StringVar(&config.Accrual, "r", "", "Accrual service")
	flag.StringVar(&config.Secret, "s", "secret", "Secret service")
	if addrValue := os.Getenv("RUN_ADDRESS"); addrValue != "" {
		config.Addr = addrValue
	}

	if databaseValue := os.Getenv("DATABASE_URI"); databaseValue != "" {
		config.Database = databaseValue
	}

	if accrualSystemaValue := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); accrualSystemaValue != "" {
		config.Accrual = accrualSystemaValue
	}

	if secretValue := os.Getenv("SECRET_KEY"); secretValue != "" {
		config.Secret = secretValue
	}
}
