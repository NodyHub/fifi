# fifi

fifi sends to a given list of url's HTTP requests, calculates on each response a signature and groups them based on the values.

Responses with the same signature may indicate similar implementation pattern, technologies and homogenious data processing. 

## Background

Recently, spring boot had a wide spreaded RCE vulnerability, known as [Spring4Shell](https://portswigger.net/daily-swig/spring4shell-microsoft-cisa-warn-of-limited-in-the-wild-exploitation) ([CVE-2022-22965](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-22965)). Due to the fact that modern web application are implemented based on micro service pattern, various paths of a domain may end up on different applications/containers/CDN. To limit the attack surface system administrator, DevOps Engineers and SRE's are highly interested in limiting the available information about a service in the public.

This tool provides help to identify differences in the response headers from a given list of urls.

## Installation
Ether download it from [the release page](https://github.com/NodyHub/fifi/releases) or complie it by yourself:
```
go install github.com/NodyHub/fifi@latest
```

## Usage and example output

```shell
[~]% fifi -h
usage: fifi [files]
fifi sends to a given list of url's HTTP requests, calculates on each response a signature and groups them based on the values.

Default reads from stdin

Options:
--------
[files] provide the urls in files.
  -C	Crash on error
  -H string
    	Host
  -X string
    	Method (default "GET")
  -a string
    	Authorization
  -c string
    	Cookie
  -j	Result as json
  -m int
    	Maximum retries for request (default 3)
  -r	Include HTTP response code in signature calculation
  -s	Include 'Server' response header in signature calculation
  -t int
    	Threads (default 1)
  -u string
    	User-Agent (default GoLang default)
  -v	Verbose output
  -w int
    	Wait ms between requests
  -x int
    	Timeout seconds (default 1)

github.com/NodyHub/fifi@v0.2.17
[~]% cat uber.url.lst | fifi -v -t 4 -s
2022/05/03 08:29:06 reading from stdin...
2022/05/03 08:29:06 Collected 11 different urls, starting analysis
2022/05/03 08:29:06 parsedArgs.ParallelRequests: 4
2022/05/03 08:29:06 Thread 3 starts
2022/05/03 08:29:06 Thread 1 starts
2022/05/03 08:29:06 Thread 0 starts
2022/05/03 08:29:06 Thread 2 starts
2022/05/03 08:29:07 1667219945 https://auth.uber.com/
2022/05/03 08:29:07 2898507639 https://auth.uber.com/login/social
2022/05/03 08:29:07 1705792451 https://auth.uber.com/login/?breeze_local_zone=dca1&state=0A-OdN1vuv_FDbpofRZqJg9maKASCY4k0kCRVEiSDGw%3D&uber_client_name=riderSignUp&uclick_id=840a8ddd-ac10-47e6-aec4-e492968acc42
2022/05/03 08:29:07 2898507639 https://auth.uber.com/login/social/
2022/05/03 08:29:07 Thread 3 finished
2022/05/03 08:29:07 1705792451 https://auth.uber.com/login/?breeze_local_zone=dca11&next_url=https%3A%2F%2Fm.uber.com%2F&state=NUUybaiHU9SIaKz56QjyvtJTz5CJC25zhhyocPV9guM%3D
2022/05/03 08:29:07 1705792451 https://auth.uber.com/login/
2022/05/03 08:29:07 1705792451 https://auth.uber.com/login/session
2022/05/03 08:29:07 Thread 0 finished
2022/05/03 08:29:07 1705792451 https://auth.uber.com/login/?next_url=https%3A%2F%2Fm.uber.com%2F&privileged_op_url=https%3A%2F%2Fm.uber.com%2F
2022/05/03 08:29:07 Thread 2 finished
2022/05/03 08:29:07 1705792451 https://auth.uber.com/login/social/?from=facebook&state=%7B%22query%22%3A%22%3Fnext_url%3Dhttps%253A%252F%252Fm.uber.com%252F%26privileged_op_url%3Dhttps%253A%252F%252Fm.uber.com%252F%26uber_client_name%3Dm2%22%2C%22csrfToken%22%3A%221650443852-01-FNOsAwdU4I8HWkiFZuimbrTHjauX146ik_Hq9h7k1Ew%22%2C%22app%22%3A%22%22%7D&response_type=token
2022/05/03 08:29:07 1705792451 https://auth.uber.com/login/social/?next_url=https%3A%2F%2Fm.uber.com%2F&privileged_op_url=https%3A%2F%2Fm.uber.com%2F&uber_client_name=m2
2022/05/03 08:29:08 ERROR (0): Get "https://auth.uber.com/login": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
2022/05/03 08:29:09 ERROR (1): Get "https://auth.uber.com/login": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
2022/05/03 08:29:11 ERROR (2): Get "https://auth.uber.com/login": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
2022/05/03 08:29:14 ERROR: maxRetry(3) reached, go to next url
2022/05/03 08:29:14 Thread 1 finished

Summary:
===================================
Headers received in every response:
===================================
 - Alt-Svc
 - X-Content-Type-Options
 - Strict-Transport-Security
 - X-Frame-Options
 - X-Xss-Protection
 - Content-Type
 - Server
 - Vary
 - Cache-Control
 - Date
 - Via
 - X-Envoy-Upstream-Service-Time
 - X-Uber-Edge
===================================

-----------------------------------
Signature: 1667219945 ; URLs: 1
Additional headers:
 - Server: ufe

Urls:
[404] https://auth.uber.com/
-----------------------------------

-----------------------------------
Signature: 1705792451 ; URLs: 7
Additional headers:
 - Content-Security-Policy
 - Etag
 - Server: ufe
 - Set-Cookie
 - Set-Cookie
 - Timing-Allow-Origin
 - X-Content-Security-Policy
 - X-Csrf-Token
 - X-Webkit-Csp

Urls:
[200] https://auth.uber.com/login/
[200] https://auth.uber.com/login/?breeze_local_zone=dca1&state=0A-OdN1vuv_FDbpofRZqJg9maKASCY4k0kCRVEiSDGw%3D&uber_client_name=riderSignUp&uclick_id=840a8ddd-ac10-47e6-aec4-e492968acc42
[200] https://auth.uber.com/login/?breeze_local_zone=dca11&next_url=https%3A%2F%2Fm.uber.com%2F&state=NUUybaiHU9SIaKz56QjyvtJTz5CJC25zhhyocPV9guM%3D
[200] https://auth.uber.com/login/?next_url=https%3A%2F%2Fm.uber.com%2F&privileged_op_url=https%3A%2F%2Fm.uber.com%2F
[200] https://auth.uber.com/login/session
[200] https://auth.uber.com/login/social/?from=facebook&state=%7B%22query%22%3A%22%3Fnext_url%3Dhttps%253A%252F%252Fm.uber.com%252F%26privileged_op_url%3Dhttps%253A%252F%252Fm.uber.com%252F%26uber_client_name%3Dm2%22%2C%22csrfToken%22%3A%221650443852-01-FNOsAwdU4I8HWkiFZuimbrTHjauX146ik_Hq9h7k1Ew%22%2C%22app%22%3A%22%22%7D&response_type=token
[200] https://auth.uber.com/login/social/?next_url=https%3A%2F%2Fm.uber.com%2F&privileged_op_url=https%3A%2F%2Fm.uber.com%2F&uber_client_name=m2
-----------------------------------

-----------------------------------
Signature: 2898507639 ; URLs: 2
Additional headers:
 - Content-Security-Policy
 - Etag
 - Server: ufe
 - Set-Cookie
 - Set-Cookie
 - X-Content-Security-Policy
 - X-Csrf-Token
 - X-Webkit-Csp

Urls:
[404] https://auth.uber.com/login/social
[404] https://auth.uber.com/login/social/
-----------------------------------

```

# Application behaviour

```mermaid
graph TD
    A[User] -->|all url's| B(fifi)
    B --> C{For all url's}
    C --> D[Send HTTP request]
    D --> E[Calculate signature]
    E --> F[Store response, based on signature]
    F --> C
    C --> G[Generate output]

```


# Similar or related projects

* https://github.com/rverton/wonitor
* https://github.com/dgtlmoon/changedetection.io
