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

//todo Add Round for 1,10,20,30,40 stacks
//todo sistemare i trade delle currency: "Aggiungere On/Off mult,gain,rounde etc per currency"
//todo sistemare in ordine alfabetico list currency
//todo aggiungere tutte le currency e fixare le tab
