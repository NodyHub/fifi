package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
	"strings"
	"time"
)

type headerEntry struct {
	Key   string
	Value string
}

type urlResponse struct {
	Url           string
	HeaderEntries []headerEntry
	StatusCode    int
	Status        string
}

// hash calulates the hash value of a given string
func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))

	return fmt.Sprint(h.Sum32())
}

func getResponseSignature(url urlResponse) string {
	raw := "(none)"
	for i, h := range url.HeaderEntries {
		if i == 0 {
			raw = h.Key
		} else {
			raw = fmt.Sprintf("%s;%s", raw, h.Key)
		}
	}

	return hash(raw)
}

func getHeaders(response *http.Response) ([]headerEntry, error) {
	res := make([]headerEntry, 0)
	// Prepare counter if a header is contained multiple times in response
	hdrCounter := make(map[string]int, 0)
	// Iterate over raw response line-wise
	rawResponse, err := httputil.DumpResponse(response, false)
	if err != nil {
		return nil, err
	}
	for _, row := range strings.Split(string(rawResponse), "\n") {
		if len(row) > 0 {
			field := strings.TrimSpace(strings.Split(row, ":")[0])
			for h, v := range response.Header {
				if strings.EqualFold(field, h) {
					idx := hdrCounter[h]
					res = append(res, headerEntry{h, v[idx]})
					hdrCounter[h] = idx + 1
				}
			}
		}
	}

	return res, nil
}

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

func getSignature(crash, verbose bool, maxRetry, timeout, wait int, authorization, cookie, host, method, useragent string, urls map[string]struct{}) (map[string][]urlResponse, error) {
	// headerMap := make(map[string][]string)
	result := make(map[string][]urlResponse)

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	for url := range urls {
		// Declare HTTP Method and Url
		req, err := http.NewRequest(method, url, nil)
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
		// first at all, check for crash
		if err != nil && crash {
			return nil, err
		}
		//  the other error handling
		retry := 0
		for retry < maxRetry && err != nil {
			log.Printf("ERROR (%v): %s\n", retry, err.Error())
			if os.IsTimeout(err) || resp.StatusCode == 429 {
				time.Sleep(time.Second * time.Duration(retry+1))
			} else {
				retry = maxRetry
			}
			retry++
		}
		if retry == maxRetry {
			log.Printf("maxRetry(%v) reached, go to next url\n", maxRetry)
			continue
		}

		// Handle response and evaluate
		headers, err := getHeaders(resp)
		if err != nil {
			return nil, err
		}

		parsedResponse := urlResponse{url, headers, resp.StatusCode, resp.Status}
		sig := getResponseSignature(parsedResponse)
		if _, exist := result[sig]; exist {
			result[sig] = append(result[sig], parsedResponse)
		} else {
			result[sig] = []urlResponse{parsedResponse}
		}

		if verbose {
			log.Printf("%s %s\n", sig, url)
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
	crash := flag.Bool("C", false, "Crash on error")
	method := flag.String("X", "GET", "Method")
	host := flag.String("H", "", "Host")
	jsonOutput := flag.Bool("j", false, "Result as json")
	maxRetry := flag.Int("r", 3, "Maximum retries for request")
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
	res, err := getSignature(*crash, *verbose, *maxRetry, *timeout, *wait, *authorization, *cookie, *host, *method, *useragent, unifiedUrls)
	if err != nil {
		log.Fatal(fmt.Sprintf("ERROR! %s", err.Error()))
	}

	// Output result
	if *jsonOutput {
		b, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(string(b))
	} else {
		fmt.Println("\nSummary:")
		// Iterate over all Signatures
		for sig, responses := range res {
			fmt.Printf("Signature: %s ; URLs: %v\n", sig, len(responses))

			if *verbose {
				// Iterate over response header
				for _, h := range responses[0].HeaderEntries {
					fmt.Printf(" - %s\n", h.Key)
				}
				fmt.Println("-----")
				fmt.Println("Urls: ")
				fmt.Println("-----")

				// Iterate over sorted list of urls
				urls := make([]string, len(responses))
				for i, r := range responses {
					urls[i] = fmt.Sprintf("[%v] %s", r.StatusCode, r.Url)
				}
				sort.Strings(urls)
				for _, u := range urls {
					fmt.Println(u)
				}
				fmt.Println("-----")
			}
		}
	}
}
