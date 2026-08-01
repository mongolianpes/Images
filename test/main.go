package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func fetchImage(url string, saveAs string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("сервер вернул статус %d", resp.StatusCode)
	}

	// сохраняем ответ в файл, чтобы потом сравнить
	out, err := os.Create(saveAs)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}

	fmt.Println("Скачано:", saveAs)
	return nil
}

func main() {
	// пример: запрос существующего файла
	err := fetchImage("http://localhost:8080/images/2026/7/22/63abfd67-58f4-4a87-aa9e-23590e203e33/img/63abfd67-58f4-4a87-aa9e-23590e203e33.webp", "existing_result.webp")
	if err != nil {
		fmt.Println("Ошибка:", err)
	}

	// пример: запрос несуществующего файла
	err = fetchImage("http://localhost:8080/images/notfound.webp", "notfound_result.webp")
	if err != nil {
		fmt.Println("Ошибка:", err)
	}
}
