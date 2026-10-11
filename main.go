package main

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
)

type Experience struct {
	Company    string   `json:"company"`
	Role       string   `json:"role"`
	Period     string   `json:"period"`
	Highlights []string `json:"highlights"`
}

type Resume struct {
	Name       string       `json:"name"`
	Title      string       `json:"title"`
	Email      string       `json:"email"`
	Location   string       `json:"location"`
	Github     string       `json:"github"`
	Linkedin   string       `json:"linkedin"`
	Summary    string       `json:"summary"`
	Experience []Experience `json:"experience"`
	Skills     []string     `json:"skills"`
}

func main() {
	// 1. Read JSON data
	data, err := os.ReadFile("data/resume.json")
	if err != nil {
		log.Fatalf("Failed to read resume.json: %v", err)
	}

	var resume Resume
	if err := json.Unmarshal(data, &resume); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// 2. Ensure dist directory exists
	if err := os.MkdirAll("dist", 0755); err != nil {
		log.Fatalf("Failed to create dist directory: %v", err)
	}

	// 3. Generate HTML
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

	// 4. Generate PDF (Placeholder / Asset copy step)
	// Tip: You can generate your PDF using a library like github.com/johnfercher/maroto
	// or drop a pre-compiled clean resume.pdf into your repo for the build script to copy over.
	copyStaticPDF("data/resume.pdf", "dist/resume.pdf")

	log.Println("Successfully built static site and assets in /dist 🚀")
}

func copyStaticPDF(src, dst string) {
	input, err := os.ReadFile(src)
	if err != nil {
		log.Println("Warning: No static PDF found to copy. Make sure dist/resume.pdf can be served.")
		return
	}
	_ = os.WriteFile(dst, input, 0644)
}
