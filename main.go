package main

import (
	"topsort/lib"
)

func main() {
	logger := lib.NewLogger()
	logger.Log("hello")

	api := lib.NewApi(logger)
	server := lib.NewHttpServer(logger, api)
	err := server.Start(":8080")
	if err != nil {
		logger.Log()
	}
}
