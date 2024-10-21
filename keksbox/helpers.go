package keksbox

import (
	"net/http"
	"net/url"
	"strings"
)

// returns the index where cookie is located in cookies, -1 if not found
func findCookie(cookies []*http.Cookie, cookie *http.Cookie) int {
	for i, c := range cookies {
		if c.Name == cookie.Name {
			return i
		}
	}
	return -1
}

func domainOf(u *url.URL) string {
	host := u.Host
	colon := strings.LastIndexByte(host, ':')
	if colon != -1 {
		host = host[:colon]
	}
	return strings.ToLower(host)
}

func isSubdomainOf(subdomain, domain string) bool {
	if subdomain == domain {
		return true
	}
	leftover, found := strings.CutSuffix(subdomain, domain)
	if !found {
		return false
	}
	if leftover[len(leftover)-1] != '.' {
		return false
	}
	return true
}

func domainsMatch(subdomain, domain string) bool {
	// we assume that domain is not an IP address
	if strings.HasSuffix(subdomain, ".") {
		return false
	}
	if strings.HasSuffix(domain, ".") {
		return false
	}
	return isSubdomainOf(subdomain, domain)
}

func pathsMatch(parentDir, childDir string) bool {
	parentDir, _ = strings.CutSuffix(parentDir, "/")
	childDir, _ = strings.CutSuffix(childDir, "/")
	if parentDir == childDir {
		return true
	}
	return strings.HasPrefix(childDir, parentDir)
}

