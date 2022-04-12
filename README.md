# Fifi

Fifi collects server response header (or the [`Server`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Server) banner) from a given list of urls and groups them.

## Background

Recently, spring boot had a wide spreaded RCE vulnerability, known as [Spring4Shell](https://portswigger.net/daily-swig/spring4shell-microsoft-cisa-warn-of-limited-in-the-wild-exploitation) ([CVE-2022-22965](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-22965)). Due to the fact that modern web application are implemented based on micro service pattern, various paths of a domain may end up on different applications/containers/CDN. To limit the attack surface system administrator, DevOps Engineers and SRE's are highly interested in limiting the available information about a service in the public.

This tool provides help to identify differences in the response headers from a given list of urls.

## Installation

```
go install github.com/NodyHub/fifi@latest
```

## Usage

```shell
% fifi -h
usage: fifi [files]
Parse urls and fetch Server banners.

Options:
[files] provide the urls in files.
  -H string
    	Host
  -X string
    	Method (default "GET")
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

```shell
% fifi -v -t 4 -w 500
2022/04/12 10:00:55 reading from stdin...

https://ap.www.namecheap.com/
https://ap.www.namecheap.com/
https://ap.www.namecheap.com/api
https://ap.www.namecheap.com/api/v1
https://ap.www.namecheap.com/api/v1/ncpl
https://ap.www.namecheap.com/JavaScriptResourceHandler.axd
[...]
https://ap.www.namecheap.com/JavaScriptResourceHandler.axd?ResourceFilePathKey=AddressResFilePath&LocaleId=en-US&VarName=AddressClientRes&ResourceType=resx&ResourceMode=1;
https://ap.www.namecheap.com/JavaScriptResourceHandler.axd?ResourceFilePathKey=DomainInfoFilePath&LocaleId=en-US&VarName=DomainInfoRes&ResourceType=resx&ResourceMode=1;
https://ap.www.namecheap.com/JavaScriptResourceHandler.axd?ResourceFilePathKey=DomainTransferResFilePath&LocaleId=en-US&VarName=DomainTransferRes&ResourceType=resx&ResourceMode=1
https://ap.www.namecheap.com/siteservices/navigationscript?fromCMS=true&fromCMSIdentity=6bdfe6b8-a8f1-45c6-81d3-b2a0dfb85786

2022/04/12 10:01:22 Collected 41 different urls, starting analysis
2022/04/12 10:01:24 3721722846 https://ap.www.namecheap.com/api/v1/ncpl/onepager
2022/04/12 10:01:25 3476117615 https://ap.www.namecheap.com/api/v1/ncpl/usermanagement/uiuser/isAdminMode
2022/04/12 10:01:26 3721722846 https://ap.www.namecheap.com/api/v1/ncpl/usermessages
2022/04/12 10:01:28 3721722846 https://ap.www.namecheap.com/Domains/DomainList
2022/04/12 10:01:31 3721722846 https://ap.www.namecheap.com/SiteServices/GetRecentMessagesJson
2022/04/12 10:01:33 3721722846 https://ap.www.namecheap.com/api/v1
[...]
2022/04/12 10:02:12 304226790 https://ap.www.namecheap.com/siteservices/navigationscript
2022/04/12 10:02:13 304226790 https://ap.www.namecheap.com/siteservices/navigationscript?_=1649748131855
2022/04/12 10:02:14 2999348701 https://ap.www.namecheap.com/api/v1/ncpl/expiringsoon/getexpiringitemcount

Summary:
ID: 304226790 ; URLs: 3 ; Server: cloudflare
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
 - X-Content-Type-Options
 - X-Frame-Options
 - X-Inst
 - X-Xss-Protection
https://ap.www.namecheap.com/siteservices/navigationscript
https://ap.www.namecheap.com/siteservices/navigationscript?_=1649748131855
https://ap.www.namecheap.com/siteservices/navigationscript?fromCMS=true&fromCMSIdentity=6bdfe6b8-a8f1-45c6-81d3-b2a0dfb85786

ID: 3721722846 ; URLs: 22 ; Server: cloudflare
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
https://ap.www.namecheap.com/
https://ap.www.namecheap.com/?_gl=1*5eunvn*_ga*ODQyMzkyNzg5LjE2NDk3NDgwOTE.*_ga_7DMJMG20P8*MTY0OTc0ODA5MC4xLjAuMTY0OTc0ODA5MC42MA..&_ga=2.76907415.360072847.1649748091-842392789.1649748091
https://ap.www.namecheap.com/Common
https://ap.www.namecheap.com/Common/ncSplitButton
https://ap.www.namecheap.com/Domains
https://ap.www.namecheap.com/Domains/DomainList
https://ap.www.namecheap.com/Domains/DomainList/DomainCategoryList
https://ap.www.namecheap.com/Domains/DomainOnly
https://ap.www.namecheap.com/Domains/GetDomainList
https://ap.www.namecheap.com/SiteServices/GetRecentMessagesJson
https://ap.www.namecheap.com/api
https://ap.www.namecheap.com/api/v1
https://ap.www.namecheap.com/api/v1/ncpl
https://ap.www.namecheap.com/api/v1/ncpl/onepager
https://ap.www.namecheap.com/api/v1/ncpl/usermanagement
https://ap.www.namecheap.com/api/v1/ncpl/usermessages
https://ap.www.namecheap.com/dashboard
https://ap.www.namecheap.com/dashboard/GetBulkModifications
https://ap.www.namecheap.com/dashboard/GetBulkModifications/
https://ap.www.namecheap.com/domains/list
https://ap.www.namecheap.com/domains/list/
https://ap.www.namecheap.com/siteservices

ID: 3476117615 ; URLs: 2 ; Server: cloudflare
 - Cache-Control
 - Cf-Cache-Status
 - Cf-Ray
 - Content-Type
 - Date
 - Expect-Ct
 - Pragma
 - Server
 - Set-Cookie
 - Strict-Transport-Security
 - Vary
 - X-Frame-Options
 - X-Xss-Protection
https://ap.www.namecheap.com/api/v1/ncpl/usermanagement/uiuser
https://ap.www.namecheap.com/api/v1/ncpl/usermanagement/uiuser/isAdminMode

ID: 2999348701 ; URLs: 5 ; Server: cloudflare
 - Cf-Cache-Status
 - Cf-Ray
 - Content-Type
 - Date
 - Expect-Ct
 - Server
 - Set-Cookie
 - Strict-Transport-Security
 - Www-Authenticate
https://ap.www.namecheap.com/api/v1/ncpl/expiringsoon/getexpiringitemcount
https://ap.www.namecheap.com/api/v1/ncpl/gatewaydomainlist/CheckSyncDomainList
https://ap.www.namecheap.com/api/v1/ncpl/gatewaydomainlist/getdomainsonly
https://ap.www.namecheap.com/api/v1/ncpl/onepager/subscription/getonepagersubscriptions
https://ap.www.namecheap.com/api/v1/ncpl/usermessages/user/getMessages

[...]

```
