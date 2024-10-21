// This package is NOT fully compliant with RFC 6265 as it sticks more to Go's
// own data structures and representations of cookies

// TODO: integrate PublicSuffixList
// TODO: thread safety through mutex's
// TODO: HttpOnly attribute

package keksbox

import (
	"fmt"
	"net/http"
	"net/url"
)

type Keksbox struct {
	Entries *[]*http.Cookie
}

func New() http.CookieJar {
	k := Keksbox{}
	factory := make([]*http.Cookie, 0)
	k.Entries = &factory
	return k
}

// Cookies implements http.CookieJar.
func (k Keksbox) Cookies(u *url.URL) (result []*http.Cookie) {
	if k.Entries == nil {
		return
	}
	// TODO:
	// Delete cookie when demanded but expired or max-age is <=0
	// Delete cookie when value is empty in SetCookies[thisCookie]
	// If a cookie has both the Max-Age and the Expires attribute, the Max-Age attribute has precedence and controls the expiration date of the cookie.
	scheme := u.Scheme
	domain := domainOf(u)
	path := u.Path
	// check all cookies from parent domains as well
	// check if cookies match security checks
	for _, c := range *k.Entries {
		// This is highly insecure because cookies with SameSite=None from this jar will be sent to _every_ domain.
		// But currently I can't see any option to restrict access. Just don't be stupid and make requests to third-party web pages with the same jar you use for other stuff.
		if c.Secure && scheme != "https" {
			continue
		}
		if !pathsMatch(c.Path, path) {
			continue
		}
		switch c.SameSite {
		case http.SameSiteStrictMode:
			if !domainsMatch(c.Domain, domain) {
				continue
			}
		case http.SameSiteLaxMode:
			// TODO
		}
		index := findCookie(result, c)
		if index == -1 {
			result = append(result, c)
			continue
		}
		if c.Domain == domain {
			// c has higher precedence
			result[index] = c
		} else if result[index].Domain == domain {
			// the existing element has higher precedence
			continue
		} else {
			// they have the same precedence, so just add the cookie
			fmt.Println("adding cookie with same name twice")
			result = append(result, c)
		}
	}
	return result
}

// SetCookies implements http.CookieJar.
func (k Keksbox) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if len(cookies) == 0 {
		return
	}
	// for simplicity and because I'm lazy as fuck
	if u.Scheme != "http" && u.Scheme != "https" {
		return
	}
	domain := domainOf(u)
	for _, c := range cookies {
		if c.Domain == "" {
			c.Domain = domain
			fmt.Println("Set domain of", c.Name, "to", domain)
		}
		// find cookie if it exists already and replace it, otherwise add it
		// cookies are the same if they share their domain, path and name
		found := false
		for i, e := range *k.Entries {
			if e.Name == c.Name && domainsMatch(e.Domain, c.Domain) && pathsMatch(e.Path, c.Path) {
				fmt.Println("overwriting cookie", e.Name)
				(*k.Entries)[i] = c
				found = true
				break
			}
		}
		if !found {
			*k.Entries = append(*k.Entries, c)
		}
	}
}
