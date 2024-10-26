package logger

import (
	"context"
	"net/http"
	"os"

	"github.com/withmandala/go-log"
)

var stdLogger *log.Logger
var fileLogger *log.Logger

func Init(std *os.File, filename string) {
	if std == nil {
		std = os.Stderr
	}
	stdLogger = log.New(std)
	_, debugFlag := os.LookupEnv("DEBUG")
	if debugFlag {
		stdLogger = stdLogger.WithDebug().WithTimestamp()
	}
	if filename != "" {
		file, err := os.Create(filename)
		if err != nil {
			stdLogger.Warn("could not create logging file", filename)
			return
		}
		fileLogger = log.New(file).WithDebug().WithTimestamp()
	}
}

func LogRequest(req *http.Request, client *http.Client) {
	ctx := context.Background()
	r := req.Clone(ctx)
	for _, c := range client.Jar.Cookies(req.URL) {
		r.AddCookie(c)
	}
	r.URL.Host = "localhost:8080"
	r.URL.Scheme = "http"
	client.Do(r)
}
