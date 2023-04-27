package main

import (
	"log"
	"time"

	"github.com/kovalyov-valentin/simple-ads-server/internal/ads"
	"github.com/kovalyov-valentin/simple-ads-server/internal/stats"
	"github.com/kovalyov-valentin/simple-ads-server/internal/stats/clickhouse"
	met "github.com/kovalyov-valentin/simple-ads-server/internal/metrics"

	// "github.com/kovalyov-valentin/simple-ads-server/mysql"
	"github.com/oschwald/geoip2-golang"
)

func main() {
	geoip, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}


	// // MySql Writer
	// mw, err := mysql.NewMySqlWriter("127.0.0.1", 13306, "rotator", "statistics", "root", "qwerty123")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// ClickHouse Writer
	cw, err := clickhouse.NewClickhouseWriter("127.0.0.1", 19000, "rotator", "statistics", "default", "qwerty123")
	if err != nil {
		log.Fatal(err)
	}

	statsManager := stats.NewManager(cw, time.Second * 10)
	statsManager.Start()

	go func() {
		_ = met.Listen("127.0.0.1:8081")
	}()


	s := ads.NewServer(geoip, statsManager)
	if err := s.Listen(); err != nil {
		log.Fatal(err)
	}
}