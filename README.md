# Fifi

Fifi fingerprints [`Server`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Server) response headers from a given list of urls and groups them.

## Background

Recently, spring boot had a wide spreaded RCE vulnerability, known as [Spring4Shell](https://portswigger.net/daily-swig/spring4shell-microsoft-cisa-warn-of-limited-in-the-wild-exploitation) ([CVE-2022-22965](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-22965)). Due to the fact that modern web application are implemented based on micro service pattern, various paths of a domain may end up on different applications/containers/CDN. To limit the attack surface system administrator, DevOps Engineers and SRE's are highly interested in limiting the available information about a service in the public.

This tool provides help to identify irregulations in the [`Server`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Server) response headers from a given list of urls.

## Installation

```
go install github.com/NodyHub/fifi@latest
```

## Usage

```
% fifi -h
usage: fifi [files]
Parse urls from stdin and fetch server banners.

Options:
[files] provide the urls in files.
 -H string
     Host
 -a string
     Authorization
 -c string
     Cookie
 -t int
     Timeout seconds (default 1)
 -u string
     User-Agent (default GoLang default)
 -v  Verbose output
 -w int
     Wait ms between requests
```

## Example output

This example gives an indication about the hosting infrastructure from Snapchat

```
% cat urls.lst | fifi -v -t 4 -w 500
2022/04/12 09:30:01 reading from stdin...
2022/04/12 09:30:01 Collected 60 different urls, starting analysis
2276290871 https://www.namecheap.com/cart
511284460 https://www.namecheap.com/cart/ajax/SessionHandler.ashx?_=1649748094813
3632238580 https://www.namecheap.com/cart/json/operation.aspx
511284460 https://www.namecheap.com/cart/json/operation.aspx/GetALLIteminCart
3721722846 https://www.namecheap.com/myaccount/login.aspx?ReturnUrl=%2f%3f_gl%3d1*5eunvn*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjAuMTY0OTc0ODA5MC42MA..%26_ga%3d2.76907415.360072847.1649748091-842392789.1649748091&_gl=1*5eunvn*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjAuMTY0OTc0ODA5MC42MA..&_ga=2.76907415.360072847.1649748091-842392789.1649748091
2276290871 https://www.namecheap.com/myaccount/login/
511284460 https://www.namecheap.com/api/v1/ncpl/cart
1254281072 https://www.namecheap.com/api/v1/ncpl/cart/user/getshoppingcartsummary
2571465822 https://www.namecheap.com/api/v1/ncpl/twofactorauth/uiauthenticate/getDeviceCodeStatus
511284460 https://www.namecheap.com/cart/json/operation.aspx/GetCartTotalAndItemCount
1254281072 https://www.namecheap.com/api/v1/ncpl/cart/user/get
1254281072 https://www.namecheap.com/api/v1/ncpl/cart/user/getShoppingCartSummary
511284460 https://www.namecheap.com/api/v1/ncpl/twofactorauth
3721722846 https://www.namecheap.com/myaccount
2571465822 https://www.namecheap.com/api/v1/ncpl/twofactorauth/uiauthenticate/getStatus
976437910 https://www.namecheap.com/domains
16678959 https://www.namecheap.com/domains/registration
976437910 https://www.namecheap.com/myaccount/twofa/secondauth.aspx?ReturnUrl=https%3a%2f%2fap.www.namecheap.com
2276290871 https://www.namecheap.com/legal/general/privacy-policy/?_gl=1*1vo5dcl*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjEuMTY0OTc0ODE0MS45&_ga=2.85306139.360072847.1649748091-842392789.1649748091
511284460 https://www.namecheap.com/api/v1
1254281072 https://www.namecheap.com/api/v1/ncpl/cart/user
1002554738 https://www.namecheap.com/cart/ajax
2276290871 https://www.namecheap.com/legal/general/privacy-policy/
2116670402 https://www.namecheap.com/
511284460 https://www.namecheap.com/api/v1/ncpl
976437910 https://www.namecheap.com/domains/registration/results?domain=foobar&_gl=1*twy0qe*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjEuMTY0OTc0ODE1NC42MA..
976437910 https://www.namecheap.com/twofa
1002554738 https://www.namecheap.com/datalayer
2276290871 https://www.namecheap.com/legal/general
1254281072 https://www.namecheap.com/api/v1/ncpl/twofactorauth/uiauthenticate
1254281072 https://www.namecheap.com/api/v1/ncpl/usermessages/user
3721722846 https://www.namecheap.com/myaccount/login
976437910 https://www.namecheap.com/twofa/device
511284460 https://www.namecheap.com/api/v1/ncpl/usermessages
2571465822 https://www.namecheap.com/api/v1/ncpl/usermessages/user/getMessages
976437910 https://www.namecheap.com/domains/registration/results
976437910 https://www.namecheap.com/twofa/device?ReturnUrl=https%3a%2f%2fap.www.namecheap.com
511284460 https://www.namecheap.com/cart/ajax/SessionHandler.ashx
1002554738 https://www.namecheap.com/cart/json
2276290871 https://www.namecheap.com/legal/general/privacy-policy
1002554738 https://www.namecheap.com/myaccount/twofa
2276290871 https://www.namecheap.com/legal/general/privacy-policy.aspx?_gl=1*1vo5dcl*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjEuMTY0OTc0ODE0MS45&_ga=2.85306139.360072847.1649748091-842392789.1649748091
2142470851 https://www.namecheap.com/myaccount/twofa/SecondAuthComplete.ashx
2142470851 https://www.namecheap.com/myaccount/twofa/SecondAuthComplete.ashx?ReturnUrl=%2f%3f_gl%3d1*5eunvn*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjAuMTY0OTc0ODA5MC42MA..%26_ga%3d2.76907415.360072847.1649748091-842392789.1649748091&_gl=1*5eunvn*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjAuMTY0OTc0ODA5MC42MA..&_ga=2.76907415.360072847.1649748091-842392789.1649748091
976437910 https://www.namecheap.com/twofa/device?ReturnUrl=%2f%3f_gl%3d1*5eunvn*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjAuMTY0OTc0ODA5MC42MA..%26_ga%3d2.76907415.360072847.1649748091-842392789.1649748091&_gl=1*5eunvn*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjAuMTY0OTc0ODA5MC42MA..&_ga=2.76907415.360072847.1649748091-842392789.1649748091
1254281072 https://www.namecheap.com/api/v1/ncpl/cart/user/refid
3926519757 https://www.namecheap.com/DataLayer/UserHashDataLayer.ashx
976437910 https://www.namecheap.com/domains/registration/results/?domain=foobar&_gl=1*twy0qe*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjEuMTY0OTc0ODE1NC42MA..
2276290871 https://www.namecheap.com/legal/general/privacy-policy.aspx
3721722846 https://www.namecheap.com/myaccount/login/?ReturnUrl=%2f%3f_gl%3d1*5eunvn*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjAuMTY0OTc0ODA5MC42MA..%26_ga%3d2.76907415.360072847.1649748091-842392789.1649748091&_gl=1*5eunvn*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjAuMTY0OTc0ODA5MC42MA..&_ga=2.76907415.360072847.1649748091-842392789.1649748091
511284460 https://www.namecheap.com/namecheap-opensearch.xml
2276290871 https://www.namecheap.com/api
2571465822 https://www.namecheap.com/api/v1/ncpl/twofactorauth/uiauthenticate/verifyDeviceCode
1712249365 https://www.namecheap.com/domains/tlds.ashx
3721722846 https://www.namecheap.com/myaccount/login.aspx
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
ID: 2116670402 ; URLs: 1 ; Server: cloudflare
 - Cache-Control
 - Cf-Cache-Status
 - Cf-Ray
 - Content-Type
 - Date
 - Expect-Ct
 - Expires
 - Last-Modified
 - Server
 - Set-Cookie
 - Strict-Transport-Security
 - Vary
 - X-Frame-Options
 - X-Xss-Protection
ID: 1712249365 ; URLs: 1 ; Server: cloudflare
 - Access-Control-Allow-Credentials
 - Access-Control-Allow-Headers
 - Access-Control-Allow-Methods
 - Cache-Control
 - Cf-Cache-Status
 - Cf-Ray
 - Content-Type
 - Date
 - Expect-Ct
 - Expires
 - Pragma
 - Server
 - Set-Cookie
 - Strict-Transport-Security
 - Vary
 - X-Frame-Options
 - X-Inst
 - X-Xss-Protection
ID: 1002554738 ; URLs: 4 ; Server: cloudflare
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
ID: 2142470851 ; URLs: 2 ; Server: cloudflare
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
 - X-Frame-Options
 - X-Inst
 - X-Xss-Protection
ID: 3926519757 ; URLs: 1 ; Server: cloudflare
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
ID: 2276290871 ; URLs: 10 ; Server: cloudflare
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
ID: 3632238580 ; URLs: 1 ; Server: cloudflare
 - Access-Control-Allow-Credentials
 - Access-Control-Allow-Headers
 - Access-Control-Allow-Methods
 - Cache-Control
 - Cf-Cache-Status
 - Cf-Ray
 - Content-Length
 - Content-Type
 - Date
 - Expect-Ct
 - Server
 - Set-Cookie
 - Strict-Transport-Security
 - X-Frame-Options
 - X-Inst
 - X-Xss-Protection
ID: 3721722846 ; URLs: 5 ; Server: cloudflare
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
ID: 1254281072 ; URLs: 7 ; Server: cloudflare
 - Cf-Cache-Status
 - Cf-Ray
 - Content-Type
 - Date
 - Expect-Ct
 - Server
 - Set-Cookie
 - Strict-Transport-Security
 - Vary
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
```

