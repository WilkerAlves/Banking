package main

import (
	"github.com/WilkerAlves/banking-lib/logger"
	"github.com/wilker/banking/app"
)

func main() {
	logger.Info("Starting the Application")
	app.Start()
}
