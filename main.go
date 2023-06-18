package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/go-rod/rod"
)

func main() {
	// Run the `terraform providers` command and capture the output
	output, err := runTerraformCommand("providers")
	if err != nil {
		log.Fatalf("Failed to execute 'terraform providers' command: %v", err)
	}

	// Use the command output as the text to be parsed
	text := string(output)

	// Define the regular expression pattern
	pattern := regexp.MustCompile(`provider\[registry\.terraform\.io/[^\]]+\]`)

	// Find all matches in the text
	matches := pattern.FindAllString(text, -1)

	// Deduplicate the matches
	uniqueProviders := make(map[string]bool)
	for _, match := range matches {
		provider := strings.TrimPrefix(match, "provider[registry.terraform.io/")
		provider = strings.TrimSuffix(provider, "]")
		uniqueProviders[provider] = true
	}

	// Check deprecation status for each provider
	for provider := range uniqueProviders {
		targetURL := fmt.Sprintf("https://registry.terraform.io/providers/hashicorp/%s/latest/docs", provider)
		isDeprecated, err := checkDeprecationStatus(targetURL)
		if err != nil {
			log.Fatalf("An error occurred while checking the deprecation status for %s: %v", provider, err)
		}

		if isDeprecated {
			fmt.Printf("Provider %s is deprecated.\n", provider)
		} else {
			fmt.Printf("Provider %s is not deprecated.\n", provider)
		}
	}
}

func runTerraformCommand(args ...string) ([]byte, error) {
	// Specify the command to execute
	command := exec.Command("terraform", args...)

	// Execute the command and capture the output
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %v", err)
	}

	return output, nil
}

func checkDeprecationStatus(targetURL string) (bool, error) {
	// Launch a headless browser
	browser := rod.New().MustConnect()

	// Create a new page
	page := browser.MustPage()

	// Navigate to the target URL
	page.MustNavigate(targetURL)

	// Wait for the page to load
	page.MustWaitLoad()

	// Find the deprecation message on the page
	elem := page.MustElement("body")
	deprecationText, err := elem.Text()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve the deprecation message: %v", err)
	}

	// Check if the deprecation message is present
	isDeprecated := strings.Contains(deprecationText, "This provider is deprecated.")

	// Close the browser
	browser.MustClose()

	return isDeprecated, nil
}
