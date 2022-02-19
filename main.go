package main

import (
	"github.com/WilkerAlves/banking-lib/logger"
	"github.com/WilkerAlves/banking/app"
)

func main() {
	logger.Info("Starting the Application")
	app.Start()
}
