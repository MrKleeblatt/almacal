package calendar

import (
	"almacal/auth"
	"almacal/logger"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrIcalFetch = errors.New("ical file fetch error")
	ErrNoSuchFile = errors.New("no ical file available for that date")
	ErrScraping  = errors.New("web scraping error")
)

func IcalFile(au *auth.AuthUser, date string) (string, error) {
	reqBody := url.Values{
		"month":     {date},
		"week":      {"0"},
		"APPNAME":   {"CampusNet"},
		"PRGNAME":   {"SCHEDULER_EXPORT_START"},
		"ARGUMENTS": {"sessionno,menuid,date"},
		"sessionno": {au.Sessionno},
		"menuid":    {au.Menuid},
		"date":      {date},
	}
	payload := strings.NewReader(reqBody.Encode())
	req, err := http.NewRequest("POST", "https://almaweb.uni-leipzig.de/scripts/mgrqispi.dll", payload)
	if err != nil {
		logger.Fatal("error during post request creation")
		panic("unreachable")
	}
	res, err := au.RedirectClient.Do(req)
	if err != nil {
		logger.Error("error during POST request when trying to export ical file", err)
		return "", ErrIcalFetch
	}
	// scrape download href from document
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logger.Error("error during web scraping", err)
		return "", ErrScraping
	}
	var href string
	doc.Find("table tr td a").Each(func(_ int, selection *goquery.Selection) {
		href, _ = selection.Attr("href")
		href = "https://almaweb.uni-leipzig.de" + href
	})
	if href == "" {
		return "", ErrNoSuchFile
	}
	req, err = http.NewRequest("GET", href, nil)
	if err != nil {
		logger.Fatal("error during get request creation")
		panic("unreachable")
	}
	res, err = au.RedirectClient.Do(req)
	if err != nil {
		logger.Error("error during GET request when trying to fetch ical file", err)
		return "", ErrIcalFetch
	}
	defer res.Body.Close()
	// Who THE FUCK uses something different than UTF-8?!?
	decoder := charmap.ISO8859_15.NewDecoder()
	reader := transform.NewReader(res.Body, decoder)
	by, err := io.ReadAll(reader)
	if err != nil {
		logger.Error("error reading body of ical file")
		return "", ErrIcalFetch
	}
	var result strings.Builder
	for _, c := range by {
		// for some reason this character encoding uses a null character after every character
		if c == 0 {
			continue
		}
		result.WriteByte(c)
	}
	return result.String(), nil
}
