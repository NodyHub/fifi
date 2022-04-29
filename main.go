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
	"strconv"
	"strings"
	"sync"
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

type cliParameter struct {
	Authorization    string
	Cookie           string
	Crash            bool
	Method           string
	Host             string
	JsonOutput       bool
	MaxRetry         int
	ParallelRequests bool
	ResponseCode     bool
	ServerHeader     bool
	Timeout          int
	Useragent        string
	Verbose          bool
	Wait             int
}

// hash calulates the hash value of a given string
func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))

	return fmt.Sprint(h.Sum32())
}

func getResponseSignature(parsedArgs *cliParameter, url *urlResponse) string {
	var raw string = ""
	// Add Response code to signature
	if parsedArgs.ResponseCode {
		raw = strconv.Itoa(url.StatusCode)
	}
	// Add server header value to signature
	if parsedArgs.ServerHeader {
		var srvHdrId = -1
		for i, h := range url.HeaderEntries {
			if strings.EqualFold("server", h.Key) {
				srvHdrId = i
				break
			}
		}
		if srvHdrId > -1 {
			raw = fmt.Sprintf("%s;%s", raw, url.HeaderEntries[srvHdrId].Value)
		} else {
			raw = fmt.Sprintf("%s;%s", raw, "(none)")
		}
	}
	// Concat all response header key's
	for _, h := range url.HeaderEntries {
		raw = fmt.Sprintf("%s;%s", raw, h.Key)
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
			// extract and preserver the original header value
			field := strings.TrimSpace(strings.Split(row, ":")[0])
			for h, v := range response.Header {
				// Compare with http lib extracted headers and store
				if strings.EqualFold(field, h) {
					idx := hdrCounter[field]
					res = append(res, headerEntry{field, v[idx]})
					hdrCounter[field] = idx + 1
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

func performRequest(parsedArgs *cliParameter, url string) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Duration(parsedArgs.Timeout) * time.Second,
	}

	// Declare HTTP Method and Url
	req, err := http.NewRequest(parsedArgs.Method, url, nil)
	if err != nil {
		return nil, err
	}

	// Set Auth
	if len(parsedArgs.Authorization) > 0 {
		req.Header.Add("Authorization", parsedArgs.Authorization)
	}

	// Set Cookie
	if len(parsedArgs.Cookie) > 0 {
		req.Header.Add("Cookie", parsedArgs.Cookie)
	}

	// Set Host
	if len(parsedArgs.Host) > 0 {
		req.Host = parsedArgs.Host
	}

	// Set UserAgent
	if len(parsedArgs.Useragent) > 0 {
		req.Header.Add("User-Agent", parsedArgs.Useragent)
	}
	// Perform get request
	resp, err := client.Do(req)

	// first at all, check for crash
	if err != nil && parsedArgs.Crash {
		return nil, err
	}
	//  the other error handling
	retry := 0
	for retry < parsedArgs.MaxRetry && err != nil {
		log.Printf("ERROR (%v): %s\n", retry, err.Error())
		if os.IsTimeout(err) || resp.StatusCode == 429 {
			time.Sleep(time.Second * time.Duration(retry+1))
		} else {
			retry = parsedArgs.MaxRetry
		}
		retry++
	}
	if retry == parsedArgs.MaxRetry {
		return nil, fmt.Errorf("maxRetry(%v) reached, go to next url\n", parsedArgs.MaxRetry)
	}

	return resp, nil
}

// storeResult evaluates the response from the http request and stores it
func storeResult(mtx *sync.RWMutex, parsedArgs *cliParameter, resp *http.Response, result *map[string][]urlResponse, url string) error {
	headers, err := getHeaders(resp)
	if err != nil {
		return err
	}
	parsedResponse := urlResponse{url, headers, resp.StatusCode, resp.Status}
	sig := getResponseSignature(parsedArgs, &parsedResponse)

	if mtx != nil {
		mtx.RLock()
	}
	if _, exist := (*result)[sig]; exist {
		(*result)[sig] = append((*result)[sig], parsedResponse)
	} else {
		(*result)[sig] = []urlResponse{parsedResponse}
	}
	if parsedArgs.Verbose {
		log.Printf("%s %s\n", sig, url)
	}
	if mtx != nil {
		mtx.RUnlock()
	}

	return nil
}

func getAllSignatures(parsedArgs *cliParameter, urls *map[string]struct{}) (map[string][]urlResponse, error) {
	result := make(map[string][]urlResponse)
	var mtx *sync.RWMutex = nil
	var wg sync.WaitGroup

	if parsedArgs.ParallelRequests {
		mtx = new(sync.RWMutex)
		urlsLen := len(*urls)
		wg.Add(urlsLen)
	}

	for url := range *urls {

		// Perform requests in parallel
		if parsedArgs.ParallelRequests {
			go func(mtx *sync.RWMutex, parsedArgs *cliParameter, result *map[string][]urlResponse, url string) {
				defer wg.Done()
				// Perform get request
				resp, err := performRequest(parsedArgs, url)
				if err != nil {
					log.Printf("ERROR: %s", err)
					return
				}

				err = storeResult(mtx, parsedArgs, resp, result, url)
				if err != nil {
					log.Printf("ERROR: %s", err)
					return
				}
			}(mtx, parsedArgs, &result, url)
		} else {
			// Perform get request sequentiel
			resp, err := performRequest(parsedArgs, url)
			if err != nil {
				log.Printf("ERROR: %s", err)
				continue
			}

			err = storeResult(mtx, parsedArgs, resp, &result, url)
			if err != nil {
				log.Printf("ERROR: %s", err)
				continue
			}
		}

	}

	if parsedArgs.ParallelRequests {
		wg.Wait()
	}

	return result, nil
}

func getSimilarHeaders(collectedResponses map[string][]urlResponse) map[string]bool {
	// Collect all headers for all responses
	headerMap := make(map[string]bool)
	for _, responses := range collectedResponses {
		resp := responses[0]
		for _, header := range resp.HeaderEntries {
			headerMap[header.Key] = true
		}
	}

	// iterate over found headers and check if they are existend in every response
	for header := range headerMap {
		found := 0
		for _, responses := range collectedResponses {
			resp := responses[0]
			for _, entry := range resp.HeaderEntries {
				if header == entry.Key {
					found++
					break
				}
			}
		}

		// Delete header if not found in every response
		if found < len(collectedResponses) {
			delete(headerMap, header)
		}
	}

	return headerMap
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
	maxRetry := flag.Int("m", 3, "Maximum retries for request")
	parallelRequests := flag.Bool("p", false, "Perform requests in parallel")
	responseCode := flag.Bool("r", false, "Include HTTP response code in signature calculation")
	serverHeader := flag.Bool("s", false, "Include 'Server' response header in signature calculation")
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

	parsedArgs := cliParameter{
		*authorization,
		*cookie,
		*crash,
		*method,
		*host,
		*jsonOutput,
		*maxRetry,
		*parallelRequests,
		*responseCode,
		*serverHeader,
		*timeout,
		*useragent,
		*verbose,
		*wait}

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
	res, err := getAllSignatures(&parsedArgs, &unifiedUrls)
	if err != nil {
		log.Fatal(fmt.Sprintf("ERROR! %s", err.Error()))
	}

	// Get all headers that are in every response
	similarHeaders := getSimilarHeaders(res)

	// Output result
	if *jsonOutput {
		b, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(string(b))
	} else {
		fmt.Println("\nSummary:")
		// Iterate over headers that are existent in every request
		if *verbose {
			fmt.Println("===================================")
			fmt.Println("Headers received in every response:")
			fmt.Println("===================================")
			for header := range similarHeaders {
				fmt.Printf(" - %s\n", header)
			}
			fmt.Println("===================================")
		}

		// Iterate over all Signatures
		fmt.Println("")
		signatures := []string{}
		for sig := range res {
			signatures = append(signatures, sig)
		}
		sort.Strings(signatures)
		for _, sig := range signatures {
			responses := res[sig]
			fmt.Println("-----------------------------------")
			fmt.Printf("Signature: %s ; URLs: %v\n", sig, len(responses))

			if *verbose {
				// Iterate over response header
				fmt.Println("Additional headers:")
				for _, h := range responses[0].HeaderEntries {
					if *serverHeader && strings.EqualFold(h.Key, "server") {
						fmt.Printf(" - %s: %s\n", h.Key, h.Value)
					} else if found := similarHeaders[h.Key]; !found {
						fmt.Printf(" - %s\n", h.Key)
					}
				}
				fmt.Println("")
				fmt.Println("Urls: ")

				// Iterate over sorted list of urls
				urls := make([]string, len(responses))
				for i, r := range responses {
					urls[i] = fmt.Sprintf("[%v] %s", r.StatusCode, r.Url)
				}
				sort.Strings(urls)
				for _, u := range urls {
					fmt.Println(u)
				}
				fmt.Println("-----------------------------------\n")
			}
		}
	}
}
