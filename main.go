package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// filter function for slices
func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}

	return
}

// getUrls reads file line-by-line and assume that they are all urls
func getUrls(inputFile string) ([]string, error) {
	// Read while file content
	log.Printf("Reading %s", inputFile)
	urls, err := os.ReadFile(inputFile)
	if err != nil {
		return []string{}, err
	}

	// Filter empty lines
	return filter(strings.Split(string(urls), "\n"), func(url string) bool {
		return url != ""
	}), nil
}

// readFromStdin reads from stdin until eol
func readFromStdin() ([]string, error) {
	var urls []string
	in := bufio.NewReader(os.Stdin)
	for {
		s, err := in.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		url := strings.TrimSpace(s)
		if url != "" {
			urls = append(urls, url)
		}
	}

	return urls, nil
}

func getSignature(verbose bool, timeout, wait int, authorization, cookie, host, useragent string, urls map[string]struct{}) (map[string][]string, error) {
	result := make(map[string][]string)

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	for url := range urls {
		// Declare HTTP Method and Url
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		// Set Auth
		if len(authorization) > 0 {
			req.Header.Add("Authorization", authorization)
		}

		// Set Cookie
		if len(cookie) > 0 {
			req.Header.Add("Cookie", cookie)
		}

		// Set Host
		if len(host) > 0 {
			req.Host = host
		}

		// Set UserAgent
		if len(useragent) > 0 {
			req.Header.Add("User-Agent", useragent)
		}

		// Perform get request
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		// Handle response and evaluate
		srv := resp.Header.Get("Server")
		if srv == "" {
			srv = "(none)"
		}

		if _, exist := result[srv]; exist {
			result[srv] = append(result[srv], url)
		} else {
			result[srv] = []string{url}
		}

		if verbose {
			fmt.Printf("%s %s\n", srv, url)
		}

		time.Sleep(time.Duration(wait) * time.Millisecond)
	}

	return result, nil
}

const (
	usage = `usage: %s [files]
Parse urls and fetch Server banners.

Options:
[files] provide the urls in files.
`
)

func main() {
	log.SetOutput(flag.CommandLine.Output())

	// Read cli param
	authorization := flag.String("a", "", "Authorization")
	cookie := flag.String("c", "", "Cookie")
	host := flag.String("H", "", "Host")
	timeout := flag.Int("t", 1, "Timeout seconds")
	useragent := flag.String("u", "", "User-Agent (default GoLang default)")
	verbose := flag.Bool("v", false, "Verbose output")
	wait := flag.Int("w", 0, "Wait ms between requests")
	flag.Usage = func() {
		log.Printf(usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	input := flag.Args()

	// Get URLS
	var err error
	var urls []string
	if len(input) == 0 {
		log.Println("reading from stdin...")
		urls, err = readFromStdin()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Read files
		for _, ifile := range input {
			newURLs, err := getUrls(ifile)
			if err != nil {
				log.Fatal(err)
			}

			urls = append(urls, newURLs...)
		}
	}

	// unify urls
	unifiedUrls := make(map[string]struct{}) // New empty set
	for _, url := range urls {
		unifiedUrls[url] = struct{}{} // Add
	}

	log.Printf("Collected %v different urls, starting analysis\n", len(unifiedUrls))
	res, err := getSignature(*verbose, *timeout, *wait, *authorization, *cookie, *host, *useragent, unifiedUrls)
	if err != nil {
		log.Fatal(err)
	}

	// Output result
	if *verbose {
		fmt.Println("\nSummary:")
	}
	for srv, subset := range res {
		fmt.Printf("%s %v urls\n", srv, len(subset))
	}
}
