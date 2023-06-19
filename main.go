package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/go-rod/rod"
)

const (
	targetURL     = "https://registry.terraform.io/providers/hashicorp/template/latest/docs"
	xpath         = "/html/body/div[3]/div[2]/div/div/article/div[1]/div/div"
	deprecationMsg = "This provider is deprecated."
	errorMsg      = "An error occurred while checking the page: %v"
)

func main() {
	// Run the `terraform providers` command and capture the output
	output, err := runTerraformCommand()
	if err != nil {
		log.Fatalf("Failed to execute 'terraform providers' command: %v", err)
	}

	// Use the command output as the text to be parsed
	text := string(output)

	// Define the regular expression pattern
	regex := regexp.MustCompile(`provider\[registry\.terraform\.io/hashicorp/([^/\]]+)`)

    lines := strings.Split(text, "\n")

// Deduplicate the matches
    uniqueProviders := make(map[string]bool)

    for _, line := range lines {
        if !strings.Contains(line, "module.") {
            match := regex.FindStringSubmatch(line)
            if len(match) > 1 {
                trimmed := match[1]
                uniqueProviders[trimmed] = true
            }
        }
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
	command := exec.Command("terraform", "providers")

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

	// Evaluate the XPath expression to retrieve the deprecation message
	elem := page.MustElementX(xpath)
	deprecationText, err := elem.Text()
	if err != nil {
		log.Fatalf("Failed to retrieve the deprecation message: %v", err)
	}

	// Check if the deprecation message indicates deprecation
	isDeprecated := strings.Contains(deprecationText, deprecationMsg)

	// Close the browser
	browser.MustClose()

	return isDeprecated, nil
}
