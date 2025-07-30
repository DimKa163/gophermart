package gophermart

import "github.com/DimKa163/gophermart/internal/shared/auth"

type Config struct {
	Addr     string
	Database string
	Accrual  string
	Secret   string
	LogLevel string
	Argon    auth.ArgonConfig
}
