package main

import (
	"example/wspinapp-backend/internal/common"
	"fmt"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func generateMigration() {
	common.InitDbWithConfig(&gorm.Config{
		DryRun: true,
	})
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run script.go /path/to/folder name_of_new_migration")
		return
	}

	inputPath := os.Args[1]
	migrationName := os.Args[2]

	files, err := os.ReadDir(inputPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	highestNumber := 0
	pattern := regexp.MustCompile(`^(\d{3})-`)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		match := pattern.FindStringSubmatch(file.Name())
		if len(match) == 2 {
			number, err := strconv.Atoi(match[1])
			if err == nil && number > highestNumber {
				highestNumber = number
			}
		}
	}

	nextNumber := highestNumber + 1
	newFileName := fmt.Sprintf("%03d-%s", nextNumber, migrationName)
	newFilePath := filepath.Join(inputPath, newFileName)

	file, err := os.Create(newFilePath)
	if err != nil {
		fmt.Println("Error creating new file:", err)
		return
	}

	_, err = file.WriteString(generateMigration())
	if err != nil {
		fmt.Println("Error creating new file:", err)
		return
	}
	fmt.Printf("Created empty file: %s\n", newFilePath)
}
