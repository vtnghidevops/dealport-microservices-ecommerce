package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// LoadEnv tải biến môi trường từ file .env nếu file tồn tại
// Nếu biến đã tồn tại trong môi trường, giá trị hiện tại sẽ được giữ nguyên
func LoadEnv(filePath string) error {
	// Kiểm tra xem file có tồn tại không
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("File %s không tồn tại, bỏ qua việc đọc biến môi trường", filePath)
		return nil
	}

	// Mở file .env
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Đọc từng dòng trong file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Bỏ qua comment và dòng trống
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Tách key và value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Chỉ đặt biến môi trường nếu nó chưa tồn tại
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
