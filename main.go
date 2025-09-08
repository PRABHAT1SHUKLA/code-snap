package main

import (
	"flag"          // For parsing CLI flags
	"fmt"           // For printing
	"io"            // For writing to files
	"os"            // For file operations
	"path/filepath" // For handling file paths
	"strings"       // For string manipulation
)

func main() {
	// Define flags
	outputFile := flag.String("o", "snapshot.txt", "Output file path")
	includeTree := flag.Bool("tree", false, "Include directory tree in output")
	flag.Parse() // Parse the flags

	// Get the input path (positional argument, default to current dir)
	inputPath := "."
	if len(flag.Args()) > 0 {
		inputPath = flag.Args()[0]
	}

	// Check if input exists
	info, err := os.Stat(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Prepare output content
	var content strings.Builder

	// If --tree is set and it's a directory, add tree
	if *includeTree && info.IsDir() {
		content.WriteString("Directory Tree:\n")
		tree, err := generateTree(inputPath)
		if err != nil {
			fmt.Printf("Error generating tree: %v\n", err)
			os.Exit(1)
		}
		content.WriteString(tree)
		content.WriteString("\n\n")
	}

	// Add file contents
	content.WriteString("File Contents:\n")
	if info.IsDir() {
		// It's a directory: Walk and dump files
		err = filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fileContent, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				content.WriteString(fmt.Sprintf("--- %s ---\n%s\n\n", path, string(fileContent)))
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error walking directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		// It's a file: Just read it
		fileContent, err := os.ReadFile(inputPath)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		content.WriteString(fmt.Sprintf("--- %s ---\n%s\n\n", inputPath, string(fileContent)))
	}

	// Write to output file
	file, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	_, err = io.WriteString(file, content.String())
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Snapshot saved to %s\n", *outputFile)
}

// Simple function to generate a directory tree (recursive)
func generateTree(root string) (string, error) {
	var builder strings.Builder
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		depth := strings.Count(relPath, string(os.PathSeparator))
		if relPath == "." {
			builder.WriteString(root + "\n")
		} else {
			prefix := strings.Repeat("  ", depth) + "- "
			if info.IsDir() {
				prefix += "[DIR] "
			}
			builder.WriteString(prefix + filepath.Base(path) + "\n")
		}
		return nil
	})
	return builder.String(), err
}
