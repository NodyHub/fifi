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
% cat bild.urls | fifi -v
2022/04/14 18:02:41 reading from stdin...
2022/04/14 18:02:41 Collected 19 different urls, starting analysis
2022/04/14 18:02:41 1749586943 https://a.bildstatic.de/img/bild-spielt.93c47b6.svg
2022/04/14 18:02:41 298926734 https://a.bildstatic.de/img/club-bremen.d75c8c1.svg
2022/04/14 18:02:41 1749586943 https://a.bildstatic.de/img/bild-gutscheine.360142b.svg
2022/04/14 18:02:41 1749586943 https://a.bildstatic.de/img/bild-vpn.51b780c.svg
2022/04/14 18:02:41 1749586943 https://a.bildstatic.de/img/club-aue.20b5c70.svg
2022/04/14 18:02:41 1749586943 https://a.bildstatic.de/img/club-augsburg.19aa74e.svg
2022/04/14 18:02:41 1749586943 https://a.bildstatic.de/img/club-dortmund.1940fa3.svg
2022/04/14 18:02:41 1749586943 https://a.bildstatic.de/img/club-dresden.caa901d.svg
2022/04/14 18:02:41 298926734 https://a.bildstatic.de/breakingnews/index.json
2022/04/14 18:02:42 3080683410 https://a.bildstatic.de/breakingnews
2022/04/14 18:02:42 3080683410 https://a.bildstatic.de/img
2022/04/14 18:02:42 1749586943 https://a.bildstatic.de/img/bild-deals.5e10a5e.svg
2022/04/14 18:02:42 1749586943 https://a.bildstatic.de/img/bild-vergleich.b0589b9.svg
2022/04/14 18:02:42 1749586943 https://a.bildstatic.de/img/club-darmstadt.bfca4d7.svg
2022/04/14 18:02:42 1749586943 https://a.bildstatic.de/img/club-duesseldorf.7210a4c.svg
2022/04/14 18:02:42 3080683410 https://a.bildstatic.de/
2022/04/14 18:02:42 1749586943 https://a.bildstatic.de/img/club-bielefeld.9bd0726.svg
2022/04/14 18:02:42 1749586943 https://a.bildstatic.de/img/club-bochum.0bba830.svg
2022/04/14 18:02:42 1749586943 https://a.bildstatic.de/img/bild-jobs.d0f1b16.svg

Summary:
Signature: 1749586943 ; URLs: 14
 - Accept-Ranges
 - Access-Control-Allow-Headers
 - Access-Control-Allow-Methods
 - Access-Control-Allow-Origin
 - Access-Control-Expose-Headers
 - Access-Control-Max-Age
 - Cache-Control
 - Content-Type
 - Date
 - Etag
 - Expires
 - Last-Modified
 - Server
 - Vary
-----
Urls:
-----
[200] https://a.bildstatic.de/img/bild-deals.5e10a5e.svg
[200] https://a.bildstatic.de/img/bild-gutscheine.360142b.svg
[200] https://a.bildstatic.de/img/bild-jobs.d0f1b16.svg
[200] https://a.bildstatic.de/img/bild-spielt.93c47b6.svg
[200] https://a.bildstatic.de/img/bild-vergleich.b0589b9.svg
[200] https://a.bildstatic.de/img/bild-vpn.51b780c.svg
[200] https://a.bildstatic.de/img/club-aue.20b5c70.svg
[200] https://a.bildstatic.de/img/club-augsburg.19aa74e.svg
[200] https://a.bildstatic.de/img/club-bielefeld.9bd0726.svg
[200] https://a.bildstatic.de/img/club-bochum.0bba830.svg
[200] https://a.bildstatic.de/img/club-darmstadt.bfca4d7.svg
[200] https://a.bildstatic.de/img/club-dortmund.1940fa3.svg
[200] https://a.bildstatic.de/img/club-dresden.caa901d.svg
[200] https://a.bildstatic.de/img/club-duesseldorf.7210a4c.svg
-----
Signature: 298926734 ; URLs: 2
 - Content-Length
 - Accept-Ranges
 - Access-Control-Allow-Headers
 - Access-Control-Allow-Methods
 - Access-Control-Allow-Origin
 - Access-Control-Expose-Headers
 - Access-Control-Max-Age
 - Cache-Control
 - Content-Type
 - Date
 - Etag
 - Expires
 - Last-Modified
 - Server
 - Vary
-----
Urls:
-----
[200] https://a.bildstatic.de/breakingnews/index.json
[200] https://a.bildstatic.de/img/club-bremen.d75c8c1.svg
-----
Signature: 3080683410 ; URLs: 3
 - Content-Length
 - Accept-Ranges
 - Access-Control-Allow-Headers
 - Access-Control-Allow-Methods
 - Access-Control-Allow-Origin
 - Access-Control-Expose-Headers
 - Access-Control-Max-Age
 - Cache-Control
 - Date
 - Expires
 - Server
 - Vary
-----
Urls:
-----
[404] https://a.bildstatic.de/
[404] https://a.bildstatic.de/breakingnews
[404] https://a.bildstatic.de/img
-----
```
