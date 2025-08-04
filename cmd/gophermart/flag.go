package main

import (
	"flag"
	"github.com/DimKa163/gophermart/app/gophermart"
	"github.com/DimKa163/gophermart/internal/env"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"os"
)

func ParseFlags(config *gophermart.Config) {
	var argonConfig auth.ArgonConfig
	var argonMemory uint
	var argonIterations uint
	var argonParallelism uint
	var argonSaltLength uint
	var argonKeyLength uint
	flag.StringVar(&config.Addr, "a", ":8080", "The address to listen on")
	flag.StringVar(&config.Database, "d", "", "The database to connect to")
	flag.StringVar(&config.Accrual, "r", "", "Accrual service")
	flag.StringVar(&config.Secret, "s", "secret", "Secret service")
	flag.StringVar(&config.Secret, "l", "info", "Log level")
	flag.StringVar(&config.Schedule, "sh", "*/10 * * * * *", "schedule")
	flag.UintVar(&argonMemory, "m", 64, "argon memory")
	flag.UintVar(&argonIterations, "i", 3, "argon iteration")
	flag.UintVar(&argonParallelism, "pr", 2, "argon parallelism")
	flag.UintVar(&argonSaltLength, "sl", 16, "argon salt length")
	flag.UintVar(&argonKeyLength, "kl", 32, "argon key length")
	flag.Parse()

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

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		config.LogLevel = envLogLevel
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		config.LogLevel = envLogLevel
	}
	if envScheduleLog := os.Getenv("WORKER_SCHEDULE"); envScheduleLog != "" {
		config.Schedule = envScheduleLog
	}
	env.ParseUIntEnv("ARGON_MEMORY", &argonMemory)
	env.ParseUIntEnv("ARGON_ITERATION", &argonIterations)
	env.ParseUIntEnv("ARGON_PARALLELISM", &argonParallelism)
	env.ParseUIntEnv("ARGON_SALT_LENGTH", &argonSaltLength)
	env.ParseUIntEnv("ARGON_KEY_LENGTH", &argonKeyLength)
	argonConfig.Memory = uint32(argonMemory * 1024)
	argonConfig.Iterations = uint32(argonIterations)
	argonConfig.Parallelism = uint32(argonParallelism)
	argonConfig.SaltLength = uint32(argonSaltLength)
	argonConfig.KeyLength = uint32(argonKeyLength)
	config.Argon = argonConfig
}
