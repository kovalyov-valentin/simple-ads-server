package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kovalyov-valentin/simple-ads-server/internal/stats"
)

const (
	insertQuery = `INSERT INTO %s (ts, country, os, browser, campaign_id, requests, impressions) VALUES (?, ?, ?, ?, ?, ?)`
)

type writer struct {
	db        *sql.DB
	tableName string
}

func NewMySqlWriter(host string, port uint16, database, table, user, password string) (*writer, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, database))
	if err != nil {
		return nil, err
	}

	return &writer{
		db: db,
		tableName: table,
	}, nil


}

func (w *writer) Insert(rows stats.Rows) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 30)
	defer cancel()

	tx, err := w.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, fmt.Sprintf(insertQuery, w.tableName))
	if err != nil {
		return err
	}

	for k, v := range rows {
		ts := time.Unix(k.Timestamp, 0).Format("2006-01-02 15:04:05")

		if _, err := stmt.Exec(ts, k.Country, k.Browser, k.CampaignId, v.Requests, v.Impressions); err != nil {
			return err
		} 
	}

	return tx.Commit()
}