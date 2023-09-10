package main

import (
	"bufio"
	"example/wspinapp-backend/internal/common"
	"fmt"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func setEnvFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2) // Split the string into at most 2 parts
		if len(parts) != 2 {
			fmt.Printf("Ignoring invalid line: %s\n", line)
			continue
		}

		key := parts[0]
		value := parts[1]

		if err := os.Setenv(key, value); err != nil {
			fmt.Printf("Failed to set environment variable for key: %s, error: %v\n", key, err)
		}
	}

	return scanner.Err()
}

func generateMigration() {
	os.Setenv("POSTGRES_HOST", "wspinapp-backend.ddns.net")
	common.InitDbWithConfig(&gorm.Config{
		DryRun: true,
	})
}

func createFile(inputPath string, migrationName string) (*os.File, error) {
	files, err := os.ReadDir(inputPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return nil, err
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

	return os.Create(newFilePath)
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run script.go /path/to/folder name_of_new_migration")
		return
	}

	inputPath := os.Args[1]
	migrationName := os.Args[2] + ".sql"
	envFilePath := os.Args[3]

	err := setEnvFromFile(envFilePath)
	if err != nil {
		fmt.Println("Error loading environment from file:", err)
		return
	}

	file, err := createFile(inputPath, migrationName)
	if err != nil {
		fmt.Println("Error creating new file:", err)
		return
	}

	func() {
		out := os.Stdout
		err := os.Stderr
		os.Stdout = file
		os.Stderr = file
		generateMigration()
		os.Stdout = out
		os.Stderr = err
	}()

	fmt.Println("Created migration file:", migrationName)
}
