package main

import (
	"almacal/auth"
	"almacal/calendar"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Content-Disposition", "inline")
		credentialsArgument := r.URL.Query().Get("credentials")
	bytes, err := base64.StdEncoding.DecodeString(credentialsArgument)
	if err != nil {
		log.Fatalln(err)
	}
	credentials := strings.Split(string(bytes), ":")
	if len(credentials) != 2 {
		return
	}
	username := credentials[0]
	password := credentials[1]
	sessionno, menuid := auth.Login(username, password)
	date := time.Now()
	var result strings.Builder
	for range 12 {
		_, err := result.WriteString(calendar.IcalFile(sessionno, menuid, date.Format("Y2006M01")))
		if err != nil {
			log.Fatalln(err)
		}
		date.AddDate(0, 1, 0)
	}
	fmt.Fprint(w, strings.ReplaceAll(result.String(), "END:VCALENDAR\nBEGIN:VCALENDAR", ""))

}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8010", nil)
}
