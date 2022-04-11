package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Filter function for slices
func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

// Read file line-by-line and assume that they are all urls
func get_urls(input_file string) []string {

	// Read while file content
	fmt.Fprintf(os.Stderr, "Reading %s\n", input_file)
	urls, err := ioutil.ReadFile(input_file)
	if err != nil {
		log.Fatal(err)
		return []string{}
	}

	// Filter empty lines
	return filter(strings.Split(string(urls), "\n"), func(url string) bool {
		return url != ""
	})
}

// read from strdin until eol
func read_from_stdin() []string {
	var urls []string
	in := bufio.NewReader(os.Stdin)
	for {
		s, err := in.ReadString('\n')
		if err != nil {

			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		urls = append(urls, strings.TrimSpace(s))
	}
	// Filter empty lines
	return filter(urls, func(url string) bool {
		return url != ""
	})
}

func get_signature(verbose bool, timeout, wait *int, authorization, cookie, host, useragent *string, urls map[string]bool) map[string][]string {

	result := make(map[string][]string)

	client := http.Client{
		Timeout: time.Duration(*timeout) * time.Second,
	}

	for url := range urls {

		// Declare HTTP Method and Url
		req, err := http.NewRequest("GET", url, nil)

		// Set Auth
		if len(*authorization) > 0 {
			req.Header.Add("Authorization", *authorization)
		}

		// Set Cookie
		if len(*cookie) > 0 {
			req.Header.Add("Cookie", *cookie)
		}

		// Set Host
		if len(*host) > 0 {
			req.Host = *host
		}

		// Set UserAgent
		if len(*useragent) > 0 {
			req.Header.Add("User-Agent", *useragent)
		}

		// Perform get request
		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
			// return result
		}

		// Handle response and evaluate
		srv := resp.Header.Get("Server")
		if srv == "" {
			srv = "(none)"
		}

		_, exist := result[srv]
		if exist {
			result[srv] = append(result[srv], url)
		} else {
			result[srv] = []string{url}
		}

		if verbose {
			fmt.Printf("%s %s\n", srv, url)
		}

		time.Sleep(time.Duration(*wait) * time.Millisecond)

	}

	return result

}

const (
	usage = `usage: %s [files]
Parse urls and fetch Server banners.

Options:
[files] provide the urls in files.
`
)

func main() {
	// Read cli param
	authorization := flag.String("a", "", "Authorization")
	cookie := flag.String("c", "", "Cookie")
	host := flag.String("H", "", "Host")
	timeout := flag.Int("t", 1, "Timeout seconds")
	useragent := flag.String("u", "", "User-Agent (default GoLang default)")
	verbose := flag.Bool("v", false, "Verbose output")
	wait := flag.Int("w", 0, "Wait ms between requests")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	input := flag.Args()

	// Get URLS
	var urls []string
	if len(input) == 0 {
		fmt.Fprintf(os.Stderr, "reading from stdin...\n")
		urls = read_from_stdin()
	} else {
		// Read files
		for _, ifile := range input {
			urls = append(urls, get_urls(ifile)...)
		}
	}

	// unify urls
	unified_urls := make(map[string]bool) // New empty set
	for _, url := range urls {
		unified_urls[url] = true // Add
	}

	fmt.Fprintf(os.Stderr, "Collected %v different urls, starting analysis\n", len(unified_urls))
	res := get_signature(*verbose, timeout, wait, authorization, cookie, host, useragent, unified_urls)

	// Output result
	if *verbose {
		fmt.Println("\nSummary:")
	}
	for srv, subset := range res {
		fmt.Printf("%s %v urls\n", srv, len(subset))
	}

}
