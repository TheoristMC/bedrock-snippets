package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	os.RemoveAll("./build/")
	os.MkdirAll("./build/", os.ModePerm)

	err := copyPublicFolder()
	if err != nil {
		panic(err)
	}

	if err = os.Mkdir("./build/snippets/", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	snippets, err := os.ReadDir("./snippets")
	if err != nil {
		panic(err)
	}

	startTime := time.Now()
	fmt.Println("generating content pages...")
	for _, e := range snippets {
		if !e.IsDir() {
			continue
		}
		name := e.Name()

		generatePagesForSnippet(name)
	}

	elapsed := time.Since(startTime)
	fmt.Println("generated content pages in", elapsed)

	generateHomepage()

	fmt.Println("generating css...")
	generateTailwindCSS()

	if len(os.Args) > 1 && os.Args[1] == "-dev" {
		runDevServer()
	}
}

func copyPublicFolder() error {
	src, dest := "public", "build"
	// thanks chatgpt
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Construct destination path
		destPath := filepath.Join(dest, path[len(src):])

		if info.IsDir() {
			// Create directory at destination
			return os.MkdirAll(destPath, os.ModePerm)
		}

		// Copy file to destination
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		return err
	})
}

func generateTailwindCSS() {
	cmd := exec.Command("npx", "@tailwindcss/cli", "-i", "./main.css", "-o", "./build/main.css")
	_, err := cmd.Output()

	if err != nil {
		fmt.Println("Failed to generate CSS:", err.Error())
		return
	}
}
