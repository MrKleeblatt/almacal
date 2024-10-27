package auth

import (
	"almacal/logger"
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrAuthFetch       = errors.New("error fetching cookies for authentication")
	ErrScrapingRvToken = errors.New("error scraping request verification token")
	ErrUnauthorized    = errors.New("unauthorized")
)

func (au *AuthUser) getKwdCookie() error {
	req, err := http.NewRequest("GET", "https://almaweb.uni-leipzig.de", nil)
	if err != nil {
		logger.Fatal("error during get request creation")
		panic("unreachable")
	}
	_, err = au.RedirectClient.Do(req)
	if err != nil {
		logger.Error("could not get kwd cookie", err)
		return ErrAuthFetch
	}
	return nil
}

func (au *AuthUser) getAuthCookieAndRVToken() (string, error) {
	uri, err := url.Parse("https://dsf.almaweb.uni-leipzig.de/IdentityServer/connect/authorize")
	if err != nil {
		logger.Fatal("error parsing constant url", err)
		panic("unreachable")
	}
	query := url.Values{
		"client_id":     {"ClassicWeb"},
		"scope":         {"openid DSF email"},
		"response_mode": {"query"},
		"response_type": {"code"},
		"redirect_uri":  {REDIRECT_URL},
	}
	uri.RawQuery = query.Encode()
	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		logger.Fatal("error during get request creation")
		panic("unreachable")
	}
	res, err := au.RedirectClient.Do(req)
	if err != nil {
		logger.Error("error getting auth cookie and rvtoken")
		return "", ErrAuthFetch
	}
	// scrape token from html
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logger.Error("error scraping rvtoken")
		return "", ErrScrapingRvToken
	}
	var token string
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		attr, ex := s.Attr("name")
		if ex && attr == "__RequestVerificationToken" {
			attr, ex = s.Attr("value")
			if ex {
				token = attr
			}
		}
	})
	if token == "" {
		return "", ErrScrapingRvToken
	}
	return token, nil
}

func (au *AuthUser) postLoginForm(username, password, rvtoken string) error {
	uri, err := url.Parse("https://dsf.almaweb.uni-leipzig.de/IdentityServer/Account/Login")
	if err != nil {
		logger.Fatal("error parsing constant url")
		panic("unreachable")
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	err = writer.WriteField("ReturnUrl", strings.Split(buildReturnUrl(), "dsf.almaweb.uni-leipzig.de")[1])
	err = writer.WriteField("CancelUrl", "")
	err = writer.WriteField("Username", username)
	err = writer.WriteField("Password", password)
	err = writer.WriteField("button", "login")
	err = writer.WriteField("__RequestVerificationToken", rvtoken)
	err = writer.WriteField("RememberLogin", "false")
	err = writer.Close()
	if err != nil {
		logger.Fatal("can't write to local writer object")
		panic("unreachable")
	}

	req, err := http.NewRequest("POST", uri.String(), payload)
	if err != nil {
		logger.Fatal("error during post request creation", err)
		panic("unreachable")
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	logger.LogRequest(req, au.NoRedirectClient)
	_, err = au.NoRedirectClient.Do(req)
	if err != nil {
		logger.Error("can't post login form", err)
		return ErrAuthFetch
	}
	return nil
}

func buildReturnUrl() string {
	result, err := url.Parse("https://dsf.almaweb.uni-leipzig.de/IdentityServer/connect/authorize/callback")
	if err != nil {
		logger.Fatal("error parsing constant url")
		panic("unreachable")
	}
	query := url.Values{
		"client_id":     {"ClassicWeb"},
		"scope":         {"openid DSF email"},
		"response_mode": {"query"},
		"response_type": {"code"},
		"redirect_uri":  {REDIRECT_URL},
	}
	result.RawQuery = query.Encode()
	return result.String()
}

func (au *AuthUser) loginCheck() (string, error) {
	req, err := http.NewRequest("GET", buildReturnUrl(), nil)
	if err != nil {
		logger.Fatal("error during get request creation")
		panic("unreachable")
	}
	logger.LogRequest(req, au.NoRedirectClient)
	res, err := au.NoRedirectClient.Do(req)
	if err != nil {
		logger.Error("error during login check")
		return "", ErrAuthFetch
	}
	location, err := res.Location()
	if err != nil {
		return "", ErrUnauthorized
	}
	return location.String(), nil
}

func (au *AuthUser) loginCheckRedirect(redirectLocation string) error {
	req, err := http.NewRequest("GET", redirectLocation, nil)
	if err != nil {
		logger.Fatal("error during get request creation")
		panic("unreachable")
	}
	logger.LogRequest(req, au.RedirectClient)
	res, err := au.RedirectClient.Do(req)
	if err != nil {
		logger.Error("can't follow login check redirect", err)
		return ErrAuthFetch
	}
	temparr := strings.Split(res.Header.Get("REFRESH"), "ARGUMENTS=-N")
	if len(temparr) < 2 {
		logger.Debug("login check redirect got result headers:", res.Header)
		return ErrUnauthorized
	}
	temp := temparr[1]
	arguments := strings.Split(temp, ",-N")
	if len(arguments) < 2 {
		return ErrUnauthorized
	}
	au.Sessionno, au.Menuid = arguments[0], arguments[1]
	return nil
}
