package main

import (
	"log"

	"github.com/kovalyov-valentin/simple-ads-server/internal/ads"
	"github.com/oschwald/geoip2-golang"
)

func main() {
	geoip, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}


	s := ads.NewServer(geoip)
	if err := s.Listen(); err != nil {
		log.Fatal(err)
	}
}