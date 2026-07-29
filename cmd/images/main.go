package main

import (
	"image-converter/internal/app"
	"log/slog"
	"net/http"
)

func main() {
	go func() {
		imageServer := app.NewImageServer()
		grpcServer := app.NewGrpcServer()
		err := grpcServer.GrpcServeServer(imageServer, ":8086")
		if err != nil {
			slog.Warn("Server shutdown with error", "error", err.Error())
		}
	}()

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("files"))))
		if err := http.ListenAndServe(":8080", mux); err != nil {
			slog.Warn("Server shutdown with error", "error", err.Error())
		}
	}()

	select {}
}
