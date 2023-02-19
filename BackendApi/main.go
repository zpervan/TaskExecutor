package main

import (
	"net/http"
	"taskexecutor/backend/server"
	"taskexecutor/logger"
)

func main() {
	logger.Info("Starting API service server")
	r := server.InitializeRoutes()
	http.ListenAndServe(":3500", r.HandleRequests())
}
