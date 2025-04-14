package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

func main() {
	// Check if file path is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path")
		os.Exit(1)
	}

	filePath := os.Args[1]

	// Exit if file doesn't have .rb extension
	if !strings.HasSuffix(filePath, ".rb") {
		fmt.Println("File must have .rb extension")
		os.Exit(1)
	}

	// Split the path into segments
	pathSegments := strings.Split(filePath, "/")

	// Check if the file is a spec file or a code file
	if strings.Contains(filePath, "_spec.rb") && pathSegments[0] == "spec" {
		// It's a spec file, find the corresponding code file
		pathSegments[0] = "app"
		codeFilePath := strings.Join(pathSegments, "/")
		codeFilePath = strings.Replace(codeFilePath, "_spec.rb", ".rb", 1)
		cmd := exec.Command("zed", codeFilePath)
		cmd.Start()
	} else if pathSegments[0] == "app" {
		// It's a code file, find or create the corresponding spec file
		pathSegments[0] = "spec"
		specFilePath := strings.Join(pathSegments, "/")
		specFilePath = strings.Replace(specFilePath, ".rb", "_spec.rb", 1)

		// Create the spec file if it doesn't exist
		if _, err := os.Stat(specFilePath); os.IsNotExist(err) {
			// Ensure the directory exists
			specDir := filepath.Dir(specFilePath)
			os.MkdirAll(specDir, 0755)

			// Extract class name from source file
			className := extractClassName(filePath)
			if className == "" {
				// Fallback to capitalized filename
				baseFileName := filepath.Base(filePath)
				baseFileName = strings.TrimSuffix(baseFileName, ".rb")
				className = capitalize(baseFileName)
			}

			// Create the spec file
			createSpecFile(specFilePath, className)
			fmt.Printf("Created spec file: %s\n", specFilePath)
		}

		// Open the spec file with Zed
		cmd := exec.Command("zed", specFilePath)
		cmd.Start()
	} else {
		fmt.Println("File must be in the app directory to create a spec for it")
		os.Exit(1)
	}
}

// extractClassName reads a file and extracts the first class name found
func extractClassName(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	// Create a regex to match class definitions
	classRegex := regexp.MustCompile(`class\s+([A-Z][A-Za-z0-9_:]*)\b`)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "class ") {
			matches := classRegex.FindStringSubmatch(line)
			if len(matches) > 1 {
				return matches[1]
			}
		}
	}

	return ""
}

// createSpecFile creates a new RSpec spec file with basic structure
func createSpecFile(filePath, className string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write the basic RSpec template
	writer.WriteString("# frozen_string_literal: true\n\n")
	writer.WriteString("require \"rails_helper\"\n\n")
	writer.WriteString(fmt.Sprintf("RSpec.describe %s do\n", className))
	writer.WriteString("  # Your specs here\n")
	writer.WriteString("end\n")

	return writer.Flush()
}

// capitalize returns a string with the first letter capitalized
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}