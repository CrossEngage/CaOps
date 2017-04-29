package main

import (
	"log"
	"net/http"
	"net/url"

	"bitbucket.org/crossengage/athena/jolokia"
)

func main() {
	jolokiaURL, _ := url.Parse("http://10.2.2.10:8778/jolokia")

	client := jolokia.Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    *jolokiaURL,
	}

	resp, err := client.Version()
	if err != nil {
		log.Printf("%#v\n", err)
	}
	log.Printf("%#v\n", resp)
}
