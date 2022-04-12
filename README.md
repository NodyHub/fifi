# Fifi

Fifi fingerprints [`Server`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Server) response headers from a given list of urls and groups them.

## Background

Recently, spring boot had a wide spreaded RCE vulnerability, known as [Spring4Shell](https://portswigger.net/daily-swig/spring4shell-microsoft-cisa-warn-of-limited-in-the-wild-exploitation) ([CVE-2022-22965](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-22965)). Due to the fact that modern web application are implemented based on micro service pattern, various paths of a domain may end up on different applications/containers/CDN. To limit the attack surface system administrator, DevOps Engineers and SRE's are highly interested in limiting the available information about a service in the public.

This tool provides help to identify irregulations in the response headers from a given list of urls.

## Installation

```
go install github.com/NodyHub/fifi@latest
```

## Usage

```
% fifi -h
usage: fifi [files]
Parse urls and fetch Server banners.

Options:
[files] provide the urls in files.
  -H string
    	Host
  -a string
    	Authorization
  -c string
    	Cookie
  -s	Server banner only grouping
  -t int
    	Timeout seconds (default 1)
  -u string
    	User-Agent (default GoLang default)
  -v	Verbose output
  -w int
    	Wait ms between requests
```

## Example output

This example gives an indication about the hosting infrastructure from Snapchat

```
% cat urls.lst | fifi -v -t 4 -w 500 | tee result.out
2022/04/12 09:30:01 reading from stdin...
2022/04/12 09:30:01 Collected 60 different urls, starting analysis
2276290871 https://www.namecheap.com/cart
511284460 https://www.namecheap.com/cart/ajax/SessionHandler.ashx?_=1649748094813
3632238580 https://www.namecheap.com/cart/json/operation.aspx
511284460 https://www.namecheap.com/cart/json/operation.aspx/GetALLIteminCart
2116670402 https://www.namecheap.com/
511284460 https://www.namecheap.com/api/v1/ncpl

[...]

976437910 https://www.namecheap.com/myaccount/twofa/secondauth.aspx
2571465822 https://www.namecheap.com/api/v1/ncpl/twofactorauth/uiauthenticate/checkAndUpdateUserLockStatus
511284460 https://www.namecheap.com/cart/ajax/messagehandler.ashx
976437910 https://www.namecheap.com/domains/registration/results/
2276290871 https://www.namecheap.com/legal

Summary:
ID: 511284460 ; URLs: 11 ; Server: cloudflare
 - Access-Control-Allow-Credentials
 - Access-Control-Allow-Headers
 - Access-Control-Allow-Methods
 - Cache-Control
 - Cf-Cache-Status
 - Cf-Ray
 - Content-Type
 - Date
 - Expect-Ct
 - Server
 - Set-Cookie
 - Strict-Transport-Security
 - Vary
 - X-Frame-Options
 - X-Inst
 - X-Xss-Protection
ID: 2571465822 ; URLs: 5 ; Server: cloudflare
 - Cf-Cache-Status
 - Cf-Ray
 - Content-Type
 - Date
 - Expect-Ct
 - Server
 - Set-Cookie
 - Strict-Transport-Security
 - Www-Authenticate
ID: 976437910 ; URLs: 11 ; Server: cloudflare
 - Cf-Cache-Status
 - Cf-Ray
 - Content-Type
 - Date
 - Expect-Ct
 - Server
 - Set-Cookie
 - Strict-Transport-Security
 - Vary
 - X-Frame-Options
 - X-Xss-Protection

[...]

ID: 16678959 ; URLs: 1 ; Server: cloudflare
 - Cf-Cache-Status
 - Cf-Ray
 - Content-Type
 - Date
 - Expect-Ct
 - Server
 - Set-Cookie
 - Strict-Transport-Security
 - Vary
 - X-Frame-Options
 - X-Xss-Protection

% cat result.out | grep 511284460
511284460 https://www.namecheap.com/cart/ajax/SessionHandler.ashx?_=1649748094813
511284460 https://www.namecheap.com/cart/json/operation.aspx/GetALLIteminCart
511284460 https://www.namecheap.com/api/v1/ncpl/cart
511284460 https://www.namecheap.com/cart/json/operation.aspx/GetCartTotalAndItemCount
511284460 https://www.namecheap.com/api/v1/ncpl/twofactorauth
511284460 https://www.namecheap.com/api/v1
511284460 https://www.namecheap.com/api/v1/ncpl
511284460 https://www.namecheap.com/api/v1/ncpl/usermessages
511284460 https://www.namecheap.com/cart/ajax/SessionHandler.ashx
511284460 https://www.namecheap.com/namecheap-opensearch.xml
511284460 https://www.namecheap.com/cart/ajax/messagehandler.ashx
ID: 511284460 ; URLs: 11 ; Server: cloudflare

% cat README.md | grep 2571465822
2571465822 https://www.namecheap.com/api/v1/ncpl/twofactorauth/uiauthenticate/getDeviceCodeStatus
2571465822 https://www.namecheap.com/api/v1/ncpl/twofactorauth/uiauthenticate/getStatus
2571465822 https://www.namecheap.com/api/v1/ncpl/usermessages/user/getMessages
2571465822 https://www.namecheap.com/api/v1/ncpl/twofactorauth/uiauthenticate/verifyDeviceCode
2571465822 https://www.namecheap.com/api/v1/ncpl/twofactorauth/uiauthenticate/checkAndUpdateUserLockStatus
ID: 2571465822 ; URLs: 5 ; Server: cloudflare

% cat README.md | grep 976437910
976437910 https://www.namecheap.com/domains
976437910 https://www.namecheap.com/myaccount/twofa/secondauth.aspx?ReturnUrl=https%3a%2f%2fap.www.namecheap.com
976437910 https://www.namecheap.com/domains/registration/results?domain=foobar&_gl=1*twy0qe*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjEuMTY0OTc0ODE1NC42MA..
976437910 https://www.namecheap.com/twofa
976437910 https://www.namecheap.com/twofa/device
976437910 https://www.namecheap.com/domains/registration/results
976437910 https://www.namecheap.com/twofa/device?ReturnUrl=https%3a%2f%2fap.www.namecheap.com
976437910 https://www.namecheap.com/twofa/device?ReturnUrl=%2f%3f_gl%3d1*5eunvn*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjAuMTY0OTc0ODA5MC42MA..%26_ga%3d2.76907415.360072847.1649748091-842392789.1649748091&_gl=1*5eunvn*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjAuMTY0OTc0ODA5MC42MA..&_ga=2.76907415.360072847.1649748091-842392789.1649748091
976437910 https://www.namecheap.com/domains/registration/results/?domain=foobar&_gl=1*twy0qe*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjEuMTY0OTc0ODE1NC42MA..
976437910 https://www.namecheap.com/myaccount/twofa/secondauth.aspx
976437910 https://www.namecheap.com/domains/registration/results/
ID: 976437910 ; URLs: 11 ; Server: cloudflare
```

