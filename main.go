package main

import (
	"almacal/auth"
	"almacal/calendar"
	"fmt"
	"os"
)

func main() {
	sessionno, menuid := auth.Login(os.Args[1], os.Args[2])
	var ical = calendar.DownloadIcalFile(sessionno, menuid)
	fmt.Println("ICAL:\n", ical)
}



