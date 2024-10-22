package main

import (
	"almacal/keksbox"
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var jar = keksbox.New()

func printJar() {
	kekse := jar.(keksbox.Keksbox)
	for _, c := range *kekse.Entries {
		fmt.Println(c)
		fmt.Println()
	}
}
func halt() {
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
}

var redirectClient = &http.Client{
	Jar: jar,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// logRequest(req, &http.Client{Jar: jar})
		return nil
	},
}
var noRedirectClient = &http.Client{
	Jar:           jar,
	CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
}

func main() {
	sessionno, menuid := login(os.Args[1], os.Args[2])
	var ical = downloadIcalFile(sessionno, menuid)
	fmt.Println("ICAL:\n", ical)
}

func downloadIcalFile(sessionno, menuid string) string {
	reqBody := url.Values{
		"month":     {"0"},
		"week":      {"Y2024W42"},
		"APPNAME":   {"CampusNet"},
		"PRGNAME":   {"SCHEDULER_EXPORT_START"},
		"ARGUMENTS": {"sessionno,menuid,date"},
		"sessionno": {sessionno},
		"menuid":    {menuid},
		"date":      {"Y2024W42"},
	}
	payload := strings.NewReader(reqBody.Encode())
	req, err := http.NewRequest("POST", "https://almaweb.uni-leipzig.de/scripts/mgrqispi.dll", payload)
	if err != nil {
		log.Fatalln(err)
	}
	res, err := redirectClient.Do(req)
	fmt.Println(res.Status)
	return ""
}

func login(username, password string) (string, string) {
	getKwdCookie()
	rvtoken := getAuthCookieAndRVToken()
	location := postLoginForm(username, password, rvtoken)
	redirectLocation := loginCheck(location)
	println("redirectLocation", redirectLocation)
	printJar()
	halt()
	return loginCheckRedirect(redirectLocation)
}

func getKwdCookie() {
	_, err := redirectClient.Get("https://almaweb.uni-leipzig.de")
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
		"redirect_uri":  {buildRedirectUrl()},
	}
	uri.RawQuery = query.Encode()
	res, err := redirectClient.Get(uri.String())
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
	println(buildRedirectUrlForBody())
	writer := multipart.NewWriter(payload)

	_ = writer.WriteField("ReturnUrl", buildRedirectUrlForBody())
	_ = writer.WriteField("CancelUrl", "")
	_ = writer.WriteField("Username", username)
	_ = writer.WriteField("Password", password)
	_ = writer.WriteField("button", "login")
	_ = writer.WriteField("__RequestVerificationToken", rvtoken)
	_ = writer.WriteField("RememberLogin", "false")
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	req, err := http.NewRequest("POST", uri.String(), payload)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := noRedirectClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	return res.Header.Get("Location")
}

func buildReturnUrl() string {
	result, err := url.Parse("/IdentityServer/connect/authorize/callback")
	if err != nil {
		log.Fatalln(err)
	}
	query := url.Values{
		"client_id":     {"ClassicWeb"},
		"scope":         {"openid DSF email"},
		"response_mode": {"query"},
		"response_type": {"code"},
		"redirect_uri":  {buildRedirectUrl()},
	}
	result.RawQuery = query.Encode()
	return result.String()
}

func buildRedirectUrl() string {
	result, err := url.Parse("https://almaweb.uni-leipzig.de/scripts/mgrqispi.dll")
	if err != nil {
		log.Fatalln(err)
	}
	query := url.Values{
		"APPNAME":   {"CampusNet"},
		"PRGNAME":   {"LOGINCHECK"},
		"ARGUMENTS": {"-N000000000000001,ids_mode"},
		"ids_mode":  {"Y"},
	}
	result.RawQuery = query.Encode()
	return result.String()
}

func buildRedirectUrlForBody() string {
	result, err := url.Parse("https://dsf.almaweb.uni-leipzig.de/IdentityServer/connect/authorize/callback")
	if err != nil {
		log.Fatalln(err)
	}
	query := url.Values{
		"client_id":     {"ClassicWeb"},
		"scope":         {"openid DSF email"},
		"response_mode": {"query"},
		"response_type": {"code"},
		"redirect_uri":  {"https://almaweb.uni-leipzig.de/scripts/mgrqispi.dll?APPNAME=CampusNet&PRGNAME=LOGINCHECK&ARGUMENTS=-N000000000000001,ids_mode&ids_mode=Y"},
	}
	result.RawQuery = query.Encode()
	return strings.Split(result.String(), "dsf.almaweb.uni-leipzig.de")[1]
}

func loginCheck(location string) string {
	uri, err := url.Parse("https://dsf.almaweb.uni-leipzig.de" + location)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("GET", uri.String(), nil)
	logRequest(req, noRedirectClient)

	res, err := noRedirectClient.Do(req)
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
	logRequest(req, redirectClient)
	res, err := redirectClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	/*defer res.Body.Close()
	by, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	str := string(by)
	fmt.Println(str)*/
	fmt.Println("REFRESH:", res.Header.Get("REFRESH"))
	return "", ""
}

func logRequest(req *http.Request, client *http.Client) {
	ctx := context.Background()
	r := req.Clone(ctx)
	for _, c := range client.Jar.Cookies(req.URL) {
		r.AddCookie(c)
	}
	r.URL.Host = "localhost:8080"
	r.URL.Scheme = "http"
	client.Do(r)
}
