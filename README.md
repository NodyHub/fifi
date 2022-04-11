# Fifi

Fifi fingerprints [`Server`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Server) response headers from a given list of urls and groups them.

## Background

Recently, spring boot had a wide spreaded RCE vulnerability, known as [Spring4Shell](https://portswigger.net/daily-swig/spring4shell-microsoft-cisa-warn-of-limited-in-the-wild-exploitation) ([CVE-2022-22965](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-22965)).

Due to the fact that modern web application are implemented based on micro service pattern,
various paths of a domain may end up on different applications/containers/CDN.

To limit the attack surface system administrator, DevOps Engineers and SRE's are highly interested in limiting the available information about a service in the public.

This tool provide a help to identify irregulation in the [`Server`](Fifi fingerprints [`Server`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Server) response headers from a given list of urls and groups them.

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

## Installation

```
go install github.com/NodyHub/fifi@latest
```


