package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/kovalyov-valentin/simple-ads-server/internal/stats"
)

const (
	insertQuery = `INSERT INTO %s (ts, country, os, browser, campaign_id, requests, impressions)`
)

type writer struct {
	conn      driver.Conn
	tableName string
}

func NewClickhouseWriter(host string, port uint16, database, table, user, password string) (*writer, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", host, port)},
		Auth: clickhouse.Auth{
			Database: database,
			Username: user,
			Password: password,
		},
		Debug:           true,
		DialTimeout:     time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})

	if err != nil {
		return nil, err
	}

	return &writer{
		conn:      conn,
		tableName: table,
	}, nil
}

func (w *writer) Insert(rows stats.Rows) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	batch, err := w.conn.PrepareBatch(ctx, fmt.Sprintf(insertQuery, w.tableName))
	if err != nil {
		return err
	}

	for k, v := range rows {
		err := batch.Append(
			time.Unix(k.Timestamp, 0),
			k.Country,
			k.Os,
			k.Browser,
			k.CampaignId,
			v.Requests,
			v.Impressions,
		)

		if err != nil {
			return err
		}
	}

	return batch.Send()
}
