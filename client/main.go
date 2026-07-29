package main

import (
	"fmt"
	"os"

	"client/images"
)

func main() {
	images.InitService()

	image, err := os.ReadFile("C:/Users/Администратор/Downloads/48x48.webp")
	if err != nil {
		fmt.Println(err.Error())
	}

	paths, err := images.SaveImage(image)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(paths)
}
