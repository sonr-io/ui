//go:build ignore
// +build ignore

// build.go - Build script for the UI module
// Usage: go run build.go [options]
//
// This script orchestrates the build process:
// 1. Runs TailwindCSS CLI to generate CSS
// 2. Runs templ generate to create Go files from templates  
// 3. Builds the Go module
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	skipCSS      = flag.Bool("skip-css", false, "Skip TailwindCSS generation")
	skipTempl    = flag.Bool("skip-templ", false, "Skip templ generation")
	skipBuild    = flag.Bool("skip-build", false, "Skip Go build")
	skipTidy     = flag.Bool("skip-tidy", false, "Skip go mod tidy")
	watch        = flag.Bool("watch", false, "Watch mode for development")
	cssOnly      = flag.Bool("css-only", false, "Only generate CSS")
	templOnly    = flag.Bool("templ-only", false, "Only generate templ files")
	localAssets  = flag.Bool("local-assets", false, "Build with localassets tag")
	installDeps  = flag.Bool("install-deps", false, "Install required dependencies")
	verbose      = flag.Bool("v", false, "Verbose output")
)

func main() {
	flag.Parse()

	if *installDeps {
		installDependencies()
		return
	}

	if *cssOnly {
		generateCSS()
		return
	}

	if *templOnly {
		generateTempl()
		return
	}

	if *watch {
		runWatchMode()
		return
	}

	// Full build process
	if !*skipCSS {
		generateCSS()
	}

	if !*skipTempl {
		generateTempl()
	}

	if !*skipTidy {
		tidyModules()
	}

	if !*skipBuild {
		buildModule()
	}

	fmt.Println("✅ Build completed successfully!")
}

func installDependencies() {
	fmt.Println("📦 Installing build dependencies...")

	// Check and install templ
	if !commandExists("templ") {
		fmt.Println("Installing templ...")
		runCommand("go", "install", "github.com/a-h/templ/cmd/templ@latest")
	} else {
		fmt.Println("✓ templ is already installed")
	}

	// Check and install tailwindcss
	if !commandExists("tailwindcss") {
		fmt.Println("Installing tailwindcss standalone CLI...")
		installTailwindCLI()
	} else {
		fmt.Println("✓ tailwindcss is already installed")
	}

	// Check and install templui
	if !commandExists("templui") {
		fmt.Println("Installing templui...")
		runCommand("go", "install", "github.com/axzilla/templui/cmd/templui@latest")
	} else {
		fmt.Println("✓ templui is already installed")
	}

	fmt.Println("✅ All dependencies installed!")
}

func installTailwindCLI() {
	// Detect OS and architecture
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	
	// Map Go arch to Tailwind arch names
	archMap := map[string]string{
		"amd64": "x64",
		"arm64": "arm64",
		"386":   "x86",
	}
	
	arch, ok := archMap[goarch]
	if !ok {
		log.Fatalf("Unsupported architecture: %s", goarch)
	}

	// Construct download URL
	version := "v3.4.17" // Latest stable version
	var url, filename string
	
	switch goos {
	case "linux":
		filename = fmt.Sprintf("tailwindcss-linux-%s", arch)
		url = fmt.Sprintf("https://github.com/tailwindlabs/tailwindcss/releases/download/%s/%s", version, filename)
	case "darwin":
		filename = fmt.Sprintf("tailwindcss-macos-%s", arch)
		url = fmt.Sprintf("https://github.com/tailwindlabs/tailwindcss/releases/download/%s/%s", version, filename)
	case "windows":
		filename = fmt.Sprintf("tailwindcss-windows-%s.exe", arch)
		url = fmt.Sprintf("https://github.com/tailwindlabs/tailwindcss/releases/download/%s/%s", version, filename)
	default:
		log.Fatalf("Unsupported OS: %s", goos)
	}

	// Create bin directory
	binDir := filepath.Join(".", "bin")
	os.MkdirAll(binDir, 0755)
	
	// Download tailwindcss
	targetPath := filepath.Join(binDir, "tailwindcss")
	if goos == "windows" {
		targetPath += ".exe"
	}
	
	fmt.Printf("Downloading TailwindCSS from %s...\n", url)
	runCommand("curl", "-L", "-o", targetPath, url)
	
	// Make executable on Unix systems
	if goos != "windows" {
		os.Chmod(targetPath, 0755)
	}
	
	fmt.Printf("TailwindCSS installed to %s\n", targetPath)
	fmt.Println("Add ./bin to your PATH to use tailwindcss command")
}

func generateCSS() {
	fmt.Println("🎨 Generating CSS with TailwindCSS...")
	
	// Check if input.css exists
	inputCSS := "./assets/css/input.css"
	if _, err := os.Stat(inputCSS); os.IsNotExist(err) {
		// Create a default input.css if it doesn't exist
		os.MkdirAll("./assets/css", 0755)
		defaultCSS := `@tailwind base;
@tailwind components;
@tailwind utilities;`
		os.WriteFile(inputCSS, []byte(defaultCSS), 0644)
	}
	
	outputCSS := "./assets/css/output.css"
	
	// Try different tailwindcss command locations
	tailwindCmd := findTailwindCommand()
	if tailwindCmd == "" {
		log.Fatal("tailwindcss not found. Run 'go run build.go -install-deps' to install it")
	}
	
	args := []string{"-i", inputCSS, "-o", outputCSS, "--minify"}
	if *watch {
		args = append(args, "--watch")
	}
	
	runCommand(tailwindCmd, args...)
	
	if !*watch {
		fmt.Println("✓ CSS generated successfully")
	}
}

func findTailwindCommand() string {
	// Check common locations
	locations := []string{
		"tailwindcss",                           // System PATH
		"./bin/tailwindcss",                     // Local bin
		"./node_modules/.bin/tailwindcss",       // npm install
	}
	
	if runtime.GOOS == "windows" {
		locations = append(locations, 
			"tailwindcss.exe",
			"./bin/tailwindcss.exe",
			"./node_modules/.bin/tailwindcss.exe",
		)
	}
	
	for _, loc := range locations {
		if commandExists(loc) || fileExists(loc) {
			return loc
		}
	}
	
	// Try with npx as fallback
	if commandExists("npx") {
		return "npx tailwindcss"
	}
	
	return ""
}

func generateTempl() {
	fmt.Println("🔧 Generating Go code from templ files...")
	
	if !commandExists("templ") {
		log.Fatal("templ not found. Run 'go run build.go -install-deps' to install it")
	}
	
	args := []string{"generate"}
	if *watch {
		args = append(args, "--watch")
	}
	
	runCommand("templ", args...)
	
	if !*watch {
		fmt.Println("✓ Templ files generated successfully")
	}
}

func tidyModules() {
	fmt.Println("📦 Tidying Go modules...")
	runCommand("go", "mod", "tidy")
	fmt.Println("✓ Modules tidied")
}

func buildModule() {
	fmt.Println("🔨 Building Go module...")
	
	args := []string{"build"}
	if *localAssets {
		args = append(args, "-tags", "localassets")
	}
	args = append(args, "./...")
	
	runCommand("go", args...)
	fmt.Println("✓ Module built successfully")
}

func runWatchMode() {
	fmt.Println("👀 Starting watch mode...")
	fmt.Println("Press Ctrl+C to stop")
	
	// Run both tailwind and templ in watch mode
	// These will run until interrupted
	done := make(chan bool)
	
	go func() {
		*watch = true
		generateCSS()
		done <- true
	}()
	
	go func() {
		*watch = true
		generateTempl()
		done <- true
	}()
	
	// Wait for both to complete (they won't in watch mode)
	<-done
	<-done
}

func runCommand(name string, args ...string) {
	if *verbose {
		fmt.Printf("Running: %s %s\n", name, strings.Join(args, " "))
	}
	
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	if err := cmd.Run(); err != nil {
		log.Fatalf("Command failed: %s %s\nError: %v", name, strings.Join(args, " "), err)
	}
}

func commandExists(cmd string) bool {
	// Handle compound commands like "npx tailwindcss"
	parts := strings.Fields(cmd)
	if len(parts) > 0 {
		cmd = parts[0]
	}
	
	_, err := exec.LookPath(cmd)
	return err == nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}