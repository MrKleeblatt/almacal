package calendar

import (
	"almacal/auth"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
)

func DownloadIcalFile(sessionno, menuid string) string {
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
	res, err := auth.RedirectClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	// scrape download href from document
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var href string
	doc.Find("table tr td a").Each(func(_ int, selection *goquery.Selection) {
		href, _ = selection.Attr("href")
		href = "https://almaweb.uni-leipzig.de" + href
	})
	res, err = auth.RedirectClient.Get(href)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	// Who THE FUCK uses something different than UTF-8?!?
	decoder := charmap.ISO8859_15.NewDecoder()
	reader := transform.NewReader(res.Body, decoder)
	by, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalln(err)
	}
	return string(by)
}
