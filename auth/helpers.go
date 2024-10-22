package auth

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getKwdCookie() {
	_, err := RedirectClient.Get("https://almaweb.uni-leipzig.de")
	if err != nil {
		log.Fatalln(err)
	}
}
func getAuthCookieAndRVToken() string {
	uri, err := url.Parse("https://dsf.almaweb.uni-leipzig.de/IdentityServer/connect/authorize")
	if err != nil {
		log.Fatalln(err)
	}
	query := url.Values{
		"client_id":     {"ClassicWeb"},
		"scope":         {"openid DSF email"},
		"response_mode": {"query"},
		"response_type": {"code"},
		"redirect_uri":  {REDIRECT_URL},
	}
	uri.RawQuery = query.Encode()
	res, err := RedirectClient.Get(uri.String())
	if err != nil {
		log.Fatalln(err)
	}
	// scrape token from html
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
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
	return token
}

func postLoginForm(username, password, rvtoken string) string {
	uri, err := url.Parse("https://dsf.almaweb.uni-leipzig.de/IdentityServer/Account/Login")
	if err != nil {
		log.Fatalln(err)
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
		log.Fatalln(err)
	}

	req, err := http.NewRequest("POST", uri.String(), payload)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := NoRedirectClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	return res.Header.Get("Location")
}


func buildReturnUrl() string {
	result, err := url.Parse("https://dsf.almaweb.uni-leipzig.de/IdentityServer/connect/authorize/callback")
	if err != nil {
		log.Fatalln(err)
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

func loginCheck() string {
	res, err := NoRedirectClient.Get(buildReturnUrl())
	if err != nil {
		log.Fatalln(err)
	}
	return res.Header.Get("Location")
}

func loginCheckRedirect(redirectLocation string) (string, string) {
	req, err := http.NewRequest("GET", redirectLocation, nil)
	if err != nil {
		log.Fatalln(err)
	}
	res, err := RedirectClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	temp := strings.Split(res.Header.Get("REFRESH"), "ARGUMENTS=-N")[1]
	arguments := strings.Split(temp, ",-N")
	return arguments[0], arguments[1]
}


