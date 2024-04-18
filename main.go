package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var lastContent string

func main() {
	url := "https://raw.githubusercontent.com/alex2ale/hello/main/m.txt"

	for {
		content, err := readURLContent(url)
		if err != nil {
			fmt.Println("error:", err)
			time.Sleep(10 * time.Second)
			continue
		}
		if content == lastContent {
			fmt.Println("URL content hasn't changed. Going to sleep.")
			time.Sleep(10 * time.Second)
			continue
		}
		lastContent = content
		lines := parseContent(content)
		var wg sync.WaitGroup
		for _, line := range lines {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				makeRequest(url)
			}(line)
		}
		wg.Wait()
	}
}

func readURLContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func parseContent(content string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

func makeRequest(url string) {
	if url == "wait" {
		time.Sleep(10 * time.Second)
		fmt.Println("wait dasht")
	} else {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error making request to %s: %v\n", url, err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body from %s: %v\n", url, err)
			return
		}
		fmt.Printf("Response from %s: %s\n", url, string(body))
	}
}
