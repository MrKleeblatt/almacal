package auth

import (
	"almacal/keksbox"
	"almacal/logger"
	"net/http"
)

// can't let this be built by code because "," would end up URL-encrypted, but needs to be as-is
const REDIRECT_URL = "https://almaweb.uni-leipzig.de/scripts/mgrqispi.dll?APPNAME=CampusNet&PRGNAME=LOGINCHECK&ARGUMENTS=-N000000000000001,ids_mode&ids_mode=Y"

type AuthUser struct {
	NoRedirectClient *http.Client
	RedirectClient   *http.Client
	Sessionno        string
	Menuid           string
}

func Login(username, password string) (*AuthUser, error) {
	au := &AuthUser{}
	var Jar = keksbox.New()
	au.RedirectClient = &http.Client{Jar: Jar}
	au.NoRedirectClient = &http.Client{
		Jar:           Jar,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}

	err := au.getKwdCookie()
	if err != nil {
		return nil, err
	}
	logger.Debug("got kwd cookie for", username, password)
	rvtoken, err := au.getAuthCookieAndRVToken()
	if err != nil {
		return nil, err
	}
	logger.Debug("got rvtoken", rvtoken)
	err = au.postLoginForm(username, password, rvtoken)
	if err != nil {
		return nil, err
	}
	logger.Debug("did post of login form")
	redirectLocation, err := au.loginCheck()
	if err != nil {
		return nil, err
	}
	logger.Debug("got redirect location", redirectLocation)
	err = au.loginCheckRedirect(redirectLocation)
	if err != nil {
		return nil, err
	}
	logger.Debug("did the redirect and got", au.Sessionno, au.Menuid)

	if au.Sessionno != "" && au.Menuid != "" {
		return au, nil
	}
	return nil, ErrUnauthorized
}
