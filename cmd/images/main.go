package main

import (
	"image-converter/internal/app"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

const defaultImagePath = "./files/d.webp"

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
		mux.HandleFunc("/images/", imagesHandle)
		if err := http.ListenAndServe(":8080", mux); err != nil {
			slog.Warn("Server shutdown with error", "error", err.Error())
		}
	}()

	select {}
}

func imagesHandle(w http.ResponseWriter, r *http.Request) {
	imagePath := strings.Replace(r.URL.Path, "/images/", "./files/", 1)
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		http.ServeFile(w, r, defaultImagePath)
	}
	http.ServeFile(w, r, imagePath)
}
