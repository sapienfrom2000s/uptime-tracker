package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
)

var (
	db  *sql.DB
	ctx context.Context
	err error
)

func main() {
	db, err = sql.Open("sqlite3", "../db/sqlite.db")
	if err != nil {
		log.Fatal("Unable to connect to DB", err)
	}
	defer db.Close()

	err = trackWebsites()
	if err != nil {
		log.Fatal("Unable to track Websites", err)
	}
}

func trackWebsites() error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var rows *sql.Rows
	rows, err = db.QueryContext(ctx, "SELECT id, url from Trackers")
	if err != nil {
		return err
	}
	defer rows.Close()

	c := cron.New()
	for rows.Next() {
		var id, url string
		if err = rows.Scan(&id, &url); err != nil {
			return err
		}

		if err = generateJob(url, id, c); err != nil {
			return err
		}
	}
	c.Start()
	time.Sleep(900 * time.Second)
	c.Stop()

	return nil
}

func generateJob(url string, id string, c *cron.Cron) error {
	c.AddFunc("@every 1s", func() {
		var res *http.Response
		if res, err = http.Get(url); err != nil {
			fmt.Printf("Unable to fetch url\n")
		}
		status := "FALSE"
		if res.StatusCode == http.StatusOK {
			status = "TRUE"
		}
		fmt.Print(url)
		_, err = db.Exec("INSERT INTO UptimeStatuses (tracker_id, up) VALUES (?, ?)", id, status)
		fmt.Printf("Added entry to DB\n %q", status)
		if err != nil {
			fmt.Printf("Unable to execute")
			return
		}
	})
	return nil
}
