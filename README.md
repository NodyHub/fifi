# Fifi

Fifi collects server response header (or the [`Server`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Server) banner) from a given list of urls and groups them.

## Background

Recently, spring boot had a wide spreaded RCE vulnerability, known as [Spring4Shell](https://portswigger.net/daily-swig/spring4shell-microsoft-cisa-warn-of-limited-in-the-wild-exploitation) ([CVE-2022-22965](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-22965)). Due to the fact that modern web application are implemented based on micro service pattern, various paths of a domain may end up on different applications/containers/CDN. To limit the attack surface system administrator, DevOps Engineers and SRE's are highly interested in limiting the available information about a service in the public.

This tool provides help to identify differences in the response headers from a given list of urls.

## Installation

```
go install github.com/NodyHub/fifi@latest
```

## Usage and example output

```shell
[~/git/fifi]% fifi -h
2022/04/20 10:42:48 usage: fifi [files]
Parse urls and fetch Server banners.

Options:
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
  -r int
    	Maximum retries for request (default 3)
  -t int
    	Timeout seconds (default 1)
  -u string
    	User-Agent (default GoLang default)
  -v	Verbose output
  -w int
    	Wait ms between requests
[~/git/fifi]% cat uber.url.lst  | fifi -v
cat:  : No such file or directory
2022/04/20 10:42:52 reading from stdin...
2022/04/20 10:42:52 Collected 11 different urls, starting analysis
2022/04/20 10:42:53 1128600947 https://auth.uber.com/login/?breeze_local_zone=dca1&state=0A-OdN1vuv_FDbpofRZqJg9maKASCY4k0kCRVEiSDGw%3D&uber_client_name=riderSignUp&uclick_id=840a8ddd-ac10-47e6-aec4-e492968acc42
2022/04/20 10:42:53 148156071 https://auth.uber.com/login/social/
2022/04/20 10:42:53 1128600947 https://auth.uber.com/login/session
2022/04/20 10:42:53 148156071 https://auth.uber.com/login/social
2022/04/20 10:42:54 ERROR (0): Get "https://auth.uber.com/login/social/?from=facebook&state=%257B%2522query%2522%3A%2522%3Fnext_url%3Dhttps%253A%252F%252Fm.uber.com%252F%26privileged_op_url%3Dhttps%253A%252F%252Fm.uber.com%252F%26uber_client_name%3Dm2%2522%252C%2522csrfToken%2522%253A%25221650443852-01-FNOsAwdU4I8HWkiFZuimbrTHjauX146ik_Hq9h7k1Ew%2522%252C%2522app%2522%253A%2522%2522%257D&response_type=token": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
2022/04/20 10:42:55 ERROR (1): Get "https://auth.uber.com/login/social/?from=facebook&state=%257B%2522query%2522%3A%2522%3Fnext_url%3Dhttps%253A%252F%252Fm.uber.com%252F%26privileged_op_url%3Dhttps%253A%252F%252Fm.uber.com%252F%26uber_client_name%3Dm2%2522%252C%2522csrfToken%2522%253A%25221650443852-01-FNOsAwdU4I8HWkiFZuimbrTHjauX146ik_Hq9h7k1Ew%2522%252C%2522app%2522%253A%2522%2522%257D&response_type=token": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
2022/04/20 10:42:57 ERROR (2): Get "https://auth.uber.com/login/social/?from=facebook&state=%257B%2522query%2522%3A%2522%3Fnext_url%3Dhttps%253A%252F%252Fm.uber.com%252F%26privileged_op_url%3Dhttps%253A%252F%252Fm.uber.com%252F%26uber_client_name%3Dm2%2522%252C%2522csrfToken%2522%253A%25221650443852-01-FNOsAwdU4I8HWkiFZuimbrTHjauX146ik_Hq9h7k1Ew%2522%252C%2522app%2522%253A%2522%2522%257D&response_type=token": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
2022/04/20 10:43:00 maxRetry(3) reached, go to next url
2022/04/20 10:43:00 264694073 https://auth.uber.com/
2022/04/20 10:43:00 1128600947 https://auth.uber.com/login
2022/04/20 10:43:01 1128600947 https://auth.uber.com/login/
2022/04/20 10:43:01 1128600947 https://auth.uber.com/login/?breeze_local_zone=dca11&next_url=https%3A%2F%2Fm.uber.com%2F&state=NUUybaiHU9SIaKz56QjyvtJTz5CJC25zhhyocPV9guM%3D
2022/04/20 10:43:01 1128600947 https://auth.uber.com/login/?next_url=https%3A%2F%2Fm.uber.com%2F&privileged_op_url=https%3A%2F%2Fm.uber.com%2F
2022/04/20 10:43:01 1128600947 https://auth.uber.com/login/social/?next_url=https%3A%2F%2Fm.uber.com%2F&privileged_op_url=https%3A%2F%2Fm.uber.com%2F&uber_client_name=m2

Summary:
Signature: 1128600947 ; URLs: 7
 - Alt-Svc
 - Cache-Control
 - Content-Security-Policy
 - Content-Type
 - Date
 - Etag
 - Server
 - Set-Cookie
 - Set-Cookie
 - Strict-Transport-Security
 - Timing-Allow-Origin
 - Vary
 - Via
 - X-Content-Security-Policy
 - X-Content-Type-Options
 - X-Csrf-Token
 - X-Envoy-Upstream-Service-Time
 - X-Frame-Options
 - X-Uber-Edge
 - X-Webkit-Csp
 - X-Xss-Protection
-----
Urls:
-----
[200] https://auth.uber.com/login
[200] https://auth.uber.com/login/
[200] https://auth.uber.com/login/?breeze_local_zone=dca1&state=0A-OdN1vuv_FDbpofRZqJg9maKASCY4k0kCRVEiSDGw%3D&uber_client_name=riderSignUp&uclick_id=840a8ddd-ac10-47e6-aec4-e492968acc42
[200] https://auth.uber.com/login/?breeze_local_zone=dca11&next_url=https%3A%2F%2Fm.uber.com%2F&state=NUUybaiHU9SIaKz56QjyvtJTz5CJC25zhhyocPV9guM%3D
[200] https://auth.uber.com/login/?next_url=https%3A%2F%2Fm.uber.com%2F&privileged_op_url=https%3A%2F%2Fm.uber.com%2F
[200] https://auth.uber.com/login/session
[200] https://auth.uber.com/login/social/?next_url=https%3A%2F%2Fm.uber.com%2F&privileged_op_url=https%3A%2F%2Fm.uber.com%2F&uber_client_name=m2
-----
Signature: 148156071 ; URLs: 2
 - Alt-Svc
 - Cache-Control
 - Content-Security-Policy
 - Content-Type
 - Date
 - Etag
 - Server
 - Set-Cookie
 - Set-Cookie
 - Strict-Transport-Security
 - Vary
 - Via
 - X-Content-Security-Policy
 - X-Content-Type-Options
 - X-Csrf-Token
 - X-Envoy-Upstream-Service-Time
 - X-Frame-Options
 - X-Uber-Edge
 - X-Webkit-Csp
 - X-Xss-Protection
-----
Urls:
-----
[404] https://auth.uber.com/login/social
[404] https://auth.uber.com/login/social/
-----
Signature: 264694073 ; URLs: 1
 - Alt-Svc
 - Cache-Control
 - Content-Type
 - Date
 - Server
 - Strict-Transport-Security
 - Vary
 - Via
 - X-Content-Type-Options
 - X-Envoy-Upstream-Service-Time
 - X-Frame-Options
 - X-Uber-Edge
 - X-Xss-Protection
-----
Urls:
-----
[404] https://auth.uber.com/
-----
```

# Similar or related projects

* https://github.com/rverton/wonitor
* https://github.com/dgtlmoon/changedetection.io
