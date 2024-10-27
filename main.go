package main

import (
	"almacal/auth"
	"almacal/calendar"
	"almacal/logger"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	credentialsArgument := r.URL.Query().Get("credentials")
	bytes, err := base64.StdEncoding.DecodeString(credentialsArgument)
	if err != nil {
		logger.Info("could not decode from", r.RemoteAddr, "credentials base 64:", err)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "could not decode credentials")
		logger.Info("login try but can't decode base 64")
		return
	}
	credentials := strings.Split(string(bytes), ":")
	if len(credentials) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "wrong credentials")
		logger.Info("login try with wrong formatted credentials")
		return
	}
	username := credentials[0]
	password := credentials[1]
	authuser, err := auth.Login(username, password)
	if err != nil {
		switch err {
		case auth.ErrUnauthorized:
			fmt.Fprint(w, "unauthorized")
		case auth.ErrScrapingRvToken:
			fmt.Fprint(w, "Can't web scrape AlmaWeb. Please contact the developer.\nhttps://github.com/MrKleeblatt/almacal")
		case auth.ErrAuthFetch:
			fmt.Fprint(w, "Can't fetch necessary sites from AlmaWeb. Please contact the developer.\nhttps://github.com/MrKleeblatt/almacal")
		}
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	date := time.Now()
	var result strings.Builder
	for range 12 {
		ical, err := calendar.IcalFile(authuser, date.Format("Y2006M01"))
		if err == calendar.ErrNoSuchFile {
			goto nextMonth
		}
		if err != nil {
			logger.Warn(err)
			goto nextMonth
		}
		_, err = result.WriteString(ical)
		if err != nil {
			logger.Error("could not write calendar data to string builder", err)
		}
	nextMonth:
		date = date.AddDate(0, 1, 0)
	}
	// str := strings.ReplaceAll(result.String(), "END:VCALENDAR\nBEGIN:VCALENDAR", "")
	str := result.String()
	regex := regexp.MustCompile("\nEND:VCALENDAR\r\nBEGIN:VCALENDAR\r\nVERSION:2\\.0\r\nPRODID:-\\/\\/Datenlotsen Informationssysteme GmbH\\/\\/CampusNet\\/\\/DE\r\nMETHOD:PUBLISH\r\n\r\nBEGIN:VTIMEZONE\r\nTZID:CampusNetZeit\r\nBEGIN:STANDARD\r\nDTSTART:[\\dT]+\r\nRRULE:FREQ=YEARLY;BYDAY=.+;BYMONTH=\\d+\r\nTZOFFSETFROM:[\\+\\d]+\r\nTZOFFSETTO:[\\+\\d]+\r\nEND:STANDARD\r\nBEGIN:DAYLIGHT\r\nDTSTART:[T\\d]+\r\nRRULE:FREQ=YEARLY;BYDAY=.+;BYMONTH=\\d+\r\nTZOFFSETFROM:[\\+\\d]+\r\nTZOFFSETTO:[\\+\\d]+\r\nEND:DAYLIGHT\r\nEND:VTIMEZONE")
	str = regex.ReplaceAllString(str, "")
	fmt.Fprint(w, str)
}

func main() {
	logger.Init(os.Stdout, "")
	logger.Debug("Debug messages enabled.")
	http.HandleFunc("/", handler)
	host := "localhost:8010"
	fmt.Println("listening on http://" + host)
	http.ListenAndServe(host, nil)
}
