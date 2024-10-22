package main

import (
	"almacal/auth"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	sessionno, menuid := auth.Login(os.Args[1], os.Args[2])
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
	res, err := auth.RedirectClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	by, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	str := string(by)

}

