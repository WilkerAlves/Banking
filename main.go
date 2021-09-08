package main

import (
	"github.com/wilker/banking/app"
	"github.com/wilker/banking/logger"
)

func main() {
	logger.Info("Starting the Application")
	app.Start()
}
