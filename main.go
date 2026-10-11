package main

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

type Experience struct {
	Company    string   `json:"company"`
	Role       string   `json:"role"`
	Period     string   `json:"period"`
	Location   string   `json:"location"`
	Highlights []string `json:"highlights"`
}

type Education struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Period      string `json:"period"`
	Details     string `json:"details"`
}

type Resume struct {
	Name       string              `json:"name"`
	Title      string              `json:"title"`
	Email      string              `json:"email"`
	Location   string              `json:"location"`
	Phone      string              `json:"phone"`
	Summary    string              `json:"summary"`
	Experience []Experience        `json:"experience"`
	Education  []Education         `json:"education"`
	Skills     map[string][]string `json:"skills"`
	Hobbies    []string            `json:"hobbies"`
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(targetPath, data, info.Mode())
	})
}

func copyStaticPDF(src, dst string) {
	input, err := os.ReadFile(src)
	if err != nil {
		log.Println("Note: No data/resume.pdf found. Skipping PDF copy step.")
		return
	}
	_ = os.WriteFile(dst, input, 0644)
}

func main() {
	data, err := os.ReadFile("data/resume.json")
	if err != nil {
		log.Fatalf("Failed to read resume.json: %v", err)
	}

	var resume Resume
	if err := json.Unmarshal(data, &resume); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	if err := os.MkdirAll("dist", 0755); err != nil {
		log.Fatalf("Failed to create dist directory: %v", err)
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("Failed to parse HTML template: %v", err)
	}

	htmlFile, err := os.Create("dist/index.html")
	if err != nil {
		log.Fatalf("Failed to create index.html: %v", err)
	}
	defer htmlFile.Close()

	if err := tmpl.Execute(htmlFile, resume); err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	copyStaticPDF("data/resume.pdf", "dist/resume.pdf")
	_ = copyDir("templates/static", "dist/static")
	log.Println("Build successful! Output generated in /dist ✅️")
}
