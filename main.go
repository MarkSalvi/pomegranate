package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
)

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	client = http.Client{
		Jar: jar,
	}

}

func main() {
	request()
	createWindow()
}
