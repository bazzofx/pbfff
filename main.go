package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)
// Color codes
const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

// Print colored text
func printColoredText(text, color string) {
	fmt.Printf("%s%s%s\n", color, text, Reset)
}
// ColorWord function colors a specific word in the string
func ColorWord(text, word, color string) string {
	// Replace the target word with the colored word
	coloredWord := color + word + Reset
	return fmt.Sprintf(strings.Replace(text, word, coloredWord, 1))
}

// Global variables for output directory, header, and body options
var outputDir string
var headerOnly bool
var bodyOnly bool
var numJobs int
var publicIP string // Store the public IP

// Create an HTTP client that ignores SSL verification
func createInsecureClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
}

// FetchPublicIP retrieves the public IP by querying an external service
func FetchPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text") // Use an IP discovery service
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

// FetchHeaders fetches the headers for a given URL and includes the status code
func FetchHeaders(url string) (string, error) {
	client := createInsecureClient()
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	headers := fmt.Sprintf("HTTP Status Code: %d %s\n", resp.StatusCode, resp.Status)
	for key, values := range resp.Header {
		headers += fmt.Sprintf("%s: %s\n", key, strings.Join(values, ", "))
	}
	return headers, nil
}

// FetchBody fetches the body for a given URL
func FetchBody(url string) (string, error) {
	client := createInsecureClient()
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

// SaveContent saves content to a specified file
func SaveContent(content, filepath string) error {
	return ioutil.WriteFile(filepath, []byte(content), 0644)
}

// ProcessURL processes a single URL, fetching headers and/or body
func ProcessURL(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	urlFolderName := strings.TrimPrefix(strings.TrimPrefix(url, "http://"), "https://")
	urlDir := filepath.Join(outputDir, SanitizeFilename(urlFolderName))
	os.MkdirAll(urlDir, os.ModePerm)

	if !bodyOnly {
		headers, err := FetchHeaders(url)
		if err == nil {
			headerFile := filepath.Join(urlDir, "header.txt")
			SaveContent(headers, headerFile)
			// Create the info string with the actual values
			infostring := fmt.Sprintf("[%s] Sucessfully fetched header for %s",publicIP, url )
			// Color the word "header" in green
			result := ColorWord(infostring, publicIP, Green)			
			// Print the result
			fmt.Println(result)
			
		} else {
			fmt.Printf("Failed to fetch header for %s: %v\n", url, err)
		}
	}

	if !headerOnly {
		body, err := FetchBody(url)
		if err == nil {
			bodyFile := filepath.Join(urlDir, "body.txt")
			SaveContent(body, bodyFile)
			infostring := fmt.Sprintf("[%s] Sucessfully fetched body for %s",publicIP, url )
			// Color the word "header" in green
			result := ColorWord(infostring, publicIP, Green)			
			// Print the result
			fmt.Println(result)
			
		} else {
			fmt.Printf("Failed to fetch body for %s: %v\n", url, err)
		}
	}
}

// SanitizeFilename converts a URL into a safe directory name
func SanitizeFilename(url string) string {
	return strings.ReplaceAll(url, "/", "_")
}

// ReadURLs reads URLs from a file or stdin
func ReadURLs(reader io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(reader)
	var urls []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			urls = append(urls, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func main() {
	urlFile := flag.String("f", "", "File containing URLs, one per line (optional, will use STDIN if not provided)")
	flag.StringVar(&outputDir, "o", "output", "Output directory for saved files")
	flag.BoolVar(&headerOnly, "h", false, "Only download headers")
	flag.BoolVar(&bodyOnly, "b", false, "Only download body")
	flag.IntVar(&numJobs, "n", 4, "Number of parallel jobs to run (default is 4)")
	flag.Parse()

	// Fetch public IP once at the beginning
	var err error
	publicIP, err = FetchPublicIP()
	if err != nil {
		fmt.Printf("Failed to fetch public IP: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Public IP: %s\n", publicIP)

	var urls []string
	if *urlFile != "" {
		file, err := os.Open(*urlFile)
		if err != nil {
			fmt.Printf("Failed to open file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		urls, err = ReadURLs(file)
	} else {
		urls, err = ReadURLs(os.Stdin)
	}

	if err != nil {
		fmt.Printf("Failed to read URLs: %v\n", err)
		os.Exit(1)
	}

	urlsChan := make(chan string, len(urls))
	var wg sync.WaitGroup

	for i := 0; i < numJobs; i++ {
		go func() {
			for url := range urlsChan {
				ProcessURL(url, &wg)
			}
		}()
	}

	for _, url := range urls {
		wg.Add(1)
		urlsChan <- url
	}

	close(urlsChan)
	wg.Wait()
	fmt.Println("Finished processing all URLs.")
}
