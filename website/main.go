package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
)

func main() {
	os.RemoveAll("./build/")
	os.MkdirAll("./build/", os.ModePerm)

	os.RemoveAll("./.tmp/snippet_repo/")
	os.MkdirAll("./.tmp/snippet_repo/", os.ModePerm)

	err := cloneGitRepo("./.tmp/snippet_repo/")
	if err != nil {
		panic(err)
	}

	err = copyPublicFolder()
	if err != nil {
		panic(err)
	}

	if err = os.Mkdir("./build/snippets/", os.ModePerm); err != nil {
		panic(err)
	}

	snippets, err := os.ReadDir(SNIPPET_DIRECTORY)
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

func cloneGitRepo(destination string) error {
	_, err := git.PlainClone(destination, false, &git.CloneOptions{
		URL:      "https://github.com/bedrock-oss/bedrock-examples",
		Progress: os.Stdout,
	})
	return err
}

// Contributor represents a GitHub contributor
type Contributor struct {
	Login         string `json:"login"`
	HTMLURL       string `json:"html_url"`
	AvatarURL     string `json:"avatar_url"`
	Contributions int    `json:"contributions"`
}

func GetGitHubContributorsAPI(owner, repo string, token string) ([]Contributor, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors", owner, repo)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// Optional: Use token for higher rate limits
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var contributors []Contributor
	err = json.NewDecoder(resp.Body).Decode(&contributors)
	if err != nil {
		return nil, err
	}

	return contributors, nil
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

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatalln("Failed to generate CSS:", err.Error(), "\n", string(output))
	}
}
