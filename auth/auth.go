package auth

import (
	"almacal/keksbox"
	"net/http"
)


var Jar = keksbox.New()
// can't let this be built by code because "," would end up URL-encrypted, but needs to be as-is
const REDIRECT_URL = "https://almaweb.uni-leipzig.de/scripts/mgrqispi.dll?APPNAME=CampusNet&PRGNAME=LOGINCHECK&ARGUMENTS=-N000000000000001,ids_mode&ids_mode=Y"

var RedirectClient = &http.Client{
	Jar:           Jar,
	// CheckRedirect: func(req *http.Request, via []*http.Request) error { return nil },
}
var NoRedirectClient = &http.Client{
	Jar:           Jar,
	CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
}


func Login(username, password string) (string, string) {
	getKwdCookie()
	rvtoken := getAuthCookieAndRVToken()
	postLoginForm(username, password, rvtoken)
	redirectLocation := loginCheck()
	return loginCheckRedirect(redirectLocation)
}
